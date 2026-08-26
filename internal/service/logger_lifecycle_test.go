package service

import (
	"testing"

	"datamasking/internal/config"
	"datamasking/internal/model"
	"datamasking/internal/store"
)

func TestDataClassCreationAllowsOptionalLogger(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("optional logger panicked after save: %v", recovered)
		}
	}()
	svc := New(store.NewMemoryStore(), nil, &config.Config{MaxPageSize: 100})
	created, err := svc.CreateDataClass(model.DataClass{Name: "phone", Level: "public"})
	if err != nil {
		t.Fatalf("create class: %v", err)
	}
	if created == nil || created.Name != "phone" {
		t.Fatalf("create result missing: %+v", created)
	}
}
