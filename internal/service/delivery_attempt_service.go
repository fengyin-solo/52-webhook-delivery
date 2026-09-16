package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/internal/store"
	"webhook/pkg/idgen"
)

func (s *Service) CreateDeliveryAttempt(input model.DeliveryAttempt) (*model.DeliveryAttempt, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDelivery(input.DeliveryID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("delivery_id", "投递记录不存在")
		}
		return nil, err
	}
	da := &model.DeliveryAttempt{
		ID:           idgen.Hex(),
		DeliveryID:   input.DeliveryID,
		AttemptNo:    input.AttemptNo,
		RequestAt:    input.RequestAt,
		ResponseCode: input.ResponseCode,
		DurationMs:   input.DurationMs,
		Success:      input.Success,
		ResponseBody: input.ResponseBody,
	}
	if da.RequestAt.IsZero() {
		da.RequestAt = time.Now()
	}
	if err := s.store.CreateDeliveryAttempt(da); err != nil {
		return nil, err
	}
	return da, nil
}

func (s *Service) GetDeliveryAttempt(id string) (*model.DeliveryAttempt, error) {
	return s.store.GetDeliveryAttempt(id)
}

func (s *Service) ListDeliveryAttempts(filter model.DeliveryAttemptFilter, page, size int) ([]*model.DeliveryAttempt, int, error) {
	all := s.store.ListDeliveryAttempts()
	matched := make([]*model.DeliveryAttempt, 0, len(all))
	for _, da := range all {
		if filter.Match(da) {
			matched = append(matched, da)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].RequestAt.After(matched[j].RequestAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.DeliveryAttempt{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateDeliveryAttempt(id string, input model.DeliveryAttempt) (*model.DeliveryAttempt, error) {
	existing, err := s.store.GetDeliveryAttempt(id)
	if err != nil {
		return nil, err
	}
	existing.ResponseCode = input.ResponseCode
	existing.DurationMs = input.DurationMs
	existing.Success = input.Success
	existing.ResponseBody = input.ResponseBody
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateDeliveryAttempt(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteDeliveryAttempt(id string) error {
	return s.store.DeleteDeliveryAttempt(id)
}
