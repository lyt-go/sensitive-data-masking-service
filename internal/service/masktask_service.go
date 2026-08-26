package service

import (
	"fmt"
	"sort"
	"time"

	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/idgen"
)

func (s *Service) CreateMaskTask(input model.MaskTask) (*model.MaskTask, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	for _, rid := range input.RuleIDs {
		if _, err := s.store.GetMaskRule(rid); err != nil {
			if err == store.ErrNotFound {
				return nil, fmt.Errorf("mask task rule lookup: %w", model.NewValidationError("rule_ids", "规则不存在: "+rid))
			}
			return nil, err
		}
	}
	t := &model.MaskTask{
		ID:               idgen.Hex(),
		Name:             input.Name,
		SourceType:       input.SourceType,
		TotalRecords:     input.TotalRecords,
		ProcessedRecords: 0,
		Status:           model.MaskTaskStatusPending,
		RuleIDs:          input.RuleIDs,
		CreatedAt:        time.Now(),
	}
	if err := s.store.CreateMaskTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) ListMaskTasks(filter model.MaskTaskFilter, page, size int) ([]*model.MaskTask, int, error) {
	all := s.store.ListMaskTasks()
	matched := make([]*model.MaskTask, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.MaskTask{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetMaskTask(id string) (*model.MaskTask, error) {
	return s.store.GetMaskTask(id)
}

func (s *Service) UpdateMaskTask(id string, input model.MaskTask) (*model.MaskTask, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	t, err := s.store.GetMaskTask(id)
	if err != nil {
		return nil, err
	}
	t.Name = input.Name
	t.SourceType = input.SourceType
	t.TotalRecords = input.TotalRecords
	t.RuleIDs = input.RuleIDs
	if err := s.store.UpdateMaskTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteMaskTask(id string) error {
	return s.store.DeleteMaskTask(id)
}

func (s *Service) TransitionMaskTaskStatus(id string, toStatus string) (*model.MaskTask, error) {
	t, err := s.store.GetMaskTask(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionMaskTask(t.Status, toStatus) {
		return nil, model.NewValidationError("status", "状态流转不合法: "+t.Status+" -> "+toStatus)
	}
	t.Status = toStatus
	if err := s.store.UpdateMaskTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) AdvanceMaskTaskProgress(id string, delta int) (*model.MaskTask, error) {
	t, err := s.store.GetMaskTask(id)
	if err != nil {
		return nil, err
	}
	if err := t.AdvanceProgress(delta); err != nil {
		return nil, err
	}
	if err := s.store.UpdateMaskTask(t); err != nil {
		return nil, err
	}
	return t, nil
}
