package vo

import "github.com/gofiber/fiber/v2/utils"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type SetErrorOption func(*Error) *Error

func SetErrorMessage(message string) SetErrorOption {
	return func(r *Error) *Error {
		r.Message = message
		return r
	}
}
func SetErrorData(data any) SetErrorOption {
	return func(r *Error) *Error {
		r.Data = data
		return r
	}
}
func NewError(code int, sets ...SetErrorOption) *Error {
	e := &Error{
		Code:    code,
		Message: utils.StatusMessage(code),
	}
	for _, set := range sets {
		e = set(e)
	}
	return e
}

func (e *Error) SetData(data any) *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Data:    data,
	}
}

func (e *Error) SetCode(code int) *Error {
	return &Error{
		Code:    code,
		Message: e.Message,
		Data:    e.Data,
	}
}

func (e *Error) SetMessage(message string) *Error {
	return &Error{
		Code:    e.Code,
		Message: message,
		Data:    e.Data,
	}
}

// Error makes it compatible with the `error` interface.
func (e *Error) Error() string {
	return e.Message
}
