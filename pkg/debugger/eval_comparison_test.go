package debugger

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type debugComparisonValue struct {
	err   error
	label string
}

func (v debugComparisonValue) String() string      { return v.label }
func (v debugComparisonValue) Hash() uint64        { return 17 }
func (v debugComparisonValue) Copy() runtime.Value { return v }

func (v debugComparisonValue) Equal(_ context.Context, other runtime.Value) (bool, error) {
	if v.err != nil {
		return false, v.err
	}

	o, ok := other.(debugComparisonValue)
	return ok && v.label == o.label, nil
}

func (v debugComparisonValue) Compare(_ context.Context, other runtime.Value) (runtime.Ordering, error) {
	if v.err != nil {
		return runtime.Equal, v.err
	}

	o, ok := other.(debugComparisonValue)
	if !ok {
		return runtime.Equal, runtime.Error(runtime.ErrInvalidOperation, "incompatible debugger values")
	}
	if v.label < o.label {
		return runtime.Less, nil
	}
	if v.label > o.label {
		return runtime.Greater, nil
	}

	return runtime.Equal, nil
}

func TestEvaluateDebugComparisonUsesFallibleHostCapabilities(t *testing.T) {
	scope := evalScope{
		locals: map[string]runtime.Value{
			"left":  debugComparisonValue{label: "a"},
			"right": debugComparisonValue{label: "b"},
		},
		params: runtime.NewParams(),
		values: vm.NewDebugValueAccess(),
	}

	value, err := evaluateExpression(context.Background(), "left < right", scope)
	if err != nil {
		t.Fatalf("order comparison: %v", err)
	}
	if value != runtime.True {
		t.Fatalf("expected true, got %v", value)
	}

	value, err = evaluateExpression(context.Background(), "left == right", scope)
	if err != nil {
		t.Fatalf("equality comparison: %v", err)
	}
	if value != runtime.False {
		t.Fatalf("expected false, got %v", value)
	}

	value, err = evaluateExpression(context.Background(), "left != right", scope)
	if err != nil {
		t.Fatalf("inequality comparison: %v", err)
	}
	if value != runtime.True {
		t.Fatalf("expected true, got %v", value)
	}

	value, err = evaluateExpression(context.Background(), "left == left", scope)
	if err != nil {
		t.Fatalf("reflexive equality comparison: %v", err)
	}
	if value != runtime.True {
		t.Fatalf("expected true, got %v", value)
	}

	sentinel := errors.New("comparison failed")
	scope.locals["left"] = debugComparisonValue{label: "a", err: sentinel}
	_, err = evaluateExpression(context.Background(), "left == right", scope)
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected host equality error, got %v", err)
	}

	_, err = evaluateExpression(context.Background(), "left < right", scope)
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected host comparison error, got %v", err)
	}
}

func TestEvaluateDebugEqualityAcceptsIncompatibleTypes(t *testing.T) {
	scope := evalScope{
		locals: map[string]runtime.Value{
			"duration": runtime.NewDuration(1),
			"text":     runtime.NewString("1ns"),
		},
		params: runtime.NewParams(),
		values: vm.NewDebugValueAccess(),
	}

	value, err := evaluateExpression(context.Background(), "duration == text", scope)
	if err != nil || value != runtime.False {
		t.Fatalf("incompatible equality = %v, %v; want false, nil", value, err)
	}

	value, err = evaluateExpression(context.Background(), "duration != text", scope)
	if err != nil || value != runtime.True {
		t.Fatalf("incompatible inequality = %v, %v; want true, nil", value, err)
	}
}
