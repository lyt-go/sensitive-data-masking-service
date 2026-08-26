package model

import (
	"strings"
	"time"
)

const (
	FieldMaskTypeFull    = "full"
	FieldMaskTypePartial = "partial"
	FieldMaskTypeHash    = "hash"
	FieldMaskTypeReplace = "replace"
)

var validFieldMaskTypes = map[string]bool{
	FieldMaskTypeFull:    true,
	FieldMaskTypePartial: true,
	FieldMaskTypeHash:    true,
	FieldMaskTypeReplace: true,
}

type FieldConfig struct {
	ID          string    `json:"id"`
	FieldName   string    `json:"field_name"`
	DataClassID string    `json:"data_class_id"`
	MaskType    string    `json:"mask_type"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
}

func (f *FieldConfig) Validate() error {
	f.FieldName = strings.TrimSpace(f.FieldName)
	if f.FieldName == "" {
		return NewValidationError("field_name", "字段名称不能为空")
	}
	f.DataClassID = strings.TrimSpace(f.DataClassID)
	if f.DataClassID == "" {
		return NewValidationError("data_class_id", "数据分类 ID 不能为空")
	}
	f.MaskType = strings.TrimSpace(f.MaskType)
	if f.MaskType == "" {
		return NewValidationError("mask_type", "脱敏类型不能为空")
	}
	if !validFieldMaskTypes[f.MaskType] {
		return NewValidationError("mask_type", "脱敏类型不合法")
	}
	return nil
}

type FieldConfigFilter struct {
	DataClassID string
	MaskType    string
	Enabled     *bool
	Keyword     string
}

func (f FieldConfigFilter) Match(fc *FieldConfig) bool {
	if f.DataClassID != "" && fc.DataClassID != f.DataClassID {
		return false
	}
	if f.MaskType != "" && fc.MaskType != f.MaskType {
		return false
	}
	if f.Enabled != nil && fc.Enabled != *f.Enabled {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(fc.FieldName), k) {
			return false
		}
	}
	return true
}
