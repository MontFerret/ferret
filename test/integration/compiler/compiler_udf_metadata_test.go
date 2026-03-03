package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func TestUdfOnlyHostCallIsPresentInHostMetadata(t *testing.T) {
	expr := `
FUNC f() => TEST_FN(1)
RETURN f()
`

	prog := compileWithLevel(t, compiler.O0, expr)
	if len(prog.Functions.Host) != 1 {
		t.Fatalf("expected exactly 1 host function, got %d", len(prog.Functions.Host))
	}

	if got := prog.Functions.Host["TEST_FN"]; got != 1 {
		t.Fatalf("expected TEST_FN arity 1, got %d", got)
	}
}

func TestUdfHostCallUpdatesMaxHostArityAcrossScopes(t *testing.T) {
	expr := `
FUNC f() => TEST_FN(1, 2)
RETURN TEST_FN(1)
`

	prog := compileWithLevel(t, compiler.O0, expr)
	if got := prog.Functions.Host["TEST_FN"]; got != 2 {
		t.Fatalf("expected TEST_FN max arity 2, got %d", got)
	}
}

func TestUdfOnlyParamIsPresentInProgramParams(t *testing.T) {
	expr := `
FUNC f() => @foo
RETURN f()
`

	prog := compileWithLevel(t, compiler.O0, expr)
	if len(prog.Params) != 1 {
		t.Fatalf("expected exactly 1 param, got %d (%v)", len(prog.Params), prog.Params)
	}

	if prog.Params[0] != "foo" {
		t.Fatalf("expected param foo, got %q", prog.Params[0])
	}
}
