package service

import (
	"testing"

	"datamasking/internal/config"
	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/logger"
)

func TestBatchMaskRecordsIsAtomic(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
	task, err := svc.CreateMaskTask(model.MaskTask{Name: "daily", SourceType: "file", TotalRecords: 2})
	if err != nil { t.Fatalf("create task: %v", err) }
	_, err = svc.BatchCreateMaskRecords([]model.MaskRecord{
		{MaskTaskID: task.ID, FieldName: "email", RuleID: "r1"},
		{MaskTaskID: "missing-task", FieldName: "phone", RuleID: "r2"},
	})
	if err == nil { t.Fatal("expected batch failure") }
	items, total, err := svc.ListMaskRecords(model.MaskRecordFilter{}, 1, 20)
	if err != nil { t.Fatalf("list records: %v", err) }
	if total != 0 || len(items) != 0 { t.Fatalf("partial records remained after failure: total=%d items=%d", total, len(items)) }
}
