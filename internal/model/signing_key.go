package model

import (
	"strings"
	"time"
)

const (
	SigningKeyActive   = "active"
	SigningKeyRotated  = "rotated"
)

const (
	AlgorithmHMACSHA256 = "hmac-sha256"
)

type SigningKey struct {
	ID         string    `json:"id"`
	EndpointID string    `json:"endpoint_id"`
	KeyID      string    `json:"key_id"`
	Algorithm  string    `json:"algorithm"`
	KeyValue   string    `json:"key_value"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (sk *SigningKey) Validate() error {
	sk.EndpointID = strings.TrimSpace(sk.EndpointID)
	sk.KeyID = strings.TrimSpace(sk.KeyID)
	sk.KeyValue = strings.TrimSpace(sk.KeyValue)
	if sk.EndpointID == "" {
		return NewValidationError("endpoint_id", "端点 ID 不能为空")
	}
	if sk.KeyID == "" {
		return NewValidationError("key_id", "密钥 ID 不能为空")
	}
	if sk.KeyValue == "" {
		return NewValidationError("key_value", "密钥值不能为空")
	}
	if sk.Algorithm == "" {
		sk.Algorithm = AlgorithmHMACSHA256
	}
	if sk.Algorithm != AlgorithmHMACSHA256 {
		return NewValidationError("algorithm", "签名算法不合法")
	}
	if sk.Status == "" {
		sk.Status = SigningKeyActive
	}
	if sk.Status != SigningKeyActive && sk.Status != SigningKeyRotated {
		return NewValidationError("status", "密钥状态不合法")
	}
	return nil
}

type SigningKeyFilter struct {
	EndpointID string
	Status     string
}

func (f SigningKeyFilter) Match(sk *SigningKey) bool {
	if f.EndpointID != "" && sk.EndpointID != f.EndpointID {
		return false
	}
	if f.Status != "" && sk.Status != f.Status {
		return false
	}
	return true
}
