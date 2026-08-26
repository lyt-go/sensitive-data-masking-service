package handler

import (
	"net/http"
	"strconv"

	"datamasking/internal/model"
	"datamasking/pkg/httpx"
)

func (s *Server) registerFieldConfigRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/field-configs", s.createFieldConfig)
	mux.HandleFunc("GET /api/field-configs", s.listFieldConfigs)
	mux.HandleFunc("GET /api/field-configs/{id}", s.getFieldConfig)
	mux.HandleFunc("PUT /api/field-configs/{id}", s.updateFieldConfig)
	mux.HandleFunc("DELETE /api/field-configs/{id}", s.deleteFieldConfig)
}

type createFieldConfigRequest struct {
	FieldName   string `json:"field_name"`
	DataClassID string `json:"data_class_id"`
	MaskType    string `json:"mask_type"`
	Enabled     bool   `json:"enabled"`
}

func (s *Server) createFieldConfig(w http.ResponseWriter, r *http.Request) {
	var req createFieldConfigRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.CreateFieldConfig(model.FieldConfig{FieldName: req.FieldName, DataClassID: req.DataClassID, MaskType: req.MaskType, Enabled: req.Enabled})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, f)
}

func (s *Server) listFieldConfigs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FieldConfigFilter{
		DataClassID: r.URL.Query().Get("data_class_id"),
		MaskType:    r.URL.Query().Get("mask_type"),
		Keyword:     r.URL.Query().Get("keyword"),
	}
	if v := r.URL.Query().Get("enabled"); v != "" {
		b, _ := strconv.ParseBool(v)
		filter.Enabled = &b
	}
	items, total, err := s.svc.ListFieldConfigs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getFieldConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	f, err := s.svc.GetFieldConfig(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

type updateFieldConfigRequest struct {
	FieldName   string `json:"field_name"`
	DataClassID string `json:"data_class_id"`
	MaskType    string `json:"mask_type"`
	Enabled     bool   `json:"enabled"`
}

func (s *Server) updateFieldConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateFieldConfigRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.UpdateFieldConfig(id, model.FieldConfig{FieldName: req.FieldName, DataClassID: req.DataClassID, MaskType: req.MaskType, Enabled: req.Enabled})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

func (s *Server) deleteFieldConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteFieldConfig(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
