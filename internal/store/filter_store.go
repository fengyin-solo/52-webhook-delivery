package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateFilter(f *model.Filter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.filters {
		if exist.Name == f.Name {
			return ErrConflict
		}
	}
	s.filters[f.ID] = f
	return nil
}

func (s *MemoryStore) GetFilter(id string) (*model.Filter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.filters[id]
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

func (s *MemoryStore) ListFilters() []*model.Filter {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Filter, 0, len(s.filters))
	for _, f := range s.filters {
		list = append(list, f)
	}
	return list
}

func (s *MemoryStore) UpdateFilter(f *model.Filter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.filters[f.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.filters {
		if exist.ID != f.ID && exist.Name == f.Name {
			return ErrConflict
		}
	}
	s.filters[f.ID] = f
	return nil
}

func (s *MemoryStore) DeleteFilter(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.filters[id]; !ok {
		return ErrNotFound
	}
	delete(s.filters, id)
	return nil
}
