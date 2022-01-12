package server

import (
	"net/http"

	"github.com/maxmoehl/tt"
)

var StatusCodeMapping = map[tt.Error]int{
	tt.ErrInvalidData:   http.StatusBadRequest,
	tt.ErrInternalError: http.StatusInternalServerError,
	tt.ErrNotFound:      http.StatusNotFound,
}
