package datetime

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Add adds an integer number of units to a date. Subday units are elapsed time;
// calendar units use Go AddDate in the input's location, normalizing month ends
// rather than clamping them. A calendar day can span 23 or 25 elapsed hours.
// @param date {DateTime} Source date.
// @param amount {Int} Number of units to add; may be negative.
// @param unit {String} Millisecond, second, minute, hour, day, week, month, or year; plurals are accepted.
// @return {DateTime} Calculated date, retaining the input's location.
// @throws {RangeError} The shift or resulting date cannot be represented safely.
func Add(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	date, amount, u, err := shiftArguments(arg1, arg2, arg3)
	if err != nil {
		return runtime.None, err
	}

	result, err := shiftDate(ctx, date, amount, u, false)
	if err != nil {
		return runtime.None, runtime.ArgError(err, 1)
	}

	return result, nil
}

// Subtract subtracts an integer number of units from a date. Subday units are
// elapsed time; calendar units use Go AddDate in the input's location, normalizing
// month ends rather than clamping them. A calendar day can span 23 or 25 hours.
// @param date {DateTime} Source date.
// @param amount {Int} Number of units to subtract; may be negative.
// @param unit {String} Millisecond, second, minute, hour, day, week, month, or year; plurals are accepted.
// @return {DateTime} Calculated date, retaining the input's location.
// @throws {RangeError} The shift or resulting date cannot be represented safely.
func Subtract(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	date, amount, u, err := shiftArguments(arg1, arg2, arg3)
	if err != nil {
		return runtime.None, err
	}

	result, err := shiftDate(ctx, date, amount, u, true)
	if err != nil {
		return runtime.None, runtime.ArgError(err, 1)
	}

	return result, nil
}

func shiftArguments(arg1, arg2, arg3 runtime.Value) (runtime.DateTime, runtime.Int, unit, error) {
	date, amount, err := runtime.CastArgs2[runtime.DateTime, runtime.Int](arg1, arg2)
	if err != nil {
		return runtime.ZeroDateTime, 0, 0, err
	}

	u, err := parseUnit(arg3, 2)

	return date, amount, u, err
}

func shiftDate(ctx context.Context, date runtime.DateTime, amount runtime.Int, u unit, subtract bool) (runtime.Value, error) {
	if duration, ok := fixedDuration(u); ok {
		delta, err := runtime.Multiply(ctx, runtime.Duration(duration), amount)
		if err != nil {
			return runtime.None, err
		}

		if subtract {
			return runtime.Subtract(ctx, date, delta)
		}

		return runtime.Add(ctx, date, delta)
	}

	if subtract {
		if amount == math.MinInt64 {
			return runtime.None, calendarRangeError()
		}

		amount = -amount
	}

	result, err := addCalendar(ctx, date.Time, amount, u)
	if err != nil {
		return runtime.None, err
	}

	return runtime.NewDateTime(result), nil
}
