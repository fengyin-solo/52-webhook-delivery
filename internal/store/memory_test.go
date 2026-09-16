package store

import (
	"testing"
	"time"

	"webhook/internal/model"
)

func TestMemoryStore_EndpointCRUD(t *testing.T) {
	s := NewMemoryStore()
	e := &model.Endpoint{ID: "e1", URL: "https://a.com", Status: model.EndpointActive, MaxRetries: 3, TimeoutMs: 5000, CreatedAt: time.Now()}
	if err := s.CreateEndpoint(e); err != nil {
		t.Fatalf("create endpoint: %v", err)
	}
	if err := s.CreateEndpoint(&model.Endpoint{ID: "e2", URL: "https://a.com", Status: model.EndpointActive, MaxRetries: 3, TimeoutMs: 5000, CreatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetEndpoint("e1")
	if err != nil {
		t.Fatalf("get endpoint: %v", err)
	}
	if got.ID != "e1" {
		t.Fatalf("expected e1, got %s", got.ID)
	}
	if _, err := s.GetEndpoint("e999"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	list := s.ListEndpoints()
	if len(list) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(list))
	}
	e.URL = "https://b.com"
	if err := s.UpdateEndpoint(e); err != nil {
		t.Fatalf("update endpoint: %v", err)
	}
	if err := s.DeleteEndpoint("e1"); err != nil {
		t.Fatalf("delete endpoint: %v", err)
	}
	if _, err := s.GetEndpoint("e1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStore_SubscriptionCRUD(t *testing.T) {
	s := NewMemoryStore()
	sub := &model.Subscription{ID: "s1", EndpointID: "e1", EventType: "order.created", Status: model.SubscriptionActive, CreatedAt: time.Now()}
	if err := s.CreateSubscription(sub); err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	got, err := s.GetSubscription("s1")
	if err != nil {
		t.Fatalf("get subscription: %v", err)
	}
	if got.ID != "s1" {
		t.Fatalf("expected s1, got %s", got.ID)
	}
	list := s.ListSubscriptions()
	if len(list) != 1 {
		t.Fatalf("expected 1 subscription, got %d", len(list))
	}
	byEndpoint := s.ListSubscriptionsByEndpointID("e1")
	if len(byEndpoint) != 1 {
		t.Fatalf("expected 1 subscription by endpoint, got %d", len(byEndpoint))
	}
	byType := s.ListSubscriptionsByEventType("order.created")
	if len(byType) != 1 {
		t.Fatalf("expected 1 subscription by type, got %d", len(byType))
	}
	sub.EventType = "order.updated"
	if err := s.UpdateSubscription(sub); err != nil {
		t.Fatalf("update subscription: %v", err)
	}
	if err := s.DeleteSubscription("s1"); err != nil {
		t.Fatalf("delete subscription: %v", err)
	}
	if _, err := s.GetSubscription("s1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStore_EventTypeCRUD(t *testing.T) {
	s := NewMemoryStore()
	et := &model.EventType{ID: "et1", Name: "order.created", Status: model.EventTypeActive, CreatedAt: time.Now()}
	if err := s.CreateEventType(et); err != nil {
		t.Fatalf("create event type: %v", err)
	}
	if err := s.CreateEventType(&model.EventType{ID: "et2", Name: "order.created", Status: model.EventTypeActive, CreatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetEventType("et1")
	if err != nil {
		t.Fatalf("get event type: %v", err)
	}
	if got.Name != "order.created" {
		t.Fatalf("expected order.created, got %s", got.Name)
	}
	byName, err := s.GetEventTypeByName("order.created")
	if err != nil {
		t.Fatalf("get by name: %v", err)
	}
	if byName.ID != "et1" {
		t.Fatalf("expected et1, got %s", byName.ID)
	}
	et.Name = "order.updated"
	if err := s.UpdateEventType(et); err != nil {
		t.Fatalf("update event type: %v", err)
	}
	if err := s.DeleteEventType("et1"); err != nil {
		t.Fatalf("delete event type: %v", err)
	}
	if _, err := s.GetEventType("et1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStore_EventCRUD(t *testing.T) {
	s := NewMemoryStore()
	e := &model.Event{ID: "ev1", Type: "order.created", Source: "api", Payload: "{}", Status: model.EventCreated, CreatedAt: time.Now()}
	if err := s.CreateEvent(e); err != nil {
		t.Fatalf("create event: %v", err)
	}
	if err := s.CreateEvent(&model.Event{ID: "ev2", Type: "order.created", Source: "api", Payload: "{}", IdempotencyKey: "key1", Status: model.EventCreated, CreatedAt: time.Now()}); err != nil {
		t.Fatalf("create event with key: %v", err)
	}
	if err := s.CreateEvent(&model.Event{ID: "ev3", Type: "order.created", Source: "api", Payload: "{}", IdempotencyKey: "key1", Status: model.EventCreated, CreatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict on idempotency key, got %v", err)
	}
	got, err := s.GetEvent("ev1")
	if err != nil {
		t.Fatalf("get event: %v", err)
	}
	if got.ID != "ev1" {
		t.Fatalf("expected ev1, got %s", got.ID)
	}
	list := s.ListEvents()
	if len(list) != 2 {
		t.Fatalf("expected 2 events, got %d", len(list))
	}
	byType := s.ListEventsByType("order.created")
	if len(byType) != 2 {
		t.Fatalf("expected 2 events by type, got %d", len(byType))
	}
	e.Payload = "{\"a\":1}"
	if err := s.UpdateEvent(e); err != nil {
		t.Fatalf("update event: %v", err)
	}
	if err := s.DeleteEvent("ev1"); err != nil {
		t.Fatalf("delete event: %v", err)
	}
	if _, err := s.GetEvent("ev1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStore_DeliveryCRUD(t *testing.T) {
	s := NewMemoryStore()
	d := &model.Delivery{ID: "d1", EventID: "ev1", EndpointID: "e1", Status: model.DeliveryPending, Attempts: 0, CreatedAt: time.Now()}
	if err := s.CreateDelivery(d); err != nil {
		t.Fatalf("create delivery: %v", err)
	}
	got, err := s.GetDelivery("d1")
	if err != nil {
		t.Fatalf("get delivery: %v", err)
	}
	if got.ID != "d1" {
		t.Fatalf("expected d1, got %s", got.ID)
	}
	list := s.ListDeliveries()
	if len(list) != 1 {
		t.Fatalf("expected 1 delivery, got %d", len(list))
	}
	byEvent := s.ListDeliveriesByEventID("ev1")
	if len(byEvent) != 1 {
		t.Fatalf("expected 1 delivery by event, got %d", len(byEvent))
	}
	byEndpoint := s.ListDeliveriesByEndpointID("e1")
	if len(byEndpoint) != 1 {
		t.Fatalf("expected 1 delivery by endpoint, got %d", len(byEndpoint))
	}
	d.Attempts = 1
	if err := s.UpdateDelivery(d); err != nil {
		t.Fatalf("update delivery: %v", err)
	}
	if err := s.DeleteDelivery("d1"); err != nil {
		t.Fatalf("delete delivery: %v", err)
	}
	if _, err := s.GetDelivery("d1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStore_DeliveryAttemptCRUD(t *testing.T) {
	s := NewMemoryStore()
	da := &model.DeliveryAttempt{ID: "da1", DeliveryID: "d1", AttemptNo: 1, RequestAt: time.Now(), ResponseCode: 200, DurationMs: 100, Success: true, ResponseBody: "ok"}
	if err := s.CreateDeliveryAttempt(da); err != nil {
		t.Fatalf("create delivery attempt: %v", err)
	}
	got, err := s.GetDeliveryAttempt("da1")
	if err != nil {
		t.Fatalf("get delivery attempt: %v", err)
	}
	if got.ID != "da1" {
		t.Fatalf("expected da1, got %s", got.ID)
	}
	list := s.ListDeliveryAttempts()
	if len(list) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(list))
	}
	byDelivery := s.ListDeliveryAttemptsByDeliveryID("d1")
	if len(byDelivery) != 1 {
		t.Fatalf("expected 1 attempt by delivery, got %d", len(byDelivery))
	}
	da.ResponseCode = 500
	if err := s.UpdateDeliveryAttempt(da); err != nil {
		t.Fatalf("update delivery attempt: %v", err)
	}
	if err := s.DeleteDeliveryAttempt("da1"); err != nil {
		t.Fatalf("delete delivery attempt: %v", err)
	}
	if _, err := s.GetDeliveryAttempt("da1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStore_RetryPolicyCRUD(t *testing.T) {
	s := NewMemoryStore()
	rp := &model.RetryPolicy{ID: "rp1", Name: "default", MaxAttempts: 3, BackoffType: model.BackoffFixed, InitialDelayMs: 1000, MaxDelayMs: 30000, Status: model.RetryPolicyActive, CreatedAt: time.Now()}
	if err := s.CreateRetryPolicy(rp); err != nil {
		t.Fatalf("create retry policy: %v", err)
	}
	if err := s.CreateRetryPolicy(&model.RetryPolicy{ID: "rp2", Name: "default", MaxAttempts: 3, BackoffType: model.BackoffFixed, InitialDelayMs: 1000, MaxDelayMs: 30000, Status: model.RetryPolicyActive, CreatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetRetryPolicy("rp1")
	if err != nil {
		t.Fatalf("get retry policy: %v", err)
	}
	if got.Name != "default" {
		t.Fatalf("expected default, got %s", got.Name)
	}
	rp.Name = "updated"
	if err := s.UpdateRetryPolicy(rp); err != nil {
		t.Fatalf("update retry policy: %v", err)
	}
	if err := s.DeleteRetryPolicy("rp1"); err != nil {
		t.Fatalf("delete retry policy: %v", err)
	}
	if _, err := s.GetRetryPolicy("rp1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStore_SigningKeyCRUD(t *testing.T) {
	s := NewMemoryStore()
	sk := &model.SigningKey{ID: "sk1", EndpointID: "e1", KeyID: "k1", Algorithm: model.AlgorithmHMACSHA256, KeyValue: "secret", Status: model.SigningKeyActive, CreatedAt: time.Now()}
	if err := s.CreateSigningKey(sk); err != nil {
		t.Fatalf("create signing key: %v", err)
	}
	if err := s.CreateSigningKey(&model.SigningKey{ID: "sk2", EndpointID: "e1", KeyID: "k1", Algorithm: model.AlgorithmHMACSHA256, KeyValue: "secret2", Status: model.SigningKeyActive, CreatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetSigningKey("sk1")
	if err != nil {
		t.Fatalf("get signing key: %v", err)
	}
	if got.KeyID != "k1" {
		t.Fatalf("expected k1, got %s", got.KeyID)
	}
	list := s.ListSigningKeysByEndpointID("e1")
	if len(list) != 1 {
		t.Fatalf("expected 1 signing key by endpoint, got %d", len(list))
	}
	sk.KeyValue = "newsecret"
	if err := s.UpdateSigningKey(sk); err != nil {
		t.Fatalf("update signing key: %v", err)
	}
	if err := s.DeleteSigningKey("sk1"); err != nil {
		t.Fatalf("delete signing key: %v", err)
	}
	if _, err := s.GetSigningKey("sk1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStore_FilterCRUD(t *testing.T) {
	s := NewMemoryStore()
	f := &model.Filter{ID: "f1", Name: "order-filter", Field: "status", Operator: model.OpEQ, Value: "paid", Status: model.FilterActive, CreatedAt: time.Now()}
	if err := s.CreateFilter(f); err != nil {
		t.Fatalf("create filter: %v", err)
	}
	if err := s.CreateFilter(&model.Filter{ID: "f2", Name: "order-filter", Field: "status", Operator: model.OpEQ, Value: "paid", Status: model.FilterActive, CreatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetFilter("f1")
	if err != nil {
		t.Fatalf("get filter: %v", err)
	}
	if got.Name != "order-filter" {
		t.Fatalf("expected order-filter, got %s", got.Name)
	}
	f.Value = "shipped"
	if err := s.UpdateFilter(f); err != nil {
		t.Fatalf("update filter: %v", err)
	}
	if err := s.DeleteFilter("f1"); err != nil {
		t.Fatalf("delete filter: %v", err)
	}
	if _, err := s.GetFilter("f1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStore_AuditLogCRUD(t *testing.T) {
	s := NewMemoryStore()
	a := &model.AuditLog{ID: "a1", Operator: "admin", Action: "create", TargetType: "endpoint", TargetID: "e1", Detail: "created endpoint", CreatedAt: time.Now()}
	if err := s.CreateAuditLog(a); err != nil {
		t.Fatalf("create audit log: %v", err)
	}
	got, err := s.GetAuditLog("a1")
	if err != nil {
		t.Fatalf("get audit log: %v", err)
	}
	if got.Operator != "admin" {
		t.Fatalf("expected admin, got %s", got.Operator)
	}
	list := s.ListAuditLogs()
	if len(list) != 1 {
		t.Fatalf("expected 1 audit log, got %d", len(list))
	}
}

func TestMemoryStore_WebhookLogCRUD(t *testing.T) {
	s := NewMemoryStore()
	wl := &model.WebhookLog{ID: "wl1", DeliveryID: "d1", EndpointID: "e1", URL: "https://a.com", Result: model.WebhookLogSuccess, DurationMs: 100, CreatedAt: time.Now()}
	if err := s.CreateWebhookLog(wl); err != nil {
		t.Fatalf("create webhook log: %v", err)
	}
	got, err := s.GetWebhookLog("wl1")
	if err != nil {
		t.Fatalf("get webhook log: %v", err)
	}
	if got.ID != "wl1" {
		t.Fatalf("expected wl1, got %s", got.ID)
	}
	list := s.ListWebhookLogs()
	if len(list) != 1 {
		t.Fatalf("expected 1 webhook log, got %d", len(list))
	}
	byDelivery := s.ListWebhookLogsByDeliveryID("d1")
	if len(byDelivery) != 1 {
		t.Fatalf("expected 1 webhook log by delivery, got %d", len(byDelivery))
	}
	wl.DurationMs = 200
	if err := s.UpdateWebhookLog(wl); err != nil {
		t.Fatalf("update webhook log: %v", err)
	}
	if err := s.DeleteWebhookLog("wl1"); err != nil {
		t.Fatalf("delete webhook log: %v", err)
	}
	if _, err := s.GetWebhookLog("wl1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStore_EndpointGroupCRUD(t *testing.T) {
	s := NewMemoryStore()
	eg := &model.EndpointGroup{ID: "eg1", Name: "prod", Description: "production", Status: model.EndpointGroupActive, CreatedAt: time.Now()}
	if err := s.CreateEndpointGroup(eg); err != nil {
		t.Fatalf("create endpoint group: %v", err)
	}
	if err := s.CreateEndpointGroup(&model.EndpointGroup{ID: "eg2", Name: "prod", Status: model.EndpointGroupActive, CreatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetEndpointGroup("eg1")
	if err != nil {
		t.Fatalf("get endpoint group: %v", err)
	}
	if got.Name != "prod" {
		t.Fatalf("expected prod, got %s", got.Name)
	}
	byName, err := s.GetEndpointGroupByName("prod")
	if err != nil {
		t.Fatalf("get by name: %v", err)
	}
	if byName.ID != "eg1" {
		t.Fatalf("expected eg1, got %s", byName.ID)
	}
	eg.Name = "staging"
	if err := s.UpdateEndpointGroup(eg); err != nil {
		t.Fatalf("update endpoint group: %v", err)
	}
	if err := s.DeleteEndpointGroup("eg1"); err != nil {
		t.Fatalf("delete endpoint group: %v", err)
	}
	if _, err := s.GetEndpointGroup("eg1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStore_MessageTemplateCRUD(t *testing.T) {
	s := NewMemoryStore()
	mt := &model.MessageTemplate{ID: "mt1", Name: "order-tpl", EventType: "order.created", ContentType: "application/json", Body: "{}", Status: model.TemplateActive, CreatedAt: time.Now()}
	if err := s.CreateMessageTemplate(mt); err != nil {
		t.Fatalf("create message template: %v", err)
	}
	got, err := s.GetMessageTemplate("mt1")
	if err != nil {
		t.Fatalf("get message template: %v", err)
	}
	if got.Name != "order-tpl" {
		t.Fatalf("expected order-tpl, got %s", got.Name)
	}
	list := s.ListMessageTemplates()
	if len(list) != 1 {
		t.Fatalf("expected 1 message template, got %d", len(list))
	}
	byType := s.ListMessageTemplatesByEventType("order.created")
	if len(byType) != 1 {
		t.Fatalf("expected 1 message template by type, got %d", len(byType))
	}
	mt.Body = "{\"a\":1}"
	if err := s.UpdateMessageTemplate(mt); err != nil {
		t.Fatalf("update message template: %v", err)
	}
	if err := s.DeleteMessageTemplate("mt1"); err != nil {
		t.Fatalf("delete message template: %v", err)
	}
	if _, err := s.GetMessageTemplate("mt1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStore_CircuitBreakerCRUD(t *testing.T) {
	s := NewMemoryStore()
	cb := &model.CircuitBreaker{ID: "cb1", EndpointID: "e1", FailureThreshold: 5, RecoveryTimeoutMs: 30000, Status: model.CircuitClosed, CreatedAt: time.Now()}
	if err := s.CreateCircuitBreaker(cb); err != nil {
		t.Fatalf("create circuit breaker: %v", err)
	}
	if err := s.CreateCircuitBreaker(&model.CircuitBreaker{ID: "cb2", EndpointID: "e1", FailureThreshold: 3, RecoveryTimeoutMs: 10000, Status: model.CircuitClosed, CreatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetCircuitBreaker("cb1")
	if err != nil {
		t.Fatalf("get circuit breaker: %v", err)
	}
	if got.EndpointID != "e1" {
		t.Fatalf("expected e1, got %s", got.EndpointID)
	}
	byEndpoint, err := s.GetCircuitBreakerByEndpointID("e1")
	if err != nil {
		t.Fatalf("get by endpoint: %v", err)
	}
	if byEndpoint.ID != "cb1" {
		t.Fatalf("expected cb1, got %s", byEndpoint.ID)
	}
	cb.FailureThreshold = 10
	if err := s.UpdateCircuitBreaker(cb); err != nil {
		t.Fatalf("update circuit breaker: %v", err)
	}
	if err := s.DeleteCircuitBreaker("cb1"); err != nil {
		t.Fatalf("delete circuit breaker: %v", err)
	}
	if _, err := s.GetCircuitBreaker("cb1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStore_RateLimitRuleCRUD(t *testing.T) {
	s := NewMemoryStore()
	rlr := &model.RateLimitRule{ID: "rlr1", Name: "global", MaxRequests: 100, WindowMs: 60000, Status: model.RateLimitRuleActive, CreatedAt: time.Now()}
	if err := s.CreateRateLimitRule(rlr); err != nil {
		t.Fatalf("create rate limit rule: %v", err)
	}
	got, err := s.GetRateLimitRule("rlr1")
	if err != nil {
		t.Fatalf("get rate limit rule: %v", err)
	}
	if got.Name != "global" {
		t.Fatalf("expected global, got %s", got.Name)
	}
	list := s.ListRateLimitRules()
	if len(list) != 1 {
		t.Fatalf("expected 1 rate limit rule, got %d", len(list))
	}
	rlr.MaxRequests = 200
	if err := s.UpdateRateLimitRule(rlr); err != nil {
		t.Fatalf("update rate limit rule: %v", err)
	}
	if err := s.DeleteRateLimitRule("rlr1"); err != nil {
		t.Fatalf("delete rate limit rule: %v", err)
	}
	if _, err := s.GetRateLimitRule("rlr1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}
