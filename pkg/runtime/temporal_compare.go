package runtime

import "context"

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
