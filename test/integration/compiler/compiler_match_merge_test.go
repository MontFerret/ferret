package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestMatchMerge_PureLiteralResults(t *testing.T) {
	src := `
LET x = @x
RETURN MATCH x (
  1 => "same",
  2 => "same",
  _ => "other",
)
`
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
	prog, err := c.Compile(file.NewSource("match_merge_pure", src))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	if got := countLoadConstValue(prog, runtime.NewString("same")); got != 1 {
		t.Fatalf("expected 1 load of \"same\", got %d", got)
	}
}

func TestMatchMerge_NoMerge_ImpureResult(t *testing.T) {
	src := `
LET x = @x
RETURN MATCH x (
  1 => LENGTH([1,2]),
  2 => LENGTH([1,2]),
  _ => 0,
)
`
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
	prog, err := c.Compile(file.NewSource("match_merge_impure", src))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	if got := countOpcode(prog, bytecode.OpLength); got != 2 {
		t.Fatalf("expected 2 LENGTH ops, got %d", got)
	}
}

func TestMatchMerge_NoMerge_GuardedArm(t *testing.T) {
	src := `
LET x = @x
RETURN MATCH x (
  1 => "same",
  2 WHEN 1 < 2 => "same",
  _ => "other",
)
`
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
	prog, err := c.Compile(file.NewSource("match_merge_guarded", src))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	if got := countLoadConstValue(prog, runtime.NewString("same")); got != 2 {
		t.Fatalf("expected 2 loads of \"same\", got %d", got)
	}
}

func countLoadConstValue(prog *bytecode.Program, val runtime.Value) int {
	if prog == nil {
		return 0
	}

	count := 0
	for _, inst := range prog.Bytecode {
		if inst.Opcode != bytecode.OpLoadConst {
			continue
		}
		constOp := inst.Operands[1]
		if !constOp.IsConstant() {
			continue
		}
		idx := constOp.Constant()
		if idx < 0 || idx >= len(prog.Constants) {
			continue
		}
		if runtime.CompareValues(prog.Constants[idx], val) == 0 {
			count++
		}
	}

	return count
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
