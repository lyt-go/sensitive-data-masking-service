package model

import (
	"strings"
	"time"
)

const (
	MaskTaskStatusPending   = "pending"
	MaskTaskStatusRunning   = "running"
	MaskTaskStatusCompleted = "completed"
	MaskTaskStatusFailed    = "failed"
)

var maskTaskTransitions = map[string]map[string]bool{
	MaskTaskStatusPending:   {MaskTaskStatusRunning: true, MaskTaskStatusFailed: true},
	MaskTaskStatusRunning:   {MaskTaskStatusCompleted: true, MaskTaskStatusFailed: true},
	MaskTaskStatusCompleted: {},
	MaskTaskStatusFailed:    {},
}

func CanTransitionMaskTask(from, to string) bool {
	if m, ok := maskTaskTransitions[from]; ok {
		return m[to]
	}
	return false
}

type MaskTask struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	SourceType       string    `json:"source_type"`
	TotalRecords     int       `json:"total_records"`
	ProcessedRecords int       `json:"processed_records"`
	Status           string    `json:"status"`
	RuleIDs          []string  `json:"rule_ids"`
	CreatedAt        time.Time `json:"created_at"`
}

func (t *MaskTask) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" {
		return NewValidationError("name", "任务名称不能为空")
	}
	t.SourceType = strings.TrimSpace(t.SourceType)
	if t.SourceType == "" {
		return NewValidationError("source_type", "来源类型不能为空")
	}
	if t.TotalRecords < 0 {
		return NewValidationError("total_records", "总记录数不能为负数")
	}
	if t.ProcessedRecords < 0 {
		return NewValidationError("processed_records", "已处理记录数不能为负数")
	}
	if t.Status == "" {
		t.Status = MaskTaskStatusPending
	}
	if t.Status != MaskTaskStatusPending && t.Status != MaskTaskStatusRunning &&
		t.Status != MaskTaskStatusCompleted && t.Status != MaskTaskStatusFailed {
		return NewValidationError("status", "任务状态不合法")
	}
	return nil
}

func (t *MaskTask) AdvanceProgress(delta int) error {
	if delta < 0 {
		return NewValidationError("delta", "进度增量不能为负数")
	}
	if t.Status != MaskTaskStatusRunning {
		return NewValidationError("status", "只有运行中的任务才能推进进度")
	}
	t.ProcessedRecords += delta
	if t.ProcessedRecords >= t.TotalRecords {
		t.ProcessedRecords = t.TotalRecords
		t.Status = MaskTaskStatusCompleted
	}
	return nil
}

type MaskTaskFilter struct {
	Status     string
	SourceType string
	Keyword    string
}

func (f MaskTaskFilter) Match(t *MaskTask) bool {
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.SourceType != "" && t.SourceType != f.SourceType {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Name), k) {
			return false
		}
	}
	return true
}
