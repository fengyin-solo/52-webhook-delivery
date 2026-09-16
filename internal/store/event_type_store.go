package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateEventType(et *model.EventType) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.eventTypes {
		if exist.Name == et.Name {
			return ErrConflict
		}
	}
	s.eventTypes[et.ID] = et
	return nil
}

func (s *MemoryStore) GetEventType(id string) (*model.EventType, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	et, ok := s.eventTypes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return et, nil
}

func (s *MemoryStore) GetEventTypeByName(name string) (*model.EventType, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, et := range s.eventTypes {
		if et.Name == name {
			return et, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListEventTypes() []*model.EventType {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.EventType, 0, len(s.eventTypes))
	for _, et := range s.eventTypes {
		list = append(list, et)
	}
	return list
}

func (s *MemoryStore) UpdateEventType(et *model.EventType) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.eventTypes[et.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.eventTypes {
		if exist.ID != et.ID && exist.Name == et.Name {
			return ErrConflict
		}
	}
	s.eventTypes[et.ID] = et
	return nil
}

func (s *MemoryStore) DeleteEventType(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.eventTypes[id]; !ok {
		return ErrNotFound
	}
	delete(s.eventTypes, id)
	return nil
}
