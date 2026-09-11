package datetime_test

import (
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/datetime"
)

func TestCalendarMonths(t *testing.T) {
	for _, tc := range []struct {
		year int
		leap bool
	}{
		{1900, false}, {2000, true}, {2023, false}, {2024, true}, {2100, false},
	} {
		lengths := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
		if tc.leap {
			lengths[1] = 29
		}

		for index, days := range lengths {
			date := runtime.NewDateTime(time.Date(tc.year, time.Month(index+1), 1, 0, 0, 0, 0, time.FixedZone("offset", 14*60*60)))
			got, err := datetime.DaysInMonth(t.Context(), date)
			if err != nil || got != runtime.Int(days) {
				t.Fatalf("%d/%d = %v, %v; want %d days", tc.year, index+1, got, err, days)
			}

			got, err = datetime.IsLeapYear(t.Context(), date)
			if err != nil || got != runtime.Boolean(tc.leap) {
				t.Fatalf("leap %d = %v, %v", tc.year, got, err)
			}
		}
	}
}

func TestCalendarNormalization(t *testing.T) {
	for _, tc := range []struct {
		start, unit, want string
		amount            runtime.Int
		subtract          bool
	}{
		{"2023-01-31T12:00:00Z", "month", "2023-03-03T12:00:00Z", 1, false},
		{"2024-01-31T12:00:00Z", "month", "2024-03-02T12:00:00Z", 1, false},
		{"2024-02-29T12:00:00Z", "year", "2025-03-01T12:00:00Z", 1, false},
		{"2024-02-29T12:00:00Z", "year", "2023-03-01T12:00:00Z", 1, true},
		{"2023-03-31T12:00:00Z", "month", "2023-03-03T12:00:00Z", 1, true},
		{"2024-12-31T12:00:00Z", "day", "2025-01-01T12:00:00Z", 1, false},
		{"2024-12-31T12:00:00Z", "week", "2025-01-07T12:00:00Z", 1, false},
		{"2024-03-01T12:00:00Z", "days", "2024-02-29T12:00:00Z", -1, false},
	} {
		fn := datetime.Add
		if tc.subtract {
			fn = datetime.Subtract
		}

		got, err := fn(t.Context(), mustDefaultLayoutDt(tc.start), tc.amount, runtime.String(tc.unit))
		date, ok := got.(runtime.DateTime)
		if err != nil || !ok || !date.Equal(mustDefaultLayoutDt(tc.want).Time) {
			t.Fatalf("%+v = %v, %v", tc, got, err)
		}
	}
}

func TestCalendarDST(t *testing.T) {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		month time.Month
		day   int
		hours runtime.Float
	}{
		{time.March, 9, 23}, {time.November, 2, 25},
	} {
		start := runtime.NewDateTime(time.Date(2024, tc.month, tc.day, 12, 0, 0, 0, location))
		next, err := datetime.Add(t.Context(), start, runtime.Int(1), runtime.String("day"))
		if err != nil {
			t.Fatal(err)
		}

		date := next.(runtime.DateTime)
		if date.Hour() != 12 || date.Day() != tc.day+1 || date.Location() != location {
			t.Fatalf("calendar shift = %v", date)
		}

		elapsed, err := datetime.Diff(t.Context(), start, next, runtime.String("hours"))
		if err != nil || elapsed != tc.hours {
			t.Fatalf("DST elapsed = %v, %v; want %v", elapsed, err, tc.hours)
		}

		restored, err := datetime.Subtract(t.Context(), next, runtime.Int(1), runtime.String("day"))
		if err != nil || !restored.(runtime.DateTime).Equal(start.Time) {
			t.Fatalf("DST subtraction = %v, %v", restored, err)
		}

		fixed, err := datetime.Add(t.Context(), start, runtime.Int(24), runtime.String("hours"))
		if err != nil || !fixed.(runtime.DateTime).Equal(start.Add(24*time.Hour)) {
			t.Fatalf("elapsed shift = %v, %v", fixed, err)
		}
	}

	first := runtime.NewDateTime(mustDefaultLayoutDt("2024-11-03T05:30:00Z").In(location))
	second := runtime.NewDateTime(mustDefaultLayoutDt("2024-11-03T06:30:00Z").In(location))
	same, err := datetime.Same(t.Context(), first, second, runtime.String("millisecond"))
	if err != nil || same != runtime.True {
		t.Fatalf("repeated local time = %v, %v", same, err)
	}
}
