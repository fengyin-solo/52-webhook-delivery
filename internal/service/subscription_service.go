package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/internal/store"
	"webhook/pkg/idgen"
)

func (s *Service) CreateSubscription(input model.Subscription) (*model.Subscription, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetEndpoint(input.EndpointID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("endpoint_id", "端点不存在")
		}
		return nil, err
	}
	if input.FilterID != "" {
		if _, err := s.store.GetFilter(input.FilterID); err != nil {
			if err == store.ErrNotFound {
				return nil, model.NewValidationError("filter_id", "过滤器不存在")
			}
			return nil, err
		}
	}
	sub := &model.Subscription{
		ID:         idgen.Hex(),
		EndpointID: input.EndpointID,
		EventType:  input.EventType,
		FilterID:   input.FilterID,
		Status:     input.Status,
		CreatedAt:  time.Now(),
	}
	if sub.Status == "" {
		sub.Status = model.SubscriptionActive
	}
	if err := s.store.CreateSubscription(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *Service) GetSubscription(id string) (*model.Subscription, error) {
	return s.store.GetSubscription(id)
}

func (s *Service) ListSubscriptions(filter model.SubscriptionFilter, page, size int) ([]*model.Subscription, int, error) {
	all := s.store.ListSubscriptions()
	matched := make([]*model.Subscription, 0, len(all))
	for _, sub := range all {
		if filter.Match(sub) {
			matched = append(matched, sub)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Subscription{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateSubscription(id string, input model.Subscription) (*model.Subscription, error) {
	existing, err := s.store.GetSubscription(id)
	if err != nil {
		return nil, err
	}
	existing.EventType = input.EventType
	existing.FilterID = input.FilterID
	if input.Status != "" && input.Status != existing.Status {
		if !model.SubscriptionCanTransition(existing.Status, input.Status) {
			return nil, model.NewValidationError("status", "订阅状态不允许此转换")
		}
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if existing.FilterID != "" {
		if _, err := s.store.GetFilter(existing.FilterID); err != nil {
			if err == store.ErrNotFound {
				return nil, model.NewValidationError("filter_id", "过滤器不存在")
			}
			return nil, err
		}
	}
	if err := s.store.UpdateSubscription(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteSubscription(id string) error {
	return s.store.DeleteSubscription(id)
}
