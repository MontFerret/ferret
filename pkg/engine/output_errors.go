package engine

import (
	"context"
	"errors"
)

// Preserve a lone diagnostic's native formatting; only aggregates need joining.
func joinOutputErrors(errs ...error) error {
	var result error
	for _, err := range errs {
		if err == nil {
			continue
		}

		if result != nil {
			return errors.Join(errs...)
		}

		result = err
	}

	return result
}

func outputContextError(ctx context.Context) error {
	err := ctx.Err()
	if err == nil {
		return nil
	}

	cause := context.Cause(ctx)
	if cause != nil && !errors.Is(err, cause) {
		return errors.Join(err, cause)
	}

	return err
}
