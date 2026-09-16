package model

import (
	"strings"
	"time"
)

const (
	TemplateActive   = "active"
	TemplateInactive = "inactive"
)

type MessageTemplate struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	EventType   string    `json:"event_type"`
	ContentType string    `json:"content_type"`
	Body        string    `json:"body"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (mt *MessageTemplate) Validate() error {
	mt.Name = strings.TrimSpace(mt.Name)
	mt.EventType = strings.TrimSpace(mt.EventType)
	mt.ContentType = strings.TrimSpace(mt.ContentType)
	mt.Body = strings.TrimSpace(mt.Body)
	if mt.Name == "" {
		return NewValidationError("name", "模板名称不能为空")
	}
	if mt.EventType == "" {
		return NewValidationError("event_type", "事件类型不能为空")
	}
	if mt.Body == "" {
		return NewValidationError("body", "模板内容不能为空")
	}
	if mt.ContentType == "" {
		mt.ContentType = "application/json"
	}
	if mt.Status == "" {
		mt.Status = TemplateActive
	}
	if mt.Status != TemplateActive && mt.Status != TemplateInactive {
		return NewValidationError("status", "模板状态不合法")
	}
	return nil
}

type MessageTemplateFilter struct {
	EventType string
	Status    string
	Keyword   string
}

func (f MessageTemplateFilter) Match(mt *MessageTemplate) bool {
	if f.EventType != "" && mt.EventType != f.EventType {
		return false
	}
	if f.Status != "" && mt.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(mt.Name), k) {
			return false
		}
	}
	return true
}
