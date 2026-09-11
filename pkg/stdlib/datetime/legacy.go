package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Forwarders attach deprecation metadata without deprecating canonical functions.
func registerLegacy(ns runtime.Namespace) {
	ns.Function().A0().
		Add("now", legacyNow)

	ns.Function().A1().
		Add("date", legacyParse1).
		Add("date_dayofweek", legacyDayOfWeek).
		Add("date_year", legacyYear).
		Add("date_month", legacyMonth).
		Add("date_day", legacyDay).
		Add("date_hour", legacyHour).
		Add("date_minute", legacyMinute).
		Add("date_second", legacySecond).
		Add("date_millisecond", legacyMillisecond).
		Add("date_dayofyear", legacyDayOfYear).
		Add("date_leapyear", legacyIsLeapYear).
		Add("date_quarter", legacyQuarter).
		Add("date_days_in_month", legacyDaysInMonth)

	ns.Function().A2().
		Add("date", legacyParse2).
		Add("date_format", legacyFormat)

	ns.Function().A3().
		Add("date_add", legacyAdd).
		Add("date_subtract", legacySubtract)

	ns.Function().A3().Add("date_compare", legacyCompare3).Add("date_diff", legacyDiff3)
	ns.Function().A4().Add("date_compare", legacyCompare4).Add("date_diff", legacyDiff4)
}

// Now returns new DateTime object with Time equal to time.Now().
// @return {DateTime} New DateTime object.
// @deprecated Use datetime::now instead.
func legacyNow(ctx context.Context) (runtime.Value, error) {
	return Now(ctx)
}

// Parse parses a formatted string and returns DateTime object it represents.
// @param time {String} String representation of DateTime.
// @return {DateTime} Parsed DateTime, preserving the parsed offset.
// @deprecated Use datetime::parse instead.
func legacyParse1(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return parse1(ctx, arg1)
}

// DayOfWeek returns number of the weekday from the date. Sunday is the 0th day of week.
// @param date {DateTime} Source DateTime.
// @return {Int} Number of the weekday.
// @deprecated Use datetime::day_of_week instead.
func legacyDayOfWeek(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return DayOfWeek(ctx, arg1)
}

// Year returns the year extracted from the given date.
// @param date {DateTime} Source DateTime.
// @return {Int} A year number.
// @deprecated Use datetime::year instead.
func legacyYear(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Year(ctx, arg1)
}

// Month returns the month of date as a number.
// @param date {DateTime} Source DateTime.
// @return {Int} A month number.
// @deprecated Use datetime::month instead.
func legacyMonth(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Month(ctx, arg1)
}

// Day returns the day of date as a number.
// @param date {DateTime} Source DateTime.
// @return {Int} A day number.
// @deprecated Use datetime::day instead.
func legacyDay(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Day(ctx, arg1)
}

// Hour returns the hour of date as a number.
// @param date {DateTime} Source DateTime.
// @return {Int} An hour number.
// @deprecated Use datetime::hour instead.
func legacyHour(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Hour(ctx, arg1)
}

// Minute returns the minute of date as a number.
// @param date {DateTime} Source DateTime.
// @return {Int} A minute number.
// @deprecated Use datetime::minute instead.
func legacyMinute(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Minute(ctx, arg1)
}

// Second returns the second of date as a number.
// @param date {DateTime} Source DateTime.
// @return {Int} A second number.
// @deprecated Use datetime::second instead.
func legacySecond(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Second(ctx, arg1)
}

// Millisecond returns the millisecond of date as a number.
// @param date {DateTime} Source DateTime.
// @return {Int} A millisecond number.
// @deprecated Use datetime::millisecond instead.
func legacyMillisecond(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Millisecond(ctx, arg1)
}

// DayOfYear returns the day of year number of date.
// The return value range from 1 to 365 (366 in a leap year).
// @param date {DateTime} Source DateTime.
// @return {Int} A day of year number.
// @deprecated Use datetime::day_of_year instead.
func legacyDayOfYear(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return DayOfYear(ctx, arg1)
}

// IsLeapYear returns true if date is in a leap year else false.
// @param date {DateTime} Source DateTime.
// @return {Boolean} Whether the date is in a leap year.
// @deprecated Use datetime::is_leap_year instead.
func legacyIsLeapYear(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return IsLeapYear(ctx, arg1)
}

// Quarter returns which quarter date belongs to.
// @param date {DateTime} Source DateTime.
// @return {Int} A quarter number.
// @deprecated Use datetime::quarter instead.
func legacyQuarter(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Quarter(ctx, arg1)
}

// DaysInMonth returns the number of days in the input's local calendar month.
// @param date {DateTime} Source date.
// @return {Int} Number of calendar days in the month, including leap days.
// @deprecated Use datetime::days_in_month instead.
func legacyDaysInMonth(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return DaysInMonth(ctx, arg1)
}

// Parse parses a formatted string and returns DateTime object it represents.
// @param time {String} String representation of DateTime.
// @param layout {String} Go time layout.
// @return {DateTime} Parsed DateTime, preserving the parsed offset.
// @deprecated Use datetime::parse instead.
func legacyParse2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return parse2(ctx, arg1, arg2)
}

// Format renders a DateTime using a Go time layout.
// @param date {DateTime} Source DateTime object.
// @param format {String} Go time layout.
// @return {String} Formatted date.
// @deprecated Use datetime::format instead.
func legacyFormat(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return Format(ctx, arg1, arg2)
}

// Add adds an integer number of units to a date. Subday units are elapsed time;
// calendar units use Go AddDate in the input's location, normalizing month ends
// rather than clamping them. A calendar day can span 23 or 25 elapsed hours.
// @param date {DateTime} Source date.
// @param amount {Int} Number of units to add; may be negative.
// @param unit {String} Millisecond, second, minute, hour, day, week, month, or year; plurals are accepted.
// @return {DateTime} Calculated date, retaining the input's location.
// @deprecated Use datetime::add instead.
func legacyAdd(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return Add(ctx, arg1, arg2, arg3)
}

// Subtract subtracts an integer number of units from a date. Subday units are
// elapsed time; calendar units use Go AddDate in the input's location, normalizing
// month ends rather than clamping them. A calendar day can span 23 or 25 hours.
// @param date {DateTime} Source date.
// @param amount {Int} Number of units to subtract; may be negative.
// @param unit {String} Millisecond, second, minute, hour, day, week, month, or year; plurals are accepted.
// @return {DateTime} Calculated date, retaining the input's location.
// @deprecated Use datetime::subtract instead.
func legacySubtract(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return Subtract(ctx, arg1, arg2, arg3)
}
