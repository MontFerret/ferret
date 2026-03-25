package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func compileWithLevel(t *testing.T, level compiler.OptimizationLevel, expr string) *bytecode.Program {
	t.Helper()

	c := compiler.New(compiler.WithOptimizationLevel(level))
	prog, err := c.Compile(file.NewAnonymousSource(expr))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	return prog
}

func countOpcode(prog *bytecode.Program, op bytecode.Opcode) int {
	if prog == nil {
		return 0
	}

	count := 0
	for _, inst := range prog.Bytecode {
		if inst.Opcode == op {
			count++
		}
	}

	return count
}

func findFirstOpcodeIndex(code []bytecode.Instruction, op bytecode.Opcode) (int, bool) {
	for i, inst := range code {
		if inst.Opcode == op {
			return i, true
		}
	}

	return -1, false
}

func hasOpcode(code []bytecode.Instruction, op bytecode.Opcode) bool {
	_, ok := findFirstOpcodeIndex(code, op)
	return ok
}

func programHasOpcode(prog *bytecode.Program, op bytecode.Opcode) bool {
	if prog == nil {
		return false
	}

	return hasOpcode(prog.Bytecode, op)
}

func lastRegisterDefOpcodeBefore(code []bytecode.Instruction, before int, reg int) (bytecode.Opcode, bool) {
	for i := before - 1; i >= 0; i-- {
		inst := code[i]
		if !inst.Operands[0].IsRegister() {
			continue
		}
		if inst.Operands[0].Register() == reg {
			return inst.Opcode, true
		}
	}

	return bytecode.OpMove, false
}

func firstCompilationError(err error) *diagnostics.Diagnostic {
	switch e := err.(type) {
	case *diagnostics.Diagnostic:
		return e
	case *diagnostics.DiagnosticSet:
		if e.Size() == 0 {
			return nil
		}

		return e.First()
	default:
		return nil
	}
}
