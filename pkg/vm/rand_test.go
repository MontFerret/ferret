package vm

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/rnd"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestRandOpcodeUsesContextSource(t *testing.T) {
	program := newTestProgram(1, nil,
		bytecode.NewInstruction(bytecode.OpRand, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpRand, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)
	instance, err := New(program)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = instance.Close() })
	env := NewDefaultEnvironment()
	if _, err := instance.Run(t.Context(), env); !errors.Is(err, runtime.ErrUnexpected) {
		t.Fatalf("missing source = %v, want ErrUnexpected", err)
	}

	src, control := rnd.NewSeed(42), rnd.NewSeed(42)
	ctx := rnd.WithContext(t.Context(), src)
	for range 3 {
		control.Float64()
		want := runtime.Float(control.Float64())
		result, err := instance.Run(ctx, env)
		if err != nil {
			t.Fatal(err)
		}

		if result.Root() != want {
			t.Fatalf("opcode result = %v, want %v", result.Root(), want)
		}

		if err := result.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
