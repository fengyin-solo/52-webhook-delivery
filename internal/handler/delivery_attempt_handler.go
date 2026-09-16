package handler

import (
	"net/http"
	"time"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerDeliveryAttemptRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/delivery-attempts", s.createDeliveryAttempt)
	mux.HandleFunc("GET /api/delivery-attempts", s.listDeliveryAttempts)
	mux.HandleFunc("GET /api/delivery-attempts/{id}", s.getDeliveryAttempt)
	mux.HandleFunc("PUT /api/delivery-attempts/{id}", s.updateDeliveryAttempt)
	mux.HandleFunc("DELETE /api/delivery-attempts/{id}", s.deleteDeliveryAttempt)
}

type createDeliveryAttemptRequest struct {
	DeliveryID   string    `json:"delivery_id"`
	AttemptNo    int       `json:"attempt_no"`
	RequestAt    time.Time `json:"request_at"`
	ResponseCode int       `json:"response_code"`
	DurationMs   int64     `json:"duration_ms"`
	Success      bool      `json:"success"`
	ResponseBody string    `json:"response_body"`
}

func (s *Server) createDeliveryAttempt(w http.ResponseWriter, r *http.Request) {
	var req createDeliveryAttemptRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	da, err := s.svc.CreateDeliveryAttempt(model.DeliveryAttempt{DeliveryID: req.DeliveryID, AttemptNo: req.AttemptNo, RequestAt: req.RequestAt, ResponseCode: req.ResponseCode, DurationMs: req.DurationMs, Success: req.Success, ResponseBody: req.ResponseBody})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, da)
}

func (s *Server) listDeliveryAttempts(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DeliveryAttemptFilter{
		DeliveryID: r.URL.Query().Get("delivery_id"),
	}
	items, total, err := s.svc.ListDeliveryAttempts(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDeliveryAttempt(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	da, err := s.svc.GetDeliveryAttempt(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, da)
}

type updateDeliveryAttemptRequest struct {
	ResponseCode int    `json:"response_code"`
	DurationMs   int64  `json:"duration_ms"`
	Success      bool   `json:"success"`
	ResponseBody string `json:"response_body"`
}

func (s *Server) updateDeliveryAttempt(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateDeliveryAttemptRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	da, err := s.svc.UpdateDeliveryAttempt(id, model.DeliveryAttempt{ResponseCode: req.ResponseCode, DurationMs: req.DurationMs, Success: req.Success, ResponseBody: req.ResponseBody})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, da)
}

func (s *Server) deleteDeliveryAttempt(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteDeliveryAttempt(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
