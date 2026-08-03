package runtime

import (
	"context"
	"errors"
)

// EqualChecked applies language equality semantics, including contextual Duration
// coercion. DateTime equality remains strict. Duration conversion failures compare
// unequal, while operational errors propagate.
func EqualChecked(ctx context.Context, left, right Value) (Boolean, error) {
	leftKind := nativeComparisonKindOf(left)
	rightKind := nativeComparisonKindOf(right)

	if leftKind == builtinComparisonUnknown || rightKind == builtinComparisonUnknown {
		if result, handled, err := dispatchCompatibleEquality(ctx, left, right); handled || err != nil {
			return result, err
		}
	}

	leftIsDuration := false
	rightIsDuration := false

	switch left.(type) {
	case DateTime:
		return equalStrictValues(ctx, left, right, leftKind, rightKind)
	case Duration:
		leftIsDuration = true
	}

	switch right.(type) {
	case DateTime:
		return equalStrictValues(ctx, left, right, leftKind, rightKind)
	case Duration:
		rightIsDuration = true
	}

	if !leftIsDuration && !rightIsDuration {
		return equalStrictValues(ctx, left, right, leftKind, rightKind)
	}

	if !nativeDurationComparisonValue(left) || !nativeDurationComparisonValue(right) {
		return equalStrictValues(ctx, left, right, leftKind, rightKind)
	}

	leftDuration, err := ToDuration(ctx, left)
	if err != nil {
		if isConversionErrorTo(err, TypeDuration) {
			return False, nil
		}

		return False, err
	}

	rightDuration, err := ToDuration(ctx, right)
	if err != nil {
		if isConversionErrorTo(err, TypeDuration) {
			return False, nil
		}

		return False, err
	}

	return leftDuration == rightDuration, nil
}

// CompareChecked applies contextual Duration coercion for language comparison
// operators. DateTime comparison and CompareValues structural ordering remain strict.
func CompareChecked(ctx context.Context, left, right Value) (Ordering, error) {
	leftKind := nativeComparisonKindOf(left)
	rightKind := nativeComparisonKindOf(right)

	if leftKind == builtinComparisonUnknown || rightKind == builtinComparisonUnknown {
		if result, handled, err := dispatchCompatibleOrdering(ctx, left, right); handled || err != nil {
			return result, err
		}
	}

	leftIsDuration := false
	rightIsDuration := false

	switch left.(type) {
	case DateTime:
		return compareStrictValues(ctx, left, right, leftKind, rightKind)
	case Duration:
		leftIsDuration = true
	}

	switch right.(type) {
	case DateTime:
		return compareStrictValues(ctx, left, right, leftKind, rightKind)
	case Duration:
		rightIsDuration = true
	}

	if !leftIsDuration && !rightIsDuration {
		return compareStrictValues(ctx, left, right, leftKind, rightKind)
	}

	if !nativeDurationComparisonValue(left) || !nativeDurationComparisonValue(right) {
		return Equal, incompatibleComparisonError(left, right)
	}

	leftDuration, err := ToDuration(ctx, left)
	if err != nil {
		return Equal, err
	}

	rightDuration, err := ToDuration(ctx, right)
	if err != nil {
		return Equal, err
	}

	return compareOrdered(leftDuration, rightDuration), nil
}

func equalStrictValues(
	ctx context.Context,
	left, right Value,
	leftKind, rightKind builtinComparison,
) (Boolean, error) {
	if leftKind != builtinComparisonUnknown && rightKind != builtinComparisonUnknown {
		return equalNativeValues(ctx, left, right, leftKind, rightKind)
	}

	return EqualValues(ctx, left, right)
}

func compareStrictValues(
	ctx context.Context,
	left, right Value,
	leftKind, rightKind builtinComparison,
) (Ordering, error) {
	if leftKind != builtinComparisonUnknown && rightKind != builtinComparisonUnknown {
		return compareNativeValues(ctx, left, right, leftKind, rightKind)
	}

	return CompareValues(ctx, left, right)
}

func nativeDurationComparisonValue(value Value) bool {
	if value == nil || value == None {
		return true
	}

	switch value := value.(type) {
	case Duration, Int, Float, String, Boolean:
		return true
	case *Array:
		for _, item := range value.data {
			if !nativeDurationComparisonValue(item) {
				return false
			}
		}

		return true
	default:
		return false
	}
}

func isConversionErrorTo(err error, target Type) bool {
	var conversionErr *conversionError
	return errors.As(err, &conversionErr) && conversionErr.targets(target)
}
