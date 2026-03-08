package vm

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestOpLoadParam_UsesBoundSlots(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  4,
		Params:     []string{"foo", "bar"},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(1), bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(2), bytecode.Operand(2)),
			bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
		},
	}

	env := NewDefaultEnvironment()
	env.Params["foo"] = runtime.NewInt(1)
	env.Params["bar"] = runtime.NewInt(2)

	out, err := New(program).Run(context.Background(), env)
	if err != nil {
		t.Fatalf("expected successful run, got %v", err)
	}

	if out != runtime.NewInt(3) {
		t.Fatalf("unexpected result: got %v, want %v", out, runtime.NewInt(3))
	}
}

func TestOpLoadParam_MissingParamsPreserveRuntimeError(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Params:     []string{"foo", "bar"},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(1), bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	env := NewDefaultEnvironment()
	env.Params["foo"] = runtime.NewInt(1)

	_, err := New(program).Run(context.Background(), env)
	if err == nil {
		t.Fatal("expected missing parameter error")
	}

	if !strings.Contains(err.Error(), "Missing parameter") {
		t.Fatalf("unexpected error: %v", err)
	}

	cause := errors.Unwrap(err)
	if cause == nil {
		cause = err
	}

	if !strings.Contains(cause.Error(), "@bar") {
		t.Fatalf("expected missing parameter name in error, got %v (cause: %v)", err, cause)
	}
}
