package handler

import (
	"net/http"

	"webhook/internal/model"
	"webhook/pkg/httpx"
)

func (s *Server) registerEndpointGroupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/endpoint-groups", s.createEndpointGroup)
	mux.HandleFunc("GET /api/endpoint-groups", s.listEndpointGroups)
	mux.HandleFunc("GET /api/endpoint-groups/{id}", s.getEndpointGroup)
	mux.HandleFunc("PUT /api/endpoint-groups/{id}", s.updateEndpointGroup)
	mux.HandleFunc("DELETE /api/endpoint-groups/{id}", s.deleteEndpointGroup)
}

type createEndpointGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (s *Server) createEndpointGroup(w http.ResponseWriter, r *http.Request) {
	var req createEndpointGroupRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	eg, err := s.svc.CreateEndpointGroup(model.EndpointGroup{Name: req.Name, Description: req.Description, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, eg)
}

func (s *Server) listEndpointGroups(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.EndpointGroupFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListEndpointGroups(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getEndpointGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	eg, err := s.svc.GetEndpointGroup(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, eg)
}

type updateEndpointGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (s *Server) updateEndpointGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateEndpointGroupRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	eg, err := s.svc.UpdateEndpointGroup(id, model.EndpointGroup{Name: req.Name, Description: req.Description, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, eg)
}

func (s *Server) deleteEndpointGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteEndpointGroup(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
