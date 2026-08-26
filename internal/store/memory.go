package store

import (
	"sync"

	"datamasking/internal/model"
)

type MemoryStore struct {
	mu           sync.RWMutex
	dataClasses  map[string]*model.DataClass
	fieldConfigs map[string]*model.FieldConfig
	maskRules    map[string]*model.MaskRule
	maskTasks    map[string]*model.MaskTask
	maskRecords  map[string]*model.MaskRecord
	policies     map[string]*model.Policy
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		dataClasses:  make(map[string]*model.DataClass),
		fieldConfigs: make(map[string]*model.FieldConfig),
		maskRules:    make(map[string]*model.MaskRule),
		maskTasks:    make(map[string]*model.MaskTask),
		maskRecords:  make(map[string]*model.MaskRecord),
		policies:     make(map[string]*model.Policy),
	}
}

var _ Store = (*MemoryStore)(nil)
