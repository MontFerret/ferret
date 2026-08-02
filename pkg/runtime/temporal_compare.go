package runtime

import (
	"context"
	"errors"
)

// EqualChecked applies language equality semantics, including temporal coercion.
// Duration conversion failures compare unequal, while operational errors propagate.
func EqualChecked(ctx context.Context, left, right Value) (Boolean, error) {
	result, err := CompareChecked(ctx, left, right)
	if err == nil {
		return result == 0, nil
	}

	if isDurationConversionError(err) {
		return False, nil
	}

	return False, err
}

// CompareChecked applies contextual temporal coercion for language comparison
// operators while preserving CompareValues for strict structural ordering.
func CompareChecked(ctx context.Context, left, right Value) (int, error) {
	leftIsDuration := false
	rightIsDuration := false

	switch left.(type) {
	case DateTime:
		return CompareValues(left, right), nil
	case Duration:
		leftIsDuration = true
	}

	switch right.(type) {
	case DateTime:
		return CompareValues(left, right), nil
	case Duration:
		rightIsDuration = true
	}

	if !leftIsDuration && !rightIsDuration {
		return CompareValues(left, right), nil
	}

	leftDuration, err := ToDuration(ctx, left)
	if err != nil {
		return 0, err
	}

	rightDuration, err := ToDuration(ctx, right)
	if err != nil {
		return 0, err
	}

	return leftDuration.Compare(rightDuration), nil
}

func isDurationConversionError(err error) bool {
	var conversionErr *durationConversionError
	return errors.As(err, &conversionErr)
}
