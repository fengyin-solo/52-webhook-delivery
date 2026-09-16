package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/pkg/idgen"
)

func (s *Service) CreateRetryPolicy(input model.RetryPolicy) (*model.RetryPolicy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	rp := &model.RetryPolicy{
		ID:             idgen.Hex(),
		Name:           input.Name,
		MaxAttempts:    input.MaxAttempts,
		BackoffType:    input.BackoffType,
		InitialDelayMs: input.InitialDelayMs,
		MaxDelayMs:     input.MaxDelayMs,
		Status:         input.Status,
		CreatedAt:      time.Now(),
	}
	if rp.Status == "" {
		rp.Status = model.RetryPolicyActive
	}
	if err := s.store.CreateRetryPolicy(rp); err != nil {
		return nil, err
	}
	return rp, nil
}

func (s *Service) GetRetryPolicy(id string) (*model.RetryPolicy, error) {
	return s.store.GetRetryPolicy(id)
}

func (s *Service) ListRetryPolicies(filter model.RetryPolicyFilter, page, size int) ([]*model.RetryPolicy, int, error) {
	all := s.store.ListRetryPolicies()
	matched := make([]*model.RetryPolicy, 0, len(all))
	for _, rp := range all {
		if filter.Match(rp) {
			matched = append(matched, rp)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RetryPolicy{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRetryPolicy(id string, input model.RetryPolicy) (*model.RetryPolicy, error) {
	existing, err := s.store.GetRetryPolicy(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.MaxAttempts = input.MaxAttempts
	existing.BackoffType = input.BackoffType
	existing.InitialDelayMs = input.InitialDelayMs
	existing.MaxDelayMs = input.MaxDelayMs
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateRetryPolicy(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteRetryPolicy(id string) error {
	return s.store.DeleteRetryPolicy(id)
}
