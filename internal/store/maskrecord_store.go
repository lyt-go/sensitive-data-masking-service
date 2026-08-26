package store

import "datamasking/internal/model"

func (s *MemoryStore) CreateMaskRecord(m *model.MaskRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.maskTasks[m.MaskTaskID]; !ok {
		return ErrNotFound
	}
	s.maskRecords[m.ID] = m
	return nil
}

func (s *MemoryStore) GetMaskRecord(id string) (*model.MaskRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.maskRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

func (s *MemoryStore) ListMaskRecords() []*model.MaskRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.MaskRecord, 0, len(s.maskRecords))
	for _, m := range s.maskRecords {
		list = append(list, m)
	}
	return list
}

func (s *MemoryStore) UpdateMaskRecord(m *model.MaskRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.maskRecords[m.ID]; !ok {
		return ErrNotFound
	}
	if _, ok := s.maskTasks[m.MaskTaskID]; !ok {
		return ErrNotFound
	}
	s.maskRecords[m.ID] = m
	return nil
}

func (s *MemoryStore) DeleteMaskRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.maskRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.maskRecords, id)
	return nil
}

func (s *MemoryStore) BatchCreateMaskRecords(records []*model.MaskRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, m := range records {
		if _, ok := s.maskTasks[m.MaskTaskID]; !ok {
			return ErrNotFound
		}
	}
	for _, m := range records {
		s.maskRecords[m.ID] = m
	}
	return nil
}
