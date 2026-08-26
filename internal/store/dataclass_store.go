package store

import "datamasking/internal/model"

func (s *MemoryStore) CreateDataClass(d *model.DataClass) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.dataClasses {
		if exist.Name == d.Name {
			return ErrConflict
		}
	}
	s.dataClasses[d.ID] = d
	return nil
}

func (s *MemoryStore) GetDataClass(id string) (*model.DataClass, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.dataClasses[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) ListDataClasses() []*model.DataClass {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.DataClass, 0, len(s.dataClasses))
	for _, d := range s.dataClasses {
		list = append(list, d)
	}
	return list
}

func (s *MemoryStore) UpdateDataClass(d *model.DataClass) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dataClasses[d.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.dataClasses {
		if exist.ID != d.ID && exist.Name == d.Name {
			return ErrConflict
		}
	}
	s.dataClasses[d.ID] = d
	return nil
}

func (s *MemoryStore) DeleteDataClass(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dataClasses[id]; !ok {
		return ErrNotFound
	}
	for _, fc := range s.fieldConfigs {
		if fc.DataClassID == id {
			return ErrConflict
		}
	}
	delete(s.dataClasses, id)
	return nil
}
