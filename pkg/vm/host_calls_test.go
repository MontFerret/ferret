package vm

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func TestCallCachedHostFunction_VarargSixArgsPreservesOrderAndCount(t *testing.T) {
	cases := []int{1, 4, 6, 8, 9}

	for _, argCount := range cases {
		t.Run(runtime.NewInt(argCount).String(), func(t *testing.T) {
			var seen []runtime.Value

			cacheFn := &mem.CachedHostFunction{
				Bound: true,
				FnV: func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
					seen = append([]runtime.Value(nil), args...)
					return runtime.NewInt(len(args)), nil
				},
			}

			reg := make([]runtime.Value, argCount+1)
			reg[0] = runtime.None
			for i := 0; i < argCount; i++ {
				reg[i+1] = runtime.NewInt(i + 1)
			}

			scratch := mem.NewScratch(0)
			out, err := callCachedHostFunction(
				context.Background(),
				cacheFn,
				reg,
				scratch,
				runtime.NewString("TEST"),
				bytecode.NewRegister(1),
				bytecode.NewRegister(argCount),
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got, want := out, runtime.NewInt(argCount); got != want {
				t.Fatalf("unexpected return value: got %v, want %v", got, want)
			}

			if got, want := len(seen), argCount; got != want {
				t.Fatalf("unexpected arg count: got %d, want %d", got, want)
			}

			for i := 0; i < argCount; i++ {
				if got, want := seen[i], runtime.NewInt(i+1); got != want {
					t.Fatalf("unexpected arg[%d]: got %v, want %v", i, got, want)
				}
			}
		})
	}
}

func TestCallCachedHostFunction_VarargArgsSliceMutationDoesNotMutateRegisters(t *testing.T) {
	cases := []int{1, 4, 6, 8, 9}

	for _, argCount := range cases {
		t.Run(runtime.NewInt(argCount).String(), func(t *testing.T) {
			cacheFn := &mem.CachedHostFunction{
				Bound: true,
				FnV: func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
					args[0] = runtime.NewInt(777)
					args[len(args)-1] = runtime.NewInt(999)
					return runtime.True, nil
				},
			}

			reg := make([]runtime.Value, argCount+1)
			reg[0] = runtime.None
			for i := 0; i < argCount; i++ {
				reg[i+1] = runtime.NewInt(i + 1)
			}

			scratch := mem.NewScratch(0)
			_, err := callCachedHostFunction(
				context.Background(),
				cacheFn,
				reg,
				scratch,
				runtime.NewString("TEST"),
				bytecode.NewRegister(1),
				bytecode.NewRegister(argCount),
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for i := 0; i < argCount; i++ {
				if got, want := reg[i+1], runtime.NewInt(i+1); got != want {
					t.Fatalf("register arg[%d] mutated: got %v, want %v", i, got, want)
				}
			}
		})
	}
}
