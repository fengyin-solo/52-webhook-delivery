package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateEvent(e *model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e.IdempotencyKey != "" {
		for _, exist := range s.events {
			if exist.IdempotencyKey == e.IdempotencyKey {
				return ErrConflict
			}
		}
	}
	s.events[e.ID] = e
	return nil
}

func (s *MemoryStore) GetEvent(id string) (*model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.events[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) ListEvents() []*model.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Event, 0, len(s.events))
	for _, e := range s.events {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) ListEventsByType(eventType string) []*model.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Event, 0)
	for _, e := range s.events {
		if e.Type == eventType {
			list = append(list, e)
		}
	}
	return list
}

func (s *MemoryStore) UpdateEvent(e *model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[e.ID]; !ok {
		return ErrNotFound
	}
	s.events[e.ID] = e
	return nil
}

func (s *MemoryStore) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[id]; !ok {
		return ErrNotFound
	}
	delete(s.events, id)
	return nil
}
