package runtime_test

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestArithmeticPreservesNonTemporalCoercion(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	tests := []struct {
		operation func() (runtime.Value, error)
		expected  runtime.Value
		name      string
	}{
		{name: "integer addition", operation: func() (runtime.Value, error) {
			return runtime.Add(ctx, runtime.NewInt(2), runtime.NewInt(3))
		}, expected: runtime.NewInt(5)},
		{name: "mixed numeric addition", operation: func() (runtime.Value, error) {
			return runtime.Add(ctx, runtime.NewInt(2), runtime.NewFloat(0.5))
		}, expected: runtime.NewFloat(2.5)},
		{name: "string concatenation", operation: func() (runtime.Value, error) {
			return runtime.Add(ctx, runtime.NewString("value="), runtime.NewInt(2))
		}, expected: runtime.NewString("value=2")},
		{name: "boolean concatenation", operation: func() (runtime.Value, error) {
			return runtime.Add(ctx, runtime.True, runtime.NewInt(2))
		}, expected: runtime.NewString("true2")},
		{name: "numeric string subtraction", operation: func() (runtime.Value, error) {
			return runtime.Subtract(ctx, runtime.NewString("5"), runtime.NewInt(2))
		}, expected: runtime.NewInt(3)},
		{name: "collection subtraction", operation: func() (runtime.Value, error) {
			return runtime.Subtract(ctx, runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)), runtime.NewInt(1))
		}, expected: runtime.NewInt(2)},
		{name: "numeric string multiplication", operation: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, runtime.NewString("2.5"), runtime.NewInt(2))
		}, expected: runtime.NewFloat(5)},
		{name: "fractional division", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewString("5"), runtime.NewString("2"))
		}, expected: runtime.NewFloat(2.5)},
		{name: "integer division", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewInt(6), runtime.NewInt(2))
		}, expected: runtime.NewInt(3)},
		{name: "numeric string modulo", operation: func() (runtime.Value, error) {
			return runtime.Modulo(ctx, runtime.NewString("5"), runtime.NewString("2"))
		}, expected: runtime.NewInt(1)},
		{name: "boolean increment", operation: func() (runtime.Value, error) {
			return runtime.Increment(ctx, runtime.True)
		}, expected: runtime.NewFloat(2)},
		{name: "numeric string decrement", operation: func() (runtime.Value, error) {
			return runtime.Decrement(ctx, runtime.NewString("2"))
		}, expected: runtime.NewInt(1)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.operation()
			if err != nil {
				t.Fatal(err)
			}

			equal, err := runtime.EqualValues(ctx, actual, test.expected)
			if err != nil || !equal {
				t.Fatalf("result = %v (%s), want %v (%s): %v", actual, runtime.TypeOf(actual), test.expected, runtime.TypeOf(test.expected), err)
			}
		})
	}
}

func TestArithmeticReportsNumericFailures(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	tests := []struct {
		operation func() (runtime.Value, error)
		target    error
		name      string
	}{
		{name: "integer addition overflow", operation: func() (runtime.Value, error) {
			return runtime.Add(ctx, runtime.NewInt64(math.MaxInt64), runtime.NewInt(1))
		}, target: runtime.ErrRange},
		{name: "integer subtraction overflow", operation: func() (runtime.Value, error) {
			return runtime.Subtract(ctx, runtime.NewInt64(math.MinInt64), runtime.NewInt(1))
		}, target: runtime.ErrRange},
		{name: "integer multiplication overflow", operation: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, runtime.NewInt64(math.MaxInt64), runtime.NewInt(2))
		}, target: runtime.ErrRange},
		{name: "minimum integer division overflow", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewInt64(math.MinInt64), runtime.NewInt(-1))
		}, target: runtime.ErrRange},
		{name: "integer division by zero", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewInt(1), runtime.ZeroInt)
		}, target: runtime.ErrInvalidOperation},
		{name: "float division by zero", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewFloat(1), runtime.ZeroFloat)
		}, target: runtime.ErrInvalidOperation},
		{name: "modulo by zero", operation: func() (runtime.Value, error) {
			return runtime.Modulo(ctx, runtime.NewInt(1), runtime.ZeroInt)
		}, target: runtime.ErrInvalidOperation},
		{name: "non-finite float operand", operation: func() (runtime.Value, error) {
			return runtime.Subtract(ctx, runtime.NewFloat(math.Inf(1)), runtime.NewFloat(1))
		}, target: runtime.ErrRange},
		{name: "non-finite float result", operation: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, runtime.NewFloat(math.MaxFloat64), runtime.NewFloat(2))
		}, target: runtime.ErrRange},
		{name: "invalid numeric string", operation: func() (runtime.Value, error) {
			return runtime.Subtract(ctx, runtime.NewString("not-a-number"), runtime.NewInt(1))
		}},
		{name: "invalid collection item", operation: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, runtime.NewArrayWith(runtime.NewInt(1), runtime.NewString("not-a-number")), runtime.NewInt(2))
		}},
		{name: "increment overflow", operation: func() (runtime.Value, error) {
			return runtime.Increment(ctx, runtime.NewInt64(math.MaxInt64))
		}, target: runtime.ErrRange},
		{name: "decrement overflow", operation: func() (runtime.Value, error) {
			return runtime.Decrement(ctx, runtime.NewInt64(math.MinInt64))
		}, target: runtime.ErrRange},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.operation()
			if err == nil || (test.target != nil && !errors.Is(err, test.target)) {
				t.Fatalf("error = %v, want %v", err, test.target)
			}
		})
	}
}

func TestDurationArithmeticContract(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	tests := []struct {
		operation func() (runtime.Value, error)
		expected  runtime.Value
		name      string
	}{
		{name: "addition", operation: func() (runtime.Value, error) {
			return runtime.Add(ctx, runtime.NewDuration(2*time.Second), runtime.NewDuration(500*time.Millisecond))
		}, expected: runtime.NewDuration(2500 * time.Millisecond)},
		{name: "subtraction", operation: func() (runtime.Value, error) {
			return runtime.Subtract(ctx, runtime.NewDuration(2*time.Second), runtime.NewDuration(500*time.Millisecond))
		}, expected: runtime.NewDuration(1500 * time.Millisecond)},
		{name: "integer multiplication", operation: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, runtime.NewDuration(time.Second), runtime.NewInt(2))
		}, expected: runtime.NewDuration(2 * time.Second)},
		{name: "reverse integer multiplication", operation: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, runtime.NewInt(2), runtime.NewDuration(time.Second))
		}, expected: runtime.NewDuration(2 * time.Second)},
		{name: "float multiplication truncates", operation: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, runtime.NewDuration(3), runtime.NewFloat(0.5))
		}, expected: runtime.NewDuration(1)},
		{name: "integer division truncates", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewDuration(3), runtime.NewInt(2))
		}, expected: runtime.NewDuration(1)},
		{name: "exact Duration ratio", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewDuration(4*time.Second), runtime.NewDuration(2*time.Second))
		}, expected: runtime.NewInt(2)},
		{name: "fractional Duration ratio", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewDuration(5*time.Second), runtime.NewDuration(2*time.Second))
		}, expected: runtime.NewFloat(2.5)},
		{name: "negative Duration", operation: func() (runtime.Value, error) {
			return runtime.Negative(runtime.NewDuration(time.Second))
		}, expected: runtime.NewDuration(-time.Second)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.operation()
			if err != nil {
				t.Fatal(err)
			}

			equal, err := runtime.EqualValues(ctx, actual, test.expected)
			if err != nil || !equal {
				t.Fatalf("result = %v (%s), want %v (%s): %v", actual, runtime.TypeOf(actual), test.expected, runtime.TypeOf(test.expected), err)
			}
		})
	}

	failures := []struct {
		operation func() (runtime.Value, error)
		target    error
		name      string
	}{
		{name: "addition overflow", operation: func() (runtime.Value, error) {
			return runtime.Add(ctx, runtime.Duration(math.MaxInt64), runtime.NewDuration(1))
		}, target: runtime.ErrRange},
		{name: "subtraction overflow", operation: func() (runtime.Value, error) {
			return runtime.Subtract(ctx, runtime.Duration(math.MinInt64), runtime.NewDuration(1))
		}, target: runtime.ErrRange},
		{name: "multiplication overflow", operation: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, runtime.Duration(math.MaxInt64), runtime.NewInt(2))
		}, target: runtime.ErrRange},
		{name: "division overflow", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.Duration(math.MinInt64), runtime.NewInt(-1))
		}, target: runtime.ErrRange},
		{name: "division by zero", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewDuration(time.Second), runtime.ZeroInt)
		}, target: runtime.ErrInvalidOperation},
		{name: "ratio division by zero", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewDuration(time.Second), runtime.ZeroDuration)
		}, target: runtime.ErrInvalidOperation},
		{name: "non-finite multiplier", operation: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, runtime.NewDuration(time.Second), runtime.NewFloat(math.Inf(1)))
		}, target: runtime.ErrInvalidOperation},
		{name: "minimum Duration negation", operation: func() (runtime.Value, error) {
			return runtime.Negative(runtime.Duration(math.MinInt64))
		}, target: runtime.ErrRange},
	}

	for _, test := range failures {
		t.Run(test.name, func(t *testing.T) {
			if _, err := test.operation(); !errors.Is(err, test.target) {
				t.Fatalf("error = %v, want %v", err, test.target)
			}
		})
	}
}

func TestTemporalArithmeticRejectsIncompatiblePairs(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	dateTime := runtime.NewDateTime(time.Unix(0, 0).UTC())
	duration := runtime.NewDuration(time.Second)
	stringValue := runtime.NewString("1s")
	integer := runtime.NewInt(1)
	array := runtime.NewArrayWith(runtime.NewString("1s"))

	tests := []struct {
		operation func() (runtime.Value, error)
		expected  string
		name      string
	}{
		{name: "DateTime plus DateTime", operation: func() (runtime.Value, error) { return runtime.Add(ctx, dateTime, dateTime) }, expected: "invalid operation: operator '+' cannot be applied to DateTime and DateTime"},
		{name: "DateTime plus String", operation: func() (runtime.Value, error) { return runtime.Add(ctx, dateTime, stringValue) }, expected: "invalid operation: operator '+' cannot be applied to DateTime and String"},
		{name: "String plus DateTime", operation: func() (runtime.Value, error) { return runtime.Add(ctx, stringValue, dateTime) }, expected: "invalid operation: operator '+' cannot be applied to String and DateTime"},
		{name: "Duration plus Int", operation: func() (runtime.Value, error) { return runtime.Add(ctx, duration, integer) }, expected: "invalid operation: operator '+' cannot be applied to Duration and Int"},
		{name: "Int plus Duration", operation: func() (runtime.Value, error) { return runtime.Add(ctx, integer, duration) }, expected: "invalid operation: operator '+' cannot be applied to Int and Duration"},
		{name: "Duration minus DateTime", operation: func() (runtime.Value, error) { return runtime.Subtract(ctx, duration, dateTime) }, expected: "invalid operation: operator '-' cannot be applied to Duration and DateTime"},
		{name: "String minus Duration", operation: func() (runtime.Value, error) { return runtime.Subtract(ctx, stringValue, duration) }, expected: "invalid operation: operator '-' cannot be applied to String and Duration"},
		{name: "DateTime times Int", operation: func() (runtime.Value, error) { return runtime.Multiply(ctx, dateTime, integer) }, expected: "invalid operation: operator '*' cannot be applied to DateTime and Int"},
		{name: "Int times DateTime", operation: func() (runtime.Value, error) { return runtime.Multiply(ctx, integer, dateTime) }, expected: "invalid operation: operator '*' cannot be applied to Int and DateTime"},
		{name: "Duration times Duration", operation: func() (runtime.Value, error) { return runtime.Multiply(ctx, duration, duration) }, expected: "invalid operation: operator '*' cannot be applied to Duration and Duration"},
		{name: "Duration times String", operation: func() (runtime.Value, error) { return runtime.Multiply(ctx, duration, stringValue) }, expected: "invalid operation: operator '*' cannot be applied to Duration and String"},
		{name: "String times Duration", operation: func() (runtime.Value, error) { return runtime.Multiply(ctx, stringValue, duration) }, expected: "invalid operation: operator '*' cannot be applied to String and Duration"},
		{name: "DateTime divided by Int", operation: func() (runtime.Value, error) { return runtime.Divide(ctx, dateTime, integer) }, expected: "invalid operation: operator '/' cannot be applied to DateTime and Int"},
		{name: "Int divided by DateTime", operation: func() (runtime.Value, error) { return runtime.Divide(ctx, integer, dateTime) }, expected: "invalid operation: operator '/' cannot be applied to Int and DateTime"},
		{name: "Int divided by Duration", operation: func() (runtime.Value, error) { return runtime.Divide(ctx, integer, duration) }, expected: "invalid operation: operator '/' cannot be applied to Int and Duration"},
		{name: "Duration divided by String", operation: func() (runtime.Value, error) { return runtime.Divide(ctx, duration, stringValue) }, expected: "invalid operation: operator '/' cannot be applied to Duration and String"},
		{name: "Duration divided by Array", operation: func() (runtime.Value, error) { return runtime.Divide(ctx, duration, array) }, expected: "invalid operation: operator '/' cannot be applied to Duration and Array"},
		{name: "DateTime modulo Int", operation: func() (runtime.Value, error) { return runtime.Modulo(ctx, dateTime, integer) }, expected: "invalid operation: operator '%' cannot be applied to DateTime and Int"},
		{name: "Int modulo Duration", operation: func() (runtime.Value, error) { return runtime.Modulo(ctx, integer, duration) }, expected: "invalid operation: operator '%' cannot be applied to Int and Duration"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.operation()
			if !errors.Is(err, runtime.ErrInvalidOperation) || err.Error() != test.expected {
				t.Fatalf("error = %v, want %q", err, test.expected)
			}
		})
	}
}

func TestTemporalArithmeticDoesNotInspectIncompatibleOperands(t *testing.T) {
	t.Parallel()

	opaque := &opaqueHostValue{}
	if _, err := runtime.Add(t.Context(), runtime.NewDuration(time.Second), opaque); !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("Duration + opaque error = %v, want ErrInvalidOperation", err)
	}
	if _, err := runtime.Subtract(t.Context(), runtime.NewDateTime(time.Unix(0, 0)), opaque); !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("DateTime - opaque error = %v, want ErrInvalidOperation", err)
	}
}

func TestUnaryArithmeticContract(t *testing.T) {
	t.Parallel()

	if actual, err := runtime.Not(runtime.True); err != nil || actual != runtime.False {
		t.Fatalf("Not(true) = %v, %v", actual, err)
	}
	if actual, err := runtime.Positive(runtime.NewFloat(1.5)); err != nil || actual != runtime.NewFloat(1.5) {
		t.Fatalf("Positive(1.5) = %v, %v", actual, err)
	}
	if actual, err := runtime.Positive(runtime.NewDuration(time.Second)); err != nil || actual != runtime.NewDuration(time.Second) {
		t.Fatalf("Positive(Duration) = %v, %v", actual, err)
	}
	if actual, err := runtime.Negative(runtime.NewInt(2)); err != nil || actual != runtime.NewInt(-2) {
		t.Fatalf("Negative(2) = %v, %v", actual, err)
	}

	tests := []struct {
		operation func() (runtime.Value, error)
		expected  string
		name      string
	}{
		{name: "not String", operation: func() (runtime.Value, error) { return runtime.Not(runtime.NewString("")) }, expected: "invalid operation: operator '!' cannot be applied to String"},
		{name: "positive String", operation: func() (runtime.Value, error) { return runtime.Positive(runtime.NewString("1")) }, expected: "invalid operation: operator '+' cannot be applied to String"},
		{name: "negative Boolean", operation: func() (runtime.Value, error) { return runtime.Negative(runtime.True) }, expected: "invalid operation: operator '-' cannot be applied to Boolean"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.operation()
			if !errors.Is(err, runtime.ErrInvalidOperation) || err.Error() != test.expected {
				t.Fatalf("error = %v, want %q", err, test.expected)
			}
		})
	}

	if _, err := runtime.Negative(runtime.NewInt64(math.MinInt64)); !errors.Is(err, runtime.ErrRange) {
		t.Fatalf("minimum Int negation error = %v, want ErrRange", err)
	}
	if _, err := runtime.Negative(runtime.NewFloat(math.Inf(1))); !errors.Is(err, runtime.ErrRange) {
		t.Fatalf("infinite Float negation error = %v, want ErrRange", err)
	}
}
