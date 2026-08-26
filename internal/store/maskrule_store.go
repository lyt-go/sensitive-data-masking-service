package store

import "datamasking/internal/model"

func (s *MemoryStore) CreateMaskRule(r *model.MaskRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.maskRules {
		if exist.Name == r.Name {
			return ErrConflict
		}
	}
	s.maskRules[r.ID] = r
	return nil
}

func (s *MemoryStore) GetMaskRule(id string) (*model.MaskRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.maskRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListMaskRules() []*model.MaskRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.MaskRule, 0, len(s.maskRules))
	for _, r := range s.maskRules {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateMaskRule(r *model.MaskRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.maskRules[r.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.maskRules {
		if exist.ID != r.ID && exist.Name == r.Name {
			return ErrConflict
		}
	}
	s.maskRules[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteMaskRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.maskRules[id]; !ok {
		return ErrNotFound
	}
	for _, t := range s.maskTasks {
		for _, rid := range t.RuleIDs {
			if rid == id {
				return ErrConflict
			}
		}
	}
	for _, p := range s.policies {
		for _, rid := range p.RuleIDs {
			if rid == id {
				return ErrConflict
			}
		}
	}
	delete(s.maskRules, id)
	return nil
}
