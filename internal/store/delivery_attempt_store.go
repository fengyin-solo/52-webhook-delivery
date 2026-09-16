package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateDeliveryAttempt(da *model.DeliveryAttempt) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deliveryAttempts[da.ID] = da
	return nil
}

func (s *MemoryStore) GetDeliveryAttempt(id string) (*model.DeliveryAttempt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	da, ok := s.deliveryAttempts[id]
	if !ok {
		return nil, ErrNotFound
	}
	return da, nil
}

func (s *MemoryStore) ListDeliveryAttempts() []*model.DeliveryAttempt {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.DeliveryAttempt, 0, len(s.deliveryAttempts))
	for _, da := range s.deliveryAttempts {
		list = append(list, da)
	}
	return list
}

func (s *MemoryStore) ListDeliveryAttemptsByDeliveryID(deliveryID string) []*model.DeliveryAttempt {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.DeliveryAttempt, 0)
	for _, da := range s.deliveryAttempts {
		if da.DeliveryID == deliveryID {
			list = append(list, da)
		}
	}
	return list
}

func (s *MemoryStore) UpdateDeliveryAttempt(da *model.DeliveryAttempt) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.deliveryAttempts[da.ID]; !ok {
		return ErrNotFound
	}
	s.deliveryAttempts[da.ID] = da
	return nil
}

func (s *MemoryStore) DeleteDeliveryAttempt(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.deliveryAttempts[id]; !ok {
		return ErrNotFound
	}
	delete(s.deliveryAttempts, id)
	return nil
}
