package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Same compares each input's own local calendar through the requested precision.
// Offset and location identity are ignored. Week compares ISO week-year and week
// number; other units require all enclosing calendar fields to match.
// @param left {DateTime} First date.
// @param right {DateTime} Second date.
// @param unit {String} Precision: year, month, week, day, hour, minute, second, or millisecond; plurals are accepted.
// @return {Boolean} Whether the local calendar values match through that precision.
// @throws {TypeError} An argument has an invalid type.
// @throws {InvalidArgument} The unit is unknown.
func Same(_ context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	left, right, u, err := datePairUnit(arg1, arg2, arg3)
	if err != nil {
		return runtime.None, err
	}

	if u == week {
		return runtime.Boolean(sameComponent(left.Time, right.Time, week)), nil
	}

	for current := year; current >= u; current-- {
		if current != week && !sameComponent(left.Time, right.Time, current) {
			return runtime.False, nil
		}
	}

	return runtime.True, nil
}

func sameComponent(left, right time.Time, u unit) bool {
	switch u {
	case year:
		return left.Year() == right.Year()
	case month:
		return left.Month() == right.Month()
	case week:
		leftYear, leftWeek := left.ISOWeek()
		rightYear, rightWeek := right.ISOWeek()

		return leftYear == rightYear && leftWeek == rightWeek
	case day:
		return left.Day() == right.Day()
	case hour:
		return left.Hour() == right.Hour()
	case minute:
		return left.Minute() == right.Minute()
	case second:
		return left.Second() == right.Second()
	case millisecond:
		return left.Nanosecond()/int(time.Millisecond) == right.Nanosecond()/int(time.Millisecond)
	default:
		return false
	}
}
