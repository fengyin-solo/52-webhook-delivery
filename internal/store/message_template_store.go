package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateMessageTemplate(mt *model.MessageTemplate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messageTemplates[mt.ID] = mt
	return nil
}

func (s *MemoryStore) GetMessageTemplate(id string) (*model.MessageTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	mt, ok := s.messageTemplates[id]
	if !ok {
		return nil, ErrNotFound
	}
	return mt, nil
}

func (s *MemoryStore) ListMessageTemplates() []*model.MessageTemplate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.MessageTemplate, 0, len(s.messageTemplates))
	for _, mt := range s.messageTemplates {
		list = append(list, mt)
	}
	return list
}

func (s *MemoryStore) ListMessageTemplatesByEventType(eventType string) []*model.MessageTemplate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.MessageTemplate, 0)
	for _, mt := range s.messageTemplates {
		if mt.EventType == eventType {
			list = append(list, mt)
		}
	}
	return list
}

func (s *MemoryStore) UpdateMessageTemplate(mt *model.MessageTemplate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.messageTemplates[mt.ID]; !ok {
		return ErrNotFound
	}
	s.messageTemplates[mt.ID] = mt
	return nil
}

func (s *MemoryStore) DeleteMessageTemplate(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.messageTemplates[id]; !ok {
		return ErrNotFound
	}
	delete(s.messageTemplates, id)
	return nil
}
