package service

import (
	"sort"
	"time"

	"datamasking/internal/model"
	"datamasking/pkg/idgen"
)

func (s *Service) CreateMaskRule(input model.MaskRule) (*model.MaskRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	r := &model.MaskRule{
		ID:          idgen.Hex(),
		Name:        input.Name,
		MaskType:    input.MaskType,
		Pattern:     input.Pattern,
		PrefixKeep:  input.PrefixKeep,
		SuffixKeep:  input.SuffixKeep,
		Replacement: input.Replacement,
		Enabled:     input.Enabled,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateMaskRule(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) ListMaskRules(filter model.MaskRuleFilter, page, size int) ([]*model.MaskRule, int, error) {
	all := s.store.ListMaskRules()
	matched := make([]*model.MaskRule, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.MaskRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetMaskRule(id string) (*model.MaskRule, error) {
	return s.store.GetMaskRule(id)
}

func (s *Service) UpdateMaskRule(id string, input model.MaskRule) (*model.MaskRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	r, err := s.store.GetMaskRule(id)
	if err != nil {
		return nil, err
	}
	r.Name = input.Name
	r.MaskType = input.MaskType
	r.Pattern = input.Pattern
	r.PrefixKeep = input.PrefixKeep
	r.SuffixKeep = input.SuffixKeep
	r.Replacement = input.Replacement
	r.Enabled = input.Enabled
	if err := s.store.UpdateMaskRule(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteMaskRule(id string) error {
	return s.store.DeleteMaskRule(id)
}

func (s *Service) ApplyMaskRule(id string, input string) (string, error) {
	r, err := s.store.GetMaskRule(id)
	if err != nil {
		return "", err
	}
	return r.Apply(input), nil
}
