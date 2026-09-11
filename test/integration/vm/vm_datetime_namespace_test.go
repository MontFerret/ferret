package vm_test

import (
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestDateTimeNamespace(t *testing.T) {
	const prefix = `let d = datetime::parse("2024-02-29T23:59:58.123+05:30") `
	RunSpecs(t, []spec.Spec{
		S(`return datetime::format(datetime::parse("2026-09-11T14:30:00Z"), "2006-01-02")`, "2026-09-11"),
		S(`return datetime::format(datetime::parse("29/02/2024", "02/01/2006"), "2006-01-02")`, "2024-02-29"),
		Array(prefix+`return [datetime::year(d), datetime::month(d), datetime::day(d), datetime::hour(d), datetime::minute(d), datetime::second(d), datetime::millisecond(d), datetime::day_of_week(d), datetime::day_of_year(d), datetime::quarter(d), datetime::days_in_month(d), datetime::is_leap_year(d)]`, []any{2024, 2, 29, 23, 59, 58, 123, 4, 60, 1, 29, true}),
		S(`return DaTeTiMe::YeAr(DATETIME::PARSE("2026-09-11T14:30:00Z"))`, 2026),
		S(`let before = now() let current = datetime::now() let after = now() return is_datetime(current) && before <= current && current <= after`, true),
		S(prefix+`return datetime::same(d, datetime::parse("2024-02-28T12:00:00Z"), "month")`, true),
		S(prefix+`return datetime::same(d, datetime::parse("2023-02-28T12:00:00Z"), "month")`, false),
		S(`return datetime::same(datetime::parse("2020-12-31T12:00:00Z"),datetime::parse("2021-01-01T12:00:00Z"),"week")`, true),
		S(`return datetime::same(datetime::parse("2023-01-02T12:00:00Z"),datetime::parse("2024-01-01T12:00:00Z"),"week")`, false),
		S(prefix+`return datetime::same(d, datetime::parse("2024-02-29T23:59:58.123999Z"), "millisecond")`, true),
		S(`return datetime::diff(datetime::parse("2024-12-31T23:30:00Z"),datetime::parse("2025-01-01T01:00:00Z"),"hours")`, 1.5),
		S(`return datetime::diff(datetime::parse("2025-01-01T01:00:00Z"),datetime::parse("2024-12-31T23:30:00Z"),"hours")`, -1.5),
		S(prefix+`return is_float(datetime::diff(d,d,"second"))`, true),
		S(`return datetime::format(datetime::add(datetime::parse("2023-01-31T12:00:00Z"),1,"month"),"2006-01-02")`, "2023-03-03"),
		S(`return datetime::format(datetime::subtract(datetime::parse("2024-02-29T12:00:00Z"),1,"year"),"2006-01-02")`, "2023-03-01"),
		S(prefix+`return date("2024-02-29T23:59:58.123+05:30") == datetime::parse("2024-02-29T23:59:58.123+05:30")`, true, "date arity 1 parity"),
		S(prefix+`return date_dayofweek(d) == datetime::day_of_week(d)`, true, "date_dayofweek arity 1 parity"),
		S(prefix+`return date_year(d) == datetime::year(d)`, true, "date_year arity 1 parity"),
		S(prefix+`return date_month(d) == datetime::month(d)`, true, "date_month arity 1 parity"),
		S(prefix+`return date_day(d) == datetime::day(d)`, true, "date_day arity 1 parity"),
		S(prefix+`return date_hour(d) == datetime::hour(d)`, true, "date_hour arity 1 parity"),
		S(prefix+`return date_minute(d) == datetime::minute(d)`, true, "date_minute arity 1 parity"),
		S(prefix+`return date_second(d) == datetime::second(d)`, true, "date_second arity 1 parity"),
		S(prefix+`return date_millisecond(d) == datetime::millisecond(d)`, true, "date_millisecond arity 1 parity"),
		S(prefix+`return date_dayofyear(d) == datetime::day_of_year(d)`, true, "date_dayofyear arity 1 parity"),
		S(prefix+`return date_leapyear(d) == datetime::is_leap_year(d)`, true, "date_leapyear arity 1 parity"),
		S(prefix+`return date_quarter(d) == datetime::quarter(d)`, true, "date_quarter arity 1 parity"),
		S(prefix+`return date_days_in_month(d) == datetime::days_in_month(d)`, true, "date_days_in_month arity 1 parity"),
		S(prefix+`return date("29/02/2024", "02/01/2006") == datetime::parse("29/02/2024", "02/01/2006")`, true, "date arity 2 parity"),
		S(prefix+`return date_format(d, "2006-01-02T15:04:05.000Z07:00") == datetime::format(d, "2006-01-02T15:04:05.000Z07:00")`, true, "date_format arity 2 parity"),
		S(prefix+`return date_add(d, 2, "months") == datetime::add(d, 2, "months")`, true, "date_add arity 3 parity"),
		S(prefix+`return date_subtract(d, 2, "months") == datetime::subtract(d, 2, "months")`, true, "date_subtract arity 3 parity"),
		S(prefix+`let end = datetime::add(d,90,"minutes") return date_diff(d,end,"hours") == 1 && date_diff(end,d,"hours",false) == (-1) && date_diff(d,end,"hours",true) == datetime::diff(d,end,"hours")`, true),
		S(prefix+`return is_int(date_diff(d,d,"hour")) && is_float(date_diff(d,d,"hour",true))`, true),
		S(`return date_compare(date("1999-02-07T00:00:00Z"),date("2000-02-09T00:00:00Z"),"year","day")`, false),
		S(prefix+`return date_compare(d,d,"year")`, true),
		S(`return date_compare(date("2020-12-31T12:00:00Z"),date("2021-01-01T12:00:00Z"),"week","week")`, true),
	})
}

func TestDateTimeMonthLengths(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`return (for month in 1..12 let d = datetime::parse(concat("2023-", month < 10 ? "0" : "", to_string(month), "-01"), "2006-01-02") return datetime::days_in_month(d))`, []any{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}),
		Array(`return (for year in [1900,2000,2024,2100] let d = datetime::parse(concat(to_string(year),"-02-01"),"2006-01-02") return [datetime::is_leap_year(d),datetime::days_in_month(d)])`, []any{[]any{false, 28}, []any{true, 29}, []any{true, 29}, []any{false, 28}}),
	})
}

func TestDateTimeNamespaceErrors(t *testing.T) {
	const prefix = `let d = datetime::parse("2024-01-01T00:00:00Z") `
	var cases []spec.Spec
	for _, expression := range []string{
		`datetime::parse()`, `datetime::parse("x","y","z")`, `datetime::parse(2024)`,
		`datetime::parse("not-a-date")`, `datetime::format(d,1)`, `datetime::year("2024-01-01")`,
		`datetime::same(d,d)`, `datetime::same(d,d,"year","day")`, `datetime::same(d,d,"unknown")`,
		`datetime::diff(d,d)`, `datetime::diff(d,d,"hour",true)`, `datetime::diff(d,d,"unknown")`,
		`date_compare(d,d,"day","year")`, `date_compare(d,d)`, `date_compare(d,d,"year","day",true)`,
		`date_diff(d,d)`, `date_diff(d,d,"hour",true,true)`, `date_diff(d,d,"hour",1)`,
		`datetime::compare(d,d,"year")`, `date_now()`,
		`datetime::diff(d,datetime::parse("2500-01-01T00:00:00Z"),"hours")`,
	} {
		cases = append(cases, Error(prefix+"return "+expression))
	}

	for _, unit := range []string{"d", "day", "days", "w", "week", "weeks", "m", "month", "months", "y", "year", "years"} {
		cases = append(cases, Error(prefix+`return datetime::diff(d,d,"`+unit+`")`))
		cases = append(cases, Error(prefix+`return date_diff(d,d,"`+unit+`",true)`))
	}

	RunSpecs(t, cases)
}

func TestDateTimeNamedLocation(t *testing.T) {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		month time.Month
		day   int
		hours int
	}{
		{time.March, 9, 23}, {time.November, 2, 25},
	} {
		start := runtime.NewDateTime(time.Date(2024, tc.month, tc.day, 12, 0, 0, 0, location))
		RunSpecs(t, []spec.Spec{
			Array(`let next = datetime::add(@start,1,"day") return [datetime::hour(next),datetime::diff(@start,next,"hours"),datetime::subtract(next,1,"day") == @start]`, []any{12, tc.hours, true}),
		}, vm.WithParam("start", start))
	}
}
