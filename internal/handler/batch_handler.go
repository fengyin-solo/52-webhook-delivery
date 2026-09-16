package handler

import (
	"net/http"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerBatchRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/batch/events", s.batchCreateEvents)
	mux.HandleFunc("POST /api/batch/endpoints/disable", s.batchDisableEndpoints)
}

func (s *Server) batchCreateEvents(w http.ResponseWriter, r *http.Request) {
	var req model.BatchCreateEventRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.BatchCreateEvents(req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) batchDisableEndpoints(w http.ResponseWriter, r *http.Request) {
	var req model.BatchDisableEndpointRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.BatchDisableEndpoints(req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
