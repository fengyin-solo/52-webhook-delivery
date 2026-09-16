package handler

import (
	"net/http"
	"time"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerDeliveryRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/deliveries", s.createDelivery)
	mux.HandleFunc("GET /api/deliveries", s.listDeliveries)
	mux.HandleFunc("GET /api/deliveries/{id}", s.getDelivery)
	mux.HandleFunc("PUT /api/deliveries/{id}", s.updateDelivery)
	mux.HandleFunc("DELETE /api/deliveries/{id}", s.deleteDelivery)
	mux.HandleFunc("POST /api/deliveries/{id}/retry", s.retryDelivery)
}

type createDeliveryRequest struct {
	EventID    string `json:"event_id"`
	EndpointID string `json:"endpoint_id"`
	Status     string `json:"status"`
}

func (s *Server) createDelivery(w http.ResponseWriter, r *http.Request) {
	var req createDeliveryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateDelivery(model.Delivery{EventID: req.EventID, EndpointID: req.EndpointID, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

func (s *Server) listDeliveries(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DeliveryFilter{
		EventID:    r.URL.Query().Get("event_id"),
		EndpointID: r.URL.Query().Get("endpoint_id"),
		Status:     r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListDeliveries(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDelivery(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := s.svc.GetDelivery(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

type updateDeliveryRequest struct {
	Status           string     `json:"status"`
	LastResponseCode int        `json:"last_response_code"`
	LastResponseBody string     `json:"last_response_body"`
	NextRetryAt      *time.Time `json:"next_retry_at"`
	DeliveredAt      *time.Time `json:"delivered_at"`
}

func (s *Server) updateDelivery(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateDeliveryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.UpdateDelivery(id, model.Delivery{Status: req.Status, LastResponseCode: req.LastResponseCode, LastResponseBody: req.LastResponseBody, NextRetryAt: req.NextRetryAt, DeliveredAt: req.DeliveredAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteDelivery(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteDelivery(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) retryDelivery(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := s.svc.RetryDelivery(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}
