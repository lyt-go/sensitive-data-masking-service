package model

import "errors"

// ValidationError 表示字段校验失败。
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return e.Field + ": " + e.Message
	}
	return e.Message
}

func NewValidationError(field, message string) error {
	return &ValidationError{Field: field, Message: message}
}

func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// ErrRuleDisabled 表示规则已停用，不能继续处理数据。
// 作为请求级校验错误，对外映射为 400 BadRequest。
var ErrRuleDisabled = errors.New("规则已停用，无法应用脱敏")
