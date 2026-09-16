package handler

import (
	"net/http"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerEventTypeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/event-types", s.createEventType)
	mux.HandleFunc("GET /api/event-types", s.listEventTypes)
	mux.HandleFunc("GET /api/event-types/{id}", s.getEventType)
	mux.HandleFunc("PUT /api/event-types/{id}", s.updateEventType)
	mux.HandleFunc("DELETE /api/event-types/{id}", s.deleteEventType)
}

type createEventTypeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      string `json:"schema"`
	Status      string `json:"status"`
}

func (s *Server) createEventType(w http.ResponseWriter, r *http.Request) {
	var req createEventTypeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	et, err := s.svc.CreateEventType(model.EventType{Name: req.Name, Description: req.Description, Schema: req.Schema, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, et)
}

func (s *Server) listEventTypes(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.EventTypeFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListEventTypes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getEventType(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	et, err := s.svc.GetEventType(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, et)
}

type updateEventTypeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      string `json:"schema"`
	Status      string `json:"status"`
}

func (s *Server) updateEventType(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateEventTypeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	et, err := s.svc.UpdateEventType(id, model.EventType{Name: req.Name, Description: req.Description, Schema: req.Schema, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, et)
}

func (s *Server) deleteEventType(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteEventType(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
