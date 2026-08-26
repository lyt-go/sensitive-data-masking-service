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
	return t, nil
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
