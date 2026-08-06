package runtime_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/internal/operationerror"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestNativeNumericArithmetic(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	tests := []struct {
		operation func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error)
		left      runtime.Value
		right     runtime.Value
		expected  runtime.Value
		name      string
	}{
		{name: "Int plus Int", operation: runtime.Add, left: runtime.NewInt(2), right: runtime.NewInt(3), expected: runtime.NewInt(5)},
		{name: "Int plus Float", operation: runtime.Add, left: runtime.NewInt(2), right: runtime.NewFloat(0.5), expected: runtime.NewFloat(2.5)},
		{name: "Float plus Int", operation: runtime.Add, left: runtime.NewFloat(2.5), right: runtime.NewInt(3), expected: runtime.NewFloat(5.5)},
		{name: "Float plus Float", operation: runtime.Add, left: runtime.NewFloat(2.5), right: runtime.NewFloat(0.5), expected: runtime.NewFloat(3)},
		{name: "Int minus Int", operation: runtime.Subtract, left: runtime.NewInt(5), right: runtime.NewInt(2), expected: runtime.NewInt(3)},
		{name: "Int minus Float", operation: runtime.Subtract, left: runtime.NewInt(5), right: runtime.NewFloat(2.5), expected: runtime.NewFloat(2.5)},
		{name: "Float minus Int", operation: runtime.Subtract, left: runtime.NewFloat(5.5), right: runtime.NewInt(2), expected: runtime.NewFloat(3.5)},
		{name: "Float minus Float", operation: runtime.Subtract, left: runtime.NewFloat(5.5), right: runtime.NewFloat(2.5), expected: runtime.NewFloat(3)},
		{name: "Int times Int", operation: runtime.Multiply, left: runtime.NewInt(3), right: runtime.NewInt(2), expected: runtime.NewInt(6)},
		{name: "Int times Float", operation: runtime.Multiply, left: runtime.NewInt(3), right: runtime.NewFloat(2.5), expected: runtime.NewFloat(7.5)},
		{name: "Float times Int", operation: runtime.Multiply, left: runtime.NewFloat(2.5), right: runtime.NewInt(3), expected: runtime.NewFloat(7.5)},
		{name: "Float times Float", operation: runtime.Multiply, left: runtime.NewFloat(2.5), right: runtime.NewFloat(2), expected: runtime.NewFloat(5)},
		{name: "exact Int division", operation: runtime.Divide, left: runtime.NewInt(6), right: runtime.NewInt(2), expected: runtime.NewInt(3)},
		{name: "fractional Int division", operation: runtime.Divide, left: runtime.NewInt(5), right: runtime.NewInt(2), expected: runtime.NewFloat(2.5)},
		{name: "Int divided by Float", operation: runtime.Divide, left: runtime.NewInt(5), right: runtime.NewFloat(2), expected: runtime.NewFloat(2.5)},
		{name: "Float divided by Int", operation: runtime.Divide, left: runtime.NewFloat(5), right: runtime.NewInt(2), expected: runtime.NewFloat(2.5)},
		{name: "Float divided by Float", operation: runtime.Divide, left: runtime.NewFloat(5), right: runtime.NewFloat(2), expected: runtime.NewFloat(2.5)},
		{name: "Int modulo Int", operation: runtime.Modulo, left: runtime.NewInt(5), right: runtime.NewInt(2), expected: runtime.NewInt(1)},
		{name: "Int modulo Float", operation: runtime.Modulo, left: runtime.NewInt(5), right: runtime.NewFloat(2), expected: runtime.NewFloat(1)},
		{name: "Float modulo Int", operation: runtime.Modulo, left: runtime.NewFloat(5.5), right: runtime.NewInt(2), expected: runtime.NewFloat(1.5)},
		{name: "Float modulo Float", operation: runtime.Modulo, left: runtime.NewFloat(5.5), right: runtime.NewFloat(2), expected: runtime.NewFloat(1.5)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.operation(ctx, test.left, test.right)
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

func TestStringTriggeredConcatenation(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	prefix := runtime.NewString("value=")
	values := []runtime.Value{
		runtime.NewInt(2),
		runtime.NewFloat(2.5),
		runtime.True,
		runtime.None,
		runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
		runtime.NewObject(),
		runtime.NewBinary([]byte("bytes")),
		runtime.NewDateTime(time.Unix(0, 0).UTC()),
		runtime.NewDuration(time.Second),
	}

	for _, value := range values {
		value := value
		t.Run(runtime.TypeName(runtime.TypeOf(value)), func(t *testing.T) {
			actual, err := runtime.Add(ctx, prefix, value)
			expected := runtime.NewString(prefix.String() + value.String())
			if err != nil || actual != expected {
				t.Fatalf("String + %s = %v, %v; want %q", runtime.TypeOf(value), actual, err, expected)
			}

			actual, err = runtime.Add(ctx, value, prefix)
			expected = runtime.NewString(value.String() + prefix.String())
			if err != nil || actual != expected {
				t.Fatalf("%s + String = %v, %v; want %q", runtime.TypeOf(value), actual, err, expected)
			}
		})
	}
}

func TestArithmeticRejectsImplicitNumericCoercion(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	tests := []struct {
		operation func() (runtime.Value, error)
		expected  string
		name      string
	}{
		{name: "Boolean plus Int", operation: func() (runtime.Value, error) { return runtime.Add(ctx, runtime.True, runtime.NewInt(1)) }, expected: "invalid operation: operator '+' cannot be applied to Boolean and Int"},
		{name: "Int plus Boolean", operation: func() (runtime.Value, error) { return runtime.Add(ctx, runtime.NewInt(1), runtime.True) }, expected: "invalid operation: operator '+' cannot be applied to Int and Boolean"},
		{name: "numeric String subtraction", operation: func() (runtime.Value, error) {
			return runtime.Subtract(ctx, runtime.NewString("10"), runtime.NewInt(2))
		}, expected: "invalid operation: operator '-' cannot be applied to String and Int"},
		{name: "numeric String multiplication", operation: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, runtime.NewString("10"), runtime.NewInt(2))
		}, expected: "invalid operation: operator '*' cannot be applied to String and Int"},
		{name: "numeric String division", operation: func() (runtime.Value, error) { return runtime.Divide(ctx, runtime.NewString("10"), runtime.NewInt(2)) }, expected: "invalid operation: operator '/' cannot be applied to String and Int"},
		{name: "numeric String modulo", operation: func() (runtime.Value, error) { return runtime.Modulo(ctx, runtime.NewString("10"), runtime.NewInt(2)) }, expected: "invalid operation: operator '%' cannot be applied to String and Int"},
		{name: "Array subtraction", operation: func() (runtime.Value, error) {
			return runtime.Subtract(ctx, runtime.NewArrayWith(runtime.NewInt(1)), runtime.NewInt(1))
		}, expected: "invalid operation: operator '-' cannot be applied to Array and Int"},
		{name: "None multiplication", operation: func() (runtime.Value, error) { return runtime.Multiply(ctx, runtime.None, runtime.NewInt(2)) }, expected: "invalid operation: operator '*' cannot be applied to None and Int"},
		{name: "Binary division", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewBinary([]byte("10")), runtime.NewInt(2))
		}, expected: "invalid operation: operator '/' cannot be applied to Binary and Int"},
		{name: "Boolean modulo", operation: func() (runtime.Value, error) { return runtime.Modulo(ctx, runtime.True, runtime.NewInt(2)) }, expected: "invalid operation: operator '%' cannot be applied to Boolean and Int"},
		{name: "String increment", operation: func() (runtime.Value, error) { return runtime.Increment(ctx, runtime.NewString("10")) }, expected: "invalid operation: operator '++' cannot be applied to String"},
		{name: "Boolean decrement", operation: func() (runtime.Value, error) { return runtime.Decrement(ctx, runtime.True) }, expected: "invalid operation: operator '--' cannot be applied to Boolean"},
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

func TestUnsupportedNativeValuesNeverParticipateInNumericArithmetic(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	values := []struct {
		value    runtime.Value
		typeName string
		name     string
	}{
		{name: "numeric String", value: runtime.NewString("10"), typeName: "String"},
		{name: "Boolean", value: runtime.True, typeName: "Boolean"},
		{name: "Array", value: runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)), typeName: "Array"},
		{name: "None", value: runtime.None, typeName: "None"},
		{name: "Binary", value: runtime.NewBinary([]byte("10")), typeName: "Binary"},
	}
	binaryOperations := []struct {
		operation func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error)
		symbol    string
		name      string
	}{
		{name: "subtract", symbol: "-", operation: runtime.Subtract},
		{name: "multiply", symbol: "*", operation: runtime.Multiply},
		{name: "divide", symbol: "/", operation: runtime.Divide},
		{name: "modulo", symbol: "%", operation: runtime.Modulo},
	}
	unaryOperations := []struct {
		operation func(context.Context, runtime.Value) (runtime.Value, error)
		symbol    string
		name      string
	}{
		{name: "increment", symbol: "++", operation: runtime.Increment},
		{name: "decrement", symbol: "--", operation: runtime.Decrement},
	}

	for _, value := range values {
		for _, operation := range binaryOperations {
			t.Run(value.name+"/"+operation.name, func(t *testing.T) {
				_, err := operation.operation(ctx, value.value, runtime.NewInt(2))
				expected := fmt.Sprintf("invalid operation: operator '%s' cannot be applied to %s and Int", operation.symbol, value.typeName)
				if !errors.Is(err, runtime.ErrInvalidOperation) || err.Error() != expected {
					t.Fatalf("error = %v, want %q", err, expected)
				}
			})
		}

		for _, operation := range unaryOperations {
			t.Run(value.name+"/"+operation.name, func(t *testing.T) {
				_, err := operation.operation(ctx, value.value)
				expected := fmt.Sprintf("invalid operation: operator '%s' cannot be applied to %s", operation.symbol, value.typeName)
				if !errors.Is(err, runtime.ErrInvalidOperation) || err.Error() != expected {
					t.Fatalf("error = %v, want %q", err, expected)
				}
			})
		}
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
		{name: "float modulo by zero", operation: func() (runtime.Value, error) {
			return runtime.Modulo(ctx, runtime.NewFloat(1), runtime.ZeroFloat)
		}, target: runtime.ErrInvalidOperation},
		{name: "non-finite float operand", operation: func() (runtime.Value, error) {
			return runtime.Subtract(ctx, runtime.NewFloat(math.Inf(1)), runtime.NewFloat(1))
		}, target: runtime.ErrRange},
		{name: "non-finite float result", operation: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, runtime.NewFloat(math.MaxFloat64), runtime.NewFloat(2))
		}, target: runtime.ErrRange},
		{name: "invalid numeric string", operation: func() (runtime.Value, error) {
			return runtime.Subtract(ctx, runtime.NewString("not-a-number"), runtime.NewInt(1))
		}, target: runtime.ErrInvalidOperation},
		{name: "invalid collection item", operation: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, runtime.NewArrayWith(runtime.NewInt(1), runtime.NewString("not-a-number")), runtime.NewInt(2))
		}, target: runtime.ErrInvalidOperation},
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

func TestZeroDivisorErrorsPreserveRuntimeAndOperationIdentity(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	tests := []struct {
		operation func() (runtime.Value, error)
		identity  error
		message   string
		name      string
	}{
		{name: "Int division", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewInt(1), runtime.ZeroInt)
		}, identity: operationerror.ErrDivisionByZero, message: "invalid operation: division by zero"},
		{name: "Float division", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewFloat(1), runtime.ZeroFloat)
		}, identity: operationerror.ErrDivisionByZero, message: "invalid operation: division by zero"},
		{name: "Duration scalar division", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewDuration(time.Second), runtime.ZeroInt)
		}, identity: operationerror.ErrDivisionByZero, message: "invalid operation: division by zero"},
		{name: "Duration ratio division", operation: func() (runtime.Value, error) {
			return runtime.Divide(ctx, runtime.NewDuration(time.Second), runtime.ZeroDuration)
		}, identity: operationerror.ErrDivisionByZero, message: "invalid operation: division by zero"},
		{name: "Int modulo", operation: func() (runtime.Value, error) {
			return runtime.Modulo(ctx, runtime.NewInt(1), runtime.ZeroInt)
		}, identity: operationerror.ErrModuloByZero, message: "invalid operation: modulo by zero"},
		{name: "Float modulo", operation: func() (runtime.Value, error) {
			return runtime.Modulo(ctx, runtime.NewFloat(1), runtime.ZeroFloat)
		}, identity: operationerror.ErrModuloByZero, message: "invalid operation: modulo by zero"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.operation()
			if err == nil {
				t.Fatal("operation returned nil error")
			}
			if got := err.Error(); got != test.message {
				t.Fatalf("error = %q, want %q", got, test.message)
			}
			if !errors.Is(err, runtime.ErrInvalidOperation) {
				t.Fatalf("error = %v, want runtime.ErrInvalidOperation", err)
			}
			if !errors.Is(err, test.identity) {
				t.Fatalf("error = %v, want %v", err, test.identity)
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
