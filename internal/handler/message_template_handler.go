package handler

import (
	"net/http"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerMessageTemplateRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/message-templates", s.createMessageTemplate)
	mux.HandleFunc("GET /api/message-templates", s.listMessageTemplates)
	mux.HandleFunc("GET /api/message-templates/{id}", s.getMessageTemplate)
	mux.HandleFunc("PUT /api/message-templates/{id}", s.updateMessageTemplate)
	mux.HandleFunc("DELETE /api/message-templates/{id}", s.deleteMessageTemplate)
}

type createMessageTemplateRequest struct {
	Name        string `json:"name"`
	EventType   string `json:"event_type"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
	Status      string `json:"status"`
}

func (s *Server) createMessageTemplate(w http.ResponseWriter, r *http.Request) {
	var req createMessageTemplateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	mt, err := s.svc.CreateMessageTemplate(model.MessageTemplate{Name: req.Name, EventType: req.EventType, ContentType: req.ContentType, Body: req.Body, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, mt)
}

func (s *Server) listMessageTemplates(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MessageTemplateFilter{
		EventType: r.URL.Query().Get("event_type"),
		Status:    r.URL.Query().Get("status"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListMessageTemplates(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMessageTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	mt, err := s.svc.GetMessageTemplate(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, mt)
}

type updateMessageTemplateRequest struct {
	Name        string `json:"name"`
	EventType   string `json:"event_type"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
	Status      string `json:"status"`
}

func (s *Server) updateMessageTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateMessageTemplateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	mt, err := s.svc.UpdateMessageTemplate(id, model.MessageTemplate{Name: req.Name, EventType: req.EventType, ContentType: req.ContentType, Body: req.Body, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, mt)
}

func (s *Server) deleteMessageTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMessageTemplate(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
