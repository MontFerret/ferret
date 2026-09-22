package datetime_test

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/datetime"
)

func TestDateTimeShiftIntegerWidth(t *testing.T) {
	base := time.Date(2024, time.February, 28, 12, 30, 15, 123, time.UTC)
	for _, amount := range []runtime.Int{math.MinInt32 - 1, math.MinInt32, math.MaxInt32, math.MaxInt32 + 1, 1 << 32} {
		for _, subtract := range []bool{false, true} {
			function := datetime.Add
			delta := time.Duration(amount) * time.Millisecond
			if subtract {
				function = datetime.Subtract
				delta = -delta
			}

			got, err := function(t.Context(), runtime.NewDateTime(base), amount, runtime.String("millisecond"))
			date, ok := got.(runtime.DateTime)
			if err != nil || !ok || !date.Equal(base.Add(delta)) {
				t.Errorf("shift %d, subtract=%v: %v, %v", amount, subtract, got, err)
			}
		}
	}
}

func TestDateTimeShiftRangeErrors(t *testing.T) {
	base := runtime.NewDateTime(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
	for _, amount := range []runtime.Int{math.MinInt64, math.MaxInt64, 1 << 50, math.MaxInt64 / 7} {
		for _, unit := range []runtime.String{"millisecond", "second", "minute", "hour", "day", "week", "month", "year"} {
			for _, function := range []runtime.Function3{datetime.Add, datetime.Subtract} {
				got, err := function(t.Context(), base, amount, unit)
				position, attributed, _ := runtime.InvalidArgumentDetails(err)
				if got != runtime.None || !errors.Is(err, runtime.ErrRange) || !attributed || position != 1 {
					t.Errorf("shift(%d, %s) = %v, %v; want argument-1 range error", amount, unit, got, err)
				}
			}
		}
	}
}

func TestCalendarShiftNativeBoundary(t *testing.T) {
	base := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	amount := runtime.Int(math.MaxInt32) + 1
	got, err := datetime.Add(t.Context(), runtime.NewDateTime(base), amount, runtime.String("day"))
	if math.MaxInt == math.MaxInt32 {
		if !errors.Is(err, runtime.ErrRange) {
			t.Fatalf("oversized calendar delta = %v, %v", got, err)
		}
	} else {
		date, ok := got.(runtime.DateTime)
		if err != nil || !ok || !date.Equal(base.AddDate(0, 0, int(amount))) {
			t.Fatalf("representable calendar delta = %v, %v", got, err)
		}
	}

	// The delta itself fits, but AddDate's year addition must also fit.
	_, err = datetime.Add(t.Context(), runtime.NewDateTime(base), runtime.Int(math.MaxInt), runtime.String("year"))
	if !errors.Is(err, runtime.ErrRange) {
		t.Fatalf("calendar component overflow = %v", err)
	}
}

func TestCalendarShiftPreservesNormalization(t *testing.T) {
	for _, zone := range []string{"UTC", "America/New_York", "Pacific/Apia"} {
		location, err := time.LoadLocation(zone)
		if err != nil {
			t.Fatal(err)
		}

		for _, base := range []time.Time{
			time.Date(-10001, time.January, 31, 12, 0, 0, 123, location),
			time.Date(1000000, time.February, 29, 12, 0, 0, 123, location),
			time.Date(2024, time.January, 31, 12, 0, 0, 123, location),
			time.Date(2024, time.February, 29, 12, 0, 0, 123, location),
			time.Date(2024, time.March, 9, 2, 30, 0, 123, location),
			time.Date(2024, time.November, 2, 1, 30, 0, 123, location),
			time.Date(2011, time.December, 29, 12, 0, 0, 123, location),
		} {
			for _, amount := range []int{-3660001, -401, -13, -1, 0, 1, 13, 401, 3660001} {
				for _, unit := range []runtime.String{"day", "week", "month", "year"} {
					var years, months, days int
					switch unit {
					case "day":
						days = amount
					case "week":
						days = amount * 7
					case "month":
						months = amount
					case "year":
						years = amount
					}

					for _, subtract := range []bool{false, true} {
						function := datetime.Add
						direction := 1
						if subtract {
							function = datetime.Subtract
							direction = -1
						}

						want := base.AddDate(years*direction, months*direction, days*direction)
						got, err := function(t.Context(), runtime.NewDateTime(base), runtime.Int(amount), unit)
						date, ok := got.(runtime.DateTime)
						if err != nil || !ok || !date.Equal(want) || date.Location() != location {
							t.Fatalf("%s %v shift %d %s subtract=%v = %v, %v; want %v", zone, base, amount, unit, subtract, got, err, want)
						}
					}
				}
			}
		}
	}
}

func TestCalendarShiftRejectsNarrowedSourceYear(t *testing.T) {
	for _, seconds := range []int64{-(1 << 60), 1 << 60} {
		base := time.Unix(seconds, 123).UTC()
		for _, amount := range []runtime.Int{0, 1} {
			got, err := datetime.Add(t.Context(), runtime.NewDateTime(base), amount, runtime.String("day"))
			if math.MaxInt == math.MaxInt32 {
				position, attributed, _ := runtime.InvalidArgumentDetails(err)
				if got != runtime.None || !errors.Is(err, runtime.ErrRange) || !attributed || position != 1 {
					t.Fatalf("narrowed source year: %v, %v; want argument-1 range error", got, err)
				}
			} else {
				date, ok := got.(runtime.DateTime)
				if err != nil || !ok || !date.Equal(base.AddDate(0, 0, int(amount))) {
					t.Fatalf("representable source year: %v, %v", got, err)
				}
			}
		}
	}
}
