package vm_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestRejectsUnsupportedBytecodeVersion(t *testing.T) {
	unsupported := bytecode.Version - 1
	if unsupported < 0 {
		unsupported = 0
	}

	program := &bytecode.Program{
		ISAVersion: unsupported,
		Bytecode:   []bytecode.Instruction{bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0))},
		Registers:  1,
	}

	_, err := vm.New(program)
	if err == nil {
		t.Fatal("expected version validation error")
	}

	if !strings.Contains(err.Error(), "unsupported bytecode version; recompile query") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFailsOnUnknownOpcode(t *testing.T) {
	unknown := bytecode.Opcode(255)
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode:   []bytecode.Instruction{bytecode.NewInstruction(unknown)},
		Registers:  1,
	}

	instance, err := vm.New(program)
	if err != nil {
		t.Fatalf("unexpected constructor error: %v", err)
	}

	_, err = instance.Run(context.Background(), vm.NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected unknown opcode error")
	}

	if !strings.Contains(errors.Unwrap(err).Error(), "unknown opcode 255 at pc 0") {
		t.Fatalf("unexpected error: %v", err)
	}
}
