package runtime

import "github.com/pkg/errors"

var (
	ErrMissedParam = errors.New("missed parameter")
)

type (
	SourceErrorDetail struct {
		error
		BaseError    error
		ComputeError error
	}
)

func (e *SourceErrorDetail) Error() string {
	return e.ComputeError.Error()
}

func SourceError(src SourceMap, err error) error {
	return &SourceErrorDetail{
		BaseError:    err,
		ComputeError: errors.Errorf("%s: %s", err.Error(), src.String()),
	}
}
