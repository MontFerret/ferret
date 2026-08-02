package datetime

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A0().
		Add("NOW", Now)

	ns.Function().A1().
		Add("DATE_DAYOFWEEK", DateDayOfWeek).
		Add("DATE_YEAR", DateYear).
		Add("DATE_MONTH", DateMonth).
		Add("DATE_DAY", DateDay).
		Add("DATE_HOUR", DateHour).
		Add("DATE_MINUTE", DateMinute).
		Add("DATE_SECOND", DateSecond).
		Add("DATE_MILLISECOND", DateMillisecond).
		Add("DATE_DAYOFYEAR", DateDayOfYear).
		Add("DATE_LEAPYEAR", DateLeapYear).
		Add("DATE_QUARTER", DateQuarter).
		Add("DATE_DAYS_IN_MONTH", DateDaysInMonth)

	ns.Function().A2().
		Add("DATE_FORMAT", DateFormat)

	ns.Function().A3().
		Add("DATE_ADD", DateAdd).
		Add("DATE_SUBTRACT", DateSubtract)

	ns.Function().Var().
		Add("DATE", Date).
		Add("DATE_COMPARE", DateCompare).
		Add("DATE_DIFF", DateDiff)
}
