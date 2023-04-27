package runtime

import (
	"github.com/pkg/errors"
)

var (
	ErrMissedParam      = errors.New("missed value for parameter(s)")
	ErrValueUndefined   = errors.New("value is undefined")
	ErrFunctionNotFound = errors.New("function not found")
)
