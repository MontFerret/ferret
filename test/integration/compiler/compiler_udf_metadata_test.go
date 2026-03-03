package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func assertParamSet(t *testing.T, got []string, expected ...string) {
	t.Helper()

	if len(got) != len(expected) {
		t.Fatalf("unexpected params length: got %d (%v), want %d (%v)", len(got), got, len(expected), expected)
	}

	gotSet := make(map[string]struct{}, len(got))
	for _, name := range got {
		gotSet[name] = struct{}{}
	}

	for _, name := range expected {
		if _, ok := gotSet[name]; !ok {
			t.Fatalf("expected param %q to be present, got %v", name, got)
		}
		delete(gotSet, name)
	}

	for name := range gotSet {
		t.Fatalf("unexpected extra param %q in %v", name, got)
	}
}

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

func TestUdfNestedHostCallIsPresentInHostMetadata(t *testing.T) {
	expr := `
FUNC outer() (
  FUNC inner(x) => TEST_FN(x)
  RETURN inner(1)
)
RETURN outer()
`

	prog := compileWithLevel(t, compiler.O0, expr)
	if len(prog.Functions.Host) != 1 {
		t.Fatalf("expected exactly 1 host function, got %d (%v)", len(prog.Functions.Host), prog.Functions.Host)
	}

	if got := prog.Functions.Host["TEST_FN"]; got != 1 {
		t.Fatalf("expected TEST_FN arity 1, got %d", got)
	}
}

func TestUdfOnlyParamIsPresentInProgramParams(t *testing.T) {
	expr := `
FUNC f() => @foo
RETURN f()
`

	prog := compileWithLevel(t, compiler.O0, expr)
	assertParamSet(t, prog.Params, "foo")
}

func TestUdfNestedParamIsPresentInProgramParams(t *testing.T) {
	expr := `
FUNC outer() (
  FUNC inner() => @foo
  RETURN inner()
)
RETURN outer()
`

	prog := compileWithLevel(t, compiler.O0, expr)
	assertParamSet(t, prog.Params, "foo")
}

func TestUdfHostArityMergesAcrossTopLevelAndMultipleUdfs(t *testing.T) {
	expr := `
FUNC a() => TEST_FN(1)
FUNC b() => TEST_FN(1, 2, 3)
LET top = TEST_FN(1, 2)
RETURN [a(), b(), top]
`

	prog := compileWithLevel(t, compiler.O0, expr)
	if len(prog.Functions.Host) != 1 {
		t.Fatalf("expected exactly 1 host function, got %d (%v)", len(prog.Functions.Host), prog.Functions.Host)
	}

	if got := prog.Functions.Host["TEST_FN"]; got != 3 {
		t.Fatalf("expected TEST_FN max arity 3, got %d", got)
	}
}

func TestUdfUnusedMetadataPresentAtO0(t *testing.T) {
	expr := `
FUNC used() => 1
FUNC unused() => TEST_FN(@foo)
RETURN used()
`

	prog := compileWithLevel(t, compiler.O0, expr)
	if got := prog.Functions.Host["TEST_FN"]; got != 1 {
		t.Fatalf("expected TEST_FN arity 1 at O0, got %d", got)
	}

	assertParamSet(t, prog.Params, "foo")
}

func TestUdfUnusedMetadataPrunedAtO1(t *testing.T) {
	expr := `
FUNC used() => 1
FUNC unused() => TEST_FN(@foo)
RETURN used()
`

	prog := compileWithLevel(t, compiler.O1, expr)
	if _, ok := prog.Functions.Host["TEST_FN"]; ok {
		t.Fatalf("expected TEST_FN metadata to be pruned at O1, got %v", prog.Functions.Host)
	}

	assertParamSet(t, prog.Params)
}

func TestUdfFunctionAliasHostMetadata(t *testing.T) {
	expr := `
USE FOO::TEST_FN AS FN
FUNC f() => FN(1)
RETURN FN()
`

	prog := compileWithLevel(t, compiler.O0, expr)
	if len(prog.Functions.Host) != 1 {
		t.Fatalf("expected exactly 1 host function, got %d (%v)", len(prog.Functions.Host), prog.Functions.Host)
	}

	if got := prog.Functions.Host["FOO::TEST_FN"]; got != 1 {
		t.Fatalf("expected FOO::TEST_FN arity 1, got %d", got)
	}
}

func TestUdfNamespaceAliasHostMetadata_TODO(t *testing.T) {
	t.Skip("TODO: namespace alias in UDF currently adds extra FOO key to host metadata")
	// Desired canonical behavior for future fix:
	// USE FOO AS F
	// FUNC f() => F::TEST_FN()
	// RETURN f()
	// Host metadata should contain only FOO::TEST_FN.
}
