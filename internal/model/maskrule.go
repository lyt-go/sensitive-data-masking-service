package model

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

const (
	MaskTypeMask     = "mask"
	MaskTypeTruncate = "truncate"
	MaskTypeHash     = "hash"
	MaskTypeReplace  = "replace"
)

var validMaskTypes = map[string]bool{
	MaskTypeMask:     true,
	MaskTypeTruncate: true,
	MaskTypeHash:     true,
	MaskTypeReplace:  true,
}

type MaskRule struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	MaskType    string    `json:"mask_type"`
	Pattern     string    `json:"pattern"`
	PrefixKeep  int       `json:"prefix_keep"`
	SuffixKeep  int       `json:"suffix_keep"`
	Replacement string    `json:"replacement"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
}

func (r *MaskRule) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	r.MaskType = strings.TrimSpace(r.MaskType)
	if r.MaskType == "" {
		return NewValidationError("mask_type", "脱敏类型不能为空")
	}
	if !validMaskTypes[r.MaskType] {
		return NewValidationError("mask_type", "脱敏类型不合法，可选值为 mask/truncate/hash/replace")
	}
	if r.PrefixKeep < 0 {
		return NewValidationError("prefix_keep", "前缀保留长度不能为负数")
	}
	if r.SuffixKeep < 0 {
		return NewValidationError("suffix_keep", "后缀保留长度不能为负数")
	}
	return nil
}

// Clone 返回 r 的值副本。MaskRule 的字段全部为值类型，浅拷贝即为完整副本，
// 用于对外返回（如策略评估）时隔离调用方修改，避免污染已保存的规则。
func (r *MaskRule) Clone() *MaskRule {
	if r == nil {
		return nil
	}
	cp := *r
	return &cp
}

func (r *MaskRule) Apply(input string) string {
	if input == "" {
		return input
	}
	switch r.MaskType {
	case MaskTypeMask:
		return applyMask(input, r.PrefixKeep, r.SuffixKeep)
	case MaskTypeTruncate:
		return applyTruncate(input, r.PrefixKeep)
	case MaskTypeHash:
		return applyHash(input)
	case MaskTypeReplace:
		return applyReplace(input, r.Replacement)
	default:
		return input
	}
}

func applyMask(input string, prefixKeep, suffixKeep int) string {
	runes := []rune(input)
	n := len(runes)
	if n == 0 {
		return input
	}
	keep := prefixKeep + suffixKeep
	if keep >= n {
		return input
	}
	if keep < 0 {
		keep = 0
	}
	var out []rune
	out = append(out, runes[:prefixKeep]...)
	for i := 0; i < n-keep; i++ {
		out = append(out, '*')
	}
	out = append(out, runes[n-suffixKeep:]...)
	return string(out)
}

func applyTruncate(input string, prefixKeep int) string {
	runes := []rune(input)
	if prefixKeep <= 0 {
		return ""
	}
	if prefixKeep >= len(runes) {
		return input
	}
	return string(runes[:prefixKeep])
}

func applyHash(input string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

func applyReplace(input, replacement string) string {
	if replacement == "" {
		return "[MASKED]"
	}
	return replacement
}

type MaskRuleFilter struct {
	MaskType string
	Enabled  *bool
	Keyword  string
}

func (f MaskRuleFilter) Match(r *MaskRule) bool {
	if f.MaskType != "" && r.MaskType != f.MaskType {
		return false
	}
	if f.Enabled != nil && r.Enabled != *f.Enabled {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Name), k) {
			return false
		}
	}
	return true
}
