// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"datamasking/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	CreateDataClass(d *model.DataClass) error
	GetDataClass(id string) (*model.DataClass, error)
	ListDataClasses() []*model.DataClass
	UpdateDataClass(d *model.DataClass) error
	DeleteDataClass(id string) error

	CreateFieldConfig(f *model.FieldConfig) error
	GetFieldConfig(id string) (*model.FieldConfig, error)
	ListFieldConfigs() []*model.FieldConfig
	UpdateFieldConfig(f *model.FieldConfig) error
	DeleteFieldConfig(id string) error

	CreateMaskRule(r *model.MaskRule) error
	GetMaskRule(id string) (*model.MaskRule, error)
	ListMaskRules() []*model.MaskRule
	UpdateMaskRule(r *model.MaskRule) error
	DeleteMaskRule(id string) error

	CreateMaskTask(t *model.MaskTask) error
	GetMaskTask(id string) (*model.MaskTask, error)
	ListMaskTasks() []*model.MaskTask
	UpdateMaskTask(t *model.MaskTask) error
	DeleteMaskTask(id string) error

	CreateMaskRecord(m *model.MaskRecord) error
	GetMaskRecord(id string) (*model.MaskRecord, error)
	ListMaskRecords() []*model.MaskRecord
	UpdateMaskRecord(m *model.MaskRecord) error
	DeleteMaskRecord(id string) error
	BatchCreateMaskRecords(records []*model.MaskRecord) error

	CreatePolicy(p *model.Policy) error
	GetPolicy(id string) (*model.Policy, error)
	ListPolicies() []*model.Policy
	UpdatePolicy(p *model.Policy) error
	DeletePolicy(id string) error
}
