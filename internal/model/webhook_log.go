package model

import (
	"strings"
	"time"
)

const (
	WebhookLogSuccess = "success"
	WebhookLogFailure = "failure"
	WebhookLogTimeout = "timeout"
)

type WebhookLog struct {
	ID           string    `json:"id"`
	DeliveryID   string    `json:"delivery_id"`
	EndpointID   string    `json:"endpoint_id"`
	URL          string    `json:"url"`
	RequestBody  string    `json:"request_body"`
	ResponseBody string    `json:"response_body"`
	StatusCode   int       `json:"status_code"`
	Result       string    `json:"result"`
	ErrorMessage string    `json:"error_message,omitempty"`
	DurationMs   int64     `json:"duration_ms"`
	CreatedAt    time.Time `json:"created_at"`
}

func (wl *WebhookLog) Validate() error {
	wl.DeliveryID = strings.TrimSpace(wl.DeliveryID)
	wl.EndpointID = strings.TrimSpace(wl.EndpointID)
	wl.URL = strings.TrimSpace(wl.URL)
	if wl.DeliveryID == "" {
		return NewValidationError("delivery_id", "投递 ID 不能为空")
	}
	if wl.EndpointID == "" {
		return NewValidationError("endpoint_id", "端点 ID 不能为空")
	}
	if wl.URL == "" {
		return NewValidationError("url", "URL 不能为空")
	}
	if wl.Result == "" {
		wl.Result = WebhookLogSuccess
	}
	if wl.Result != WebhookLogSuccess && wl.Result != WebhookLogFailure && wl.Result != WebhookLogTimeout {
		return NewValidationError("result", "日志结果不合法")
	}
	if wl.DurationMs < 0 {
		return NewValidationError("duration_ms", "持续时间不能为负数")
	}
	return nil
}

type WebhookLogFilter struct {
	DeliveryID string
	EndpointID string
	Result     string
}

func (f WebhookLogFilter) Match(wl *WebhookLog) bool {
	if f.DeliveryID != "" && wl.DeliveryID != f.DeliveryID {
		return false
	}
	if f.EndpointID != "" && wl.EndpointID != f.EndpointID {
		return false
	}
	if f.Result != "" && wl.Result != f.Result {
		return false
	}
	return true
}
