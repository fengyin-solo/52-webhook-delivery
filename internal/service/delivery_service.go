package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/internal/store"
	"webhook/pkg/idgen"
)

func (s *Service) CreateDelivery(input model.Delivery) (*model.Delivery, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetEvent(input.EventID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("event_id", "事件不存在")
		}
		return nil, err
	}
	if _, err := s.store.GetEndpoint(input.EndpointID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("endpoint_id", "端点不存在")
		}
		return nil, err
	}
	d := &model.Delivery{
		ID:         idgen.Hex(),
		EventID:    input.EventID,
		EndpointID: input.EndpointID,
		Status:     input.Status,
		Attempts:   input.Attempts,
		CreatedAt:  time.Now(),
	}
	if d.Status == "" {
		d.Status = model.DeliveryPending
	}
	if err := s.store.CreateDelivery(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) GetDelivery(id string) (*model.Delivery, error) {
	return s.store.GetDelivery(id)
}

func (s *Service) ListDeliveries(filter model.DeliveryFilter, page, size int) ([]*model.Delivery, int, error) {
	all := s.store.ListDeliveries()
	matched := make([]*model.Delivery, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Delivery{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateDelivery(id string, input model.Delivery) (*model.Delivery, error) {
	existing, err := s.store.GetDelivery(id)
	if err != nil {
		return nil, err
	}
	if input.Status != "" && input.Status != existing.Status {
		if !model.DeliveryCanTransition(existing.Status, input.Status) {
			return nil, model.NewValidationError("status", "投递状态不允许此转换")
		}
		existing.Status = input.Status
	}
	existing.LastResponseCode = input.LastResponseCode
	existing.LastResponseBody = input.LastResponseBody
	if input.NextRetryAt != nil {
		existing.NextRetryAt = input.NextRetryAt
	}
	if input.DeliveredAt != nil {
		existing.DeliveredAt = input.DeliveredAt
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateDelivery(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteDelivery(id string) error {
	return s.store.DeleteDelivery(id)
}

func (s *Service) RetryDelivery(id string) (*model.Delivery, error) {
	d, err := s.store.GetDelivery(id)
	if err != nil {
		return nil, err
	}
	if d.Status != model.DeliveryFailed {
		return nil, model.NewValidationError("status", "只有失败状态的投递才能重试")
	}
	d.Status = model.DeliveryRetrying
	d.Attempts = 0
	d.NextRetryAt = nil
	now := time.Now()
	d.NextRetryAt = &now
	if err := s.store.UpdateDelivery(d); err != nil {
		return nil, err
	}
	return d, nil
}
