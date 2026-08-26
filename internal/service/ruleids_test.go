package service

import (
	"testing"

	"datamasking/internal/config"
	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/logger"
)

func TestMaskTaskRuleIDsStayIsolated(t *testing.T) {
	st := store.NewMemoryStore()
	svc := New(st, logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
	rule, err := svc.CreateMaskRule(model.MaskRule{Name: "phone", MaskType: "mask", PrefixKeep: 3, SuffixKeep: 4, Enabled: true})
	if err != nil { t.Fatalf("create rule: %v", err) }
	ids := []string{rule.ID}
	task, err := svc.CreateMaskTask(model.MaskTask{Name: "import", SourceType: "csv", TotalRecords: 1, RuleIDs: ids})
	if err != nil { t.Fatalf("create task: %v", err) }
	ids[0] = "changed-by-caller"
	got, err := svc.GetMaskTask(task.ID)
	if err != nil { t.Fatalf("get task: %v", err) }
	if len(got.RuleIDs) != 1 || got.RuleIDs[0] != rule.ID { t.Fatalf("task rule ids were changed: %#v", got.RuleIDs) }
	if err := svc.DeleteMaskRule(rule.ID); err != store.ErrConflict { t.Fatalf("expected referenced rule conflict, got %v", err) }
}
