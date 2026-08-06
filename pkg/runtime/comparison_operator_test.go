package runtime_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestEvaluateComparisonPredicates(t *testing.T) {
	tests := []struct {
		left     runtime.Value
		right    runtime.Value
		name     string
		op       operator.Binary
		expected runtime.Boolean
	}{
		{name: "equal", op: operator.Equal, left: runtime.NewInt(2), right: runtime.NewFloat(2), expected: runtime.True},
		{name: "not equal", op: operator.NotEqual, left: runtime.NewInt(2), right: runtime.NewInt(3), expected: runtime.True},
		{name: "less", op: operator.Less, left: runtime.NewInt(1), right: runtime.NewInt(2), expected: runtime.True},
		{name: "less or equal", op: operator.LessOrEqual, left: runtime.NewInt(2), right: runtime.NewInt(2), expected: runtime.True},
		{name: "greater", op: operator.Greater, left: runtime.NewInt(3), right: runtime.NewInt(2), expected: runtime.True},
		{name: "greater or equal", op: operator.GreaterOrEqual, left: runtime.NewInt(2), right: runtime.NewInt(2), expected: runtime.True},
		{name: "in list", op: operator.In, left: runtime.NewInt(2), right: runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)), expected: runtime.True},
		{name: "in map values", op: operator.In, left: runtime.NewInt(2), right: runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewInt(2)}), expected: runtime.True},
		{name: "in string", op: operator.In, left: runtime.NewString("ell"), right: runtime.NewString("hello"), expected: runtime.True},
		{name: "string membership stringifies needle", op: operator.In, left: runtime.NewInt(1), right: runtime.NewString("value 1"), expected: runtime.True},
		{name: "unsupported membership", op: operator.In, left: runtime.NewInt(1), right: runtime.NewInt(1), expected: runtime.False},
		{name: "unsupported host membership", op: operator.In, left: runtime.NewInt(1), right: &contractHostValue{}, expected: runtime.False},
		{name: "incompatible equality", op: operator.Equal, left: runtime.NewDuration(time.Second), right: runtime.NewString("1s"), expected: runtime.False},
		{name: "incompatible inequality", op: operator.NotEqual, left: runtime.NewDuration(time.Second), right: runtime.NewString("1s"), expected: runtime.True},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := runtime.EvaluateComparison(t.Context(), test.op, test.left, test.right)
			if err != nil {
				t.Fatalf("EvaluateComparison() error = %v", err)
			}
			if actual != test.expected {
				t.Fatalf("EvaluateComparison() = %v, want %v", actual, test.expected)
			}

			predicate, err := runtime.ResolveComparison(test.op)
			if err != nil {
				t.Fatalf("ResolveComparison() error = %v", err)
			}
			resolved, err := predicate(t.Context(), test.left, test.right)
			if err != nil {
				t.Fatalf("resolved comparison error = %v", err)
			}
			if resolved != test.expected {
				t.Fatalf("resolved comparison = %v, want %v", resolved, test.expected)
			}
		})
	}
}

func TestEvaluateComparisonRelationalDiagnostics(t *testing.T) {
	duration := runtime.NewDuration(time.Second)
	stringValue := runtime.NewString("1s")
	tests := []struct {
		left     runtime.Value
		right    runtime.Value
		expected string
		op       operator.Binary
	}{
		{op: operator.Greater, left: stringValue, right: duration, expected: "invalid operation: operator '>' cannot be applied to String and Duration"},
		{op: operator.LessOrEqual, left: duration, right: stringValue, expected: "invalid operation: operator '<=' cannot be applied to Duration and String"},
	}

	for _, test := range tests {
		_, err := runtime.EvaluateComparison(t.Context(), test.op, test.left, test.right)
		if !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Fatalf("EvaluateComparison() error = %v, want ErrInvalidOperation", err)
		}
		if err.Error() != test.expected {
			t.Fatalf("EvaluateComparison() error = %q, want %q", err, test.expected)
		}
	}
}

func TestEvaluateComparisonPropagatesHostErrorsAndContext(t *testing.T) {
	hostErr := errors.New("host comparison failed")
	left := &contractHostValue{equalityErr: hostErr, comparisonErr: hostErr}
	right := &contractHostValue{}

	for _, op := range []operator.Binary{operator.Equal, operator.Greater} {
		if _, err := runtime.EvaluateComparison(t.Context(), op, left, right); !errors.Is(err, hostErr) {
			t.Fatalf("EvaluateComparison(%s) error = %v, want host error", op, err)
		}
	}

	container := runtime.NewArrayWith(&contractHostValue{})
	if _, err := runtime.EvaluateComparison(t.Context(), operator.In, &contractHostValue{equalityErr: hostErr}, container); !errors.Is(err, hostErr) {
		t.Fatalf("EvaluateComparison(IN) error = %v, want host error", err)
	}

	calls := 0
	cancelled := &contractHostValue{equalityCalls: &calls}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := runtime.EvaluateComparison(ctx, operator.Equal, cancelled, right); !errors.Is(err, context.Canceled) {
		t.Fatalf("EvaluateComparison() error = %v, want context.Canceled", err)
	}
	if calls != 1 {
		t.Fatalf("host equality calls = %d, want 1", calls)
	}
}

func TestEvaluateComparisonRejectsNonComparisonOperator(t *testing.T) {
	_, err := runtime.EvaluateComparison(t.Context(), operator.Add, runtime.NewInt(1), runtime.NewInt(2))
	if !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("EvaluateComparison() error = %v, want ErrInvalidOperation", err)
	}
	if actual, expected := err.Error(), `invalid operation: operator "+" is not a comparison operator`; actual != expected {
		t.Fatalf("EvaluateComparison() error = %q, want %q", actual, expected)
	}

	if predicate, resolvedErr := runtime.ResolveComparison(operator.Add); predicate != nil || resolvedErr == nil || resolvedErr.Error() != err.Error() {
		t.Fatalf("ResolveComparison() = %T, %v; want nil, %q", predicate, resolvedErr, err)
	}
}
