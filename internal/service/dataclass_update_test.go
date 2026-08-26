package service

import (
	"testing"

	"datamasking/internal/config"
	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/logger"
)

func TestFailedDataClassRenameDoesNotLeakIntoStore(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
	first, err := svc.CreateDataClass(model.DataClass{Name: "phone", Level: "confidential"})
	if err != nil { t.Fatalf("create first class: %v", err) }
	if _, err = svc.CreateDataClass(model.DataClass{Name: "email", Level: "confidential"}); err != nil { t.Fatalf("create second class: %v", err) }
	if _, err = svc.UpdateDataClass(first.ID, model.DataClass{Name: "email", Level: "confidential"}); err == nil { t.Fatal("expected duplicate name error") }
	got, err := svc.GetDataClass(first.ID)
	if err != nil { t.Fatalf("get class: %v", err) }
	if got.Name != "phone" { t.Fatalf("failed rename leaked into stored class: %q", got.Name) }
}
