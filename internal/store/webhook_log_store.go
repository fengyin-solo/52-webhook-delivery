package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateWebhookLog(wl *model.WebhookLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.webhookLogs[wl.ID] = wl
	return nil
}

func (s *MemoryStore) GetWebhookLog(id string) (*model.WebhookLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	wl, ok := s.webhookLogs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return wl, nil
}

func (s *MemoryStore) ListWebhookLogs() []*model.WebhookLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.WebhookLog, 0, len(s.webhookLogs))
	for _, wl := range s.webhookLogs {
		list = append(list, wl)
	}
	return list
}

func (s *MemoryStore) ListWebhookLogsByDeliveryID(deliveryID string) []*model.WebhookLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.WebhookLog, 0)
	for _, wl := range s.webhookLogs {
		if wl.DeliveryID == deliveryID {
			list = append(list, wl)
		}
	}
	return list
}

func (s *MemoryStore) UpdateWebhookLog(wl *model.WebhookLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.webhookLogs[wl.ID]; !ok {
		return ErrNotFound
	}
	s.webhookLogs[wl.ID] = wl
	return nil
}

func (s *MemoryStore) DeleteWebhookLog(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.webhookLogs[id]; !ok {
		return ErrNotFound
	}
	delete(s.webhookLogs, id)
	return nil
}
