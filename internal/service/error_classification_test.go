package service

import (
	"testing"

	"datamasking/internal/config"
	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/logger"
)

func TestWrappedRuleValidationRemainsClassifiable(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
	_, err := svc.CreateMaskTask(model.MaskTask{Name: "x", SourceType: "file", TotalRecords: 1, RuleIDs: []string{"missing"}})
	if err == nil {
		t.Fatal("expected missing rule validation")
	}
	if !model.IsValidationError(err) {
		t.Fatalf("wrapped validation was not classifiable: %v", err)
	}
}
