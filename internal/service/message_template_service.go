package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/internal/store"
	"webhook/pkg/idgen"
)

func (s *Service) CreateMessageTemplate(input model.MessageTemplate) (*model.MessageTemplate, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetEventTypeByName(input.EventType); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("event_type", "事件类型不存在")
		}
		return nil, err
	}
	mt := &model.MessageTemplate{
		ID:          idgen.Hex(),
		Name:        input.Name,
		EventType:   input.EventType,
		ContentType: input.ContentType,
		Body:        input.Body,
		Status:      input.Status,
		CreatedAt:   time.Now(),
	}
	if mt.Status == "" {
		mt.Status = model.TemplateActive
	}
	if err := s.store.CreateMessageTemplate(mt); err != nil {
		return nil, err
	}
	return mt, nil
}

func (s *Service) GetMessageTemplate(id string) (*model.MessageTemplate, error) {
	return s.store.GetMessageTemplate(id)
}

func (s *Service) ListMessageTemplates(filter model.MessageTemplateFilter, page, size int) ([]*model.MessageTemplate, int, error) {
	all := s.store.ListMessageTemplates()
	matched := make([]*model.MessageTemplate, 0, len(all))
	for _, mt := range all {
		if filter.Match(mt) {
			matched = append(matched, mt)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.MessageTemplate{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateMessageTemplate(id string, input model.MessageTemplate) (*model.MessageTemplate, error) {
	existing, err := s.store.GetMessageTemplate(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.EventType = input.EventType
	existing.ContentType = input.ContentType
	existing.Body = input.Body
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if existing.EventType != "" {
		if _, err := s.store.GetEventTypeByName(existing.EventType); err != nil {
			if err == store.ErrNotFound {
				return nil, model.NewValidationError("event_type", "事件类型不存在")
			}
			return nil, err
		}
	}
	if err := s.store.UpdateMessageTemplate(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteMessageTemplate(id string) error {
	return s.store.DeleteMessageTemplate(id)
}
