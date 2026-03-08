package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
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

func TestUdfNamespaceAliasHostMetadata(t *testing.T) {
	expr := `
USE FOO AS F
FUNC f() => F::TEST_FN()
RETURN f()
`

	prog := compileWithLevel(t, compiler.O0, expr)
	if got := prog.Functions.Host["FOO::TEST_FN"]; got != 0 {
		t.Fatalf("expected FOO::TEST_FN arity 0, got %d", got)
	}

	if _, ok := prog.Functions.Host["FOO"]; ok {
		t.Fatalf("expected no bare FOO host metadata, got %v", prog.Functions.Host)
	}
}

func TestNamespaceAliasDoesNotShadowUdfCallAtO1(t *testing.T) {
	expr := `
USE FOO AS F
FUNC f() => 1
RETURN f()
`

	prog := compileWithLevel(t, compiler.O1, expr)
	found := false
	for _, udf := range prog.Functions.UserDefined {
		if udf.Name == "F" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected UDF F to remain reachable at O1, got %v", prog.Functions.UserDefined)
	}

	if _, ok := prog.Functions.Host["FOO"]; ok {
		t.Fatalf("expected no bare FOO host metadata at O1, got %v", prog.Functions.Host)
	}
}

func TestProgramParamsPreserveFirstUseOrder(t *testing.T) {
	expr := `RETURN [@beta, @alpha, @beta, @gamma]`

	prog := compileWithLevel(t, compiler.O0, expr)

	if len(prog.Params) != 3 {
		t.Fatalf("unexpected params count: got %d (%v), want %d", len(prog.Params), prog.Params, 3)
	}

	want := []string{"beta", "alpha", "gamma"}
	for i := range want {
		if prog.Params[i] != want[i] {
			t.Fatalf("unexpected param order at index %d: got %q, want %q", i, prog.Params[i], want[i])
		}
	}
}

func TestLoadParamUsesSlotOperand(t *testing.T) {
	expr := `RETURN @foo + @bar + @foo`

	prog := compileWithLevel(t, compiler.O0, expr)

	var loads []bytecode.Instruction
	for _, inst := range prog.Bytecode {
		if inst.Opcode == bytecode.OpLoadParam {
			loads = append(loads, inst)
		}
	}

	if len(loads) != 3 {
		t.Fatalf("unexpected number of LOADP instructions: got %d", len(loads))
	}

	got := []bytecode.Operand{
		loads[0].Operands[1],
		loads[1].Operands[1],
		loads[2].Operands[1],
	}
	want := []bytecode.Operand{1, 2, 1}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected slot at LOADP #%d: got %d, want %d", i, got[i], want[i])
		}

		if got[i].IsConstant() {
			t.Fatalf("expected LOADP source operand to be slot-encoded, got constant %s", got[i])
		}
	}
}
