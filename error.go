package tt

import (
	"fmt"
)

var (
	ErrInvalidData = fmt.Errorf("invalid data [00]")
	// ErrInternal indicates an error with the database
	ErrInternal          = fmt.Errorf("internal error [01]")
	ErrNotFound          = fmt.Errorf("not found [02]")
	ErrInvalidFormat     = fmt.Errorf("invalid format [03]")
	ErrNotImplemented    = fmt.Errorf("not implemented [04]")
	ErrInvalidTimer      = fmt.Errorf("invalid timer [05]")
	ErrInvalidParameter  = fmt.Errorf("invalid parameter supplied")
	ErrInvalidParameters = fmt.Errorf("invalid parameters supplied")
)
