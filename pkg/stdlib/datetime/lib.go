package datetime

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// RegisterLib registers canonical datetime functions and deprecated global adapters.
// @namespace datetime
func RegisterLib(ns runtime.Namespace) {
	canonical := ns.Namespace("datetime")
	canonical.Function().A0().
		Add("now", Now)

	canonical.Function().A1().
		Add("parse", parse1).
		Add("day_of_week", DayOfWeek).
		Add("year", Year).
		Add("month", Month).
		Add("day", Day).
		Add("hour", Hour).
		Add("minute", Minute).
		Add("second", Second).
		Add("millisecond", Millisecond).
		Add("day_of_year", DayOfYear).
		Add("is_leap_year", IsLeapYear).
		Add("quarter", Quarter).
		Add("days_in_month", DaysInMonth)

	canonical.Function().A2().
		Add("parse", parse2).
		Add("format", Format)

	canonical.Function().A3().
		Add("add", Add).
		Add("subtract", Subtract)

	canonical.Function().A3().Add("same", Same).Add("diff", Diff)

	registerLegacy(ns)
}
