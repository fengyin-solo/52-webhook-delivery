package service

import (
	"testing"

	"webhook/internal/model"
)

func TestService_ProcessPendingDeliveries(t *testing.T) {
	s := newTestService()
	ep, _ := s.CreateEndpoint(model.Endpoint{URL: "https://hook.com", Secret: "s", MaxRetries: 3, TimeoutMs: 5000})
	_, _ = s.CreateEventType(model.EventType{Name: "order.created", Description: "Order"})
	_, _ = s.CreateSubscription(model.Subscription{EndpointID: ep.ID, EventType: "order.created"})
	ev, _ := s.CreateEvent(model.Event{Type: "order.created", Source: "api", Payload: "{}"})
	ds := s.store.ListDeliveriesByEventID(ev.ID)
	if len(ds) != 1 {
		t.Fatalf("expected 1 delivery, got %d", len(ds))
	}
	processed, err := s.ProcessPendingDeliveries()
	if err != nil {
		t.Fatalf("process pending: %v", err)
	}
	if processed != 1 {
		t.Fatalf("expected 1 processed, got %d", processed)
	}
	updated, _ := s.GetDelivery(ds[0].ID)
	if updated.Status != model.DeliveryDelivered {
		t.Fatalf("expected delivered, got %s", updated.Status)
	}
}

func TestService_GetPendingDeliveryCount(t *testing.T) {
	s := newTestService()
	if cnt := s.GetPendingDeliveryCount(); cnt != 0 {
		t.Fatalf("expected 0, got %d", cnt)
	}
}

func TestService_GetOverdueDeliveries(t *testing.T) {
	s := newTestService()
	overdue := s.GetOverdueDeliveries()
	if len(overdue) != 0 {
		t.Fatalf("expected 0 overdue, got %d", len(overdue))
	}
}

func TestService_RetryAllFailedDeliveries(t *testing.T) {
	s := newTestService()
	ep, _ := s.CreateEndpoint(model.Endpoint{URL: "https://hook.com", Secret: "s", MaxRetries: 3, TimeoutMs: 5000})
	_, _ = s.CreateEventType(model.EventType{Name: "order.created", Description: "Order"})
	ev, _ := s.CreateEvent(model.Event{Type: "order.created", Source: "api", Payload: "{}"})
	d, _ := s.CreateDelivery(model.Delivery{EventID: ev.ID, EndpointID: ep.ID, Status: model.DeliveryFailed})
	count, err := s.RetryAllFailedDeliveries()
	if err != nil {
		t.Fatalf("retry all failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 retried, got %d", count)
	}
	updated, _ := s.GetDelivery(d.ID)
	if updated.Status != model.DeliveryRetrying {
		t.Fatalf("expected retrying, got %s", updated.Status)
	}
}
