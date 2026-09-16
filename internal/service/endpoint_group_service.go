package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/pkg/idgen"
)

func (s *Service) CreateEndpointGroup(input model.EndpointGroup) (*model.EndpointGroup, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	eg := &model.EndpointGroup{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Description: input.Description,
		Status:      input.Status,
		CreatedAt:   time.Now(),
	}
	if eg.Status == "" {
		eg.Status = model.EndpointGroupActive
	}
	if err := s.store.CreateEndpointGroup(eg); err != nil {
		return nil, err
	}
	return eg, nil
}

func (s *Service) GetEndpointGroup(id string) (*model.EndpointGroup, error) {
	return s.store.GetEndpointGroup(id)
}

func (s *Service) ListEndpointGroups(filter model.EndpointGroupFilter, page, size int) ([]*model.EndpointGroup, int, error) {
	all := s.store.ListEndpointGroups()
	matched := make([]*model.EndpointGroup, 0, len(all))
	for _, eg := range all {
		if filter.Match(eg) {
			matched = append(matched, eg)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.EndpointGroup{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateEndpointGroup(id string, input model.EndpointGroup) (*model.EndpointGroup, error) {
	existing, err := s.store.GetEndpointGroup(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Description = input.Description
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateEndpointGroup(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteEndpointGroup(id string) error {
	return s.store.DeleteEndpointGroup(id)
}
