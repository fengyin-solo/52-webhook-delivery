package handler

import (
	"net/http"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerFilterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/filters", s.createFilter)
	mux.HandleFunc("GET /api/filters", s.listFilters)
	mux.HandleFunc("GET /api/filters/{id}", s.getFilter)
	mux.HandleFunc("PUT /api/filters/{id}", s.updateFilter)
	mux.HandleFunc("DELETE /api/filters/{id}", s.deleteFilter)
	mux.HandleFunc("POST /api/filters/{id}/match", s.matchFilter)
}

type createFilterRequest struct {
	Name     string `json:"name"`
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
	Status   string `json:"status"`
}

func (s *Server) createFilter(w http.ResponseWriter, r *http.Request) {
	var req createFilterRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.CreateFilter(model.Filter{Name: req.Name, Field: req.Field, Operator: req.Operator, Value: req.Value, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, f)
}

func (s *Server) listFilters(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FilterFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListFilters(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getFilter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	f, err := s.svc.GetFilter(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

type updateFilterRequest struct {
	Name     string `json:"name"`
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
	Status   string `json:"status"`
}

func (s *Server) updateFilter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateFilterRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.UpdateFilter(id, model.Filter{Name: req.Name, Field: req.Field, Operator: req.Operator, Value: req.Value, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

func (s *Server) deleteFilter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteFilter(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type matchFilterRequest struct {
	Payload string `json:"payload"`
}

type matchFilterResponse struct {
	Match bool `json:"match"`
}

func (s *Server) matchFilter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req matchFilterRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.GetFilter(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, matchFilterResponse{Match: f.MatchPayload(req.Payload)})
}
