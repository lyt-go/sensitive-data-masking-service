package service

import (
	"sort"
	"time"

	"datamasking/internal/model"
	"datamasking/pkg/idgen"
)

func (s *Service) CreateDataClass(input model.DataClass) (*model.DataClass, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	d := &model.DataClass{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Level:       input.Level,
		Description: input.Description,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateDataClass(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) ListDataClasses(filter model.DataClassFilter, page, size int) ([]*model.DataClass, int, error) {
	all := s.store.ListDataClasses()
	matched := make([]*model.DataClass, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.DataClass{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetDataClass(id string) (*model.DataClass, error) {
	return s.store.GetDataClass(id)
}

func (s *Service) UpdateDataClass(id string, input model.DataClass) (*model.DataClass, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	d, err := s.store.GetDataClass(id)
	if err != nil {
		return nil, err
	}
	// 构造新对象再交给 store 校验落库，避免直接修改内存中的原记录：
	// 一旦新名称与其它分类冲突，store 会返回 ErrConflict，此时原分类内容保持不变。
	updated := &model.DataClass{
		ID:          d.ID,
		Name:        input.Name,
		Level:       input.Level,
		Description: input.Description,
		CreatedAt:   d.CreatedAt,
	}
	if err := s.store.UpdateDataClass(updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) DeleteDataClass(id string) error {
	return s.store.DeleteDataClass(id)
}
