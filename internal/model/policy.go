package model

import (
	"sort"
	"strings"
	"time"
)

type Policy struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Scope     string    `json:"scope"`
	RuleIDs   []string  `json:"rule_ids"`
	Priority  int       `json:"priority"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *Policy) Validate() error {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return NewValidationError("name", "策略名称不能为空")
	}
	p.Scope = strings.TrimSpace(p.Scope)
	if p.Scope == "" {
		return NewValidationError("scope", "作用范围不能为空")
	}
	return nil
}

func (p *Policy) SelectRules(allRules map[string]*MaskRule) []*MaskRule {
	if !p.Enabled || len(p.RuleIDs) == 0 {
		return nil
	}
	selected := make([]*MaskRule, 0, len(p.RuleIDs))
	for _, rid := range p.RuleIDs {
		if r, ok := allRules[rid]; ok && r.Enabled {
			selected = append(selected, r)
		}
	}
	sort.Slice(selected, func(i, j int) bool {
		return selected[i].CreatedAt.Before(selected[j].CreatedAt)
	})
	return selected
}

type PolicyFilter struct {
	Scope   string
	Enabled *bool
	Keyword string
}

func (f PolicyFilter) Match(p *Policy) bool {
	if f.Scope != "" && p.Scope != f.Scope {
		return false
	}
	if f.Enabled != nil && p.Enabled != *f.Enabled {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(p.Name), k) {
			return false
		}
	}
	return true
}
