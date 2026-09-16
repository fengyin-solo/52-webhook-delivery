package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/internal/store"
	"webhook/pkg/idgen"
)

func (s *Service) CreateSigningKey(input model.SigningKey) (*model.SigningKey, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetEndpoint(input.EndpointID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("endpoint_id", "端点不存在")
		}
		return nil, err
	}
	sk := &model.SigningKey{
		ID:         idgen.Hex(),
		EndpointID: input.EndpointID,
		KeyID:      input.KeyID,
		Algorithm:  input.Algorithm,
		KeyValue:   input.KeyValue,
		Status:     input.Status,
		CreatedAt:  time.Now(),
	}
	if sk.Status == "" {
		sk.Status = model.SigningKeyActive
	}
	if err := s.store.CreateSigningKey(sk); err != nil {
		return nil, err
	}
	return sk, nil
}

func (s *Service) GetSigningKey(id string) (*model.SigningKey, error) {
	return s.store.GetSigningKey(id)
}

func (s *Service) ListSigningKeys(filter model.SigningKeyFilter, page, size int) ([]*model.SigningKey, int, error) {
	all := s.store.ListSigningKeys()
	matched := make([]*model.SigningKey, 0, len(all))
	for _, sk := range all {
		if filter.Match(sk) {
			matched = append(matched, sk)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.SigningKey{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateSigningKey(id string, input model.SigningKey) (*model.SigningKey, error) {
	existing, err := s.store.GetSigningKey(id)
	if err != nil {
		return nil, err
	}
	existing.KeyValue = input.KeyValue
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSigningKey(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteSigningKey(id string) error {
	return s.store.DeleteSigningKey(id)
}

func (s *Service) ComputeHMAC(key, payload string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *Service) VerifyHMAC(key, payload, signature string) bool {
	expected := s.ComputeHMAC(key, payload)
	return hmac.Equal([]byte(expected), []byte(signature))
}
