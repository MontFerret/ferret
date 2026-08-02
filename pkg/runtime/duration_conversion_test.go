package runtime_test

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestToDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    runtime.Value
		name     string
		expected runtime.Duration
	}{
		{name: "duration", input: runtime.NewDuration(5 * time.Second), expected: runtime.NewDuration(5 * time.Second)},
		{name: "string", input: runtime.NewString("500ms"), expected: runtime.NewDuration(500 * time.Millisecond)},
		{name: "compound string", input: runtime.NewString("1m30s"), expected: runtime.NewDuration(90 * time.Second)},
		{name: "integer milliseconds", input: runtime.NewInt(500), expected: runtime.NewDuration(500 * time.Millisecond)},
		{name: "float milliseconds", input: runtime.NewFloat(1.5), expected: runtime.NewDuration(1500 * time.Microsecond)},
		{name: "fractional nanosecond truncates", input: runtime.NewFloat(0.0000005), expected: runtime.ZeroDuration},
		{name: "none", input: runtime.None, expected: runtime.ZeroDuration},
		{name: "false", input: runtime.False, expected: runtime.ZeroDuration},
		{name: "true", input: runtime.True, expected: runtime.NewDuration(time.Millisecond)},
		{name: "empty list", input: runtime.EmptyArray(), expected: runtime.ZeroDuration},
		{name: "singleton number", input: runtime.NewArrayWith(runtime.NewInt(500)), expected: runtime.NewDuration(500 * time.Millisecond)},
		{name: "singleton string", input: runtime.NewArrayWith(runtime.NewString("5s")), expected: runtime.NewDuration(5 * time.Second)},
		{name: "recursive singleton", input: runtime.NewArrayWith(runtime.NewArrayWith(runtime.NewInt(2))), expected: runtime.NewDuration(2 * time.Millisecond)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual, err := runtime.ToDuration(t.Context(), test.input)
			if err != nil {
				t.Fatalf("ToDuration(%s): %v", test.input, err)
			}
			if actual != test.expected {
				t.Fatalf("ToDuration(%s) = %s, want %s", test.input, actual, test.expected)
			}
		})
	}
}

func TestToDurationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input  runtime.Value
		target error
		name   string
	}{
		{name: "invalid string", input: runtime.NewString("five seconds"), target: runtime.ErrInvalidArgument},
		{name: "object", input: runtime.NewObject(), target: runtime.ErrInvalidType},
		{name: "multiple list items", input: runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)), target: runtime.ErrInvalidArgument},
		{name: "integer overflow", input: runtime.NewInt64(9_223_372_036_855), target: runtime.ErrRange},
		{name: "float overflow", input: runtime.NewFloat(math.MaxFloat64), target: runtime.ErrRange},
		{name: "not a number", input: runtime.NewFloat(math.NaN()), target: runtime.ErrInvalidArgument},
		{name: "infinity", input: runtime.NewFloat(math.Inf(1)), target: runtime.ErrInvalidArgument},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := runtime.ToDuration(t.Context(), test.input)
			if !errors.Is(err, test.target) {
				t.Fatalf("ToDuration(%s) error = %v, want %v", test.input, err, test.target)
			}
		})
	}
}

func TestDurationContextualCoercion(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	fiveSeconds := runtime.NewDuration(5 * time.Second)

	assertDuration := func(name string, value runtime.Value, err error, expected runtime.Duration) {
		t.Helper()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		actual, ok := value.(runtime.Duration)
		if !ok || actual != expected {
			t.Fatalf("%s = %v, want %s", name, value, expected)
		}
	}

	value, err := runtime.AddChecked(ctx, fiveSeconds, runtime.NewString("500ms"))
	assertDuration("string addition", value, err, runtime.NewDuration(5500*time.Millisecond))
	value, err = runtime.AddChecked(ctx, fiveSeconds, runtime.NewInt(2))
	assertDuration("numeric addition", value, err, runtime.NewDuration(5002*time.Millisecond))
	value, err = runtime.AddChecked(ctx, runtime.NewInt(2), fiveSeconds)
	assertDuration("reverse numeric addition", value, err, runtime.NewDuration(5002*time.Millisecond))
	value, err = runtime.SubtractChecked(ctx, fiveSeconds, runtime.NewString("500ms"))
	assertDuration("string subtraction", value, err, runtime.NewDuration(4500*time.Millisecond))
	value, err = runtime.DivideChecked(ctx, fiveSeconds, runtime.NewString("2"))
	assertDuration("numeric string division", value, err, runtime.NewDuration(2500*time.Millisecond))

	ratio, err := runtime.DivideChecked(ctx, runtime.NewDuration(10*time.Second), runtime.NewString("5s"))
	if err != nil || ratio != runtime.NewInt(2) {
		t.Fatalf("duration string ratio = %v, %v", ratio, err)
	}
}

func TestDurationCheckedComparison(t *testing.T) {
	t.Parallel()

	tests := []struct {
		left, right runtime.Value
		expected    int
	}{
		{runtime.NewDuration(5 * time.Second), runtime.NewString("5s"), 0},
		{runtime.NewDuration(5 * time.Second), runtime.NewInt(5000), 0},
		{runtime.NewInt(4999), runtime.NewDuration(5 * time.Second), -1},
	}

	for _, test := range tests {
		actual, err := runtime.CompareChecked(t.Context(), test.left, test.right)
		if err != nil || actual != test.expected {
			t.Fatalf("CompareChecked(%s, %s) = %d, %v", test.left, test.right, actual, err)
		}
	}

	if _, err := runtime.CompareChecked(t.Context(), runtime.NewDuration(time.Second), runtime.NewString("tomorrow")); !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("invalid comparison error = %v", err)
	}
	if runtime.CompareValues(runtime.NewDuration(time.Millisecond), runtime.NewInt(1)) == 0 {
		t.Fatal("strict CompareValues unexpectedly applied temporal coercion")
	}

	dateTime := runtime.NewDateTime(time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC))
	actual, err := runtime.CompareChecked(t.Context(), dateTime, runtime.NewDuration(time.Second))
	if err != nil || actual != runtime.CompareValues(dateTime, runtime.NewDuration(time.Second)) {
		t.Fatalf("mixed DateTime comparison = %d, %v", actual, err)
	}
}
