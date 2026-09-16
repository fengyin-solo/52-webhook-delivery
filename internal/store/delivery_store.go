package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateDelivery(d *model.Delivery) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deliveries[d.ID] = d
	return nil
}

func (s *MemoryStore) GetDelivery(id string) (*model.Delivery, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.deliveries[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) ListDeliveries() []*model.Delivery {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Delivery, 0, len(s.deliveries))
	for _, d := range s.deliveries {
		list = append(list, d)
	}
	return list
}

func (s *MemoryStore) ListDeliveriesByEventID(eventID string) []*model.Delivery {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Delivery, 0)
	for _, d := range s.deliveries {
		if d.EventID == eventID {
			list = append(list, d)
		}
	}
	return list
}

func (s *MemoryStore) ListDeliveriesByEndpointID(endpointID string) []*model.Delivery {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Delivery, 0)
	for _, d := range s.deliveries {
		if d.EndpointID == endpointID {
			list = append(list, d)
		}
	}
	return list
}

func (s *MemoryStore) UpdateDelivery(d *model.Delivery) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.deliveries[d.ID]; !ok {
		return ErrNotFound
	}
	s.deliveries[d.ID] = d
	return nil
}

func (s *MemoryStore) DeleteDelivery(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.deliveries[id]; !ok {
		return ErrNotFound
	}
	delete(s.deliveries, id)
	return nil
}
