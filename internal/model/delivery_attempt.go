package model

import (
	"strings"
	"time"
)

type DeliveryAttempt struct {
	ID           string    `json:"id"`
	DeliveryID   string    `json:"delivery_id"`
	AttemptNo    int       `json:"attempt_no"`
	RequestAt    time.Time `json:"request_at"`
	ResponseCode int       `json:"response_code"`
	DurationMs   int64     `json:"duration_ms"`
	Success      bool      `json:"success"`
	ResponseBody string    `json:"response_body"`
}

func (da *DeliveryAttempt) Validate() error {
	da.DeliveryID = strings.TrimSpace(da.DeliveryID)
	if da.DeliveryID == "" {
		return NewValidationError("delivery_id", "投递 ID 不能为空")
	}
	if da.AttemptNo < 1 {
		return NewValidationError("attempt_no", "尝试序号必须大于 0")
	}
	if da.DurationMs < 0 {
		return NewValidationError("duration_ms", "持续时间不能为负数")
	}
	return nil
}

type DeliveryAttemptFilter struct {
	DeliveryID string
	Success    *bool
}

func (f DeliveryAttemptFilter) Match(da *DeliveryAttempt) bool {
	if f.DeliveryID != "" && da.DeliveryID != f.DeliveryID {
		return false
	}
	if f.Success != nil && da.Success != *f.Success {
		return false
	}
	return true
}
