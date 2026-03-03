package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
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

func TestUdfUnusedEliminationO1(t *testing.T) {
	expr := `
FUNC used() => 1
FUNC unused() => 2
RETURN used()
`

	prog := compileWithLevel(t, compiler.O1, expr)
	if len(prog.Functions.UserDefined) != 1 {
		t.Fatalf("expected 1 UDF at O1, got %d", len(prog.Functions.UserDefined))
	}

	if prog.Functions.UserDefined[0].Name != "USED" {
		t.Fatalf("expected USED UDF, got %q", prog.Functions.UserDefined[0].Name)
	}
}

func TestUdfUnusedEliminationO0(t *testing.T) {
	expr := `
FUNC used() => 1
FUNC unused() => 2
RETURN used()
`

	prog := compileWithLevel(t, compiler.O0, expr)
	if len(prog.Functions.UserDefined) != 2 {
		t.Fatalf("expected 2 UDFs at O0, got %d", len(prog.Functions.UserDefined))
	}
}

func TestUdfUnusedNestedCaptureNotLifted(t *testing.T) {
	expr := `
LET base = 5
FUNC outer() (
  FUNC inner() ( RETURN base )
  RETURN 1
)
RETURN outer()
`

	prog := compileWithLevel(t, compiler.O1, expr)
	if len(prog.Functions.UserDefined) != 1 {
		t.Fatalf("expected 1 UDF at O1, got %d", len(prog.Functions.UserDefined))
	}

	var outerParams = -1
	for _, udf := range prog.Functions.UserDefined {
		if udf.Name == "OUTER" {
			outerParams = udf.Params
			break
		}
	}

	if outerParams == -1 {
		t.Fatalf("expected OUTER UDF metadata")
	}

	if outerParams != 0 {
		t.Fatalf("expected OUTER to have 0 params (no captures), got %d", outerParams)
	}
}

func TestUdfRecursionReachable(t *testing.T) {
	expr := `
FUNC fact(n) (
  RETURN MATCH n (
    0 => 1,
    _ => n * fact(n - 1),
  )
)
RETURN fact(5)
`

	prog := compileWithLevel(t, compiler.O1, expr)
	found := false
	for _, udf := range prog.Functions.UserDefined {
		if udf.Name == "FACT" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected FACT UDF to be reachable at O1")
	}
}
