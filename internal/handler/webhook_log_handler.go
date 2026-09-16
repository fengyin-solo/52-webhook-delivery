package handler

import (
	"net/http"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerWebhookLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/webhook-logs", s.createWebhookLog)
	mux.HandleFunc("GET /api/webhook-logs", s.listWebhookLogs)
	mux.HandleFunc("GET /api/webhook-logs/{id}", s.getWebhookLog)
	mux.HandleFunc("PUT /api/webhook-logs/{id}", s.updateWebhookLog)
	mux.HandleFunc("DELETE /api/webhook-logs/{id}", s.deleteWebhookLog)
}

type createWebhookLogRequest struct {
	DeliveryID   string `json:"delivery_id"`
	EndpointID   string `json:"endpoint_id"`
	URL          string `json:"url"`
	RequestBody  string `json:"request_body"`
	ResponseBody string `json:"response_body"`
	StatusCode   int    `json:"status_code"`
	Result       string `json:"result"`
	ErrorMessage string `json:"error_message"`
	DurationMs   int64  `json:"duration_ms"`
}

func (s *Server) createWebhookLog(w http.ResponseWriter, r *http.Request) {
	var req createWebhookLogRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wl, err := s.svc.CreateWebhookLog(model.WebhookLog{DeliveryID: req.DeliveryID, EndpointID: req.EndpointID, URL: req.URL, RequestBody: req.RequestBody, ResponseBody: req.ResponseBody, StatusCode: req.StatusCode, Result: req.Result, ErrorMessage: req.ErrorMessage, DurationMs: req.DurationMs})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, wl)
}

func (s *Server) listWebhookLogs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.WebhookLogFilter{
		DeliveryID: r.URL.Query().Get("delivery_id"),
		EndpointID: r.URL.Query().Get("endpoint_id"),
		Result:     r.URL.Query().Get("result"),
	}
	items, total, err := s.svc.ListWebhookLogs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getWebhookLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	wl, err := s.svc.GetWebhookLog(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wl)
}

type updateWebhookLogRequest struct {
	ResponseBody string `json:"response_body"`
	StatusCode   int    `json:"status_code"`
	Result       string `json:"result"`
	ErrorMessage string `json:"error_message"`
	DurationMs   int64  `json:"duration_ms"`
}

func (s *Server) updateWebhookLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateWebhookLogRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wl, err := s.svc.UpdateWebhookLog(id, model.WebhookLog{ResponseBody: req.ResponseBody, StatusCode: req.StatusCode, Result: req.Result, ErrorMessage: req.ErrorMessage, DurationMs: req.DurationMs})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wl)
}

func (s *Server) deleteWebhookLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteWebhookLog(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
