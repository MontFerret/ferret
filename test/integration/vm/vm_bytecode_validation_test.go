package vm_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
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

	RunProgramSpecs(t, []spec.Spec{
		spec.NewSpecWith(spec.NewProgramInput(program), "unsupported ISA version").
			Expect().ExecError(assert.NewUnaryAssertion(func(actual any) error {
			err, ok := actual.(error)
			if !ok || err == nil {
				return errors.New("expected version validation error")
			}

			if !strings.Contains(err.Error(), "unsupported bytecode version; recompile query") {
				return fmt.Errorf("unexpected error: %v", err)
			}

			return nil
		})),
	})
}

func TestFailsOnUnknownOpcode(t *testing.T) {
	unknown := bytecode.Opcode(255)
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode:   []bytecode.Instruction{bytecode.NewInstruction(unknown)},
		Registers:  1,
	}

	RunProgramSpecs(t, []spec.Spec{
		spec.NewSpecWith(spec.NewProgramInput(program), "unknown opcode").
			Expect().ExecError(assert.NewUnaryAssertion(func(actual any) error {
			err, ok := actual.(error)
			if !ok || err == nil {
				return errors.New("expected unknown opcode error")
			}

			var rtErr *vm.RuntimeError
			if !errors.As(err, &rtErr) {
				return fmt.Errorf("expected runtime error, got %T", err)
			}

			if rtErr.Note != "unknown opcode 255 at pc 0" {
				return fmt.Errorf("unexpected error: %v", err)
			}

			return nil
		})),
	})
}
