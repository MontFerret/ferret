package runtime_test

import (
	"errors"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestToDateTimeEpochUnits(t *testing.T) {
	t.Parallel()

	expected := time.Unix(1, 0).UTC()
	tests := []struct {
		value runtime.Value
		name  string
		unit  string
	}{
		{name: "seconds short", value: runtime.NewInt(1), unit: "s"},
		{name: "seconds sec", value: runtime.NewInt(1), unit: "sec"},
		{name: "seconds singular", value: runtime.NewInt(1), unit: "second"},
		{name: "seconds plural", value: runtime.NewInt(1), unit: "seconds"},
		{name: "milliseconds short", value: runtime.NewInt(1_000), unit: "ms"},
		{name: "milliseconds singular", value: runtime.NewInt(1_000), unit: "millisecond"},
		{name: "milliseconds plural", value: runtime.NewInt(1_000), unit: "milliseconds"},
		{name: "microseconds ascii short", value: runtime.NewInt(1_000_000), unit: "us"},
		{name: "microseconds micro sign", value: runtime.NewInt(1_000_000), unit: "µs"},
		{name: "microseconds greek mu", value: runtime.NewInt(1_000_000), unit: "μs"},
		{name: "microseconds singular", value: runtime.NewInt(1_000_000), unit: "microsecond"},
		{name: "microseconds plural", value: runtime.NewInt(1_000_000), unit: "microseconds"},
		{name: "nanoseconds short", value: runtime.NewInt(1_000_000_000), unit: "ns"},
		{name: "nanoseconds singular", value: runtime.NewInt(1_000_000_000), unit: "nanosecond"},
		{name: "nanoseconds plural", value: runtime.NewInt(1_000_000_000), unit: "nanoseconds"},
		{name: "case insensitive seconds", value: runtime.NewInt(1), unit: "SECONDS"},
		{name: "case insensitive milliseconds", value: runtime.NewInt(1_000), unit: "MS"},
		{name: "case insensitive microseconds", value: runtime.NewInt(1_000_000), unit: "MICROSECONDS"},
		{name: "case insensitive nanoseconds", value: runtime.NewInt(1_000_000_000), unit: "NS"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := runtime.ToDateTimeEpoch(t.Context(), test.value, runtime.NewString(test.unit))
			if err != nil || !actual.Equal(expected) || actual.Location() != time.UTC {
				t.Fatalf("ToDateTimeEpoch(%v, %q) = %v, %v; want %v UTC", test.value, test.unit, actual, err, expected)
			}
		})
	}
}

func TestToDateTimeEpochFractionsAndNegativeValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expected time.Time
		value    runtime.Value
		name     string
		unit     string
	}{
		{name: "zero", value: runtime.NewInt(0), unit: "s", expected: time.Unix(0, 0).UTC()},
		{name: "negative second", value: runtime.NewInt(-1), unit: "s", expected: time.Unix(-1, 0).UTC()},
		{name: "fractional second", value: runtime.NewFloat(1.5), unit: "s", expected: time.Unix(1, 500_000_000).UTC()},
		{name: "negative fractional second", value: runtime.NewFloat(-1.5), unit: "s", expected: time.Unix(-2, 500_000_000).UTC()},
		{name: "fractional millisecond", value: runtime.NewFloat(1.5), unit: "ms", expected: time.Unix(0, 1_500_000).UTC()},
		{name: "negative fractional millisecond", value: runtime.NewFloat(-1.5), unit: "ms", expected: time.Unix(-1, 998_500_000).UTC()},
		{name: "positive sub-nanosecond truncates", value: runtime.NewFloat(0.5), unit: "ns", expected: time.Unix(0, 0).UTC()},
		{name: "negative sub-nanosecond truncates", value: runtime.NewFloat(-0.5), unit: "ns", expected: time.Unix(0, 0).UTC()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := runtime.ToDateTimeEpoch(t.Context(), test.value, runtime.NewString(test.unit))
			if err != nil || !actual.Equal(test.expected) {
				t.Fatalf("ToDateTimeEpoch(%v, %q) = %v, %v; want %v", test.value, test.unit, actual, err, test.expected)
			}
		})
	}
}

func TestToDateTimeEpochRange(t *testing.T) {
	t.Parallel()

	const unixToInternal = int64(62_135_596_800)
	maximumSeconds := int64(math.MaxInt64 - unixToInternal)

	maximum, err := runtime.ToDateTimeEpoch(t.Context(), runtime.NewInt64(maximumSeconds), runtime.NewString("s"))
	if err != nil || maximum.Unix() != maximumSeconds {
		t.Fatalf("maximum DateTime = %v, %v", maximum, err)
	}

	minimum, err := runtime.ToDateTimeEpoch(t.Context(), runtime.NewInt64(math.MinInt64), runtime.NewString("s"))
	if err != nil || minimum.Unix() != math.MinInt64 {
		t.Fatalf("minimum DateTime = %v, %v", minimum, err)
	}

	for _, input := range []runtime.Value{
		runtime.NewInt64(maximumSeconds + 1),
		runtime.NewFloat(math.MaxFloat64),
	} {
		if _, err := runtime.ToDateTimeEpoch(t.Context(), input, runtime.NewString("s")); !errors.Is(err, runtime.ErrRange) {
			t.Fatalf("ToDateTimeEpoch(%v, s) error = %v, want ErrRange", input, err)
		}
	}
}

func TestToDateTimeEpochErrors(t *testing.T) {
	t.Parallel()

	existing := runtime.NewDateTime(time.Unix(0, 0).UTC())
	tests := []struct {
		input    runtime.Value
		unit     runtime.Value
		category error
		name     string
		contains []string
	}{
		{
			name:     "numeric string",
			input:    runtime.NewString("1"),
			unit:     runtime.NewString("s"),
			category: runtime.ErrInvalidArgument,
			contains: []string{"String", `"1"`, "Int or Float"},
		},
		{
			name:     "RFC3339 string with unit",
			input:    runtime.NewString("2026-08-02T12:00:00Z"),
			unit:     runtime.NewString("ms"),
			category: runtime.ErrInvalidArgument,
			contains: []string{"String", "epoch units are only valid"},
		},
		{
			name:     "DateTime with unit",
			input:    existing,
			unit:     runtime.NewString("s"),
			category: runtime.ErrInvalidArgument,
			contains: []string{"DateTime", "epoch units are only valid"},
		},
		{
			name:     "unsupported input",
			input:    runtime.NewObject(),
			unit:     runtime.NewString("s"),
			category: runtime.ErrInvalidType,
			contains: []string{"Object", "expected Int or Float"},
		},
		{
			name:     "NONE unit",
			input:    runtime.NewInt(1),
			unit:     runtime.None,
			category: runtime.ErrInvalidType,
			contains: []string{"None", "expected String", `"s", "ms", "us", or "ns"`},
		},
		{
			name:     "numeric unit",
			input:    runtime.NewInt(1),
			unit:     runtime.NewInt(1),
			category: runtime.ErrInvalidType,
			contains: []string{"Int", "expected String", `"s", "ms", "us", or "ns"`},
		},
		{
			name:     "unknown unit",
			input:    runtime.NewInt(1),
			unit:     runtime.NewString("minutes"),
			category: runtime.ErrInvalidArgument,
			contains: []string{"minutes", "unsupported epoch unit", `"s", "ms", "us", or "ns"`},
		},
		{
			name:     "unit whitespace is not ignored",
			input:    runtime.NewInt(1),
			unit:     runtime.NewString(" s"),
			category: runtime.ErrInvalidArgument,
			contains: []string{`" s"`, "unsupported epoch unit"},
		},
		{
			name:     "NaN",
			input:    runtime.NewFloat(math.NaN()),
			unit:     runtime.NewString("s"),
			category: runtime.ErrInvalidArgument,
			contains: []string{"Float", "finite"},
		},
		{
			name:     "positive infinity",
			input:    runtime.NewFloat(math.Inf(1)),
			unit:     runtime.NewString("s"),
			category: runtime.ErrInvalidArgument,
			contains: []string{"Float", "finite"},
		},
		{
			name:     "negative infinity",
			input:    runtime.NewFloat(math.Inf(-1)),
			unit:     runtime.NewString("s"),
			category: runtime.ErrInvalidArgument,
			contains: []string{"Float", "finite"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := runtime.ToDateTimeEpoch(t.Context(), test.input, test.unit)
			if !errors.Is(err, test.category) {
				t.Fatalf("error = %v, want %v", err, test.category)
			}
			for _, expected := range test.contains {
				if !strings.Contains(err.Error(), expected) {
					t.Fatalf("error %q does not contain %q", err, expected)
				}
			}
		})
	}
}
