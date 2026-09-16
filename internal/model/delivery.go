package model

import (
	"strings"
	"time"
)

const (
	DeliveryPending    = "pending"
	DeliveryRetrying   = "retrying"
	DeliveryDelivered  = "delivered"
	DeliveryFailed     = "failed"
)

type Delivery struct {
	ID               string    `json:"id"`
	EventID          string    `json:"event_id"`
	EndpointID       string    `json:"endpoint_id"`
	Status           string    `json:"status"`
	Attempts         int       `json:"attempts"`
	LastResponseCode int       `json:"last_response_code"`
	LastResponseBody string    `json:"last_response_body"`
	NextRetryAt      *time.Time `json:"next_retry_at,omitempty"`
	DeliveredAt      *time.Time `json:"delivered_at,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

func (d *Delivery) Validate() error {
	d.EventID = strings.TrimSpace(d.EventID)
	d.EndpointID = strings.TrimSpace(d.EndpointID)
	if d.EventID == "" {
		return NewValidationError("event_id", "事件 ID 不能为空")
	}
	if d.EndpointID == "" {
		return NewValidationError("endpoint_id", "端点 ID 不能为空")
	}
	if d.Status == "" {
		d.Status = DeliveryPending
	}
	if d.Status != DeliveryPending && d.Status != DeliveryRetrying && d.Status != DeliveryDelivered && d.Status != DeliveryFailed {
		return NewValidationError("status", "投递状态不合法")
	}
	if d.Attempts < 0 {
		return NewValidationError("attempts", "尝试次数不能为负数")
	}
	return nil
}

type DeliveryFilter struct {
	EventID    string
	EndpointID string
	Status     string
}

func (f DeliveryFilter) Match(d *Delivery) bool {
	if f.EventID != "" && d.EventID != f.EventID {
		return false
	}
	if f.EndpointID != "" && d.EndpointID != f.EndpointID {
		return false
	}
	if f.Status != "" && d.Status != f.Status {
		return false
	}
	return true
}

var deliveryTransitions = map[string]map[string]bool{
	DeliveryPending:   {DeliveryRetrying: true, DeliveryDelivered: true, DeliveryFailed: true},
	DeliveryRetrying:  {DeliveryRetrying: true, DeliveryDelivered: true, DeliveryFailed: true},
	DeliveryDelivered: {},
	DeliveryFailed:    {DeliveryRetrying: true},
}

func DeliveryCanTransition(from, to string) bool {
	if m, ok := deliveryTransitions[from]; ok {
		return m[to]
	}
	return false
}
