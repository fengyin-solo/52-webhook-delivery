package model

import (
	"strings"
	"time"
)

const (
	EventTypeActive   = "active"
	EventTypeDisabled = "disabled"
)

type EventType struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Schema      string    `json:"schema"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (et *EventType) Validate() error {
	et.Name = strings.TrimSpace(et.Name)
	et.Description = strings.TrimSpace(et.Description)
	et.Schema = strings.TrimSpace(et.Schema)
	if et.Name == "" {
		return NewValidationError("name", "事件类型名称不能为空")
	}
	if len(et.Name) > 128 {
		return NewValidationError("name", "事件类型名称不能超过 128 个字符")
	}
	if et.Status == "" {
		et.Status = EventTypeActive
	}
	if et.Status != EventTypeActive && et.Status != EventTypeDisabled {
		return NewValidationError("status", "事件类型状态不合法")
	}
	return nil
}

type EventTypeFilter struct {
	Status  string
	Keyword string
}

func (f EventTypeFilter) Match(et *EventType) bool {
	if f.Status != "" && et.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(et.Name), k) &&
			!strings.Contains(strings.ToLower(et.Description), k) {
			return false
		}
	}
	return true
}
