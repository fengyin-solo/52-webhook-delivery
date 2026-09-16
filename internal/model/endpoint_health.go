package model

import (
	"strings"
	"time"
)

const (
	EndpointHealthHealthy   = "healthy"
	EndpointHealthUnhealthy = "unhealthy"
	EndpointHealthUnknown   = "unknown"
)

type EndpointHealth struct {
	ID           string    `json:"id"`
	EndpointID   string    `json:"endpoint_id"`
	Status       string    `json:"status"`
	ResponseTime int64     `json:"response_time"`
	LastChecked  time.Time `json:"last_checked"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

func (eh *EndpointHealth) Validate() error {
	eh.EndpointID = strings.TrimSpace(eh.EndpointID)
	if eh.EndpointID == "" {
		return NewValidationError("endpoint_id", "端点 ID 不能为空")
	}
	if eh.Status == "" {
		eh.Status = EndpointHealthUnknown
	}
	if eh.Status != EndpointHealthHealthy && eh.Status != EndpointHealthUnhealthy && eh.Status != EndpointHealthUnknown {
		return NewValidationError("status", "健康状态不合法")
	}
	if eh.ResponseTime < 0 {
		return NewValidationError("response_time", "响应时间不能为负数")
	}
	return nil
}

type EndpointHealthFilter struct {
	EndpointID string
	Status     string
}

func (f EndpointHealthFilter) Match(eh *EndpointHealth) bool {
	if f.EndpointID != "" && eh.EndpointID != f.EndpointID {
		return false
	}
	if f.Status != "" && eh.Status != f.Status {
		return false
	}
	return true
}
