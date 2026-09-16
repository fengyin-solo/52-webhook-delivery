package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateEndpointHealth(eh *model.EndpointHealth) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.endpointHealths[eh.ID] = eh
	return nil
}

func (s *MemoryStore) GetEndpointHealth(id string) (*model.EndpointHealth, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	eh, ok := s.endpointHealths[id]
	if !ok {
		return nil, ErrNotFound
	}
	return eh, nil
}

func (s *MemoryStore) ListEndpointHealths() []*model.EndpointHealth {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.EndpointHealth, 0, len(s.endpointHealths))
	for _, eh := range s.endpointHealths {
		list = append(list, eh)
	}
	return list
}

func (s *MemoryStore) ListEndpointHealthsByEndpointID(endpointID string) []*model.EndpointHealth {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.EndpointHealth, 0)
	for _, eh := range s.endpointHealths {
		if eh.EndpointID == endpointID {
			list = append(list, eh)
		}
	}
	return list
}

func (s *MemoryStore) UpdateEndpointHealth(eh *model.EndpointHealth) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.endpointHealths[eh.ID]; !ok {
		return ErrNotFound
	}
	s.endpointHealths[eh.ID] = eh
	return nil
}

func (s *MemoryStore) DeleteEndpointHealth(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.endpointHealths[id]; !ok {
		return ErrNotFound
	}
	delete(s.endpointHealths, id)
	return nil
}
