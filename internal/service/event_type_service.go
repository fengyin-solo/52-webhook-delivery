package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/pkg/idgen"
)

func (s *Service) CreateEventType(input model.EventType) (*model.EventType, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	et := &model.EventType{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Description: input.Description,
		Schema:      input.Schema,
		Status:      input.Status,
		CreatedAt:   time.Now(),
	}
	if et.Status == "" {
		et.Status = model.EventTypeActive
	}
	if err := s.store.CreateEventType(et); err != nil {
		return nil, err
	}
	return et, nil
}

func (s *Service) GetEventType(id string) (*model.EventType, error) {
	return s.store.GetEventType(id)
}

func (s *Service) ListEventTypes(filter model.EventTypeFilter, page, size int) ([]*model.EventType, int, error) {
	all := s.store.ListEventTypes()
	matched := make([]*model.EventType, 0, len(all))
	for _, et := range all {
		if filter.Match(et) {
			matched = append(matched, et)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.EventType{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateEventType(id string, input model.EventType) (*model.EventType, error) {
	existing, err := s.store.GetEventType(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Description = input.Description
	existing.Schema = input.Schema
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateEventType(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteEventType(id string) error {
	return s.store.DeleteEventType(id)
}
