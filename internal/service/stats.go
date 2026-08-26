package service

import (
	"sort"

	"datamasking/internal/model"
)

type OverviewStats struct {
	DataClassCount    int                    `json:"data_class_count"`
	MaskRuleCount     int                    `json:"mask_rule_count"`
	MaskTaskByStatus  map[string]int         `json:"mask_task_by_status"`
	RecordByMaskType  map[string]int         `json:"record_by_mask_type"`
	TopFieldRecordCnt []FieldRecordCountItem `json:"top_field_record_cnt"`
}

type FieldRecordCountItem struct {
	FieldName string `json:"field_name"`
	Count     int    `json:"count"`
}

func (s *Service) GetOverviewStats() (*OverviewStats, error) {
	stats := &OverviewStats{
		MaskTaskByStatus: make(map[string]int),
		RecordByMaskType: make(map[string]int),
	}

	stats.DataClassCount = len(s.store.ListDataClasses())
	stats.MaskRuleCount = len(s.store.ListMaskRules())

	for _, t := range s.store.ListMaskTasks() {
		stats.MaskTaskByStatus[t.Status]++
	}

	allRules := make(map[string]*model.MaskRule)
	for _, r := range s.store.ListMaskRules() {
		allRules[r.ID] = r
	}

	fieldCnt := make(map[string]int)
	for _, rec := range s.store.ListMaskRecords() {
		fieldCnt[rec.FieldName]++
		if r, ok := allRules[rec.RuleID]; ok {
			stats.RecordByMaskType[r.MaskType]++
		} else {
			stats.RecordByMaskType["unknown"]++
		}
	}

	items := make([]FieldRecordCountItem, 0, len(fieldCnt))
	for fn, c := range fieldCnt {
		items = append(items, FieldRecordCountItem{FieldName: fn, Count: c})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count != items[j].Count {
			return items[i].Count > items[j].Count
		}
		return items[i].FieldName < items[j].FieldName
	})
	if len(items) > 10 {
		items = items[:10]
	}
	stats.TopFieldRecordCnt = items

	return stats, nil
}
