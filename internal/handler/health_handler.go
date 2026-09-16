package handler

import (
	"net/http"
	"runtime"
	"time"

	"webhook/pkg/httpx"
)

func (s *Server) registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/health", s.healthCheck)
	mux.HandleFunc("GET /api/health/system", s.systemInfo)
}

type healthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, healthResponse{Status: "ok", Timestamp: time.Now()})
}

type systemInfoResponse struct {
	GoVersion     string `json:"go_version"`
	NumGoroutine  int    `json:"num_goroutine"`
	NumCPU        int    `json:"num_cpu"`
	Timestamp     string `json:"timestamp"`
}

func (s *Server) systemInfo(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, systemInfoResponse{
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
		Timestamp:    time.Now().Format(time.RFC3339),
	})
}
