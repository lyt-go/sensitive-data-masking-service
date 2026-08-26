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

// IsValidationError 判断 err 是否为（或包装了）ValidationError。
// 使用 errors.As 以穿透 fmt.Errorf("%w", ...) 之类的包装，
// 保证服务层用 %w 包装后的校验错误仍能被识别为参数错误。
func IsValidationError(err error) bool {
	var target *ValidationError
	return errors.As(err, &target)
}

// AsValidationError 提取 err 中包装的 ValidationError，便于读取字段信息。
// 若 err 未包装 ValidationError 则返回 nil。
func AsValidationError(err error) *ValidationError {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve
	}
	return nil
}
