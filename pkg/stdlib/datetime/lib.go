package datetime

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A0().
		Add("now", Now)

	ns.Function().A1().
		Add("date", date1).
		Add("date_dayofweek", DateDayOfWeek).
		Add("date_year", DateYear).
		Add("date_month", DateMonth).
		Add("date_day", DateDay).
		Add("date_hour", DateHour).
		Add("date_minute", DateMinute).
		Add("date_second", DateSecond).
		Add("date_millisecond", DateMillisecond).
		Add("date_dayofyear", DateDayOfYear).
		Add("date_leapyear", DateLeapYear).
		Add("date_quarter", DateQuarter).
		Add("date_days_in_month", DateDaysInMonth)

	ns.Function().A2().
		Add("date", date2).
		Add("date_format", DateFormat)

	ns.Function().A3().
		Add("date_add", DateAdd).
		Add("date_compare", dateCompare3).
		Add("date_diff", dateDiff3).
		Add("date_subtract", DateSubtract)

	ns.Function().A4().
		Add("date_compare", dateCompare4).
		Add("date_diff", dateDiff4)
}
