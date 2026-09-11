package datetime_test

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/datetime"
)

func TestLegacyComparison(t *testing.T) {
	functions := datetimeFunctions(t)
	compare3, _ := functions.A3().Get("date_compare")
	compare4, _ := functions.A4().Get("date_compare")
	for _, tc := range []struct {
		left, right, start, end string
		want                    bool
	}{
		{"1999-02-07T00:00:00Z", "2000-02-09T00:00:00Z", "year", "day", false},
		{"1999-02-07T00:00:00Z", "2000-02-07T00:00:00Z", "year", "millisecond", false},
		{"1999-02-07T00:00:00Z", "2000-02-07T00:00:00Z", "day", "day", true},
		{"2024-06-12T10:20:30Z", "2024-06-12T11:20:30Z", "hour", "second", false},
		{"2024-06-12T10:20:30Z", "2024-06-13T10:20:30Z", "hour", "second", true},
		{"2024-06-12T10:20:30.123Z", "2024-06-12T10:20:30.123999Z", "year", "millisecond", true},
		{"2020-12-31T12:00:00Z", "2021-01-01T12:00:00Z", "week", "week", true},
		{"2023-01-02T12:00:00Z", "2024-01-01T12:00:00Z", "week", "week", false},
		{"2023-01-01T12:00:00Z", "2023-01-09T12:00:00Z", "week", "week", false},
		{"2023-01-02T12:00:00Z", "2023-01-08T12:00:00Z", "week", "week", true},
	} {
		left, right := mustDefaultLayoutDt(tc.left), mustDefaultLayoutDt(tc.right)
		got, err := compare4(t.Context(), left, right, runtime.String(tc.start), runtime.String(tc.end))
		if err != nil || got != runtime.Boolean(tc.want) {
			t.Fatalf("%+v = %v, %v", tc, got, err)
		}

		if tc.end == "millisecond" {
			got, err = compare3(t.Context(), left, right, runtime.String(tc.start))
			if err != nil || got != runtime.Boolean(tc.want) {
				t.Fatalf("default end %+v = %v, %v", tc, got, err)
			}
		}
	}

	base := mustDefaultLayoutDt("2024-01-01T00:00:00Z")
	for _, args := range [][]runtime.Value{
		{base, base, runtime.String("day"), runtime.String("year")},
		{base, base, runtime.String("bogus"), runtime.String("day")},
		{base, base, runtime.String("year"), runtime.String("bogus")},
	} {
		got, err := compare4(t.Context(), args[0], args[1], args[2], args[3])
		if got != runtime.None || !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("invalid range %v = %v, %v", args, got, err)
		}
	}

	for pos := 0; pos < 4; pos++ {
		args := []runtime.Value{base, base, runtime.String("year"), runtime.String("millisecond")}
		args[pos] = runtime.None
		got, err := compare4(t.Context(), args[0], args[1], args[2], args[3])
		arg, ok, _ := runtime.InvalidArgumentDetails(err)
		if got != runtime.None || !errors.Is(err, runtime.ErrInvalidType) || !ok || arg != pos {
			t.Fatalf("invalid argument %d = %v, %v", pos, got, err)
		}
	}
}

func TestLegacyDiff(t *testing.T) {
	functions := datetimeFunctions(t)
	diff3, _ := functions.A3().Get("date_diff")
	diff4, _ := functions.A4().Get("date_diff")
	base := mustDefaultLayoutDt("2024-01-01T00:00:00Z")
	for _, delta := range []time.Duration{0, 90 * time.Minute, -90 * time.Minute, time.Duration(math.MaxInt64), time.Duration(math.MinInt64)} {
		end := runtime.NewDateTime(base.Add(delta))
		for _, tc := range []struct {
			unit    string
			divisor time.Duration
		}{
			{"millisecond", time.Millisecond}, {"second", time.Second}, {"minute", time.Minute}, {"hour", time.Hour},
		} {
			got, err := diff3(t.Context(), base, end, runtime.String(tc.unit))
			expected := runtime.Int(delta / tc.divisor)
			if err != nil || got != expected {
				t.Fatalf("%d %s = %v, %v; want %v", delta, tc.unit, got, err, expected)
			}

			got, err = diff4(t.Context(), base, end, runtime.String(tc.unit), runtime.False)
			if err != nil || got != expected {
				t.Fatalf("false = %v, %v; want %v", got, err, expected)
			}

			canonical, err := datetime.Diff(t.Context(), base, end, runtime.String(tc.unit))
			if err != nil {
				t.Fatal(err)
			}

			got, err = diff4(t.Context(), base, end, runtime.String(tc.unit), runtime.True)
			if err != nil || got != canonical {
				t.Fatalf("true = %v, %v; want %v", got, err, canonical)
			}
		}
	}

	for _, unit := range []string{"day", "week", "month", "year", "unknown"} {
		for _, flag := range []runtime.Boolean{runtime.False, runtime.True} {
			got, err := diff4(t.Context(), base, base, runtime.String(unit), flag)
			if got != runtime.None || !errors.Is(err, runtime.ErrInvalidArgument) {
				t.Fatalf("%s = %v, %v", unit, got, err)
			}
		}
	}

	got, err := diff4(t.Context(), base, base, runtime.String("hour"), runtime.None)
	if got != runtime.None || !errors.Is(err, runtime.ErrInvalidType) {
		t.Fatalf("invalid flag = %v, %v", got, err)
	}

	// Float conversion would round across the next whole-millisecond boundary.
	for _, sign := range []int64{-1, 1} {
		end := runtime.NewDateTime(base.Add(time.Duration(sign * 9_000_000_000_000_999_999)))
		got, err := diff3(t.Context(), base, end, runtime.String("milliseconds"))
		if err != nil || got != runtime.Int(sign*9_000_000_000_000) {
			t.Fatalf("integer precision = %v, %v", got, err)
		}
	}

	far := runtime.NewDateTime(base.AddDate(500, 0, 0))
	got, err = diff3(t.Context(), base, far, runtime.String("hour"))
	if got != runtime.None || !errors.Is(err, runtime.ErrRange) {
		t.Fatalf("range = %v, %v", got, err)
	}
}

func datetimeFunctions(t *testing.T) *runtime.Functions {
	t.Helper()
	ns := runtime.NewLibrary()
	datetime.RegisterLib(ns)
	functions, err := ns.Build()
	if err != nil {
		t.Fatal(err)
	}

	return functions
}
