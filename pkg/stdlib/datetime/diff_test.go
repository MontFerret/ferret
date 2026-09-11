package datetime_test

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/datetime"
)

func TestDiffSignedFractions(t *testing.T) {
	base := mustDefaultLayoutDt("2024-02-28T12:00:00Z")
	for _, tc := range []struct {
		unit  string
		delta time.Duration
		want  float64
	}{
		{"millisecond", 1500 * time.Microsecond, 1.5},
		{"second", 1500 * time.Millisecond, 1.5},
		{"minute", 90 * time.Second, 1.5},
		{"hour", 90 * time.Minute, 1.5},
		{"hours", 48 * time.Hour, 48},
	} {
		for _, sign := range []int{-1, 0, 1} {
			end := runtime.NewDateTime(base.Add(time.Duration(sign) * tc.delta))
			got, err := datetime.Diff(t.Context(), base, end, runtime.String(tc.unit))
			if err != nil || got != runtime.Float(float64(sign)*tc.want) {
				t.Fatalf("%s sign %d = %v (%T), %v", tc.unit, sign, got, got, err)
			}

			reversed, err := datetime.Diff(t.Context(), end, base, runtime.String(tc.unit))
			if err != nil || reversed != runtime.Float(-float64(sign)*tc.want) {
				t.Fatalf("reversed %s = %v, %v", tc.unit, reversed, err)
			}
		}
	}

	for _, tc := range []struct {
		start, end string
		hours      float64
	}{
		{"2024-02-28T00:00:00Z", "2024-03-01T00:00:00Z", 48},
		{"2023-02-28T00:00:00Z", "2023-03-01T00:00:00Z", 24},
		{"2024-04-01T00:00:00Z", "2024-05-01T00:00:00Z", 720},
		{"2024-01-01T00:00:00Z", "2025-01-01T00:00:00Z", 8784},
		{"2024-12-31T23:30:00Z", "2025-01-01T00:30:00Z", 1},
		{"2024-03-10T01:30:00-05:00", "2024-03-10T03:30:00-04:00", 1},
		{"2024-11-03T01:30:00-04:00", "2024-11-03T01:30:00-05:00", 1},
		{"2024-06-12T10:00:00Z", "2024-06-12T15:30:00+05:30", 0},
	} {
		got, err := datetime.Diff(t.Context(), mustDefaultLayoutDt(tc.start), mustDefaultLayoutDt(tc.end), runtime.String("hours"))
		if err != nil || got != runtime.Float(tc.hours) {
			t.Fatalf("%s / %s = %v, %v; want %v hours", tc.start, tc.end, got, err, tc.hours)
		}
	}
}

func TestDiffRange(t *testing.T) {
	base := mustDefaultLayoutDt("2000-01-01T00:00:00Z")
	for _, delta := range []time.Duration{time.Duration(math.MaxInt64), time.Duration(math.MinInt64)} {
		end := runtime.NewDateTime(base.Add(delta))
		got, err := datetime.Diff(t.Context(), base, end, runtime.String("milliseconds"))
		if err != nil || got != runtime.Float(float64(delta)/float64(time.Millisecond)) {
			t.Fatalf("boundary %d = %v, %v", delta, got, err)
		}
	}

	for _, end := range []time.Time{
		base.Add(time.Duration(math.MaxInt64)).Add(time.Nanosecond),
		base.Add(time.Duration(math.MinInt64)).Add(-time.Nanosecond),
		time.Date(2500, 1, 1, 0, 0, 0, 0, time.UTC),
	} {
		got, err := datetime.Diff(t.Context(), base, runtime.NewDateTime(end), runtime.String("hours"))
		if got != runtime.None || !errors.Is(err, runtime.ErrRange) {
			t.Fatalf("overflow = %v, %v", got, err)
		}
	}
}

func TestDateTimeArgumentErrors(t *testing.T) {
	base := mustDefaultLayoutDt("2024-01-01T00:00:00Z")
	for _, fn := range []runtime.Function3{datetime.Same, datetime.Diff} {
		for pos := 0; pos < 3; pos++ {
			args := []runtime.Value{base, base, runtime.String("hour")}
			args[pos] = runtime.None
			got, err := fn(t.Context(), args[0], args[1], args[2])
			arg, ok, _ := runtime.InvalidArgumentDetails(err)
			if got != runtime.None || !errors.Is(err, runtime.ErrInvalidType) || !ok || arg != pos {
				t.Fatalf("argument %d = %v, %v", pos, got, err)
			}
		}

		got, err := fn(t.Context(), base, base, runtime.String("unknown"))
		if got != runtime.None || !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("unknown unit on equal dates = %v, %v", got, err)
		}
	}

	for _, name := range []string{"d", "day", "days", "w", "week", "weeks", "m", "month", "months", "y", "year", "years"} {
		for _, right := range []runtime.DateTime{base, runtime.NewDateTime(base.Add(time.Hour))} {
			got, err := datetime.Diff(t.Context(), base, right, runtime.String(name))
			pos, ok, _ := runtime.InvalidArgumentDetails(err)
			if got != runtime.None || !errors.Is(err, runtime.ErrInvalidArgument) || !ok || pos != 2 {
				t.Fatalf("%s = %v, %v", name, got, err)
			}
		}
	}
}
