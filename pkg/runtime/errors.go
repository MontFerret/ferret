package runtime

import (
	"github.com/pkg/errors"
)

var (
	ErrMissedParams = errors.New("missed parameter values")
	ErrMissedParam  = errors.New("missed value for parameter")
)
