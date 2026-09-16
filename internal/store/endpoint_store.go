package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateEndpoint(e *model.Endpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.endpoints {
		if exist.URL == e.URL {
			return ErrConflict
		}
	}
	s.endpoints[e.ID] = e
	return nil
}

func (s *MemoryStore) GetEndpoint(id string) (*model.Endpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.endpoints[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) GetEndpointByURL(url string) (*model.Endpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.endpoints {
		if e.URL == url {
			return e, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListEndpoints() []*model.Endpoint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Endpoint, 0, len(s.endpoints))
	for _, e := range s.endpoints {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) UpdateEndpoint(e *model.Endpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.endpoints[e.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.endpoints {
		if exist.ID != e.ID && exist.URL == e.URL {
			return ErrConflict
		}
	}
	s.endpoints[e.ID] = e
	return nil
}

func (s *MemoryStore) DeleteEndpoint(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.endpoints[id]; !ok {
		return ErrNotFound
	}
	delete(s.endpoints, id)
	return nil
}
