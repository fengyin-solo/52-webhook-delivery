package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateCircuitBreaker(cb *model.CircuitBreaker) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.circuitBreakers {
		if exist.EndpointID == cb.EndpointID {
			return ErrConflict
		}
	}
	s.circuitBreakers[cb.ID] = cb
	return nil
}

func (s *MemoryStore) GetCircuitBreaker(id string) (*model.CircuitBreaker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cb, ok := s.circuitBreakers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cb, nil
}

func (s *MemoryStore) GetCircuitBreakerByEndpointID(endpointID string) (*model.CircuitBreaker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, cb := range s.circuitBreakers {
		if cb.EndpointID == endpointID {
			return cb, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListCircuitBreakers() []*model.CircuitBreaker {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.CircuitBreaker, 0, len(s.circuitBreakers))
	for _, cb := range s.circuitBreakers {
		list = append(list, cb)
	}
	return list
}

func (s *MemoryStore) UpdateCircuitBreaker(cb *model.CircuitBreaker) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.circuitBreakers[cb.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.circuitBreakers {
		if exist.ID != cb.ID && exist.EndpointID == cb.EndpointID {
			return ErrConflict
		}
	}
	s.circuitBreakers[cb.ID] = cb
	return nil
}

func (s *MemoryStore) DeleteCircuitBreaker(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.circuitBreakers[id]; !ok {
		return ErrNotFound
	}
	delete(s.circuitBreakers, id)
	return nil
}
