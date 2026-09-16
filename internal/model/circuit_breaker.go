package model

import (
	"strings"
	"time"
)

const (
	CircuitClosed     = "closed"
	CircuitOpen       = "open"
	CircuitHalfOpen   = "half_open"
)

type CircuitBreaker struct {
	ID                string    `json:"id"`
	EndpointID        string    `json:"endpoint_id"`
	FailureThreshold  int       `json:"failure_threshold"`
	RecoveryTimeoutMs int       `json:"recovery_timeout_ms"`
	Status            string    `json:"status"`
	ConsecutiveFailures int     `json:"consecutive_failures"`
	LastFailureAt     *time.Time `json:"last_failure_at,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

func (cb *CircuitBreaker) Validate() error {
	cb.EndpointID = strings.TrimSpace(cb.EndpointID)
	if cb.EndpointID == "" {
		return NewValidationError("endpoint_id", "端点 ID 不能为空")
	}
	if cb.FailureThreshold < 1 {
		return NewValidationError("failure_threshold", "失败阈值至少为 1")
	}
	if cb.RecoveryTimeoutMs < 1000 {
		return NewValidationError("recovery_timeout_ms", "恢复超时至少为 1000ms")
	}
	if cb.Status == "" {
		cb.Status = CircuitClosed
	}
	if cb.Status != CircuitClosed && cb.Status != CircuitOpen && cb.Status != CircuitHalfOpen {
		return NewValidationError("status", "熔断器状态不合法")
	}
	if cb.ConsecutiveFailures < 0 {
		return NewValidationError("consecutive_failures", "连续失败次数不能为负数")
	}
	return nil
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.ConsecutiveFailures++
	now := time.Now()
	cb.LastFailureAt = &now
	if cb.ConsecutiveFailures >= cb.FailureThreshold {
		cb.Status = CircuitOpen
	}
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.ConsecutiveFailures = 0
	cb.LastFailureAt = nil
	if cb.Status == CircuitHalfOpen {
		cb.Status = CircuitClosed
	}
}

func (cb *CircuitBreaker) CanAttempt() bool {
	if cb.Status == CircuitClosed {
		return true
	}
	if cb.Status == CircuitHalfOpen {
		return true
	}
	if cb.Status == CircuitOpen && cb.LastFailureAt != nil {
		elapsed := time.Since(*cb.LastFailureAt).Milliseconds()
		if elapsed >= int64(cb.RecoveryTimeoutMs) {
			cb.Status = CircuitHalfOpen
			return true
		}
	}
	return false
}

type CircuitBreakerFilter struct {
	EndpointID string
	Status     string
}

func (f CircuitBreakerFilter) Match(cb *CircuitBreaker) bool {
	if f.EndpointID != "" && cb.EndpointID != f.EndpointID {
		return false
	}
	if f.Status != "" && cb.Status != f.Status {
		return false
	}
	return true
}
