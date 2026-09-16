package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/pkg/idgen"
)

func (s *Service) CreateEndpoint(input model.Endpoint) (*model.Endpoint, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	e := &model.Endpoint{
		ID:          idgen.Hex(),
		URL:         input.URL,
		Secret:      input.Secret,
		Status:      input.Status,
		MaxRetries:  input.MaxRetries,
		TimeoutMs:   input.TimeoutMs,
		Description: input.Description,
		CreatedAt:   time.Now(),
	}
	if e.Status == "" {
		e.Status = model.EndpointActive
	}
	if err := s.store.CreateEndpoint(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) GetEndpoint(id string) (*model.Endpoint, error) {
	return s.store.GetEndpoint(id)
}

func (s *Service) ListEndpoints(filter model.EndpointFilter, page, size int) ([]*model.Endpoint, int, error) {
	all := s.store.ListEndpoints()
	matched := make([]*model.Endpoint, 0, len(all))
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
		return []*model.Endpoint{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateEndpoint(id string, input model.Endpoint) (*model.Endpoint, error) {
	existing, err := s.store.GetEndpoint(id)
	if err != nil {
		return nil, err
	}
	existing.URL = input.URL
	existing.Secret = input.Secret
	existing.Description = input.Description
	existing.MaxRetries = input.MaxRetries
	existing.TimeoutMs = input.TimeoutMs
	if input.Status != "" && input.Status != existing.Status {
		if !model.EndpointCanTransition(existing.Status, input.Status) {
			return nil, model.NewValidationError("status", "端点状态不允许此转换")
		}
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateEndpoint(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteEndpoint(id string) error {
	return s.store.DeleteEndpoint(id)
}
