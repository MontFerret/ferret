package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/integration/base"
)

func compileNoOpt(t *testing.T, expr string) *bytecode.Program {
	t.Helper()
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
	prog, err := c.Compile(file.NewSource("template-literal-opt", expr))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}
	return prog
}

func execProgram(t *testing.T, prog *bytecode.Program) any {
	t.Helper()
	out, err := base.Exec(prog, false, vm.WithFunctions(base.Stdlib()))
	if err != nil {
		t.Fatalf("exec failed: %v", err)
	}
	return out
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

func assertOpcodeCount(t *testing.T, prog *bytecode.Program, op bytecode.Opcode, want int) {
	t.Helper()
	if got := countOpcode(prog, op); got != want {
		t.Fatalf("expected %d %s opcode(s), got %d", want, op.String(), got)
	}
}

func TestTemplateLiteral_ConstFoldedToSingleConst(t *testing.T) {
	prog := compileNoOpt(t, "RETURN `foo-${1}-bar-${true}`")
	assertOpcodeCount(t, prog, bytecode.OpAdd, 0)
	out := execProgram(t, prog)
	if out != "foo-1-bar-true" {
		t.Fatalf("expected %q, got %v", "foo-1-bar-true", out)
	}
}

func TestTemplateLiteral_FoldsConstExpressionsIntoChunks(t *testing.T) {
	prog := compileNoOpt(t, "LET x = \"X\" RETURN `a-${1}-b-${x}-c-${true}-d`")
	assertOpcodeCount(t, prog, bytecode.OpAdd, 2)
	out := execProgram(t, prog)
	if out != "a-1-b-X-c-true-d" {
		t.Fatalf("expected %q, got %v", "a-1-b-X-c-true-d", out)
	}
}
