package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateRateLimitRule(rlr *model.RateLimitRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rateLimitRules[rlr.ID] = rlr
	return nil
}

func (s *MemoryStore) GetRateLimitRule(id string) (*model.RateLimitRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rlr, ok := s.rateLimitRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return rlr, nil
}

func (s *MemoryStore) ListRateLimitRules() []*model.RateLimitRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RateLimitRule, 0, len(s.rateLimitRules))
	for _, rlr := range s.rateLimitRules {
		list = append(list, rlr)
	}
	return list
}

func (s *MemoryStore) UpdateRateLimitRule(rlr *model.RateLimitRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rateLimitRules[rlr.ID]; !ok {
		return ErrNotFound
	}
	s.rateLimitRules[rlr.ID] = rlr
	return nil
}

func (s *MemoryStore) DeleteRateLimitRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rateLimitRules[id]; !ok {
		return ErrNotFound
	}
	delete(s.rateLimitRules, id)
	return nil
}
