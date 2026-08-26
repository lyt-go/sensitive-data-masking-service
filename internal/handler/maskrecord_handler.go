package handler

import (
	"net/http"

	"datamasking/internal/model"
	"datamasking/pkg/httpx"
)

func (s *Server) registerMaskRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/mask-records", s.createMaskRecord)
	mux.HandleFunc("GET /api/mask-records", s.listMaskRecords)
	mux.HandleFunc("GET /api/mask-records/{id}", s.getMaskRecord)
	mux.HandleFunc("PUT /api/mask-records/{id}", s.updateMaskRecord)
	mux.HandleFunc("DELETE /api/mask-records/{id}", s.deleteMaskRecord)
	mux.HandleFunc("POST /api/mask-records/batch", s.batchCreateMaskRecords)
}

type createMaskRecordRequest struct {
	MaskTaskID     string `json:"mask_task_id"`
	FieldName      string `json:"field_name"`
	RuleID         string `json:"rule_id"`
	OriginalSample string `json:"original_sample"`
	MaskedSample   string `json:"masked_sample"`
}

func (s *Server) createMaskRecord(w http.ResponseWriter, r *http.Request) {
	var req createMaskRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.CreateMaskRecord(model.MaskRecord{
		MaskTaskID: req.MaskTaskID, FieldName: req.FieldName, RuleID: req.RuleID,
		OriginalSample: req.OriginalSample, MaskedSample: req.MaskedSample,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, m)
}

func (s *Server) listMaskRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MaskRecordFilter{
		MaskTaskID: r.URL.Query().Get("mask_task_id"),
		FieldName:  r.URL.Query().Get("field_name"),
		RuleID:     r.URL.Query().Get("rule_id"),
	}
	items, total, err := s.svc.ListMaskRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMaskRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, err := s.svc.GetMaskRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

type updateMaskRecordRequest struct {
	MaskTaskID     string `json:"mask_task_id"`
	FieldName      string `json:"field_name"`
	RuleID         string `json:"rule_id"`
	OriginalSample string `json:"original_sample"`
	MaskedSample   string `json:"masked_sample"`
}

func (s *Server) updateMaskRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateMaskRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.UpdateMaskRecord(id, model.MaskRecord{
		MaskTaskID: req.MaskTaskID, FieldName: req.FieldName, RuleID: req.RuleID,
		OriginalSample: req.OriginalSample, MaskedSample: req.MaskedSample,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) deleteMaskRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMaskRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchCreateMaskRecordsRequest struct {
	Records []createMaskRecordRequest `json:"records"`
}

type batchCreateMaskRecordsResponse struct {
	Count int `json:"count"`
}

func (s *Server) batchCreateMaskRecords(w http.ResponseWriter, r *http.Request) {
	var req batchCreateMaskRecordsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.MaskRecord, len(req.Records))
	for i, rec := range req.Records {
		inputs[i] = model.MaskRecord{
			MaskTaskID: rec.MaskTaskID, FieldName: rec.FieldName, RuleID: rec.RuleID,
			OriginalSample: rec.OriginalSample, MaskedSample: rec.MaskedSample,
		}
	}
	cnt, err := s.svc.BatchCreateMaskRecords(inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, batchCreateMaskRecordsResponse{Count: cnt})
}
