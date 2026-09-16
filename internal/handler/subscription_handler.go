package handler

import (
	"net/http"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerSubscriptionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/subscriptions", s.createSubscription)
	mux.HandleFunc("GET /api/subscriptions", s.listSubscriptions)
	mux.HandleFunc("GET /api/subscriptions/{id}", s.getSubscription)
	mux.HandleFunc("PUT /api/subscriptions/{id}", s.updateSubscription)
	mux.HandleFunc("DELETE /api/subscriptions/{id}", s.deleteSubscription)
}

type createSubscriptionRequest struct {
	EndpointID string `json:"endpoint_id"`
	EventType  string `json:"event_type"`
	FilterID   string `json:"filter_id"`
	Status     string `json:"status"`
}

func (s *Server) createSubscription(w http.ResponseWriter, r *http.Request) {
	var req createSubscriptionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sub, err := s.svc.CreateSubscription(model.Subscription{EndpointID: req.EndpointID, EventType: req.EventType, FilterID: req.FilterID, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sub)
}

func (s *Server) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SubscriptionFilter{
		EndpointID: r.URL.Query().Get("endpoint_id"),
		EventType:  r.URL.Query().Get("event_type"),
		Status:     r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListSubscriptions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sub, err := s.svc.GetSubscription(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

type updateSubscriptionRequest struct {
	EventType string `json:"event_type"`
	FilterID  string `json:"filter_id"`
	Status    string `json:"status"`
}

func (s *Server) updateSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateSubscriptionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sub, err := s.svc.UpdateSubscription(id, model.Subscription{EventType: req.EventType, FilterID: req.FilterID, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

func (s *Server) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSubscription(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
