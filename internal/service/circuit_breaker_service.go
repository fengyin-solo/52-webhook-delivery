package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/internal/store"
	"webhook/pkg/idgen"
)

func (s *Service) CreateCircuitBreaker(input model.CircuitBreaker) (*model.CircuitBreaker, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetEndpoint(input.EndpointID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("endpoint_id", "端点不存在")
		}
		return nil, err
	}
	cb := &model.CircuitBreaker{
		ID:                  idgen.Hex(),
		EndpointID:          input.EndpointID,
		FailureThreshold:    input.FailureThreshold,
		RecoveryTimeoutMs:   input.RecoveryTimeoutMs,
		Status:              input.Status,
		ConsecutiveFailures: 0,
		CreatedAt:           time.Now(),
	}
	if cb.Status == "" {
		cb.Status = model.CircuitClosed
	}
	if err := s.store.CreateCircuitBreaker(cb); err != nil {
		return nil, err
	}
	return cb, nil
}

func (s *Service) GetCircuitBreaker(id string) (*model.CircuitBreaker, error) {
	return s.store.GetCircuitBreaker(id)
}

func (s *Service) ListCircuitBreakers(filter model.CircuitBreakerFilter, page, size int) ([]*model.CircuitBreaker, int, error) {
	all := s.store.ListCircuitBreakers()
	matched := make([]*model.CircuitBreaker, 0, len(all))
	for _, cb := range all {
		if filter.Match(cb) {
			matched = append(matched, cb)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.CircuitBreaker{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateCircuitBreaker(id string, input model.CircuitBreaker) (*model.CircuitBreaker, error) {
	existing, err := s.store.GetCircuitBreaker(id)
	if err != nil {
		return nil, err
	}
	existing.FailureThreshold = input.FailureThreshold
	existing.RecoveryTimeoutMs = input.RecoveryTimeoutMs
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateCircuitBreaker(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteCircuitBreaker(id string) error {
	return s.store.DeleteCircuitBreaker(id)
}

func (s *Service) RecordCircuitBreakerFailure(endpointID string) (*model.CircuitBreaker, error) {
	cb, err := s.store.GetCircuitBreakerByEndpointID(endpointID)
	if err != nil {
		return nil, err
	}
	cb.RecordFailure()
	if err := s.store.UpdateCircuitBreaker(cb); err != nil {
		return nil, err
	}
	return cb, nil
}

func (s *Service) RecordCircuitBreakerSuccess(endpointID string) (*model.CircuitBreaker, error) {
	cb, err := s.store.GetCircuitBreakerByEndpointID(endpointID)
	if err != nil {
		return nil, err
	}
	cb.RecordSuccess()
	if err := s.store.UpdateCircuitBreaker(cb); err != nil {
		return nil, err
	}
	return cb, nil
}
