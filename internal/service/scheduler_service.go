package service

import (
	"time"

	"webhook/internal/model"
	"webhook/pkg/idgen"
)

// ProcessPendingDeliveries 处理待投递和重试中的投递记录。
func (s *Service) ProcessPendingDeliveries() (processed int, err error) {
	all := s.store.ListDeliveries()
	now := time.Now()
	for _, d := range all {
		if d.Status != model.DeliveryPending && d.Status != model.DeliveryRetrying {
			continue
		}
		if d.NextRetryAt != nil && d.NextRetryAt.After(now) {
			continue
		}
		ep, err := s.store.GetEndpoint(d.EndpointID)
		if err != nil || ep.Status != model.EndpointActive {
			continue
		}
		event, err := s.store.GetEvent(d.EventID)
		if err != nil {
			continue
		}
		d.Attempts++
		attempt := &model.DeliveryAttempt{
			ID:         idgen.Hex(),
			DeliveryID: d.ID,
			AttemptNo:  d.Attempts,
			RequestAt:  now,
		}
		if ep.MaxRetries > 0 && d.Attempts > ep.MaxRetries {
			d.Status = model.DeliveryFailed
			attempt.Success = false
			attempt.ResponseBody = "超过最大重试次数"
			_ = s.store.CreateDeliveryAttempt(attempt)
			_ = s.store.UpdateDelivery(d)
			processed++
			continue
		}
		signKeys := s.store.ListSigningKeysByEndpointID(ep.ID)
		var signature string
		for _, sk := range signKeys {
			if sk.Status == model.SigningKeyActive {
				signature = s.ComputeHMAC(sk.KeyValue, event.Payload)
				break
			}
		}
		_ = signature
		d.LastResponseCode = 200
		d.LastResponseBody = "simulated_ok"
		attempt.Success = true
		attempt.ResponseCode = 200
		attempt.ResponseBody = "simulated_ok"
		attempt.DurationMs = 10
		_ = s.store.CreateDeliveryAttempt(attempt)
		if d.Attempts >= 1 {
			d.Status = model.DeliveryDelivered
			d.DeliveredAt = &now
			d.NextRetryAt = nil
		} else {
			retryPolicies := s.store.ListRetryPolicies()
			var policy *model.RetryPolicy
			for _, rp := range retryPolicies {
				if rp.Status == model.RetryPolicyActive {
					policy = rp
					break
				}
			}
			if policy != nil {
				delay := policy.ComputeDelay(d.Attempts)
				next := now.Add(time.Duration(delay) * time.Millisecond)
				d.NextRetryAt = &next
			}
		}
		_ = s.store.UpdateDelivery(d)
		processed++
	}
	return processed, nil
}

// GetPendingDeliveryCount 获取待处理投递数量。
func (s *Service) GetPendingDeliveryCount() int {
	all := s.store.ListDeliveries()
	count := 0
	for _, d := range all {
		if d.Status == model.DeliveryPending || d.Status == model.DeliveryRetrying {
			count++
		}
	}
	return count
}

// GetOverdueDeliveries 获取已超时的投递记录。
func (s *Service) GetOverdueDeliveries() []*model.Delivery {
	all := s.store.ListDeliveries()
	now := time.Now()
	var result []*model.Delivery
	for _, d := range all {
		if d.Status != model.DeliveryPending && d.Status != model.DeliveryRetrying {
			continue
		}
		if d.NextRetryAt != nil && d.NextRetryAt.Before(now) {
			result = append(result, d)
		} else if d.NextRetryAt == nil && d.Status == model.DeliveryPending {
			result = append(result, d)
		}
	}
	return result
}

// RetryAllFailedDeliveries 批量重试所有失败投递。
func (s *Service) RetryAllFailedDeliveries() (int, error) {
	all := s.store.ListDeliveries()
	count := 0
	for _, d := range all {
		if d.Status == model.DeliveryFailed {
			d.Status = model.DeliveryRetrying
			d.Attempts = 0
			now := time.Now()
			d.NextRetryAt = &now
			_ = s.store.UpdateDelivery(d)
			count++
		}
	}
	return count, nil
}
