package handler

import (
	"net/http"
	"time"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerEndpointHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/endpoint-healths", s.createEndpointHealth)
	mux.HandleFunc("GET /api/endpoint-healths", s.listEndpointHealths)
	mux.HandleFunc("GET /api/endpoint-healths/{id}", s.getEndpointHealth)
	mux.HandleFunc("PUT /api/endpoint-healths/{id}", s.updateEndpointHealth)
	mux.HandleFunc("DELETE /api/endpoint-healths/{id}", s.deleteEndpointHealth)
}

type createEndpointHealthRequest struct {
	EndpointID   string    `json:"endpoint_id"`
	Status       string    `json:"status"`
	ResponseTime int64     `json:"response_time"`
	LastChecked  time.Time `json:"last_checked"`
	ErrorMessage string    `json:"error_message"`
}

func (s *Server) createEndpointHealth(w http.ResponseWriter, r *http.Request) {
	var req createEndpointHealthRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	eh, err := s.svc.CreateEndpointHealth(model.EndpointHealth{EndpointID: req.EndpointID, Status: req.Status, ResponseTime: req.ResponseTime, LastChecked: req.LastChecked, ErrorMessage: req.ErrorMessage})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, eh)
}

func (s *Server) listEndpointHealths(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.EndpointHealthFilter{
		EndpointID: r.URL.Query().Get("endpoint_id"),
		Status:     r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListEndpointHealths(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getEndpointHealth(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	eh, err := s.svc.GetEndpointHealth(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, eh)
}

type updateEndpointHealthRequest struct {
	Status       string    `json:"status"`
	ResponseTime int64     `json:"response_time"`
	LastChecked  time.Time `json:"last_checked"`
	ErrorMessage string    `json:"error_message"`
}

func (s *Server) updateEndpointHealth(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateEndpointHealthRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	eh, err := s.svc.UpdateEndpointHealth(id, model.EndpointHealth{Status: req.Status, ResponseTime: req.ResponseTime, LastChecked: req.LastChecked, ErrorMessage: req.ErrorMessage})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, eh)
}

func (s *Server) deleteEndpointHealth(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteEndpointHealth(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
