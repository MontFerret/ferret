package runtime_test

import (
	"context"
	"errors"
	"math"
	"strings"
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

	for _, test := range []struct {
		left     runtime.Value
		right    runtime.Value
		name     string
		expected runtime.Duration
	}{
		{
			name:     "integer string multiplier",
			left:     fiveSeconds,
			right:    runtime.NewString("2"),
			expected: runtime.NewDuration(10 * time.Second),
		},
		{
			name:     "reverse integer string multiplier",
			left:     runtime.NewString("2"),
			right:    fiveSeconds,
			expected: runtime.NewDuration(10 * time.Second),
		},
		{
			name:     "fractional string multiplier",
			left:     fiveSeconds,
			right:    runtime.NewString("2.5"),
			expected: runtime.NewDuration(12500 * time.Millisecond),
		},
		{
			name:     "reverse fractional string multiplier",
			left:     runtime.NewString("2.5"),
			right:    fiveSeconds,
			expected: runtime.NewDuration(12500 * time.Millisecond),
		},
	} {
		value, err := runtime.MultiplyChecked(ctx, test.left, test.right)
		assertDuration(test.name, value, err, test.expected)
	}

	for _, test := range []struct {
		left, right runtime.Value
		name        string
	}{
		{name: "invalid string multiplier", left: fiveSeconds, right: runtime.NewString("invalid")},
		{name: "reverse invalid string multiplier", left: runtime.NewString("invalid"), right: fiveSeconds},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := runtime.MultiplyChecked(ctx, test.left, test.right)
			if !errors.Is(err, runtime.ErrInvalidArgument) {
				t.Fatalf("error = %v, want ErrInvalidArgument", err)
			}

			for _, part := range []string{"cannot use String", "invalid", "numeric Duration scale"} {
				if !strings.Contains(err.Error(), part) {
					t.Fatalf("error = %q, want %q", err, part)
				}
			}
		})
	}

	value, err = runtime.DivideChecked(ctx, fiveSeconds, runtime.NewString("2"))
	assertDuration("numeric string division", value, err, runtime.NewDuration(2500*time.Millisecond))

	ratio, err := runtime.DivideChecked(ctx, runtime.NewDuration(10*time.Second), runtime.NewString("5s"))
	if err != nil || ratio != runtime.NewInt(2) {
		t.Fatalf("duration string ratio = %v, %v", ratio, err)
	}
}

func TestDurationComparisonRequiresExplicitConversion(t *testing.T) {
	t.Parallel()

	duration := runtime.NewDuration(5 * time.Second)
	converted, err := runtime.ToDuration(t.Context(), runtime.NewString("5s"))
	if err != nil {
		t.Fatal(err)
	}

	equal, err := runtime.EqualValues(t.Context(), duration, converted)
	if err != nil || !equal {
		t.Fatalf("explicitly converted equality = %v, %v, want true, nil", equal, err)
	}
	actual, err := runtime.CompareValues(t.Context(), duration, converted)
	if err != nil || actual != runtime.Equal {
		t.Fatalf("explicitly converted ordering = %d, %v, want Equal, nil", actual, err)
	}

	actual, err = runtime.CompareValues(t.Context(), runtime.NewDuration(4999*time.Millisecond), duration)
	if err != nil || actual != runtime.Less {
		t.Fatalf("native Duration ordering = %d, %v, want Less, nil", actual, err)
	}
}

func TestDurationComparisonDoesNotConvertOtherTypes(t *testing.T) {
	t.Parallel()

	duration := runtime.NewDuration(time.Second)
	tests := []struct {
		input runtime.Value
		name  string
	}{
		{name: "valid string", input: runtime.NewString("1s")},
		{name: "matching number", input: runtime.NewInt(1000)},
		{name: "boolean", input: runtime.True},
		{name: "invalid type", input: runtime.NewObject()},
		{name: "malformed string", input: runtime.NewString("tomorrow")},
		{name: "malformed list shape", input: runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))},
		{name: "overflow", input: runtime.NewInt64(9_223_372_036_855)},
		{name: "recursive conversion", input: runtime.NewArrayWith(runtime.NewArrayWith(runtime.NewString("tomorrow")))},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			equal, err := runtime.EqualValues(t.Context(), duration, test.input)
			if err != nil {
				t.Fatalf("EqualValues(%s, %s): %v", duration, test.input, err)
			}
			if equal {
				t.Fatalf("EqualValues(%s, %s) = true, want false", duration, test.input)
			}
			if _, err := runtime.CompareValues(t.Context(), duration, test.input); !errors.Is(err, runtime.ErrInvalidOperation) {
				t.Fatalf("CompareValues(%s, %s) error = %v, want ErrInvalidOperation", duration, test.input, err)
			}
		})
	}
}

func TestDurationComparisonDoesNotInspectOpaqueLists(t *testing.T) {
	t.Parallel()

	duration := runtime.NewDuration(time.Second)
	customErr := errors.New("duration list access failed")
	rangeErr := runtime.Error(runtime.ErrRange, "duration list backend failed")
	tests := []struct {
		input runtime.Value
		name  string
	}{
		{
			name:  "length cancellation",
			input: newFallibleDurationList(context.Canceled, nil),
		},
		{
			name:  "length deadline",
			input: newFallibleDurationList(context.DeadlineExceeded, nil),
		},
		{
			name:  "length range error remains operational",
			input: newFallibleDurationList(rangeErr, nil),
		},
		{
			name:  "item access error",
			input: newFallibleDurationList(nil, customErr, runtime.NewInt(1)),
		},
		{
			name:  "item range error remains operational",
			input: newFallibleDurationList(nil, rangeErr, runtime.NewInt(1)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			equal, err := runtime.EqualValues(t.Context(), duration, test.input)
			if err != nil || equal {
				t.Fatalf("EqualValues = %v, %v; want false, nil", equal, err)
			}

			_, err = runtime.CompareValues(t.Context(), duration, test.input)
			if !errors.Is(err, runtime.ErrInvalidOperation) {
				t.Fatalf("CompareValues error = %v, want ErrInvalidOperation", err)
			}
		})
	}
}
