package handler

import (
	"net/http"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerRateLimitRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rate-limit-rules", s.createRateLimitRule)
	mux.HandleFunc("GET /api/rate-limit-rules", s.listRateLimitRules)
	mux.HandleFunc("GET /api/rate-limit-rules/{id}", s.getRateLimitRule)
	mux.HandleFunc("PUT /api/rate-limit-rules/{id}", s.updateRateLimitRule)
	mux.HandleFunc("DELETE /api/rate-limit-rules/{id}", s.deleteRateLimitRule)
}

type createRateLimitRuleRequest struct {
	Name        string `json:"name"`
	EndpointID  string `json:"endpoint_id"`
	EventType   string `json:"event_type"`
	MaxRequests int    `json:"max_requests"`
	WindowMs    int    `json:"window_ms"`
	Status      string `json:"status"`
}

func (s *Server) createRateLimitRule(w http.ResponseWriter, r *http.Request) {
	var req createRateLimitRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rlr, err := s.svc.CreateRateLimitRule(model.RateLimitRule{Name: req.Name, EndpointID: req.EndpointID, EventType: req.EventType, MaxRequests: req.MaxRequests, WindowMs: req.WindowMs, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rlr)
}

func (s *Server) listRateLimitRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RateLimitRuleFilter{
		EndpointID: r.URL.Query().Get("endpoint_id"),
		EventType:  r.URL.Query().Get("event_type"),
		Status:     r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListRateLimitRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRateLimitRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rlr, err := s.svc.GetRateLimitRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rlr)
}

type updateRateLimitRuleRequest struct {
	Name        string `json:"name"`
	EndpointID  string `json:"endpoint_id"`
	EventType   string `json:"event_type"`
	MaxRequests int    `json:"max_requests"`
	WindowMs    int    `json:"window_ms"`
	Status      string `json:"status"`
}

func (s *Server) updateRateLimitRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRateLimitRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rlr, err := s.svc.UpdateRateLimitRule(id, model.RateLimitRule{Name: req.Name, EndpointID: req.EndpointID, EventType: req.EventType, MaxRequests: req.MaxRequests, WindowMs: req.WindowMs, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rlr)
}

func (s *Server) deleteRateLimitRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRateLimitRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
