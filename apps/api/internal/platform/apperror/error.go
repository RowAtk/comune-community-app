package apperror

import "errors"

type Kind string

const (
	KindNotFound     Kind = "not_found"
	KindConflict     Kind = "conflict"
	KindUnauthorized Kind = "unauthorized"
	KindValidation   Kind = "validation"
	KindInternal     Kind = "internal"
)

type Error struct {
	Kind    Kind
	Code    string
	Message string
	Err     error
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) Wrap(err error) *Error {
	e.Err = err
	return &e
}

func Validation(code, msg string, err error) *Error {
	return &Error{Kind: KindValidation, Code: code, Message: msg, Err: err}
}

func Conflict(code, msg string, err error) *Error {
	return &Error{Kind: KindConflict, Code: code, Message: msg, Err: err}
}

func NotFound(code, msg string, err error) *Error {
	return &Error{Kind: KindNotFound, Code: code, Message: msg, Err: err}
}

func Unauthorized(code, msg string, err error) *Error {
	return &Error{Kind: KindUnauthorized, Code: code, Message: msg, Err: err}
}

func Internal(code, msg string, err error) *Error {
	return &Error{Kind: KindInternal, Code: code, Message: msg, Err: err}
}

func As(err error) (*Error, bool) {
	var appErr *Error
	if !errors.As(err, &appErr) {
		return nil, false
	}

	return appErr, true
}
