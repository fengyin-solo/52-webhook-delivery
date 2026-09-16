package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/pkg/idgen"
)

func (s *Service) CreateRateLimitRule(input model.RateLimitRule) (*model.RateLimitRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	rlr := &model.RateLimitRule{
		ID:          idgen.Hex(),
		Name:        input.Name,
		EndpointID:  input.EndpointID,
		EventType:   input.EventType,
		MaxRequests: input.MaxRequests,
		WindowMs:    input.WindowMs,
		Status:      input.Status,
		CreatedAt:   time.Now(),
	}
	if rlr.Status == "" {
		rlr.Status = model.RateLimitRuleActive
	}
	if err := s.store.CreateRateLimitRule(rlr); err != nil {
		return nil, err
	}
	return rlr, nil
}

func (s *Service) GetRateLimitRule(id string) (*model.RateLimitRule, error) {
	return s.store.GetRateLimitRule(id)
}

func (s *Service) ListRateLimitRules(filter model.RateLimitRuleFilter, page, size int) ([]*model.RateLimitRule, int, error) {
	all := s.store.ListRateLimitRules()
	matched := make([]*model.RateLimitRule, 0, len(all))
	for _, rlr := range all {
		if filter.Match(rlr) {
			matched = append(matched, rlr)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RateLimitRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRateLimitRule(id string, input model.RateLimitRule) (*model.RateLimitRule, error) {
	existing, err := s.store.GetRateLimitRule(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.EndpointID = input.EndpointID
	existing.EventType = input.EventType
	existing.MaxRequests = input.MaxRequests
	existing.WindowMs = input.WindowMs
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateRateLimitRule(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteRateLimitRule(id string) error {
	return s.store.DeleteRateLimitRule(id)
}
