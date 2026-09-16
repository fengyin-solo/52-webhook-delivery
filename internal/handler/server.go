// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"webhook/internal/config"
	"webhook/internal/model"
	"webhook/internal/service"
	"webhook/internal/store"
	"webhook/pkg/httpx"
	"webhook/pkg/logger"
)

type Server struct {
	svc *service.Service
	log *logger.Logger
	cfg *config.Config
}

func NewServer(svc *service.Service, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{svc: svc, log: log, cfg: cfg}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerEndpointRoutes(mux)
	s.registerSubscriptionRoutes(mux)
	s.registerEventTypeRoutes(mux)
	s.registerEventRoutes(mux)
	s.registerDeliveryRoutes(mux)
	s.registerDeliveryAttemptRoutes(mux)
	s.registerRetryPolicyRoutes(mux)
	s.registerSigningKeyRoutes(mux)
	s.registerFilterRoutes(mux)
	s.registerAuditLogRoutes(mux)
	s.registerWebhookLogRoutes(mux)
	s.registerEndpointGroupRoutes(mux)
	s.registerMessageTemplateRoutes(mux)
	s.registerEndpointHealthRoutes(mux)
	s.registerCircuitBreakerRoutes(mux)
	s.registerRateLimitRuleRoutes(mux)
	s.registerHealthRoutes(mux)
	s.registerStatsRoutes(mux)
	s.registerBatchRoutes(mux)
	s.registerSchedulerRoutes(mux)
	s.registerExportRoutes(mux)
	mux.Handle("GET /", http.FileServer(http.Dir("web")))
	return s.loggingMiddleware(s.recoveryMiddleware(s.apiKeyMiddleware(s.rateLimitMiddleware(mux))))
}

func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.log.Errorf("panic: %v\n%s", rec, debug.Stack())
				httpx.InternalError(w, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.APIKey == "" {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			key := r.Header.Get("X-API-Key")
			if key == "" {
				key = r.URL.Query().Get("api_key")
			}
			if key != s.cfg.APIKey {
				httpx.Unauthorized(w, "无效的 API Key")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

type rateLimiter struct {
	mu      sync.Mutex
	clients map[string][]time.Time
	limit   int
	window  time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{clients: make(map[string][]time.Time), limit: limit, window: window}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-rl.window)
	times := rl.clients[ip]
	valid := make([]time.Time, 0, len(times))
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	if len(valid) >= rl.limit {
		rl.clients[ip] = valid
		return false
	}
	valid = append(valid, now)
	rl.clients[ip] = valid
	return true
}

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	limiter := newRateLimiter(100, time.Minute)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		if !limiter.allow(ip) {
			httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}
