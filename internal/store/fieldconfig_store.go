package store

import "datamasking/internal/model"

func (s *MemoryStore) CreateFieldConfig(f *model.FieldConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dataClasses[f.DataClassID]; !ok {
		return ErrNotFound
	}
	s.fieldConfigs[f.ID] = f
	return nil
}

func (s *MemoryStore) GetFieldConfig(id string) (*model.FieldConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.fieldConfigs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

func (s *MemoryStore) ListFieldConfigs() []*model.FieldConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.FieldConfig, 0, len(s.fieldConfigs))
	for _, f := range s.fieldConfigs {
		list = append(list, f)
	}
	return list
}

func (s *MemoryStore) UpdateFieldConfig(f *model.FieldConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.fieldConfigs[f.ID]; !ok {
		return ErrNotFound
	}
	if _, ok := s.dataClasses[f.DataClassID]; !ok {
		return ErrNotFound
	}
	s.fieldConfigs[f.ID] = f
	return nil
}

func (s *MemoryStore) DeleteFieldConfig(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.fieldConfigs[id]; !ok {
		return ErrNotFound
	}
	delete(s.fieldConfigs, id)
	return nil
}
