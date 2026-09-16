package model

import (
	"strings"
	"time"
)

const (
	EventCreated    = "created"
	EventDispatched = "dispatched"
	EventCompleted  = "completed"
	EventFailed     = "failed"
)

type Event struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	Source         string    `json:"source"`
	Payload        string    `json:"payload"`
	IdempotencyKey string    `json:"idempotency_key,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

func (e *Event) Validate() error {
	e.Type = strings.TrimSpace(e.Type)
	e.Source = strings.TrimSpace(e.Source)
	e.Payload = strings.TrimSpace(e.Payload)
	e.IdempotencyKey = strings.TrimSpace(e.IdempotencyKey)
	if e.Type == "" {
		return NewValidationError("type", "事件类型不能为空")
	}
	if e.Source == "" {
		return NewValidationError("source", "事件来源不能为空")
	}
	if e.Payload == "" {
		return NewValidationError("payload", "事件负载不能为空")
	}
	if e.Status == "" {
		e.Status = EventCreated
	}
	if e.Status != EventCreated && e.Status != EventDispatched && e.Status != EventCompleted && e.Status != EventFailed {
		return NewValidationError("status", "事件状态不合法")
	}
	return nil
}

type EventFilter struct {
	Type   string
	Source string
	Status string
}

func (f EventFilter) Match(e *Event) bool {
	if f.Type != "" && e.Type != f.Type {
		return false
	}
	if f.Source != "" && e.Source != f.Source {
		return false
	}
	if f.Status != "" && e.Status != f.Status {
		return false
	}
	return true
}

var eventTransitions = map[string]map[string]bool{
	EventCreated:    {EventDispatched: true, EventFailed: true},
	EventDispatched: {EventCompleted: true, EventFailed: true},
	EventCompleted:  {},
	EventFailed:     {},
}

func EventCanTransition(from, to string) bool {
	if m, ok := eventTransitions[from]; ok {
		return m[to]
	}
	return false
}
