package datetime

import (
	"github.com/MontFerret/ferret/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) error {
	return ns.RegisterFunctions(
		runtime.NewFunctionsFromMap(map[string]runtime.Function{
			"NOW":                Now,
			"DATE":               Date,
			"DATE_COMPARE":       DateCompare,
			"DATE_DAYOFWEEK":     DateDayOfWeek,
			"DATE_YEAR":          DateYear,
			"DATE_MONTH":         DateMonth,
			"DATE_DAY":           DateDay,
			"DATE_HOUR":          DateHour,
			"DATE_MINUTE":        DateMinute,
			"DATE_SECOND":        DateSecond,
			"DATE_MILLISECOND":   DateMillisecond,
			"DATE_DAYOFYEAR":     DateDayOfYear,
			"DATE_LEAPYEAR":      DateLeapYear,
			"DATE_QUARTER":       DateQuarter,
			"DATE_DAYS_IN_MONTH": DateDaysInMonth,
			"DATE_FORMAT":        DateFormat,
			"DATE_ADD":           DateAdd,
			"DATE_SUBTRACT":      DateSubtract,
			"DATE_DIFF":          DateDiff,
		}),
	)
}
