package handler

import (
	"net/http"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerEndpointRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/endpoints", s.createEndpoint)
	mux.HandleFunc("GET /api/endpoints", s.listEndpoints)
	mux.HandleFunc("GET /api/endpoints/{id}", s.getEndpoint)
	mux.HandleFunc("PUT /api/endpoints/{id}", s.updateEndpoint)
	mux.HandleFunc("DELETE /api/endpoints/{id}", s.deleteEndpoint)
}

type createEndpointRequest struct {
	URL         string `json:"url"`
	Secret      string `json:"secret"`
	Status      string `json:"status"`
	MaxRetries  int    `json:"max_retries"`
	TimeoutMs   int    `json:"timeout_ms"`
	Description string `json:"description"`
}

func (s *Server) createEndpoint(w http.ResponseWriter, r *http.Request) {
	var req createEndpointRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.CreateEndpoint(model.Endpoint{URL: req.URL, Secret: req.Secret, Status: req.Status, MaxRetries: req.MaxRetries, TimeoutMs: req.TimeoutMs, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

func (s *Server) listEndpoints(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.EndpointFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListEndpoints(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getEndpoint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.GetEndpoint(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

type updateEndpointRequest struct {
	URL         string `json:"url"`
	Secret      string `json:"secret"`
	Status      string `json:"status"`
	MaxRetries  int    `json:"max_retries"`
	TimeoutMs   int    `json:"timeout_ms"`
	Description string `json:"description"`
}

func (s *Server) updateEndpoint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateEndpointRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.UpdateEndpoint(id, model.Endpoint{URL: req.URL, Secret: req.Secret, Status: req.Status, MaxRetries: req.MaxRetries, TimeoutMs: req.TimeoutMs, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) deleteEndpoint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteEndpoint(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
