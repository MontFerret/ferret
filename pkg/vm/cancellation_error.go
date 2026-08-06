package vm

import (
	"context"
	"errors"
)

func isExecutionCancellation(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
