package model

import (
	"strings"
	"time"
)

const (
	EndpointActive   = "active"
	EndpointPaused   = "paused"
	EndpointDisabled = "disabled"
)

type Endpoint struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Secret      string    `json:"secret"`
	Status      string    `json:"status"`
	MaxRetries  int       `json:"max_retries"`
	TimeoutMs   int       `json:"timeout_ms"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (e *Endpoint) Validate() error {
	e.URL = strings.TrimSpace(e.URL)
	e.Secret = strings.TrimSpace(e.Secret)
	e.Description = strings.TrimSpace(e.Description)
	if e.URL == "" {
		return NewValidationError("url", "URL 不能为空")
	}
	if !strings.HasPrefix(e.URL, "http://") && !strings.HasPrefix(e.URL, "https://") {
		return NewValidationError("url", "URL 必须以 http:// 或 https:// 开头")
	}
	if e.Status == "" {
		e.Status = EndpointActive
	}
	if e.Status != EndpointActive && e.Status != EndpointPaused && e.Status != EndpointDisabled {
		return NewValidationError("status", "端点状态不合法")
	}
	if e.MaxRetries < 0 {
		return NewValidationError("max_retries", "最大重试次数不能为负数")
	}
	if e.MaxRetries > 100 {
		return NewValidationError("max_retries", "最大重试次数不能超过 100")
	}
	if e.TimeoutMs < 100 {
		return NewValidationError("timeout_ms", "超时时间不能小于 100ms")
	}
	if e.TimeoutMs > 300000 {
		return NewValidationError("timeout_ms", "超时时间不能超过 300000ms")
	}
	return nil
}

type EndpointFilter struct {
	Status  string
	Keyword string
}

func (f EndpointFilter) Match(e *Endpoint) bool {
	if f.Status != "" && e.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(e.URL), k) &&
			!strings.Contains(strings.ToLower(e.Description), k) {
			return false
		}
	}
	return true
}

var endpointTransitions = map[string]map[string]bool{
	EndpointActive:   {EndpointPaused: true, EndpointDisabled: true},
	EndpointPaused:   {EndpointActive: true, EndpointDisabled: true},
	EndpointDisabled: {EndpointActive: true, EndpointPaused: true},
}

func EndpointCanTransition(from, to string) bool {
	if m, ok := endpointTransitions[from]; ok {
		return m[to]
	}
	return false
}
