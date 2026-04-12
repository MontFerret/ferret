package runtime

import "errors"

// InvalidArgumentDetails returns the outermost invalid argument position and cause.
// The returned position is zero-based so callers can use it both for metadata lookups
// and for user-facing messages after converting to a one-based index.
func InvalidArgumentDetails(err error) (pos int, ok bool, cause error) {
	var invalidErr *invalidArgumentError

	if !errors.As(err, &invalidErr) || invalidErr == nil {
		return 0, false, nil
	}

	return invalidErr.position, true, invalidErr.cause
}
