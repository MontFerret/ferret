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

func findUserDefined(t *testing.T, prog *bytecode.Program, name string) bytecode.UDF {
	t.Helper()

	for _, udf := range prog.Functions.UserDefined {
		if udf.Name == name {
			return udf
		}
	}

	t.Fatalf("expected UDF %q in %v", name, prog.Functions.UserDefined)

	return bytecode.UDF{}
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

func TestUdfNestedCaptureMetadataAcrossScopes(t *testing.T) {
	expr := `
LET global = 100
FUNC outer(a) (
  LET outerLocal = 10
  FUNC middle(b) (
    FUNC inner(c) => global + a + outerLocal + b + c
    RETURN inner(1)
  )
  RETURN middle(2)
)
RETURN outer(3)
`

	assertLevel := func(level compiler.OptimizationLevel) {
		t.Helper()

		prog := compileWithLevel(t, level, expr)

		if got := findUserDefined(t, prog, "outer").Params; got != 2 {
			t.Fatalf("expected outer total params/captures to be 2, got %d", got)
		}

		if got := findUserDefined(t, prog, "middle").Params; got != 4 {
			t.Fatalf("expected middle total params/captures to be 4, got %d", got)
		}

		if got := findUserDefined(t, prog, "inner").Params; got != 5 {
			t.Fatalf("expected inner total params/captures to be 5, got %d", got)
		}
	}

	assertLevel(compiler.O0)
	assertLevel(compiler.O1)
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
		if udf.Name == "f" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected UDF f to remain reachable at O1, got %v", prog.Functions.UserDefined)
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

func TestUdfParamSlotMatchesProgramParamOrdering(t *testing.T) {
	expr := `
LET x = @alpha
FUNC f() => @beta
RETURN x + f()
`

	prog := compileWithLevel(t, compiler.O0, expr)

	wantParams := []string{"alpha", "beta"}
	if len(prog.Params) != len(wantParams) {
		t.Fatalf("unexpected params count: got %d (%v), want %d", len(prog.Params), prog.Params, len(wantParams))
	}

	for i := range wantParams {
		if prog.Params[i] != wantParams[i] {
			t.Fatalf("unexpected param at index %d: got %q, want %q", i, prog.Params[i], wantParams[i])
		}
	}

	udfEntry := -1
	for _, udf := range prog.Functions.UserDefined {
		if udf.Name == "f" {
			udfEntry = udf.Entry
			break
		}
	}

	if udfEntry < 0 {
		t.Fatalf("expected UDF f metadata, got %v", prog.Functions.UserDefined)
	}

	if udfEntry >= len(prog.Bytecode) {
		t.Fatalf("invalid UDF entry: %d (bytecode len: %d)", udfEntry, len(prog.Bytecode))
	}

	inst := prog.Bytecode[udfEntry]
	if inst.Opcode != bytecode.OpLoadParam {
		t.Fatalf("unexpected opcode at UDF entry %d: got %s, want %s", udfEntry, inst.Opcode, bytecode.OpLoadParam)
	}

	if got := inst.Operands[1]; got != bytecode.Operand(2) {
		t.Fatalf("unexpected UDF LOADP slot: got %d, want %d", got, bytecode.Operand(2))
	}

	if inst.Operands[1].IsConstant() {
		t.Fatalf("expected UDF LOADP source operand to be slot-encoded, got constant %s", inst.Operands[1])
	}
}

func TestCaseDistinctUdfNamesPreserveMetadata(t *testing.T) {
	expr := `
FUNC a() => 1
FUNC A() => 2
RETURN a() + A()
`

	prog := compileWithLevel(t, compiler.O0, expr)

	if got := findUserDefined(t, prog, "a").DisplayName; got != "a" {
		t.Fatalf("expected lowercase UDF display name, got %q", got)
	}

	if got := findUserDefined(t, prog, "A").DisplayName; got != "A" {
		t.Fatalf("expected uppercase UDF display name, got %q", got)
	}
}
