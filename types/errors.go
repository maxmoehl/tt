package types

import (
	"errors"
)

var (
	// ErrNotFound indicates that the error is related to (a) record(s)
	// not existing or not found.
	ErrNotFound = errors.New("not found")
)
