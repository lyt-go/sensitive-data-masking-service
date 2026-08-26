package store

import "datamasking/internal/model"

func (s *MemoryStore) CreateMaskTask(t *model.MaskTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, rid := range t.RuleIDs {
		if _, ok := s.maskRules[rid]; !ok {
			return ErrNotFound
		}
	}
	s.maskTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) GetMaskTask(id string) (*model.MaskTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.maskTasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	// 返回值的拷贝，避免调用方拿到 map 内部指针后跨锁并发写。
	cp := *t
	return &cp, nil
}

func (s *MemoryStore) ListMaskTasks() []*model.MaskTask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.MaskTask, 0, len(s.maskTasks))
	for _, t := range s.maskTasks {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateMaskTask(t *model.MaskTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.maskTasks[t.ID]; !ok {
		return ErrNotFound
	}
	for _, rid := range t.RuleIDs {
		if _, ok := s.maskRules[rid]; !ok {
			return ErrNotFound
		}
	}
	s.maskTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteMaskTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.maskTasks[id]; !ok {
		return ErrNotFound
	}
	for _, mr := range s.maskRecords {
		if mr.MaskTaskID == id {
			return ErrConflict
		}
	}
	delete(s.maskTasks, id)
	return nil
}

// AdvanceMaskTaskProgress 原子地推进任务进度。
// 读取当前进度、累加 delta、按业务规则判定是否完成的整个过程
// 都在 s.mu.Lock() 临界区内完成，避免并发推进时丢失更新。
// 返回推进后的任务拷贝。
func (s *MemoryStore) AdvanceMaskTaskProgress(id string, delta int) (*model.MaskTask, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.maskTasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	if err := t.AdvanceProgress(delta); err != nil {
		return nil, err
	}
	cp := *t
	return &cp, nil
}
