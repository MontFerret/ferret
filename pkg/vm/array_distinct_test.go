package vm

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type distinctCollisionValue struct {
	err   error
	label string
}

func (v distinctCollisionValue) String() string {
	return v.label
}

func (v distinctCollisionValue) Hash() uint64 {
	return 7
}

func (v distinctCollisionValue) Copy() runtime.Value {
	return v
}

func (v distinctCollisionValue) Equal(_ context.Context, other runtime.Value) (bool, error) {
	if v.err != nil {
		return false, v.err
	}

	o, ok := other.(distinctCollisionValue)
	if !ok {
		return false, nil
	}

	return v.label == o.label, nil
}

func (v distinctCollisionValue) Compare(_ context.Context, other runtime.Value) (runtime.Ordering, error) {
	if v.err != nil {
		return runtime.Equal, v.err
	}

	o, ok := other.(distinctCollisionValue)
	if !ok {
		return runtime.Equal, runtime.Error(runtime.ErrInvalidOperation, "incompatible values")
	}
	if v.label < o.label {
		return runtime.Less, nil
	}
	if v.label > o.label {
		return runtime.Greater, nil
	}

	return runtime.Equal, nil
}

func TestArrayDistinctSeparatesHashCollisions(t *testing.T) {
	ctx := context.Background()
	first := distinctCollisionValue{label: "first"}
	second := distinctCollisionValue{label: "second"}

	result, err := arrayDistinct(ctx, runtime.NewArrayWith(first, second, first))
	if err != nil {
		t.Fatalf("arrayDistinct: %v", err)
	}

	length, err := result.Length(ctx)
	if err != nil {
		t.Fatalf("result length: %v", err)
	}

	if length != 2 {
		t.Fatalf("expected 2 distinct values, got %d", length)
	}

	for idx, want := range []runtime.Value{first, second} {
		got, err := result.At(ctx, runtime.Int(idx))
		if err != nil {
			t.Fatalf("result at %d: %v", idx, err)
		}

		equal, err := runtime.EqualValues(ctx, got, want)
		if err != nil {
			t.Fatalf("compare result at %d: %v", idx, err)
		}
		if !equal {
			t.Fatalf("result at %d: expected %v, got %v", idx, want, got)
		}
	}
}

func TestArrayDistinctUsesNumericEqualityAcrossRepresentations(t *testing.T) {
	result, err := arrayDistinct(
		t.Context(),
		runtime.NewArrayWith(runtime.NewInt(1), runtime.NewFloat(1)),
	)
	if err != nil {
		t.Fatalf("arrayDistinct: %v", err)
	}

	length, err := result.Length(t.Context())
	if err != nil {
		t.Fatalf("result length: %v", err)
	}
	if length != 1 {
		t.Fatalf("expected one distinct numeric value, got %d", length)
	}
}

func TestArrayDistinctPropagatesEqualityErrors(t *testing.T) {
	sentinel := errors.New("distinct equality failed")
	first := distinctCollisionValue{label: "first", err: sentinel}
	second := distinctCollisionValue{label: "second"}

	_, err := arrayDistinct(context.Background(), runtime.NewArrayWith(first, second))
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected equality error, got %v", err)
	}
}

func TestContainsPropagatesCollectionEqualityErrors(t *testing.T) {
	sentinel := errors.New("membership equality failed")
	first := distinctCollisionValue{label: "first"}
	second := distinctCollisionValue{label: "second", err: sentinel}

	result, err := contains(context.Background(), runtime.NewArrayWith(first), second)
	if result {
		t.Fatal("failed membership lookup must be false")
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected equality error, got %v", err)
	}
}

func TestVMInPropagatesCollectionEqualityErrors(t *testing.T) {
	sentinel := errors.New("vm membership equality failed")
	first := distinctCollisionValue{label: "first"}
	second := distinctCollisionValue{label: "second", err: sentinel}
	program := newTestProgram(
		3,
		[]runtime.Value{second, runtime.NewArrayWith(first)},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
		bytecode.NewInstruction(bytecode.OpIn, bytecode.NewRegister(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
	)

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected equality error, got %v", err)
	}
}

func TestArrayMembershipComparatorsPropagateEqualityErrors(t *testing.T) {
	sentinel := errors.New("quantified membership equality failed")
	first := distinctCollisionValue{label: "first"}
	second := distinctCollisionValue{label: "second", err: sentinel}
	left := runtime.NewArrayWith(second)
	right := runtime.NewArrayWith(first)

	for _, tc := range []struct {
		run  func(context.Context, arrayComparator, runtime.Value, runtime.Value) (runtime.Boolean, error)
		name string
	}{
		{name: "all", run: arrayAll},
		{name: "any", run: arrayAny},
		{name: "none", run: arrayNone},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.run(context.Background(), IN, left, right)
			if !errors.Is(err, sentinel) {
				t.Fatalf("expected equality error, got %v", err)
			}
		})
	}
}

func TestVMComparisonInstructionsPropagateHostErrors(t *testing.T) {
	sentinel := errors.New("vm comparison failed")
	left := distinctCollisionValue{label: "left", err: sentinel}
	right := distinctCollisionValue{label: "right"}

	for _, opcode := range []bytecode.Opcode{bytecode.OpEq, bytecode.OpGt} {
		t.Run(opcode.String(), func(t *testing.T) {
			program := newTestProgram(
				3,
				[]runtime.Value{left, right},
				bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
				bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
				bytecode.NewInstruction(opcode, bytecode.NewRegister(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			)

			_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
			if !errors.Is(err, sentinel) {
				t.Fatalf("expected host comparison error, got %v", err)
			}
		})
	}
}

func TestVMFusedEqualityJumpPropagatesHostErrors(t *testing.T) {
	sentinel := errors.New("vm fused comparison failed")
	left := distinctCollisionValue{label: "left", err: sentinel}
	right := distinctCollisionValue{label: "right"}
	program := newTestProgram(
		1,
		[]runtime.Value{left, right},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpJumpIfEqConst, bytecode.Operand(3), bytecode.NewRegister(0), bytecode.NewConstant(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected host comparison error, got %v", err)
	}
}
