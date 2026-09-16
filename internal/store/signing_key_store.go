package store

import (
	"webhook/internal/model"
)

func (s *MemoryStore) CreateSigningKey(sk *model.SigningKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.signingKeys {
		if exist.EndpointID == sk.EndpointID && exist.KeyID == sk.KeyID {
			return ErrConflict
		}
	}
	s.signingKeys[sk.ID] = sk
	return nil
}

func (s *MemoryStore) GetSigningKey(id string) (*model.SigningKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sk, ok := s.signingKeys[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sk, nil
}

func (s *MemoryStore) ListSigningKeys() []*model.SigningKey {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.SigningKey, 0, len(s.signingKeys))
	for _, sk := range s.signingKeys {
		list = append(list, sk)
	}
	return list
}

func (s *MemoryStore) ListSigningKeysByEndpointID(endpointID string) []*model.SigningKey {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.SigningKey, 0)
	for _, sk := range s.signingKeys {
		if sk.EndpointID == endpointID {
			list = append(list, sk)
		}
	}
	return list
}

func (s *MemoryStore) UpdateSigningKey(sk *model.SigningKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.signingKeys[sk.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.signingKeys {
		if exist.ID != sk.ID && exist.EndpointID == sk.EndpointID && exist.KeyID == sk.KeyID {
			return ErrConflict
		}
	}
	s.signingKeys[sk.ID] = sk
	return nil
}

func (s *MemoryStore) DeleteSigningKey(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.signingKeys[id]; !ok {
		return ErrNotFound
	}
	delete(s.signingKeys, id)
	return nil
}
