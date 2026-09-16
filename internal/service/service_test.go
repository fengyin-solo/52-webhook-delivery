package service

import (
	"testing"

	"webhook/internal/config"
	"webhook/internal/model"
	"webhook/internal/store"
	"webhook/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestService_CreateEndpoint(t *testing.T) {
	s := newTestService()
	e, err := s.CreateEndpoint(model.Endpoint{URL: "https://example.com/webhook", Secret: "secret", MaxRetries: 3, TimeoutMs: 5000})
	if err != nil {
		t.Fatalf("create endpoint: %v", err)
	}
	if e.ID == "" {
		t.Fatal("expected endpoint id")
	}
	if e.Status != model.EndpointActive {
		t.Fatalf("expected active, got %s", e.Status)
	}
	_, err = s.CreateEndpoint(model.Endpoint{URL: "https://example.com/webhook", Secret: "s2", MaxRetries: 3, TimeoutMs: 5000})
	if err == nil {
		t.Fatal("expected conflict on duplicate url")
	}
}

func TestService_EndpointStatusTransition(t *testing.T) {
	s := newTestService()
	e, _ := s.CreateEndpoint(model.Endpoint{URL: "https://a.com", Secret: "s", MaxRetries: 3, TimeoutMs: 5000})
	_, err := s.UpdateEndpoint(e.ID, model.Endpoint{URL: "https://a.com", Secret: "s", Status: model.EndpointDisabled, MaxRetries: 3, TimeoutMs: 5000})
	if err != nil {
		t.Fatalf("update status: %v", err)
	}
	updated, _ := s.GetEndpoint(e.ID)
	if updated.Status != model.EndpointDisabled {
		t.Fatalf("expected disabled, got %s", updated.Status)
	}
}

func TestService_CreateSubscription(t *testing.T) {
	s := newTestService()
	_, err := s.CreateSubscription(model.Subscription{EndpointID: "nonexist", EventType: "order.created"})
	if err == nil {
		t.Fatal("expected validation error for missing endpoint")
	}
	ep, _ := s.CreateEndpoint(model.Endpoint{URL: "https://b.com", Secret: "s", MaxRetries: 3, TimeoutMs: 5000})
	sub, err := s.CreateSubscription(model.Subscription{EndpointID: ep.ID, EventType: "order.created"})
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	if sub.Status != model.SubscriptionActive {
		t.Fatalf("expected active, got %s", sub.Status)
	}
}

func TestService_CreateEventType(t *testing.T) {
	s := newTestService()
	et, err := s.CreateEventType(model.EventType{Name: "user.signup", Description: "User signed up", Schema: "{}", Status: model.EventTypeActive})
	if err != nil {
		t.Fatalf("create event type: %v", err)
	}
	if et.ID == "" {
		t.Fatal("expected event type id")
	}
	_, err = s.CreateEventType(model.EventType{Name: "user.signup", Description: "Dup"})
	if err == nil {
		t.Fatal("expected conflict on duplicate name")
	}
}

func TestService_CreateEvent(t *testing.T) {
	s := newTestService()
	_, err := s.CreateEvent(model.Event{Type: "nonexist", Source: "api", Payload: "{}"})
	if err == nil {
		t.Fatal("expected validation error for missing event type")
	}
	_, _ = s.CreateEventType(model.EventType{Name: "order.created", Description: "Order"})
	ev, err := s.CreateEvent(model.Event{Type: "order.created", Source: "api", Payload: "{\"id\":1}"})
	if err != nil {
		t.Fatalf("create event: %v", err)
	}
	if ev.Status != model.EventCreated && ev.Status != model.EventDispatched {
		t.Fatalf("unexpected event status: %s", ev.Status)
	}
}

func TestService_EventDispatch(t *testing.T) {
	s := newTestService()
	ep, _ := s.CreateEndpoint(model.Endpoint{URL: "https://hook.com", Secret: "s", MaxRetries: 3, TimeoutMs: 5000})
	_, _ = s.CreateEventType(model.EventType{Name: "order.created", Description: "Order"})
	_, _ = s.CreateSubscription(model.Subscription{EndpointID: ep.ID, EventType: "order.created"})
	ev, _ := s.CreateEvent(model.Event{Type: "order.created", Source: "api", Payload: "{\"id\":1}"})
	ds := s.store.ListDeliveriesByEventID(ev.ID)
	if len(ds) != 1 {
		t.Fatalf("expected 1 delivery, got %d", len(ds))
	}
	if ds[0].Status != model.DeliveryPending {
		t.Fatalf("expected pending delivery, got %s", ds[0].Status)
	}
}

func TestService_EventDispatchWithFilter(t *testing.T) {
	s := newTestService()
	ep, _ := s.CreateEndpoint(model.Endpoint{URL: "https://hook.com", Secret: "s", MaxRetries: 3, TimeoutMs: 5000})
	f, _ := s.CreateFilter(model.Filter{Name: "status-filter", Field: "status", Operator: model.OpEQ, Value: "paid"})
	_, _ = s.CreateEventType(model.EventType{Name: "order.created", Description: "Order"})
	_, _ = s.CreateSubscription(model.Subscription{EndpointID: ep.ID, EventType: "order.created", FilterID: f.ID})
	ev1, _ := s.CreateEvent(model.Event{Type: "order.created", Source: "api", Payload: "{\"status\":\"paid\"}"})
	ds1 := s.store.ListDeliveriesByEventID(ev1.ID)
	if len(ds1) != 1 {
		t.Fatalf("expected 1 delivery for matching payload, got %d", len(ds1))
	}
	ev2, _ := s.CreateEvent(model.Event{Type: "order.created", Source: "api", Payload: "{\"status\":\"pending\"}"})
	ds2 := s.store.ListDeliveriesByEventID(ev2.ID)
	if len(ds2) != 0 {
		t.Fatalf("expected 0 deliveries for non-matching payload, got %d", len(ds2))
	}
}

func TestService_DeliveryStatusTransition(t *testing.T) {
	s := newTestService()
	ep, _ := s.CreateEndpoint(model.Endpoint{URL: "https://hook.com", Secret: "s", MaxRetries: 3, TimeoutMs: 5000})
	_, _ = s.CreateEventType(model.EventType{Name: "order.created", Description: "Order"})
	ev, _ := s.CreateEvent(model.Event{Type: "order.created", Source: "api", Payload: "{}"})
	d, _ := s.CreateDelivery(model.Delivery{EventID: ev.ID, EndpointID: ep.ID, Status: model.DeliveryPending})
	_, err := s.UpdateDelivery(d.ID, model.Delivery{Status: model.DeliveryDelivered})
	if err != nil {
		t.Fatalf("transition to delivered: %v", err)
	}
	updated, _ := s.GetDelivery(d.ID)
	if updated.Status != model.DeliveryDelivered {
		t.Fatalf("expected delivered, got %s", updated.Status)
	}
}

func TestService_RetryDelivery(t *testing.T) {
	s := newTestService()
	ep, _ := s.CreateEndpoint(model.Endpoint{URL: "https://hook.com", Secret: "s", MaxRetries: 3, TimeoutMs: 5000})
	_, _ = s.CreateEventType(model.EventType{Name: "order.created", Description: "Order"})
	ev, _ := s.CreateEvent(model.Event{Type: "order.created", Source: "api", Payload: "{}"})
	d, _ := s.CreateDelivery(model.Delivery{EventID: ev.ID, EndpointID: ep.ID, Status: model.DeliveryFailed})
	_, err := s.RetryDelivery(d.ID)
	if err != nil {
		t.Fatalf("retry delivery: %v", err)
	}
	updated, _ := s.GetDelivery(d.ID)
	if updated.Status != model.DeliveryRetrying {
		t.Fatalf("expected retrying, got %s", updated.Status)
	}
	_, err = s.RetryDelivery(d.ID)
	if err == nil {
		t.Fatal("expected error when retrying non-failed delivery")
	}
}

func TestService_RetryPolicyBackoff(t *testing.T) {
	s := newTestService()
	rp, _ := s.CreateRetryPolicy(model.RetryPolicy{Name: "exp", MaxAttempts: 5, BackoffType: model.BackoffExponential, InitialDelayMs: 1000, MaxDelayMs: 30000})
	if d := rp.ComputeDelay(1); d != 1000 {
		t.Fatalf("expected 1000, got %d", d)
	}
	if d := rp.ComputeDelay(2); d != 2000 {
		t.Fatalf("expected 2000, got %d", d)
	}
	if d := rp.ComputeDelay(3); d != 4000 {
		t.Fatalf("expected 4000, got %d", d)
	}
	l, _ := s.CreateRetryPolicy(model.RetryPolicy{Name: "lin", MaxAttempts: 5, BackoffType: model.BackoffLinear, InitialDelayMs: 1000, MaxDelayMs: 30000})
	if d := l.ComputeDelay(3); d != 3000 {
		t.Fatalf("expected 3000, got %d", d)
	}
	f, _ := s.CreateRetryPolicy(model.RetryPolicy{Name: "fix", MaxAttempts: 5, BackoffType: model.BackoffFixed, InitialDelayMs: 1000, MaxDelayMs: 30000})
	if d := f.ComputeDelay(5); d != 1000 {
		t.Fatalf("expected 1000, got %d", d)
	}
}

func TestService_HMAC(t *testing.T) {
	s := newTestService()
	sig := s.ComputeHMAC("secret", "payload")
	if sig == "" {
		t.Fatal("expected signature")
	}
	if !s.VerifyHMAC("secret", "payload", sig) {
		t.Fatal("expected valid hmac")
	}
	if s.VerifyHMAC("secret", "payload", "bad") {
		t.Fatal("expected invalid hmac")
	}
}

func TestService_BatchCreateEvents(t *testing.T) {
	s := newTestService()
	_, _ = s.CreateEventType(model.EventType{Name: "order.created", Description: "Order"})
	result, err := s.BatchCreateEvents(model.BatchCreateEventRequest{Events: []model.Event{
		{Type: "order.created", Source: "api", Payload: "{}"},
		{Type: "order.created", Source: "api", Payload: "{}"},
	}})
	if err != nil {
		t.Fatalf("batch create events: %v", err)
	}
	if result.Success != 2 {
		t.Fatalf("expected 2 success, got %d", result.Success)
	}
}

func TestService_BatchDisableEndpoints(t *testing.T) {
	s := newTestService()
	ep, _ := s.CreateEndpoint(model.Endpoint{URL: "https://a.com", Secret: "s", MaxRetries: 3, TimeoutMs: 5000})
	result, err := s.BatchDisableEndpoints(model.BatchDisableEndpointRequest{IDs: []string{ep.ID}})
	if err != nil {
		t.Fatalf("batch disable: %v", err)
	}
	if result.Success != 1 {
		t.Fatalf("expected 1 success, got %d", result.Success)
	}
	updated, _ := s.GetEndpoint(ep.ID)
	if updated.Status != model.EndpointDisabled {
		t.Fatalf("expected disabled, got %s", updated.Status)
	}
}

func TestService_Stats(t *testing.T) {
	s := newTestService()
	overview := s.GetDeliveryOverviewStats()
	if overview.TotalCount != 0 {
		t.Fatalf("expected 0 total, got %d", overview.TotalCount)
	}
	avg, p95 := s.GetLatencyStats()
	if avg != 0 || p95 != 0 {
		t.Fatalf("expected 0 latency stats")
	}
	snapshot := s.ExportSnapshot()
	if len(snapshot.Endpoints) != 0 {
		t.Fatalf("expected 0 endpoints in snapshot")
	}
}

func TestService_FilterMatchPayload(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFilter(model.Filter{Name: "test", Field: "status", Operator: model.OpEQ, Value: "paid"})
	if !f.MatchPayload("{\"status\":\"paid\"}") {
		t.Fatal("expected match")
	}
	if f.MatchPayload("{\"status\":\"pending\"}") {
		t.Fatal("expected no match")
	}
	f2, _ := s.CreateFilter(model.Filter{Name: "test2", Field: "name", Operator: model.OpContains, Value: "test"})
	if !f2.MatchPayload("{\"name\":\"test-user\"}") {
		t.Fatal("expected contains match")
	}
	f3, _ := s.CreateFilter(model.Filter{Name: "test3", Field: "name", Operator: model.OpStartsWith, Value: "pre"})
	if !f3.MatchPayload("{\"name\":\"prefix\"}") {
		t.Fatal("expected starts_with match")
	}
	f4, _ := s.CreateFilter(model.Filter{Name: "test4", Field: "status", Operator: model.OpNE, Value: "deleted"})
	if !f4.MatchPayload("{\"status\":\"active\"}") {
		t.Fatal("expected ne match")
	}
}

func TestService_CreateDeliveryAttempt(t *testing.T) {
	s := newTestService()
	ep, _ := s.CreateEndpoint(model.Endpoint{URL: "https://a.com", Secret: "s", MaxRetries: 3, TimeoutMs: 5000})
	_, _ = s.CreateEventType(model.EventType{Name: "order.created", Description: "Order"})
	ev, _ := s.CreateEvent(model.Event{Type: "order.created", Source: "api", Payload: "{}"})
	d, _ := s.CreateDelivery(model.Delivery{EventID: ev.ID, EndpointID: ep.ID})
	_, err := s.CreateDeliveryAttempt(model.DeliveryAttempt{DeliveryID: d.ID, AttemptNo: 1, ResponseCode: 200, DurationMs: 100, Success: true})
	if err != nil {
		t.Fatalf("create delivery attempt: %v", err)
	}
}

func TestService_AuditLog(t *testing.T) {
	s := newTestService()
	a, err := s.CreateAuditLog(model.AuditLog{Operator: "admin", Action: "create", TargetType: "endpoint", TargetID: "e1", Detail: "detail"})
	if err != nil {
		t.Fatalf("create audit log: %v", err)
	}
	if a.ID == "" {
		t.Fatal("expected audit log id")
	}
}

func TestService_ListPagination(t *testing.T) {
	s := newTestService()
	for i := 0; i < 5; i++ {
		_, _ = s.CreateEndpoint(model.Endpoint{URL: "https://example" + string(rune('a'+i)) + ".com", Secret: "s", MaxRetries: 3, TimeoutMs: 5000})
	}
	items, total, err := s.ListEndpoints(model.EndpointFilter{}, 1, 2)
	if err != nil {
		t.Fatalf("list endpoints: %v", err)
	}
	if total != 5 {
		t.Fatalf("expected total 5, got %d", total)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	items2, _, _ := s.ListEndpoints(model.EndpointFilter{}, 3, 2)
	if len(items2) != 1 {
		t.Fatalf("expected 1 item on page 3, got %d", len(items2))
	}
}

func TestService_CreateEventIdempotency(t *testing.T) {
	s := newTestService()
	_, _ = s.CreateEventType(model.EventType{Name: "order.created", Description: "Order"})
	_, _ = s.CreateEvent(model.Event{Type: "order.created", Source: "api", Payload: "{}", IdempotencyKey: "key1"})
	_, err := s.CreateEvent(model.Event{Type: "order.created", Source: "api", Payload: "{}", IdempotencyKey: "key1"})
	if err == nil {
		t.Fatal("expected conflict on duplicate idempotency key")
	}
}

func TestService_EventStatusMachine(t *testing.T) {
	s := newTestService()
	_, _ = s.CreateEventType(model.EventType{Name: "order.created", Description: "Order"})
	ep, _ := s.CreateEndpoint(model.Endpoint{URL: "https://hook.com", Secret: "s", MaxRetries: 3, TimeoutMs: 5000})
	_, _ = s.CreateSubscription(model.Subscription{EndpointID: ep.ID, EventType: "order.created"})
	ev, _ := s.CreateEvent(model.Event{Type: "order.created", Source: "api", Payload: "{}"})
	if ev.Status != model.EventDispatched {
		t.Fatalf("expected dispatched after auto-dispatch, got %s", ev.Status)
	}
	_, err := s.UpdateEvent(ev.ID, model.Event{Status: model.EventCompleted})
	if err != nil {
		t.Fatalf("transition dispatched->completed: %v", err)
	}
	_, err = s.UpdateEvent(ev.ID, model.Event{Status: model.EventFailed})
	if err == nil {
		t.Fatal("expected error transitioning completed->failed")
	}
}

func TestService_SigningKeyUnique(t *testing.T) {
	s := newTestService()
	ep, _ := s.CreateEndpoint(model.Endpoint{URL: "https://a.com", Secret: "s", MaxRetries: 3, TimeoutMs: 5000})
	_, _ = s.CreateSigningKey(model.SigningKey{EndpointID: ep.ID, KeyID: "k1", KeyValue: "v1"})
	_, err := s.CreateSigningKey(model.SigningKey{EndpointID: ep.ID, KeyID: "k1", KeyValue: "v2"})
	if err == nil {
		t.Fatal("expected conflict on duplicate endpoint_id+key_id")
	}
}
