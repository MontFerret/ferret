package runtime

import (
	"github.com/pkg/errors"
)

var (
	ErrMissedParam = errors.New("missed value for parameter(s)")
)
