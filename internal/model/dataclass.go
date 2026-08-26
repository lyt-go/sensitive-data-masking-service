package model

import (
	"strings"
	"time"
)

const (
	DataClassLevelPublic      = "public"
	DataClassLevelInternal    = "internal"
	DataClassLevelConfidential = "confidential"
	DataClassLevelSecret      = "secret"
)

var validDataClassLevels = map[string]bool{
	DataClassLevelPublic:       true,
	DataClassLevelInternal:     true,
	DataClassLevelConfidential: true,
	DataClassLevelSecret:       true,
}

type DataClass struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Level       string    `json:"level"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (d *DataClass) Validate() error {
	d.Name = strings.TrimSpace(d.Name)
	if d.Name == "" {
		return NewValidationError("name", "数据分类名称不能为空")
	}
	d.Level = strings.TrimSpace(d.Level)
	if d.Level == "" {
		return NewValidationError("level", "安全级别不能为空")
	}
	if !validDataClassLevels[d.Level] {
		return NewValidationError("level", "安全级别不合法，可选值为 public/internal/confidential/secret")
	}
	return nil
}

type DataClassFilter struct {
	Level   string
	Keyword string
}

func (f DataClassFilter) Match(d *DataClass) bool {
	if f.Level != "" && d.Level != f.Level {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(d.Name), k) {
			return false
		}
	}
	return true
}
