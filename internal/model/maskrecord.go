package model

import (
	"strings"
	"time"
)

type MaskRecord struct {
	ID             string    `json:"id"`
	MaskTaskID     string    `json:"mask_task_id"`
	FieldName      string    `json:"field_name"`
	RuleID         string    `json:"rule_id"`
	OriginalSample string    `json:"original_sample"`
	MaskedSample   string    `json:"masked_sample"`
	CreatedAt      time.Time `json:"created_at"`
}

func (m *MaskRecord) Validate() error {
	m.MaskTaskID = strings.TrimSpace(m.MaskTaskID)
	if m.MaskTaskID == "" {
		return NewValidationError("mask_task_id", "脱敏任务 ID 不能为空")
	}
	m.FieldName = strings.TrimSpace(m.FieldName)
	if m.FieldName == "" {
		return NewValidationError("field_name", "字段名称不能为空")
	}
	m.RuleID = strings.TrimSpace(m.RuleID)
	if m.RuleID == "" {
		return NewValidationError("rule_id", "规则 ID 不能为空")
	}
	return nil
}

type MaskRecordFilter struct {
	MaskTaskID string
	FieldName  string
	RuleID     string
}

func (f MaskRecordFilter) Match(m *MaskRecord) bool {
	if f.MaskTaskID != "" && m.MaskTaskID != f.MaskTaskID {
		return false
	}
	if f.FieldName != "" && m.FieldName != f.FieldName {
		return false
	}
	if f.RuleID != "" && m.RuleID != f.RuleID {
		return false
	}
	return true
}
