package handler

import (
	"net/http"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerCircuitBreakerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/circuit-breakers", s.createCircuitBreaker)
	mux.HandleFunc("GET /api/circuit-breakers", s.listCircuitBreakers)
	mux.HandleFunc("GET /api/circuit-breakers/{id}", s.getCircuitBreaker)
	mux.HandleFunc("PUT /api/circuit-breakers/{id}", s.updateCircuitBreaker)
	mux.HandleFunc("DELETE /api/circuit-breakers/{id}", s.deleteCircuitBreaker)
	mux.HandleFunc("POST /api/circuit-breakers/{id}/record-failure", s.recordCircuitBreakerFailure)
	mux.HandleFunc("POST /api/circuit-breakers/{id}/record-success", s.recordCircuitBreakerSuccess)
}

type createCircuitBreakerRequest struct {
	EndpointID        string `json:"endpoint_id"`
	FailureThreshold  int    `json:"failure_threshold"`
	RecoveryTimeoutMs int    `json:"recovery_timeout_ms"`
	Status            string `json:"status"`
}

func (s *Server) createCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	var req createCircuitBreakerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	cb, err := s.svc.CreateCircuitBreaker(model.CircuitBreaker{EndpointID: req.EndpointID, FailureThreshold: req.FailureThreshold, RecoveryTimeoutMs: req.RecoveryTimeoutMs, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, cb)
}

func (s *Server) listCircuitBreakers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CircuitBreakerFilter{
		EndpointID: r.URL.Query().Get("endpoint_id"),
		Status:     r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListCircuitBreakers(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cb, err := s.svc.GetCircuitBreaker(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, cb)
}

type updateCircuitBreakerRequest struct {
	FailureThreshold  int    `json:"failure_threshold"`
	RecoveryTimeoutMs int    `json:"recovery_timeout_ms"`
	Status            string `json:"status"`
}

func (s *Server) updateCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateCircuitBreakerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	cb, err := s.svc.UpdateCircuitBreaker(id, model.CircuitBreaker{FailureThreshold: req.FailureThreshold, RecoveryTimeoutMs: req.RecoveryTimeoutMs, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, cb)
}

func (s *Server) deleteCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCircuitBreaker(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) recordCircuitBreakerFailure(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cb, err := s.svc.GetCircuitBreaker(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	cb, err = s.svc.RecordCircuitBreakerFailure(cb.EndpointID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, cb)
}

func (s *Server) recordCircuitBreakerSuccess(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cb, err := s.svc.GetCircuitBreaker(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	cb, err = s.svc.RecordCircuitBreakerSuccess(cb.EndpointID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, cb)
}
