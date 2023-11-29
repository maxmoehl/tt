package tt

import (
	"fmt"
)

var (
	ErrInvalidData = fmt.Errorf("invalid data")
	// ErrInternal indicates an error with the database
	ErrInternal              = fmt.Errorf("internal error")
	ErrNotFound              = fmt.Errorf("not found")
	ErrInvalidFormat         = fmt.Errorf("invalid format")
	ErrNotImplemented        = fmt.Errorf("not implemented")
	ErrInvalidTimer          = fmt.Errorf("invalid timer")
	ErrInvalidParameter      = fmt.Errorf("invalid parameter supplied")
	ErrInvalidParameters     = fmt.Errorf("invalid parameters supplied")
	ErrOperationNotPermitted = fmt.Errorf("operation not permitted")
)
