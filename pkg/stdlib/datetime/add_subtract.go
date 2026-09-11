package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Add adds an integer number of units to a date. Subday units are elapsed time;
// calendar units use Go AddDate in the input's location, normalizing month ends
// rather than clamping them. A calendar day can span 23 or 25 elapsed hours.
// @param date {DateTime} Source date.
// @param amount {Int} Number of units to add; may be negative.
// @param unit {String} Millisecond, second, minute, hour, day, week, month, or year; plurals are accepted.
// @return {DateTime} Calculated date, retaining the input's location.
func Add(_ context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	date, amount, u, err := shiftArguments(arg1, arg2, arg3)
	if err != nil {
		return runtime.None, err
	}

	return runtime.NewDateTime(addUnit(date.Time, int(amount), u)), nil
}

// Subtract subtracts an integer number of units from a date. Subday units are
// elapsed time; calendar units use Go AddDate in the input's location, normalizing
// month ends rather than clamping them. A calendar day can span 23 or 25 hours.
// @param date {DateTime} Source date.
// @param amount {Int} Number of units to subtract; may be negative.
// @param unit {String} Millisecond, second, minute, hour, day, week, month, or year; plurals are accepted.
// @return {DateTime} Calculated date, retaining the input's location.
func Subtract(_ context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	date, amount, u, err := shiftArguments(arg1, arg2, arg3)
	if err != nil {
		return runtime.None, err
	}

	return runtime.NewDateTime(addUnit(date.Time, -int(amount), u)), nil
}

func shiftArguments(arg1, arg2, arg3 runtime.Value) (runtime.DateTime, runtime.Int, unit, error) {
	date, amount, err := runtime.CastArgs2[runtime.DateTime, runtime.Int](arg1, arg2)
	if err != nil {
		return runtime.ZeroDateTime, 0, 0, err
	}

	u, err := parseUnit(arg3, 2)

	return date, amount, u, err
}

func addUnit(date time.Time, amount int, u unit) time.Time {
	if duration, ok := fixedDuration(u); ok {
		return date.Add(time.Duration(amount) * duration)
	}

	switch u {
	case day:
		return date.AddDate(0, 0, amount)
	case week:
		return date.AddDate(0, 0, amount*7)
	case month:
		return date.AddDate(0, amount, 0)
	case year:
		return date.AddDate(amount, 0, 0)
	default:
		return date
	}
}
