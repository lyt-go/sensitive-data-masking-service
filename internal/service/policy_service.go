package service

import (
	"sort"
	"time"

	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/idgen"
)

func (s *Service) CreatePolicy(input model.Policy) (*model.Policy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	for _, rid := range input.RuleIDs {
		if _, err := s.store.GetMaskRule(rid); err != nil {
			if err == store.ErrNotFound {
				return nil, model.NewValidationError("rule_ids", "规则不存在: "+rid)
			}
			return nil, err
		}
	}
	p := &model.Policy{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Scope:     input.Scope,
		RuleIDs:   input.RuleIDs,
		Priority:  input.Priority,
		Enabled:   input.Enabled,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreatePolicy(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) ListPolicies(filter model.PolicyFilter, page, size int) ([]*model.Policy, int, error) {
	all := s.store.ListPolicies()
	matched := make([]*model.Policy, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].Priority != matched[j].Priority {
			return matched[i].Priority > matched[j].Priority
		}
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Policy{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetPolicy(id string) (*model.Policy, error) {
	return s.store.GetPolicy(id)
}

func (s *Service) UpdatePolicy(id string, input model.Policy) (*model.Policy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	for _, rid := range input.RuleIDs {
		if _, err := s.store.GetMaskRule(rid); err != nil {
			if err == store.ErrNotFound {
				return nil, model.NewValidationError("rule_ids", "规则不存在: "+rid)
			}
			return nil, err
		}
	}
	p, err := s.store.GetPolicy(id)
	if err != nil {
		return nil, err
	}
	p.Name = input.Name
	p.Scope = input.Scope
	p.RuleIDs = input.RuleIDs
	p.Priority = input.Priority
	p.Enabled = input.Enabled
	if err := s.store.UpdatePolicy(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) DeletePolicy(id string) error {
	return s.store.DeletePolicy(id)
}

func (s *Service) EvaluatePolicy(id string) ([]*model.MaskRule, error) {
	p, err := s.store.GetPolicy(id)
	if err != nil {
		return nil, err
	}
	allRules := make(map[string]*model.MaskRule)
	for _, r := range s.store.ListMaskRules() {
		allRules[r.ID] = r
	}
	// SelectRules 返回的是指向已保存规则的指针，这里统一克隆，
	// 确保调用方修改评估结果不会污染已保存的规则（影响后续真正执行的脱敏）。
	selected := p.SelectRules(allRules)
	rules := make([]*model.MaskRule, 0, len(selected))
	for _, r := range selected {
		rules = append(rules, r.Clone())
	}
	return rules, nil
}
