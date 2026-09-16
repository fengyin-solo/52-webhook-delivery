package model

import (
	"strings"
	"time"
)

const (
	BackoffFixed       = "fixed"
	BackoffExponential = "exponential"
	BackoffLinear      = "linear"
)

const (
	RetryPolicyActive   = "active"
	RetryPolicyDisabled = "disabled"
)

type RetryPolicy struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	MaxAttempts    int       `json:"max_attempts"`
	BackoffType    string    `json:"backoff_type"`
	InitialDelayMs int       `json:"initial_delay_ms"`
	MaxDelayMs     int       `json:"max_delay_ms"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

func (rp *RetryPolicy) Validate() error {
	rp.Name = strings.TrimSpace(rp.Name)
	if rp.Name == "" {
		return NewValidationError("name", "策略名称不能为空")
	}
	if rp.MaxAttempts < 1 {
		return NewValidationError("max_attempts", "最大尝试次数至少为 1")
	}
	if rp.MaxAttempts > 100 {
		return NewValidationError("max_attempts", "最大尝试次数不能超过 100")
	}
	if rp.BackoffType == "" {
		rp.BackoffType = BackoffFixed
	}
	if rp.BackoffType != BackoffFixed && rp.BackoffType != BackoffExponential && rp.BackoffType != BackoffLinear {
		return NewValidationError("backoff_type", "退避类型不合法")
	}
	if rp.InitialDelayMs < 1 {
		return NewValidationError("initial_delay_ms", "初始延迟必须大于 0")
	}
	if rp.MaxDelayMs < rp.InitialDelayMs {
		return NewValidationError("max_delay_ms", "最大延迟不能小于初始延迟")
	}
	if rp.Status == "" {
		rp.Status = RetryPolicyActive
	}
	if rp.Status != RetryPolicyActive && rp.Status != RetryPolicyDisabled {
		return NewValidationError("status", "策略状态不合法")
	}
	return nil
}

func (rp *RetryPolicy) ComputeDelay(attempt int) int {
	if attempt < 1 {
		attempt = 1
	}
	switch rp.BackoffType {
	case BackoffFixed:
		return rp.InitialDelayMs
	case BackoffLinear:
		d := rp.InitialDelayMs * attempt
		if d > rp.MaxDelayMs {
			return rp.MaxDelayMs
		}
		return d
	case BackoffExponential:
		d := rp.InitialDelayMs
		for i := 1; i < attempt; i++ {
			d *= 2
			if d > rp.MaxDelayMs {
				return rp.MaxDelayMs
			}
		}
		return d
	default:
		return rp.InitialDelayMs
	}
}

type RetryPolicyFilter struct {
	Status string
	Name   string
}

func (f RetryPolicyFilter) Match(rp *RetryPolicy) bool {
	if f.Status != "" && rp.Status != f.Status {
		return false
	}
	if f.Name != "" {
		k := strings.ToLower(strings.TrimSpace(f.Name))
		if k != "" && !strings.Contains(strings.ToLower(rp.Name), k) {
			return false
		}
	}
	return true
}
