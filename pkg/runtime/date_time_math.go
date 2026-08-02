package runtime

import (
	"context"
	"errors"
	"math"
	"time"
)

func addDateTimeChecked(ctx context.Context, left, right Value) (Value, error) {
	leftDateTime, leftIsDateTime := left.(DateTime)
	rightDateTime, rightIsDateTime := right.(DateTime)

	if !leftIsDateTime && !rightIsDateTime {
		return nil, nil
	}

	if leftIsDateTime && rightIsDateTime {
		return None, temporalBinaryTypeError("+", left, right)
	}

	if leftIsDateTime {
		duration, err := ToDuration(ctx, right)
		if err != nil {
			return None, err
		}

		return checkedDateTimeAdd(leftDateTime, duration)
	}

	duration, err := ToDuration(ctx, left)
	if err != nil {
		return None, err
	}

	return checkedDateTimeAdd(rightDateTime, duration)
}

func subtractDateTimeChecked(ctx context.Context, left, right Value) (Value, error) {
	leftDateTime, leftIsDateTime := left.(DateTime)
	_, rightIsDateTime := right.(DateTime)

	if !leftIsDateTime && !rightIsDateTime {
		return nil, nil
	}

	if !leftIsDateTime {
		return None, temporalBinaryTypeError("-", left, right)
	}

	if rightDateTime, ok := right.(DateTime); ok {
		return checkedDateTimeDifference(leftDateTime, rightDateTime)
	}

	if rightString, ok := right.(String); ok {
		if dateTime, err := ToDateTime(ctx, rightString); err == nil {
			return checkedDateTimeDifference(leftDateTime, dateTime)
		}

		duration, err := ToDuration(ctx, rightString)
		if err == nil {
			return checkedDateTimeSubtract(leftDateTime, duration)
		}

		if errors.Is(err, ErrRange) {
			return None, err
		}

		return None, Errorf(
			ErrInvalidArgument,
			"cannot convert String %q to DateTime or Duration",
			rightString.String(),
		)
	}

	duration, err := ToDuration(ctx, right)
	if err != nil {
		return None, err
	}

	return checkedDateTimeSubtract(leftDateTime, duration)
}

func checkedDateTimeAdd(dateTime DateTime, duration Duration) (DateTime, error) {
	base := dateTime.Time.Round(0)
	candidate := base.Add(time.Duration(duration))

	if (duration > 0 && !candidate.After(base)) || (duration < 0 && !candidate.Before(base)) {
		return ZeroDateTime, dateTimeRangeError("addition")
	}

	var restored time.Time
	if duration == Duration(math.MinInt64) {
		restored = candidate.Add(time.Duration(math.MaxInt64)).Add(time.Nanosecond)
	} else {
		restored = candidate.Add(-time.Duration(duration))
	}

	if !restored.Equal(base) {
		return ZeroDateTime, dateTimeRangeError("addition")
	}

	return NewDateTime(candidate), nil
}

func checkedDateTimeSubtract(dateTime DateTime, duration Duration) (DateTime, error) {
	if duration != Duration(math.MinInt64) {
		return checkedDateTimeAdd(dateTime, -duration)
	}

	partial, err := checkedDateTimeAdd(dateTime, Duration(math.MaxInt64))
	if err != nil {
		return ZeroDateTime, dateTimeRangeError("subtraction")
	}

	result, err := checkedDateTimeAdd(partial, Duration(1))
	if err != nil {
		return ZeroDateTime, dateTimeRangeError("subtraction")
	}

	return result, nil
}

func checkedDateTimeDifference(left, right DateTime) (Duration, error) {
	leftInstant := left.Time.Round(0)
	rightInstant := right.Time.Round(0)
	difference := leftInstant.Sub(rightInstant)

	if !rightInstant.Add(difference).Equal(leftInstant) {
		return ZeroDuration, dateTimeRangeError("subtraction")
	}

	return Duration(difference), nil
}

func temporalBinaryTypeError(operator string, left, right Value) error {
	return Errorf(
		ErrInvalidOperation,
		"operator %s is not supported for %s and %s",
		operator,
		TypeName(TypeOf(left)),
		TypeName(TypeOf(right)),
	)
}

func dateTimeRangeError(operation string) error {
	return Errorf(ErrRange, "DateTime %s exceeds the supported range", operation)
}
