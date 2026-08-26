package service

import (
	"testing"

	"datamasking/internal/config"
	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/logger"
)

func TestDisabledRuleCannotMaskData(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
	rule, err := svc.CreateMaskRule(model.MaskRule{Name: "old-phone", MaskType: "mask", PrefixKeep: 3, SuffixKeep: 4, Enabled: false})
	if err != nil { t.Fatalf("create disabled rule: %v", err) }
	if _, err = svc.ApplyMaskRule(rule.ID, "13800138000"); err == nil { t.Fatal("disabled rule unexpectedly masked data") }
}
