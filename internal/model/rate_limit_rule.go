package model

import (
	"strings"
	"time"
)

const (
	RateLimitRuleActive   = "active"
	RateLimitRuleInactive = "inactive"
)

type RateLimitRule struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	EndpointID  string    `json:"endpoint_id,omitempty"`
	EventType   string    `json:"event_type,omitempty"`
	MaxRequests int       `json:"max_requests"`
	WindowMs    int       `json:"window_ms"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (rlr *RateLimitRule) Validate() error {
	rlr.Name = strings.TrimSpace(rlr.Name)
	rlr.EndpointID = strings.TrimSpace(rlr.EndpointID)
	rlr.EventType = strings.TrimSpace(rlr.EventType)
	if rlr.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	if rlr.MaxRequests < 1 {
		return NewValidationError("max_requests", "最大请求数至少为 1")
	}
	if rlr.WindowMs < 1000 {
		return NewValidationError("window_ms", "时间窗口至少为 1000ms")
	}
	if rlr.Status == "" {
		rlr.Status = RateLimitRuleActive
	}
	if rlr.Status != RateLimitRuleActive && rlr.Status != RateLimitRuleInactive {
		return NewValidationError("status", "限流规则状态不合法")
	}
	return nil
}

func (rlr *RateLimitRule) IsGlobal() bool {
	return rlr.EndpointID == "" && rlr.EventType == ""
}

func (rlr *RateLimitRule) MatchEndpoint(endpointID, eventType string) bool {
	if rlr.Status != RateLimitRuleActive {
		return false
	}
	if rlr.EndpointID != "" && rlr.EndpointID != endpointID {
		return false
	}
	if rlr.EventType != "" && rlr.EventType != eventType {
		return false
	}
	return true
}

type RateLimitRuleFilter struct {
	EndpointID string
	EventType  string
	Status     string
}

func (f RateLimitRuleFilter) Match(rlr *RateLimitRule) bool {
	if f.EndpointID != "" && rlr.EndpointID != f.EndpointID {
		return false
	}
	if f.EventType != "" && rlr.EventType != f.EventType {
		return false
	}
	if f.Status != "" && rlr.Status != f.Status {
		return false
	}
	return true
}
