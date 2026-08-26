package service

import (
	"testing"

	"datamasking/internal/config"
	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/logger"
)

func TestFailedMaskTaskUpdateDoesNotLeakState(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
	task, err := svc.CreateMaskTask(model.MaskTask{Name: "nightly", SourceType: "file", TotalRecords: 5})
	if err != nil { t.Fatalf("create task: %v", err) }
	_, err = svc.UpdateMaskTask(task.ID, model.MaskTask{Name: "broken", SourceType: "file", TotalRecords: 5, RuleIDs: []string{"missing"}})
	if err == nil { t.Fatal("expected missing rule error") }
	got, err := svc.GetMaskTask(task.ID)
	if err != nil { t.Fatalf("get task: %v", err) }
	if got.Name != "nightly" || len(got.RuleIDs) != 0 { t.Fatalf("failed task update leaked state: %+v", got) }
}
