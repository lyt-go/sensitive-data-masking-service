package service

import (
	"sync"
	"testing"

	"datamasking/internal/config"
	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/logger"
)

func TestConcurrentMaskTaskProgress(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
	task, err := svc.CreateMaskTask(model.MaskTask{Name: "batch", SourceType: "stream", TotalRecords: 32})
	if err != nil { t.Fatalf("create task: %v", err) }
	if _, err = svc.TransitionMaskTaskStatus(task.ID, model.MaskTaskStatusRunning); err != nil { t.Fatalf("transition: %v", err) }
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if _, err := svc.AdvanceMaskTaskProgress(task.ID, 1); err != nil { t.Errorf("advance: %v", err) }
		}()
	}
	close(start)
	wg.Wait()
	got, err := svc.GetMaskTask(task.ID)
	if err != nil { t.Fatalf("get task: %v", err) }
	if got.ProcessedRecords != 32 || got.Status != model.MaskTaskStatusCompleted {
		t.Fatalf("progress was not accumulated: %+v", got)
	}
}
