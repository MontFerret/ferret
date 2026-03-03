package vm

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestHandleErrorWithCatch_AppliesJumpTargetZero(t *testing.T) {
	instance := New(&bytecode.Program{
		Registers: 1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
		CatchTable: []bytecode.Catch{
			{1, 1, 0},
		},
	})

	instance.pc = 1
	called := false

	err := instance.handleErrorWithCatch(errors.New("boom"), func() {
		called = true
	})

	if err != nil {
		t.Fatalf("expected caught error to be swallowed, got %v", err)
	}

	if !called {
		t.Fatal("expected onCatch callback to be called")
	}

	if got, want := instance.pc, 0; got != want {
		t.Fatalf("expected catch jump target %d, got %d", want, got)
	}
}

func TestHandleErrorWithCatch_AppliesPositiveJumpTarget(t *testing.T) {
	instance := New(&bytecode.Program{
		Registers: 1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
		CatchTable: []bytecode.Catch{
			{1, 1, 2},
		},
	})

	instance.pc = 1

	err := instance.handleErrorWithCatch(errors.New("boom"), nil)
	if err != nil {
		t.Fatalf("expected caught error to be swallowed, got %v", err)
	}

	if got, want := instance.pc, 2; got != want {
		t.Fatalf("expected catch jump target %d, got %d", want, got)
	}
}

func TestHandleErrorWithCatch_ReturnsErrorOutsideCatchRegion(t *testing.T) {
	instance := New(&bytecode.Program{
		Registers: 1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
		CatchTable: []bytecode.Catch{
			{0, 0, 1},
		},
	})

	instance.pc = 1
	called := false
	wantErr := errors.New("boom")

	err := instance.handleErrorWithCatch(wantErr, func() {
		called = true
	})

	if err != wantErr {
		t.Fatalf("expected original error to be returned, got %v", err)
	}

	if called {
		t.Fatal("expected onCatch callback not to be called")
	}

	if got, want := instance.pc, 1; got != want {
		t.Fatalf("expected pc to stay unchanged at %d, got %d", want, got)
	}
}
