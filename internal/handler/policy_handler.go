package handler

import (
	"net/http"
	"strconv"

	"datamasking/internal/model"
	"datamasking/pkg/httpx"
)

func (s *Server) registerPolicyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/policies", s.createPolicy)
	mux.HandleFunc("GET /api/policies", s.listPolicies)
	mux.HandleFunc("GET /api/policies/{id}", s.getPolicy)
	mux.HandleFunc("PUT /api/policies/{id}", s.updatePolicy)
	mux.HandleFunc("DELETE /api/policies/{id}", s.deletePolicy)
	mux.HandleFunc("GET /api/policies/{id}/evaluate", s.evaluatePolicy)
}

type createPolicyRequest struct {
	Name     string   `json:"name"`
	Scope    string   `json:"scope"`
	RuleIDs  []string `json:"rule_ids"`
	Priority int      `json:"priority"`
	Enabled  bool     `json:"enabled"`
}

func (s *Server) createPolicy(w http.ResponseWriter, r *http.Request) {
	var req createPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreatePolicy(model.Policy{
		Name: req.Name, Scope: req.Scope, RuleIDs: req.RuleIDs,
		Priority: req.Priority, Enabled: req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listPolicies(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.PolicyFilter{
		Scope:   r.URL.Query().Get("scope"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	if v := r.URL.Query().Get("enabled"); v != "" {
		b, _ := strconv.ParseBool(v)
		filter.Enabled = &b
	}
	items, total, err := s.svc.ListPolicies(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := s.svc.GetPolicy(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

type updatePolicyRequest struct {
	Name     string   `json:"name"`
	Scope    string   `json:"scope"`
	RuleIDs  []string `json:"rule_ids"`
	Priority int      `json:"priority"`
	Enabled  bool     `json:"enabled"`
}

func (s *Server) updatePolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updatePolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdatePolicy(id, model.Policy{
		Name: req.Name, Scope: req.Scope, RuleIDs: req.RuleIDs,
		Priority: req.Priority, Enabled: req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deletePolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeletePolicy(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) evaluatePolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rules, err := s.svc.EvaluatePolicy(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rules)
}
