package util

import (
	"errors"
	"fmt"

	"github.com/ld/welfaremall/internal/constants"
)

// 哨兵错误。
var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource conflict")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrValidation        = errors.New("validation failed")
	ErrRateLimited       = errors.New("rate limited")
	ErrPointsNotEnough   = errors.New("points not enough")
	ErrProductOutOfStock = errors.New("product out of stock")
	ErrExchangeLimit     = errors.New("exchange limit reached")
	ErrActivityNotActive = errors.New("activity not active")
	ErrSoldOut           = errors.New("sold out")
)

// AppError 业务错误：携带业务码与 HTTP 状态码。
type AppError struct {
	Code       int
	HTTPStatus int
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func NewAppError(httpStatus, code int, message string, err error) *AppError {
	return &AppError{Code: code, HTTPStatus: httpStatus, Message: message, Err: err}
}

func BadRequest(message string, err error) *AppError {
	return NewAppError(400, constants.CodeBadRequest, message, err)
}

func Unauthorized(message string, err error) *AppError {
	return NewAppError(401, constants.CodeUnauthorized, message, err)
}

func Forbidden(message string, err error) *AppError {
	return NewAppError(403, constants.CodeForbidden, message, err)
}

func NotFound(message string, err error) *AppError {
	return NewAppError(404, constants.CodeNotFound, message, err)
}

func Conflict(message string, err error) *AppError {
	return NewAppError(409, constants.CodeConflict, message, err)
}

func Validation(message string, err error) *AppError {
	return NewAppError(422, constants.CodeValidationError, message, err)
}

func Internal(message string, err error) *AppError {
	return NewAppError(500, constants.CodeInternalError, message, err)
}

// AsAppError 将任意 error 归一为 AppError。
func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	switch {
	case errors.Is(err, ErrNotFound):
		return NotFound(constants.MsgNotFound, err)
	case errors.Is(err, ErrConflict):
		return Conflict(constants.MsgConflict, err)
	case errors.Is(err, ErrUnauthorized):
		return Unauthorized(constants.MsgUnauthorized, err)
	case errors.Is(err, ErrForbidden):
		return Forbidden(constants.MsgForbidden, err)
	case errors.Is(err, ErrValidation):
		return Validation(constants.MsgInvalidRequest, err)
	case errors.Is(err, ErrRateLimited):
		return NewAppError(429, constants.CodeRateLimited, constants.MsgRateLimited, err)
	case errors.Is(err, ErrPointsNotEnough):
		return BadRequest(constants.MsgPointsNotEnough, err)
	case errors.Is(err, ErrProductOutOfStock):
		return Conflict(constants.MsgProductOutOfStock, err)
	case errors.Is(err, ErrExchangeLimit):
		return Conflict(constants.MsgExchangeLimitReached, err)
	case errors.Is(err, ErrActivityNotActive):
		return Conflict(constants.MsgSeckillNotActive, err)
	case errors.Is(err, ErrSoldOut):
		return Conflict(constants.MsgSeckillSoldOut, err)
	}
	return Internal(constants.MsgInternalError, err)
}
