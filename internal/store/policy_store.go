package store

import "datamasking/internal/model"

func (s *MemoryStore) CreatePolicy(p *model.Policy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, rid := range p.RuleIDs {
		if _, ok := s.maskRules[rid]; !ok {
			return ErrNotFound
		}
	}
	s.policies[p.ID] = p
	return nil
}

func (s *MemoryStore) GetPolicy(id string) (*model.Policy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.policies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *MemoryStore) ListPolicies() []*model.Policy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Policy, 0, len(s.policies))
	for _, p := range s.policies {
		list = append(list, p)
	}
	return list
}

func (s *MemoryStore) UpdatePolicy(p *model.Policy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.policies[p.ID]; !ok {
		return ErrNotFound
	}
	for _, rid := range p.RuleIDs {
		if _, ok := s.maskRules[rid]; !ok {
			return ErrNotFound
		}
	}
	s.policies[p.ID] = p
	return nil
}

func (s *MemoryStore) DeletePolicy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.policies[id]; !ok {
		return ErrNotFound
	}
	delete(s.policies, id)
	return nil
}
