package handler

import (
	"net/http"

	"datamasking/internal/model"
	"datamasking/pkg/httpx"
)

func (s *Server) registerMaskTaskRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/mask-tasks", s.createMaskTask)
	mux.HandleFunc("GET /api/mask-tasks", s.listMaskTasks)
	mux.HandleFunc("GET /api/mask-tasks/{id}", s.getMaskTask)
	mux.HandleFunc("PUT /api/mask-tasks/{id}", s.updateMaskTask)
	mux.HandleFunc("DELETE /api/mask-tasks/{id}", s.deleteMaskTask)
	mux.HandleFunc("POST /api/mask-tasks/{id}/transition", s.transitionMaskTask)
	mux.HandleFunc("POST /api/mask-tasks/{id}/advance", s.advanceMaskTask)
}

type createMaskTaskRequest struct {
	Name         string   `json:"name"`
	SourceType   string   `json:"source_type"`
	TotalRecords int      `json:"total_records"`
	RuleIDs      []string `json:"rule_ids"`
}

func (s *Server) createMaskTask(w http.ResponseWriter, r *http.Request) {
	var req createMaskTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateMaskTask(model.MaskTask{
		Name: req.Name, SourceType: req.SourceType,
		TotalRecords: req.TotalRecords, RuleIDs: req.RuleIDs,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listMaskTasks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MaskTaskFilter{
		Status:     r.URL.Query().Get("status"),
		SourceType: r.URL.Query().Get("source_type"),
		Keyword:    r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListMaskTasks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMaskTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetMaskTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type updateMaskTaskRequest struct {
	Name         string   `json:"name"`
	SourceType   string   `json:"source_type"`
	TotalRecords int      `json:"total_records"`
	RuleIDs      []string `json:"rule_ids"`
}

func (s *Server) updateMaskTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateMaskTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateMaskTask(id, model.MaskTask{
		Name: req.Name, SourceType: req.SourceType,
		TotalRecords: req.TotalRecords, RuleIDs: req.RuleIDs,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteMaskTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMaskTask(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionMaskTaskRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionMaskTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionMaskTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.TransitionMaskTaskStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type advanceMaskTaskRequest struct {
	Delta int `json:"delta"`
}

func (s *Server) advanceMaskTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req advanceMaskTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.AdvanceMaskTaskProgress(id, req.Delta)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}
