package handler

import (
	"net/http"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerSigningKeyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/signing-keys", s.createSigningKey)
	mux.HandleFunc("GET /api/signing-keys", s.listSigningKeys)
	mux.HandleFunc("GET /api/signing-keys/{id}", s.getSigningKey)
	mux.HandleFunc("PUT /api/signing-keys/{id}", s.updateSigningKey)
	mux.HandleFunc("DELETE /api/signing-keys/{id}", s.deleteSigningKey)
	mux.HandleFunc("POST /api/signing-keys/verify", s.verifyHMAC)
}

type createSigningKeyRequest struct {
	EndpointID string `json:"endpoint_id"`
	KeyID      string `json:"key_id"`
	Algorithm  string `json:"algorithm"`
	KeyValue   string `json:"key_value"`
	Status     string `json:"status"`
}

func (s *Server) createSigningKey(w http.ResponseWriter, r *http.Request) {
	var req createSigningKeyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sk, err := s.svc.CreateSigningKey(model.SigningKey{EndpointID: req.EndpointID, KeyID: req.KeyID, Algorithm: req.Algorithm, KeyValue: req.KeyValue, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sk)
}

func (s *Server) listSigningKeys(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SigningKeyFilter{
		EndpointID: r.URL.Query().Get("endpoint_id"),
		Status:     r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListSigningKeys(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSigningKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sk, err := s.svc.GetSigningKey(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sk)
}

type updateSigningKeyRequest struct {
	KeyValue string `json:"key_value"`
	Status   string `json:"status"`
}

func (s *Server) updateSigningKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateSigningKeyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sk, err := s.svc.UpdateSigningKey(id, model.SigningKey{KeyValue: req.KeyValue, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sk)
}

func (s *Server) deleteSigningKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSigningKey(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type verifyHMACRequest struct {
	Key       string `json:"key"`
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}

type verifyHMACResponse struct {
	Valid    bool   `json:"valid"`
	Computed string `json:"computed"`
}

func (s *Server) verifyHMAC(w http.ResponseWriter, r *http.Request) {
	var req verifyHMACRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	computed := s.svc.ComputeHMAC(req.Key, req.Payload)
	valid := s.svc.VerifyHMAC(req.Key, req.Payload, req.Signature)
	httpx.OK(w, verifyHMACResponse{Valid: valid, Computed: computed})
}
