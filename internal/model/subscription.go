package model

import (
	"strings"
	"time"
)

const (
	SubscriptionActive = "active"
	SubscriptionPaused = "paused"
)

type Subscription struct {
	ID         string    `json:"id"`
	EndpointID string    `json:"endpoint_id"`
	EventType  string    `json:"event_type"`
	FilterID   string    `json:"filter_id,omitempty"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *Subscription) Validate() error {
	s.EndpointID = strings.TrimSpace(s.EndpointID)
	s.EventType = strings.TrimSpace(s.EventType)
	s.FilterID = strings.TrimSpace(s.FilterID)
	if s.EndpointID == "" {
		return NewValidationError("endpoint_id", "端点 ID 不能为空")
	}
	if s.EventType == "" {
		return NewValidationError("event_type", "事件类型不能为空")
	}
	if s.Status == "" {
		s.Status = SubscriptionActive
	}
	if s.Status != SubscriptionActive && s.Status != SubscriptionPaused {
		return NewValidationError("status", "订阅状态不合法")
	}
	return nil
}

type SubscriptionFilter struct {
	EndpointID string
	EventType  string
	Status     string
}

func (f SubscriptionFilter) Match(s *Subscription) bool {
	if f.EndpointID != "" && s.EndpointID != f.EndpointID {
		return false
	}
	if f.EventType != "" && s.EventType != f.EventType {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	return true
}

var subscriptionTransitions = map[string]map[string]bool{
	SubscriptionActive: {SubscriptionPaused: true},
	SubscriptionPaused: {SubscriptionActive: true},
}

func SubscriptionCanTransition(from, to string) bool {
	if m, ok := subscriptionTransitions[from]; ok {
		return m[to]
	}
	return false
}
