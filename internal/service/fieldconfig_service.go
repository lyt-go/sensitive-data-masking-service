package service

import (
	"sort"
	"time"

	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/idgen"
)

func (s *Service) CreateFieldConfig(input model.FieldConfig) (*model.FieldConfig, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDataClass(input.DataClassID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("data_class_id", "数据分类不存在")
		}
		return nil, err
	}
	f := &model.FieldConfig{
		ID:          idgen.Hex(),
		FieldName:   input.FieldName,
		DataClassID: input.DataClassID,
		MaskType:    input.MaskType,
		Enabled:     input.Enabled,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateFieldConfig(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Service) ListFieldConfigs(filter model.FieldConfigFilter, page, size int) ([]*model.FieldConfig, int, error) {
	all := s.store.ListFieldConfigs()
	matched := make([]*model.FieldConfig, 0, len(all))
	for _, f := range all {
		if filter.Match(f) {
			matched = append(matched, f)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.FieldConfig{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetFieldConfig(id string) (*model.FieldConfig, error) {
	return s.store.GetFieldConfig(id)
}

func (s *Service) UpdateFieldConfig(id string, input model.FieldConfig) (*model.FieldConfig, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDataClass(input.DataClassID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("data_class_id", "数据分类不存在")
		}
		return nil, err
	}
	f, err := s.store.GetFieldConfig(id)
	if err != nil {
		return nil, err
	}
	f.FieldName = input.FieldName
	f.DataClassID = input.DataClassID
	f.MaskType = input.MaskType
	f.Enabled = input.Enabled
	if err := s.store.UpdateFieldConfig(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Service) DeleteFieldConfig(id string) error {
	return s.store.DeleteFieldConfig(id)
}
