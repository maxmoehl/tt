package tt

import (
	"fmt"
)

var (
	ErrInvalidData   = NewError("invalid data")
	ErrInternalError = NewError("internal error")
	ErrNotFound      = NewError("not found")
)

type Error interface {
	Error() string
	Msg() string
	Unwrap() error
	WithCause(error) Error
}

type err struct {
	Message string `json:"message"`
	Cause   error  `json:"cause,omitempty"`
}

func (e *err) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s; reason: [%s]", e.Message, e.Cause.Error())
	}
	return e.Message
}

func (e *err) Msg() string {
	return e.Message
}

func (e *err) Unwrap() error {
	return e.Cause
}

// WithCause copies the underlying error and returns a new instance of it
// with the given cause set.
func (e *err) WithCause(cause error) Error {
	var newE = *e
	newE.Cause = cause
	return &newE
}

func NewError(msg string) Error {
	return &err{msg, nil}
}

func NewErrorf(msgFormat string, a ...interface{}) Error {
	return &err{fmt.Sprintf(msgFormat, a...), nil}
}
