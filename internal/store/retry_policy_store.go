package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateRetryPolicy(rp *model.RetryPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.retryPolicies {
		if exist.Name == rp.Name {
			return ErrConflict
		}
	}
	s.retryPolicies[rp.ID] = rp
	return nil
}

func (s *MemoryStore) GetRetryPolicy(id string) (*model.RetryPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rp, ok := s.retryPolicies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return rp, nil
}

func (s *MemoryStore) ListRetryPolicies() []*model.RetryPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RetryPolicy, 0, len(s.retryPolicies))
	for _, rp := range s.retryPolicies {
		list = append(list, rp)
	}
	return list
}

func (s *MemoryStore) UpdateRetryPolicy(rp *model.RetryPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.retryPolicies[rp.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.retryPolicies {
		if exist.ID != rp.ID && exist.Name == rp.Name {
			return ErrConflict
		}
	}
	s.retryPolicies[rp.ID] = rp
	return nil
}

func (s *MemoryStore) DeleteRetryPolicy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.retryPolicies[id]; !ok {
		return ErrNotFound
	}
	delete(s.retryPolicies, id)
	return nil
}
