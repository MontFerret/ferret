package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/integration/base"
)

func compileOptimized(t *testing.T, expr string) *vm.Program {
	t.Helper()
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O1))
	prog, err := c.Compile(file.NewSource("const-prop", expr))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}
	return prog
}

func execOptimized(t *testing.T, prog *vm.Program) any {
	t.Helper()
	out, err := base.Exec(prog, false, vm.WithFunctions(base.Stdlib()))
	if err != nil {
		t.Fatalf("exec failed: %v", err)
	}
	return out
}

func assertHasOpcode(t *testing.T, prog *vm.Program, op vm.Opcode) {
	t.Helper()
	for _, inst := range prog.Bytecode {
		if inst.Opcode == op {
			return
		}
	}
	t.Fatalf("expected opcode %s to be present", op.String())
}

func assertNoOpcode(t *testing.T, prog *vm.Program, op vm.Opcode) {
	t.Helper()
	for _, inst := range prog.Bytecode {
		if inst.Opcode == op {
			t.Fatalf("unexpected opcode %s in bytecode", op.String())
		}
	}
}

func TestConstantPropagation_FoldsArithmetic(t *testing.T) {
	prog := compileOptimized(t, `LET a = 1 + 2 RETURN a`)
	assertNoOpcode(t, prog, vm.OpAdd)
	out := execOptimized(t, prog)
	if out != float64(3) && out != 3 {
		t.Fatalf("expected 3, got %v", out)
	}
}

func TestConstantPropagation_FoldsUnary(t *testing.T) {
	prog := compileOptimized(t, `LET a = 1 + 2 RETURN -a`)
	assertNoOpcode(t, prog, vm.OpAdd)
	assertNoOpcode(t, prog, vm.OpFlipNegative)
	out := execOptimized(t, prog)
	if out != float64(-3) && out != -3 {
		t.Fatalf("expected -3, got %v", out)
	}
}

func TestConstantPropagation_FoldsChain(t *testing.T) {
	prog := compileOptimized(t, `LET a = 10 RETURN (a - 3) * 2`)
	assertNoOpcode(t, prog, vm.OpSub)
	assertNoOpcode(t, prog, vm.OpMulti)
	out := execOptimized(t, prog)
	if out != float64(14) && out != 14 {
		t.Fatalf("expected 14, got %v", out)
	}
}

func TestConstantPropagation_DivideByZeroNotFolded(t *testing.T) {
	prog := compileOptimized(t, `RETURN 1 / 0`)
	assertHasOpcode(t, prog, vm.OpDiv)
	_, err := base.Exec(prog, false, vm.WithFunctions(base.Stdlib()))
	if err == nil {
		t.Fatalf("expected divide by zero error, got nil")
	}
}

func TestConstantPropagation_DivideByZeroStringNotFolded(t *testing.T) {
	prog := compileOptimized(t, `RETURN 1 / "0"`)
	assertHasOpcode(t, prog, vm.OpDiv)
	_, err := base.Exec(prog, false, vm.WithFunctions(base.Stdlib()))
	if err == nil {
		t.Fatalf("expected divide by zero error, got nil")
	}
}
