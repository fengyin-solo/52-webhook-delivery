// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"webhook/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Endpoint
	CreateEndpoint(e *model.Endpoint) error
	GetEndpoint(id string) (*model.Endpoint, error)
	GetEndpointByURL(url string) (*model.Endpoint, error)
	ListEndpoints() []*model.Endpoint
	UpdateEndpoint(e *model.Endpoint) error
	DeleteEndpoint(id string) error

	// Subscription
	CreateSubscription(s *model.Subscription) error
	GetSubscription(id string) (*model.Subscription, error)
	ListSubscriptions() []*model.Subscription
	ListSubscriptionsByEndpointID(endpointID string) []*model.Subscription
	ListSubscriptionsByEventType(eventType string) []*model.Subscription
	UpdateSubscription(s *model.Subscription) error
	DeleteSubscription(id string) error

	// EventType
	CreateEventType(et *model.EventType) error
	GetEventType(id string) (*model.EventType, error)
	GetEventTypeByName(name string) (*model.EventType, error)
	ListEventTypes() []*model.EventType
	UpdateEventType(et *model.EventType) error
	DeleteEventType(id string) error

	// Event
	CreateEvent(e *model.Event) error
	GetEvent(id string) (*model.Event, error)
	ListEvents() []*model.Event
	ListEventsByType(eventType string) []*model.Event
	UpdateEvent(e *model.Event) error
	DeleteEvent(id string) error

	// Delivery
	CreateDelivery(d *model.Delivery) error
	GetDelivery(id string) (*model.Delivery, error)
	ListDeliveries() []*model.Delivery
	ListDeliveriesByEventID(eventID string) []*model.Delivery
	ListDeliveriesByEndpointID(endpointID string) []*model.Delivery
	UpdateDelivery(d *model.Delivery) error
	DeleteDelivery(id string) error

	// DeliveryAttempt
	CreateDeliveryAttempt(da *model.DeliveryAttempt) error
	GetDeliveryAttempt(id string) (*model.DeliveryAttempt, error)
	ListDeliveryAttempts() []*model.DeliveryAttempt
	ListDeliveryAttemptsByDeliveryID(deliveryID string) []*model.DeliveryAttempt
	UpdateDeliveryAttempt(da *model.DeliveryAttempt) error
	DeleteDeliveryAttempt(id string) error

	// RetryPolicy
	CreateRetryPolicy(rp *model.RetryPolicy) error
	GetRetryPolicy(id string) (*model.RetryPolicy, error)
	ListRetryPolicies() []*model.RetryPolicy
	UpdateRetryPolicy(rp *model.RetryPolicy) error
	DeleteRetryPolicy(id string) error

	// SigningKey
	CreateSigningKey(sk *model.SigningKey) error
	GetSigningKey(id string) (*model.SigningKey, error)
	ListSigningKeys() []*model.SigningKey
	ListSigningKeysByEndpointID(endpointID string) []*model.SigningKey
	UpdateSigningKey(sk *model.SigningKey) error
	DeleteSigningKey(id string) error

	// Filter
	CreateFilter(f *model.Filter) error
	GetFilter(id string) (*model.Filter, error)
	ListFilters() []*model.Filter
	UpdateFilter(f *model.Filter) error
	DeleteFilter(id string) error

	// AuditLog
	CreateAuditLog(a *model.AuditLog) error
	GetAuditLog(id string) (*model.AuditLog, error)
	ListAuditLogs() []*model.AuditLog

	// WebhookLog
	CreateWebhookLog(wl *model.WebhookLog) error
	GetWebhookLog(id string) (*model.WebhookLog, error)
	ListWebhookLogs() []*model.WebhookLog
	ListWebhookLogsByDeliveryID(deliveryID string) []*model.WebhookLog
	UpdateWebhookLog(wl *model.WebhookLog) error
	DeleteWebhookLog(id string) error

	// EndpointGroup
	CreateEndpointGroup(eg *model.EndpointGroup) error
	GetEndpointGroup(id string) (*model.EndpointGroup, error)
	GetEndpointGroupByName(name string) (*model.EndpointGroup, error)
	ListEndpointGroups() []*model.EndpointGroup
	UpdateEndpointGroup(eg *model.EndpointGroup) error
	DeleteEndpointGroup(id string) error

	// MessageTemplate
	CreateMessageTemplate(mt *model.MessageTemplate) error
	GetMessageTemplate(id string) (*model.MessageTemplate, error)
	ListMessageTemplates() []*model.MessageTemplate
	ListMessageTemplatesByEventType(eventType string) []*model.MessageTemplate
	UpdateMessageTemplate(mt *model.MessageTemplate) error
	DeleteMessageTemplate(id string) error

	// EndpointHealth
	CreateEndpointHealth(eh *model.EndpointHealth) error
	GetEndpointHealth(id string) (*model.EndpointHealth, error)
	ListEndpointHealths() []*model.EndpointHealth
	ListEndpointHealthsByEndpointID(endpointID string) []*model.EndpointHealth
	UpdateEndpointHealth(eh *model.EndpointHealth) error
	DeleteEndpointHealth(id string) error

	// CircuitBreaker
	CreateCircuitBreaker(cb *model.CircuitBreaker) error
	GetCircuitBreaker(id string) (*model.CircuitBreaker, error)
	GetCircuitBreakerByEndpointID(endpointID string) (*model.CircuitBreaker, error)
	ListCircuitBreakers() []*model.CircuitBreaker
	UpdateCircuitBreaker(cb *model.CircuitBreaker) error
	DeleteCircuitBreaker(id string) error

	// RateLimitRule
	CreateRateLimitRule(rlr *model.RateLimitRule) error
	GetRateLimitRule(id string) (*model.RateLimitRule, error)
	ListRateLimitRules() []*model.RateLimitRule
	UpdateRateLimitRule(rlr *model.RateLimitRule) error
	DeleteRateLimitRule(id string) error
}
