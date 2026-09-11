package datetime

import (
	"strings"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type unit int

// Order also defines the inclusive legacy comparison range.
const (
	millisecond unit = iota
	second
	minute
	hour
	day
	week
	month
	year
)

func parseUnit(value runtime.Value, position int) (unit, error) {
	name, err := runtime.CastArg[runtime.String](value, position)
	if err != nil {
		return 0, err
	}

	switch strings.ToLower(string(name)) {
	case "f", "millisecond", "milliseconds":
		return millisecond, nil
	case "s", "second", "seconds":
		return second, nil
	case "i", "minute", "minutes":
		return minute, nil
	case "h", "hour", "hours":
		return hour, nil
	case "d", "day", "days":
		return day, nil
	case "w", "week", "weeks":
		return week, nil
	case "m", "month", "months":
		return month, nil
	case "y", "year", "years":
		return year, nil
	}

	return 0, runtime.ArgError(runtime.Errorf(runtime.ErrInvalidArgument, "unknown datetime unit %q", name), position)
}

func fixedDuration(u unit) (time.Duration, bool) {
	switch u {
	case millisecond:
		return time.Millisecond, true
	case second:
		return time.Second, true
	case minute:
		return time.Minute, true
	case hour:
		return time.Hour, true
	default:
		return 0, false
	}
}

func datePairUnit(arg1, arg2, arg3 runtime.Value) (runtime.DateTime, runtime.DateTime, unit, error) {
	left, right, err := runtime.CastArgs2[runtime.DateTime, runtime.DateTime](arg1, arg2)
	if err != nil {
		return runtime.ZeroDateTime, runtime.ZeroDateTime, 0, err
	}

	u, err := parseUnit(arg3, 2)

	return left, right, u, err
}
