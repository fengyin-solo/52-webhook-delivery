package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/internal/store"
	"webhook/pkg/idgen"
)

func (s *Service) CreateEndpointHealth(input model.EndpointHealth) (*model.EndpointHealth, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetEndpoint(input.EndpointID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("endpoint_id", "端点不存在")
		}
		return nil, err
	}
	eh := &model.EndpointHealth{
		ID:           idgen.Hex(),
		EndpointID:   input.EndpointID,
		Status:       input.Status,
		ResponseTime: input.ResponseTime,
		LastChecked:  input.LastChecked,
		ErrorMessage: input.ErrorMessage,
		CreatedAt:    time.Now(),
	}
	if eh.Status == "" {
		eh.Status = model.EndpointHealthUnknown
	}
	if eh.LastChecked.IsZero() {
		eh.LastChecked = time.Now()
	}
	if err := s.store.CreateEndpointHealth(eh); err != nil {
		return nil, err
	}
	return eh, nil
}

func (s *Service) GetEndpointHealth(id string) (*model.EndpointHealth, error) {
	return s.store.GetEndpointHealth(id)
}

func (s *Service) ListEndpointHealths(filter model.EndpointHealthFilter, page, size int) ([]*model.EndpointHealth, int, error) {
	all := s.store.ListEndpointHealths()
	matched := make([]*model.EndpointHealth, 0, len(all))
	for _, eh := range all {
		if filter.Match(eh) {
			matched = append(matched, eh)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].LastChecked.After(matched[j].LastChecked)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.EndpointHealth{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateEndpointHealth(id string, input model.EndpointHealth) (*model.EndpointHealth, error) {
	existing, err := s.store.GetEndpointHealth(id)
	if err != nil {
		return nil, err
	}
	existing.Status = input.Status
	existing.ResponseTime = input.ResponseTime
	existing.LastChecked = input.LastChecked
	existing.ErrorMessage = input.ErrorMessage
	if existing.LastChecked.IsZero() {
		existing.LastChecked = time.Now()
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateEndpointHealth(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteEndpointHealth(id string) error {
	return s.store.DeleteEndpointHealth(id)
}
