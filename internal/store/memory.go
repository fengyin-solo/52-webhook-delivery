package store

import (
	"sync"

	"webhook/internal/model"
)

type MemoryStore struct {
	mu                sync.RWMutex
	endpoints         map[string]*model.Endpoint
	subscriptions     map[string]*model.Subscription
	eventTypes        map[string]*model.EventType
	events            map[string]*model.Event
	deliveries        map[string]*model.Delivery
	deliveryAttempts  map[string]*model.DeliveryAttempt
	retryPolicies     map[string]*model.RetryPolicy
	signingKeys       map[string]*model.SigningKey
	filters           map[string]*model.Filter
	auditLogs         map[string]*model.AuditLog
	webhookLogs       map[string]*model.WebhookLog
	endpointGroups    map[string]*model.EndpointGroup
	messageTemplates  map[string]*model.MessageTemplate
	endpointHealths   map[string]*model.EndpointHealth
	circuitBreakers   map[string]*model.CircuitBreaker
	rateLimitRules    map[string]*model.RateLimitRule
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		endpoints:        make(map[string]*model.Endpoint),
		subscriptions:    make(map[string]*model.Subscription),
		eventTypes:       make(map[string]*model.EventType),
		events:           make(map[string]*model.Event),
		deliveries:       make(map[string]*model.Delivery),
		deliveryAttempts: make(map[string]*model.DeliveryAttempt),
		retryPolicies:    make(map[string]*model.RetryPolicy),
		signingKeys:      make(map[string]*model.SigningKey),
		filters:          make(map[string]*model.Filter),
		auditLogs:        make(map[string]*model.AuditLog),
		webhookLogs:      make(map[string]*model.WebhookLog),
		endpointGroups:   make(map[string]*model.EndpointGroup),
		messageTemplates: make(map[string]*model.MessageTemplate),
		endpointHealths:  make(map[string]*model.EndpointHealth),
		circuitBreakers:  make(map[string]*model.CircuitBreaker),
		rateLimitRules:   make(map[string]*model.RateLimitRule),
	}
}

var _ Store = (*MemoryStore)(nil)
