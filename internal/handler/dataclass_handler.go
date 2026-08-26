package handler

import (
	"net/http"

	"datamasking/internal/model"
	"datamasking/pkg/httpx"
)

func (s *Server) registerDataClassRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/data-classes", s.createDataClass)
	mux.HandleFunc("GET /api/data-classes", s.listDataClasses)
	mux.HandleFunc("GET /api/data-classes/{id}", s.getDataClass)
	mux.HandleFunc("PUT /api/data-classes/{id}", s.updateDataClass)
	mux.HandleFunc("DELETE /api/data-classes/{id}", s.deleteDataClass)
}

type createDataClassRequest struct {
	Name        string `json:"name"`
	Level       string `json:"level"`
	Description string `json:"description"`
}

func (s *Server) createDataClass(w http.ResponseWriter, r *http.Request) {
	var req createDataClassRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateDataClass(model.DataClass{Name: req.Name, Level: req.Level, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

func (s *Server) listDataClasses(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DataClassFilter{
		Level:   r.URL.Query().Get("level"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListDataClasses(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDataClass(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := s.svc.GetDataClass(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

type updateDataClassRequest struct {
	Name        string `json:"name"`
	Level       string `json:"level"`
	Description string `json:"description"`
}

func (s *Server) updateDataClass(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateDataClassRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.UpdateDataClass(id, model.DataClass{Name: req.Name, Level: req.Level, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteDataClass(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteDataClass(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
