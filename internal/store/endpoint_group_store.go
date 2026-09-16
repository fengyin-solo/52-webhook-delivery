package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateEndpointGroup(eg *model.EndpointGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.endpointGroups {
		if exist.Name == eg.Name {
			return ErrConflict
		}
	}
	s.endpointGroups[eg.ID] = eg
	return nil
}

func (s *MemoryStore) GetEndpointGroup(id string) (*model.EndpointGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	eg, ok := s.endpointGroups[id]
	if !ok {
		return nil, ErrNotFound
	}
	return eg, nil
}

func (s *MemoryStore) GetEndpointGroupByName(name string) (*model.EndpointGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, eg := range s.endpointGroups {
		if eg.Name == name {
			return eg, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListEndpointGroups() []*model.EndpointGroup {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.EndpointGroup, 0, len(s.endpointGroups))
	for _, eg := range s.endpointGroups {
		list = append(list, eg)
	}
	return list
}

func (s *MemoryStore) UpdateEndpointGroup(eg *model.EndpointGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.endpointGroups[eg.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.endpointGroups {
		if exist.ID != eg.ID && exist.Name == eg.Name {
			return ErrConflict
		}
	}
	s.endpointGroups[eg.ID] = eg
	return nil
}

func (s *MemoryStore) DeleteEndpointGroup(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.endpointGroups[id]; !ok {
		return ErrNotFound
	}
	delete(s.endpointGroups, id)
	return nil
}
