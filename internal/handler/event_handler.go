package handler

import (
	"net/http"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerEventRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/events", s.createEvent)
	mux.HandleFunc("GET /api/events", s.listEvents)
	mux.HandleFunc("GET /api/events/{id}", s.getEvent)
	mux.HandleFunc("PUT /api/events/{id}", s.updateEvent)
	mux.HandleFunc("DELETE /api/events/{id}", s.deleteEvent)
}

type createEventRequest struct {
	Type           string `json:"type"`
	Source         string `json:"source"`
	Payload        string `json:"payload"`
	IdempotencyKey string `json:"idempotency_key"`
	Status         string `json:"status"`
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	var req createEventRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.CreateEvent(model.Event{Type: req.Type, Source: req.Source, Payload: req.Payload, IdempotencyKey: req.IdempotencyKey, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

func (s *Server) listEvents(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.EventFilter{
		Type:   r.URL.Query().Get("type"),
		Source: r.URL.Query().Get("source"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListEvents(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.GetEvent(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

type updateEventRequest struct {
	Status string `json:"status"`
}

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateEventRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.UpdateEvent(id, model.Event{Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteEvent(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
