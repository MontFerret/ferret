package vm

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func TestCallCachedHostFunction_VarargSixArgsPreservesOrderAndCount(t *testing.T) {
	var seen []runtime.Value

	cacheFn := &mem.CachedHostFunction{
		FnV: func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
			seen = append([]runtime.Value(nil), args...)
			return runtime.NewInt(len(args)), nil
		},
	}

	reg := []runtime.Value{
		runtime.None,
		runtime.NewInt(10),
		runtime.NewInt(20),
		runtime.NewInt(30),
		runtime.NewInt(40),
		runtime.NewInt(50),
		runtime.NewInt(60),
	}

	out, err := callCachedHostFunction(
		context.Background(),
		cacheFn,
		reg,
		runtime.NewString("TEST"),
		bytecode.NewRegister(1),
		bytecode.NewRegister(6),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := out, runtime.NewInt(6); got != want {
		t.Fatalf("unexpected return value: got %v, want %v", got, want)
	}

	if got, want := len(seen), 6; got != want {
		t.Fatalf("unexpected arg count: got %d, want %d", got, want)
	}

	wantVals := []runtime.Value{
		runtime.NewInt(10),
		runtime.NewInt(20),
		runtime.NewInt(30),
		runtime.NewInt(40),
		runtime.NewInt(50),
		runtime.NewInt(60),
	}
	for i := range wantVals {
		if got, want := seen[i], wantVals[i]; got != want {
			t.Fatalf("unexpected arg[%d]: got %v, want %v", i, got, want)
		}
	}
}

func TestCallCachedHostFunction_VarargArgsSliceMutationDoesNotMutateRegisters(t *testing.T) {
	cacheFn := &mem.CachedHostFunction{
		FnV: func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
			args[0] = runtime.NewInt(777)
			args[len(args)-1] = runtime.NewInt(999)
			return runtime.True, nil
		},
	}

	reg := []runtime.Value{
		runtime.None,
		runtime.NewInt(1),
		runtime.NewInt(2),
		runtime.NewInt(3),
		runtime.NewInt(4),
		runtime.NewInt(5),
		runtime.NewInt(6),
	}

	_, err := callCachedHostFunction(
		context.Background(),
		cacheFn,
		reg,
		runtime.NewString("TEST"),
		bytecode.NewRegister(1),
		bytecode.NewRegister(6),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantVals := []runtime.Value{
		runtime.NewInt(1),
		runtime.NewInt(2),
		runtime.NewInt(3),
		runtime.NewInt(4),
		runtime.NewInt(5),
		runtime.NewInt(6),
	}
	for i := range wantVals {
		if got, want := reg[i+1], wantVals[i]; got != want {
			t.Fatalf("register arg[%d] mutated: got %v, want %v", i, got, want)
		}
	}
}
