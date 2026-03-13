package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type hostCallSpec struct {
	name      string
	args      []runtime.Value
	protected bool
}

func newTestProgram(registers int, constants []runtime.Value, instructions ...bytecode.Instruction) *bytecode.Program {
	return &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  registers,
		Bytecode:   instructions,
		Constants:  constants,
	}
}

func newHostCallProgram(specs ...hostCallSpec) *bytecode.Program {
	constants := make([]runtime.Value, 0, len(specs)*2)
	instructions := make([]bytecode.Instruction, 0, len(specs)*3+1)
	nextRegister := 0
	returnReg := 0

	for _, spec := range specs {
		dst := nextRegister
		nextRegister++

		nameIdx := len(constants)
		constants = append(constants, runtime.NewString(spec.name))
		instructions = append(instructions,
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(dst), bytecode.NewConstant(nameIdx)),
		)

		if len(spec.args) == 0 {
			op := bytecode.OpHCall
			if spec.protected {
				op = bytecode.OpProtectedHCall
			}

			instructions = append(instructions, bytecode.NewInstruction(op, bytecode.NewRegister(dst)))
			returnReg = dst
			continue
		}

		argStart := nextRegister
		argEnd := argStart + len(spec.args) - 1
		for _, arg := range spec.args {
			constIdx := len(constants)
			constants = append(constants, arg)
			instructions = append(instructions,
				bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(nextRegister), bytecode.NewConstant(constIdx)),
			)
			nextRegister++
		}

		op := bytecode.OpHCall
		if spec.protected {
			op = bytecode.OpProtectedHCall
		}

		instructions = append(instructions,
			bytecode.NewInstruction(op, bytecode.NewRegister(dst), bytecode.NewRegister(argStart), bytecode.NewRegister(argEnd)),
		)
		returnReg = dst
	}

	if nextRegister == 0 {
		nextRegister = 1
	}

	instructions = append(instructions, bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(returnReg)))

	return newTestProgram(nextRegister, constants, instructions...)
}

func assertRuntimeValueEquals(t *testing.T, got, want runtime.Value) {
	t.Helper()

	if got != want {
		t.Fatalf("unexpected result: got %v, want %v", got, want)
	}
}
