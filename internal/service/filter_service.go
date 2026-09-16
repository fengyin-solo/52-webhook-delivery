package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/pkg/idgen"
)

func (s *Service) CreateFilter(input model.Filter) (*model.Filter, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	f := &model.Filter{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Field:     input.Field,
		Operator:  input.Operator,
		Value:     input.Value,
		Status:    input.Status,
		CreatedAt: time.Now(),
	}
	if f.Status == "" {
		f.Status = model.FilterActive
	}
	if err := s.store.CreateFilter(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Service) GetFilter(id string) (*model.Filter, error) {
	return s.store.GetFilter(id)
}

func (s *Service) ListFilters(filter model.FilterFilter, page, size int) ([]*model.Filter, int, error) {
	all := s.store.ListFilters()
	matched := make([]*model.Filter, 0, len(all))
	for _, f := range all {
		if filter.Match(f) {
			matched = append(matched, f)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Filter{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateFilter(id string, input model.Filter) (*model.Filter, error) {
	existing, err := s.store.GetFilter(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Field = input.Field
	existing.Operator = input.Operator
	existing.Value = input.Value
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateFilter(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteFilter(id string) error {
	return s.store.DeleteFilter(id)
}
