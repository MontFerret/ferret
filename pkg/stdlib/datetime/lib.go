package datetime

import "github.com/MontFerret/ferret/pkg/runtime/core"

func NewLib() map[string]core.Function {
	return map[string]core.Function{
		"NOW":              Now,
		"DATE":             Date,
		"DATE_DAYOFWEEK":   DateDayOfWeek,
		"DATE_YEAR":        DateYear,
		"DATE_MONTH":       DateMonth,
		"DATE_DAY":         DateDay,
		"DATE_HOUR":        DateHour,
		"DATE_MINUTE":      DateMinute,
		"DATE_SECOND":      DateSecond,
		"DATE_MILLISECOND": DateMillisecond,
		"DATE_DAYOFYEAR":   DateDayOfYear,
		"DATE_LEAPYEAR":    DateLeapYear,
		"DATE_QUARTER":     DateQuarter,
	}
}
