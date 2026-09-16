package handler

import (
	"net/http"

	"webhook/pkg/httpx"
)

func (s *Server) registerSchedulerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/scheduler/process-pending", s.processPending)
	mux.HandleFunc("GET /api/scheduler/overdue", s.listOverdue)
	mux.HandleFunc("POST /api/scheduler/retry-all-failed", s.retryAllFailed)
}

func (s *Server) processPending(w http.ResponseWriter, r *http.Request) {
	processed, err := s.svc.ProcessPendingDeliveries()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"processed": processed})
}

func (s *Server) listOverdue(w http.ResponseWriter, r *http.Request) {
	items := s.svc.GetOverdueDeliveries()
	httpx.OK(w, items)
}

func (s *Server) retryAllFailed(w http.ResponseWriter, r *http.Request) {
	count, err := s.svc.RetryAllFailedDeliveries()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"retried": count})
}
