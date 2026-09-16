package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/internal/store"
	"webhook/pkg/idgen"
)

func (s *Service) CreateEvent(input model.Event) (*model.Event, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetEventTypeByName(input.Type); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("type", "事件类型不存在")
		}
		return nil, err
	}
	e := &model.Event{
		ID:             idgen.Hex(),
		Type:           input.Type,
		Source:         input.Source,
		Payload:        input.Payload,
		IdempotencyKey: input.IdempotencyKey,
		Status:         input.Status,
		CreatedAt:      time.Now(),
	}
	if e.Status == "" {
		e.Status = model.EventCreated
	}
	if err := s.store.CreateEvent(e); err != nil {
		return nil, err
	}
	s.dispatchEvent(e)
	return e, nil
}

func (s *Service) GetEvent(id string) (*model.Event, error) {
	return s.store.GetEvent(id)
}

func (s *Service) ListEvents(filter model.EventFilter, page, size int) ([]*model.Event, int, error) {
	all := s.store.ListEvents()
	matched := make([]*model.Event, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Event{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateEvent(id string, input model.Event) (*model.Event, error) {
	existing, err := s.store.GetEvent(id)
	if err != nil {
		return nil, err
	}
	if input.Status != "" && input.Status != existing.Status {
		if !model.EventCanTransition(existing.Status, input.Status) {
			return nil, model.NewValidationError("status", "事件状态不允许此转换")
		}
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateEvent(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteEvent(id string) error {
	return s.store.DeleteEvent(id)
}

func (s *Service) dispatchEvent(e *model.Event) {
	subscriptions := s.store.ListSubscriptionsByEventType(e.Type)
	for _, sub := range subscriptions {
		if sub.Status != model.SubscriptionActive {
			continue
		}
		ep, err := s.store.GetEndpoint(sub.EndpointID)
		if err != nil || ep.Status != model.EndpointActive {
			continue
		}
		if sub.FilterID != "" {
			f, err := s.store.GetFilter(sub.FilterID)
			if err == nil && !f.MatchPayload(e.Payload) {
				continue
			}
		}
		d := &model.Delivery{
			ID:         idgen.Hex(),
			EventID:    e.ID,
			EndpointID: sub.EndpointID,
			Status:     model.DeliveryPending,
			Attempts:   0,
			CreatedAt:  time.Now(),
		}
		_ = s.store.CreateDelivery(d)
	}
	if e.Status == model.EventCreated {
		e.Status = model.EventDispatched
		_ = s.store.UpdateEvent(e)
	}
}
