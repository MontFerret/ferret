package runtime

import (
	"math"
	"time"
)

func addDateTimeDuration(dateTime DateTime, duration Duration) (Value, error) {
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

func subtractDateTimeDuration(dateTime DateTime, duration Duration) (Value, error) {
	if duration != Duration(math.MinInt64) {
		return addDateTimeDuration(dateTime, -duration)
	}

	partialValue, err := addDateTimeDuration(dateTime, Duration(math.MaxInt64))
	if err != nil {
		return ZeroDateTime, dateTimeRangeError("subtraction")
	}
	partial := partialValue.(DateTime)

	result, err := addDateTimeDuration(partial, Duration(1))
	if err != nil {
		return ZeroDateTime, dateTimeRangeError("subtraction")
	}

	return result, nil
}

func subtractDateTimes(left, right DateTime) (Value, error) {
	leftInstant := left.Time.Round(0)
	rightInstant := right.Time.Round(0)
	difference := leftInstant.Sub(rightInstant)

	if !rightInstant.Add(difference).Equal(leftInstant) {
		return ZeroDuration, dateTimeRangeError("subtraction")
	}

	return Duration(difference), nil
}

func dateTimeRangeError(operation string) error {
	return Errorf(ErrRange, "DateTime %s exceeds the supported range", operation)
}
