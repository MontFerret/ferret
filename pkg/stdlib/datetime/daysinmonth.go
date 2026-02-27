package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var daysCount = map[time.Month]int{
	time.January:   31,
	time.February:  28,
	time.March:     31,
	time.April:     30,
	time.May:       31,
	time.June:      30,
	time.July:      30,
	time.August:    31,
	time.September: 30,
	time.October:   31,
	time.November:  30,
	time.December:  31,
}

// DATE_DAYS_IN_MONTH returns the number of days in the month of date.
// @param {DateTime} date - Source DateTime.
// @return {Int} - Number of the days.
func DateDaysInMonth(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	month := dt.Month()
	count := daysCount[month]

	if month == time.February && isLeap(dt.Year()) {
		count++
	}

	return runtime.Int(count), nil
}
