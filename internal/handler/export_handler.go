package handler

import (
	"net/http"

	"webhook/pkg/httpx"
)

func (s *Server) registerExportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/export/snapshot", s.exportSnapshot)
}

func (s *Server) exportSnapshot(w http.ResponseWriter, r *http.Request) {
	snapshot := s.svc.ExportSnapshot()
	httpx.OK(w, snapshot)
}
