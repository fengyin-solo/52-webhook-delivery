package handler

import (
	"net/http"

	"webhook/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.statsOverview)
	mux.HandleFunc("GET /api/stats/endpoints", s.statsEndpoints)
	mux.HandleFunc("GET /api/stats/event-types", s.statsEventTypes)
	mux.HandleFunc("GET /api/stats/timeline", s.statsTimeline)
	mux.HandleFunc("GET /api/stats/latency", s.statsLatency)
	mux.HandleFunc("GET /api/stats/webhook-logs", s.statsWebhookLogs)
	mux.HandleFunc("GET /api/stats/templates", s.statsTemplates)
	mux.HandleFunc("GET /api/stats/pending-count", s.statsPendingCount)
}

func (s *Server) statsOverview(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetDeliveryOverviewStats()
	httpx.OK(w, stats)
}

func (s *Server) statsEndpoints(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetEndpointPerformance()
	httpx.OK(w, stats)
}

func (s *Server) statsEventTypes(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetEventTypeDistribution()
	httpx.OK(w, stats)
}

func (s *Server) statsTimeline(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetDeliveryTimeline()
	httpx.OK(w, stats)
}

func (s *Server) statsLatency(w http.ResponseWriter, r *http.Request) {
	avg, p95 := s.svc.GetLatencyStats()
	httpx.OK(w, map[string]interface{}{"avg_ms": avg, "p95_ms": p95})
}

func (s *Server) statsWebhookLogs(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetWebhookLogResultStats()
	httpx.OK(w, stats)
}

func (s *Server) statsTemplates(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetTemplateUsageStats()
	httpx.OK(w, stats)
}

func (s *Server) statsPendingCount(w http.ResponseWriter, r *http.Request) {
	count := s.svc.GetPendingDeliveryCount()
	httpx.OK(w, map[string]interface{}{"pending_count": count})
}
