package handler

import (
	"net/http"
	"strconv"

	"datamasking/internal/model"
	"datamasking/pkg/httpx"
)

func (s *Server) registerMaskRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/mask-rules", s.createMaskRule)
	mux.HandleFunc("GET /api/mask-rules", s.listMaskRules)
	mux.HandleFunc("GET /api/mask-rules/{id}", s.getMaskRule)
	mux.HandleFunc("PUT /api/mask-rules/{id}", s.updateMaskRule)
	mux.HandleFunc("DELETE /api/mask-rules/{id}", s.deleteMaskRule)
	mux.HandleFunc("POST /api/mask-rules/{id}/apply", s.applyMaskRule)
}

type createMaskRuleRequest struct {
	Name        string `json:"name"`
	MaskType    string `json:"mask_type"`
	Pattern     string `json:"pattern"`
	PrefixKeep  int    `json:"prefix_keep"`
	SuffixKeep  int    `json:"suffix_keep"`
	Replacement string `json:"replacement"`
	Enabled     bool   `json:"enabled"`
}

func (s *Server) createMaskRule(w http.ResponseWriter, r *http.Request) {
	var req createMaskRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	mr, err := s.svc.CreateMaskRule(model.MaskRule{
		Name: req.Name, MaskType: req.MaskType, Pattern: req.Pattern,
		PrefixKeep: req.PrefixKeep, SuffixKeep: req.SuffixKeep,
		Replacement: req.Replacement, Enabled: req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, mr)
}

func (s *Server) listMaskRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MaskRuleFilter{
		MaskType: r.URL.Query().Get("mask_type"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	if v := r.URL.Query().Get("enabled"); v != "" {
		b, _ := strconv.ParseBool(v)
		filter.Enabled = &b
	}
	items, total, err := s.svc.ListMaskRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMaskRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	mr, err := s.svc.GetMaskRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, mr)
}

type updateMaskRuleRequest struct {
	Name        string `json:"name"`
	MaskType    string `json:"mask_type"`
	Pattern     string `json:"pattern"`
	PrefixKeep  int    `json:"prefix_keep"`
	SuffixKeep  int    `json:"suffix_keep"`
	Replacement string `json:"replacement"`
	Enabled     bool   `json:"enabled"`
}

func (s *Server) updateMaskRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateMaskRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	mr, err := s.svc.UpdateMaskRule(id, model.MaskRule{
		Name: req.Name, MaskType: req.MaskType, Pattern: req.Pattern,
		PrefixKeep: req.PrefixKeep, SuffixKeep: req.SuffixKeep,
		Replacement: req.Replacement, Enabled: req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, mr)
}

func (s *Server) deleteMaskRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMaskRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type applyMaskRuleRequest struct {
	Input string `json:"input"`
}

type applyMaskRuleResponse struct {
	Result string `json:"result"`
}

func (s *Server) applyMaskRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req applyMaskRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.ApplyMaskRule(id, req.Input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, applyMaskRuleResponse{Result: result})
}
