package runtime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestInvalidBinaryOperationsUseOperatorDiagnostics(t *testing.T) {
	dateTime := runtime.NewDateTime(time.Date(2026, time.August, 5, 12, 0, 0, 0, time.UTC))
	duration := runtime.NewDuration(time.Second)

	tests := []struct {
		operation func() (runtime.Value, error)
		name      string
		expected  string
	}{
		{
			name:     "addition",
			expected: "invalid operation: operator '+' cannot be applied to DateTime and DateTime",
			operation: func() (runtime.Value, error) {
				return runtime.Add(t.Context(), dateTime, dateTime)
			},
		},
		{
			name:     "subtraction preserves operand order",
			expected: "invalid operation: operator '-' cannot be applied to Duration and DateTime",
			operation: func() (runtime.Value, error) {
				return runtime.Subtract(t.Context(), duration, dateTime)
			},
		},
		{
			name:     "multiplication",
			expected: "invalid operation: operator '*' cannot be applied to Duration and Duration",
			operation: func() (runtime.Value, error) {
				return runtime.Multiply(t.Context(), duration, duration)
			},
		},
		{
			name:     "division preserves operand order",
			expected: "invalid operation: operator '/' cannot be applied to Int and Duration",
			operation: func() (runtime.Value, error) {
				return runtime.Divide(t.Context(), runtime.NewInt(1), duration)
			},
		},
		{
			name:     "modulus",
			expected: "invalid operation: operator '%' cannot be applied to Duration and Int",
			operation: func() (runtime.Value, error) {
				return runtime.Modulo(t.Context(), duration, runtime.NewInt(1))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.operation()
			if err == nil || err.Error() != test.expected {
				t.Fatalf("error = %v, want %q", err, test.expected)
			}
		})
	}
}

func TestComparisonWithoutSourceOperatorUsesComparisonDiagnostic(t *testing.T) {
	_, err := runtime.CompareValues(
		t.Context(),
		runtime.NewDuration(time.Second),
		runtime.NewString("1s"),
	)
	const expected = "invalid operation: comparison cannot be applied to Duration and String"
	if err == nil || err.Error() != expected {
		t.Fatalf("error = %v, want %q", err, expected)
	}
}
