package service

import (
	"sort"
	"time"

	"webhook/internal/model"
	"webhook/internal/store"
	"webhook/pkg/idgen"
)

func (s *Service) CreateWebhookLog(input model.WebhookLog) (*model.WebhookLog, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDelivery(input.DeliveryID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("delivery_id", "投递记录不存在")
		}
		return nil, err
	}
	if _, err := s.store.GetEndpoint(input.EndpointID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("endpoint_id", "端点不存在")
		}
		return nil, err
	}
	wl := &model.WebhookLog{
		ID:           idgen.Hex(),
		DeliveryID:   input.DeliveryID,
		EndpointID:   input.EndpointID,
		URL:          input.URL,
		RequestBody:  input.RequestBody,
		ResponseBody: input.ResponseBody,
		StatusCode:   input.StatusCode,
		Result:       input.Result,
		ErrorMessage: input.ErrorMessage,
		DurationMs:   input.DurationMs,
		CreatedAt:    time.Now(),
	}
	if wl.Result == "" {
		wl.Result = model.WebhookLogSuccess
	}
	if err := s.store.CreateWebhookLog(wl); err != nil {
		return nil, err
	}
	return wl, nil
}

func (s *Service) GetWebhookLog(id string) (*model.WebhookLog, error) {
	return s.store.GetWebhookLog(id)
}

func (s *Service) ListWebhookLogs(filter model.WebhookLogFilter, page, size int) ([]*model.WebhookLog, int, error) {
	all := s.store.ListWebhookLogs()
	matched := make([]*model.WebhookLog, 0, len(all))
	for _, wl := range all {
		if filter.Match(wl) {
			matched = append(matched, wl)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.WebhookLog{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateWebhookLog(id string, input model.WebhookLog) (*model.WebhookLog, error) {
	existing, err := s.store.GetWebhookLog(id)
	if err != nil {
		return nil, err
	}
	existing.ResponseBody = input.ResponseBody
	existing.StatusCode = input.StatusCode
	existing.Result = input.Result
	existing.ErrorMessage = input.ErrorMessage
	existing.DurationMs = input.DurationMs
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateWebhookLog(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteWebhookLog(id string) error {
	return s.store.DeleteWebhookLog(id)
}
