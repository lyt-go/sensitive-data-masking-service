package service

import (
	"sort"
	"time"

	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/idgen"
)

func (s *Service) CreateMaskRecord(input model.MaskRecord) (*model.MaskRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetMaskTask(input.MaskTaskID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("mask_task_id", "脱敏任务不存在")
		}
		return nil, err
	}
	m := &model.MaskRecord{
		ID:             idgen.Hex(),
		MaskTaskID:     input.MaskTaskID,
		FieldName:      input.FieldName,
		RuleID:         input.RuleID,
		OriginalSample: input.OriginalSample,
		MaskedSample:   input.MaskedSample,
		CreatedAt:      time.Now(),
	}
	if err := s.store.CreateMaskRecord(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) ListMaskRecords(filter model.MaskRecordFilter, page, size int) ([]*model.MaskRecord, int, error) {
	all := s.store.ListMaskRecords()
	matched := make([]*model.MaskRecord, 0, len(all))
	for _, m := range all {
		if filter.Match(m) {
			matched = append(matched, m)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.MaskRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetMaskRecord(id string) (*model.MaskRecord, error) {
	return s.store.GetMaskRecord(id)
}

func (s *Service) UpdateMaskRecord(id string, input model.MaskRecord) (*model.MaskRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetMaskTask(input.MaskTaskID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("mask_task_id", "脱敏任务不存在")
		}
		return nil, err
	}
	m, err := s.store.GetMaskRecord(id)
	if err != nil {
		return nil, err
	}
	m.MaskTaskID = input.MaskTaskID
	m.FieldName = input.FieldName
	m.RuleID = input.RuleID
	m.OriginalSample = input.OriginalSample
	m.MaskedSample = input.MaskedSample
	if err := s.store.UpdateMaskRecord(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) DeleteMaskRecord(id string) error {
	return s.store.DeleteMaskRecord(id)
}

func (s *Service) BatchCreateMaskRecords(inputs []model.MaskRecord) (int, error) {
	records := make([]*model.MaskRecord, 0, len(inputs))
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			return 0, err
		}
		records = append(records, &model.MaskRecord{
			ID:             idgen.Hex(),
			MaskTaskID:     input.MaskTaskID,
			FieldName:      input.FieldName,
			RuleID:         input.RuleID,
			OriginalSample: input.OriginalSample,
			MaskedSample:   input.MaskedSample,
			CreatedAt:      time.Now(),
		})
	}
	if err := s.store.BatchCreateMaskRecords(records); err != nil {
		return 0, err
	}
	return len(records), nil
}
