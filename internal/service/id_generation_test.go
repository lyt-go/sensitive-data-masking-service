package service

import (
	"crypto/rand"
	"fmt"
	"io"
	"sync"
	"testing"

	"datamasking/internal/config"
	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/logger"
)

type coordinatedReader struct {
	mu sync.Mutex
	calls int
	firstReady chan struct{}
	secondFilled chan struct{}
}

func (r *coordinatedReader) Read(p []byte) (int, error) {
	r.mu.Lock()
	r.calls++
	call := r.calls
	r.mu.Unlock()
	if call == 1 {
		for i := range p { p[i] = 0x11 }
		close(r.firstReady)
		<-r.secondFilled
		return len(p), nil
	}
	<-r.firstReady
	for i := range p { p[i] = 0x22 }
	close(r.secondFilled)
	return len(p), nil
}

func TestConcurrentRuleCreationProducesDistinctIDs(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
	reader := &coordinatedReader{firstReady: make(chan struct{}), secondFilled: make(chan struct{})}
	originalReader := rand.Reader
	rand.Reader = io.Reader(reader)
	defer func() { rand.Reader = originalReader }()
	ids := make(chan string, 2)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done(); <-start
			r, err := svc.CreateMaskRule(model.MaskRule{Name: fmt.Sprintf("rule-%d", i), MaskType: "mask", Enabled: true})
			if err != nil { t.Errorf("create: %v", err); return }
			ids <- r.ID
		}(i)
	}
	close(start); wg.Wait(); close(ids)
	seen := map[string]bool{}
	for id := range ids {
		if seen[id] { t.Fatalf("duplicate rule id %q", id) }
		seen[id] = true
	}
	if len(seen) != 2 { t.Fatalf("created %d unique rules", len(seen)) }
}
