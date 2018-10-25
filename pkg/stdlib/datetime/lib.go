package datetime

import "github.com/MontFerret/ferret/pkg/runtime/core"

func NewLib() map[string]core.Function {
	return map[string]core.Function{
<<<<<<< HEAD
<<<<<<< HEAD
		"NOW":                Now,
		"DATE":               Date,
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
=======
		"NOW": Now,
>>>>>>> 09546ae... added lib.go
=======
		"NOW":            Now,
		"DATE":           Date,
		"DATE_DAYOFWEEK": DateDayOfWeek,
		"DATE_YEAR":      DateYear,
>>>>>>> db2ab4b... added DATE_YEAR function
	}
}
