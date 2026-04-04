package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func paramSet(got []string, expected ...string) error {
	if len(got) != len(expected) {
		return fmt.Errorf("unexpected params length: got %d (%v), want %d (%v)", len(got), got, len(expected), expected)
	}

	gotSet := make(map[string]struct{}, len(got))
	for _, name := range got {
		gotSet[name] = struct{}{}
	}

	for _, name := range expected {
		if _, ok := gotSet[name]; !ok {
			return fmt.Errorf("expected param %q to be present, got %v", name, got)
		}
		delete(gotSet, name)
	}

	for name := range gotSet {
		return fmt.Errorf("unexpected extra param %q in %v", name, got)
	}

	return nil
}

func findUserDefined(prog *bytecode.Program, name string) (bytecode.UDF, error) {
	for _, udf := range prog.Functions.UserDefined {
		if udf.Name == name {
			return udf, nil
		}
	}

	return bytecode.UDF{}, fmt.Errorf("expected UDF %q in %v", name, prog.Functions.UserDefined)
}

func hostArity(host map[string]int, name string, want int) error {
	got, ok := host[name]
	if !ok {
		return fmt.Errorf("expected host function %q in %v", name, host)
	}

	if got != want {
		return fmt.Errorf("expected %s arity %d, got %d", name, want, got)
	}

	return nil
}

func TestUdfMetadataO0(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
FUNC f() => TEST_FN(1)
RETURN f()
`, func(prog *bytecode.Program) error {
			if len(prog.Functions.Host) != 1 {
				return fmt.Errorf("expected exactly 1 host function, got %d", len(prog.Functions.Host))
			}

			return hostArity(prog.Functions.Host, "TEST_FN", 1)
		}, "udf host call included in metadata"),
		ProgramCheck(`
FUNC f() => TEST_FN(1, 2)
RETURN TEST_FN(1)
`, func(prog *bytecode.Program) error {
			return hostArity(prog.Functions.Host, "TEST_FN", 2)
		}, "max host arity merges across scopes"),
		ProgramCheck(`
FUNC outer() (
  FUNC inner(x) => TEST_FN(x)
  RETURN inner(1)
)
RETURN outer()
`, func(prog *bytecode.Program) error {
			if len(prog.Functions.Host) != 1 {
				return fmt.Errorf("expected exactly 1 host function, got %d (%v)", len(prog.Functions.Host), prog.Functions.Host)
			}

			return hostArity(prog.Functions.Host, "TEST_FN", 1)
		}, "nested host call included in metadata"),
		ProgramCheck(`
FUNC f() => @foo
RETURN f()
`, func(prog *bytecode.Program) error {
			return paramSet(prog.Params, "foo")
		}, "udf param included in program params"),
		ProgramCheck(`
FUNC outer() (
  FUNC inner() => @foo
  RETURN inner()
)
RETURN outer()
`, func(prog *bytecode.Program) error {
			return paramSet(prog.Params, "foo")
		}, "nested udf param included in program params"),
		ProgramCheck(`
FUNC a() => TEST_FN(1)
FUNC b() => TEST_FN(1, 2, 3)
LET top = TEST_FN(1, 2)
RETURN [a(), b(), top]
`, func(prog *bytecode.Program) error {
			if len(prog.Functions.Host) != 1 {
				return fmt.Errorf("expected exactly 1 host function, got %d (%v)", len(prog.Functions.Host), prog.Functions.Host)
			}

			return hostArity(prog.Functions.Host, "TEST_FN", 3)
		}, "host arity merges across multiple udfs and top level"),
		ProgramCheck(`
FUNC used() => 1
FUNC unused() => TEST_FN(@foo)
RETURN used()
`, func(prog *bytecode.Program) error {
			if err := hostArity(prog.Functions.Host, "TEST_FN", 1); err != nil {
				return err
			}

			return paramSet(prog.Params, "foo")
		}, "unused udf metadata kept at o0"),
		ProgramCheck(`
USE FOO::TEST_FN AS FN
FUNC f() => FN(1)
RETURN FN()
`, func(prog *bytecode.Program) error {
			if len(prog.Functions.Host) != 1 {
				return fmt.Errorf("expected exactly 1 host function, got %d (%v)", len(prog.Functions.Host), prog.Functions.Host)
			}

			return hostArity(prog.Functions.Host, "FOO::TEST_FN", 1)
		}, "function alias preserves host metadata"),
		ProgramCheck(`
LET upper = Foo()
LET lower = foo()
RETURN [upper, lower]
`, func(prog *bytecode.Program) error {
			expected := map[string]int{
				"Foo": 0,
				"foo": 0,
			}

			if len(prog.Functions.Host) != len(expected) {
				return fmt.Errorf("expected %d host functions, got %d (%v)", len(expected), len(prog.Functions.Host), prog.Functions.Host)
			}

			for name, arity := range expected {
				if err := hostArity(prog.Functions.Host, name, arity); err != nil {
					return err
				}
			}

			return nil
		}, "case-distinct host names preserved"),
		ProgramCheck(`
USE Foo::Test_FN AS Fn
RETURN Fn()
`, func(prog *bytecode.Program) error {
			if len(prog.Functions.Host) != 1 {
				return fmt.Errorf("expected exactly 1 host function, got %d (%v)", len(prog.Functions.Host), prog.Functions.Host)
			}

			return hostArity(prog.Functions.Host, "Foo::Test_FN", 0)
		}, "function alias preserves exact case"),
		ProgramCheck(`
USE Foo AS F
RETURN f::Test_FN()
`, func(prog *bytecode.Program) error {
			if err := hostArity(prog.Functions.Host, "f::Test_FN", 0); err != nil {
				return err
			}

			if _, ok := prog.Functions.Host["Foo::Test_FN"]; ok {
				return fmt.Errorf("expected no exact-case alias rewrite on mismatch, got %v", prog.Functions.Host)
			}

			return nil
		}, "namespace alias mismatch does not rewrite metadata"),
		ProgramCheck(`
USE FOO AS F
FUNC f() => F::TEST_FN()
RETURN f()
`, func(prog *bytecode.Program) error {
			if err := hostArity(prog.Functions.Host, "FOO::TEST_FN", 0); err != nil {
				return err
			}

			if _, ok := prog.Functions.Host["FOO"]; ok {
				return fmt.Errorf("expected no bare FOO host metadata, got %v", prog.Functions.Host)
			}

			return nil
		}, "namespace alias preserves fully qualified host name"),
		ProgramCheck(`RETURN [@beta, @alpha, @beta, @gamma]`, func(prog *bytecode.Program) error {
			want := []string{"beta", "alpha", "gamma"}

			if len(prog.Params) != len(want) {
				return fmt.Errorf("unexpected params count: got %d (%v), want %d", len(prog.Params), prog.Params, len(want))
			}

			for i := range want {
				if prog.Params[i] != want[i] {
					return fmt.Errorf("unexpected param order at index %d: got %q, want %q", i, prog.Params[i], want[i])
				}
			}

			return nil
		}, "program params preserve first use order"),
		ProgramCheck(`RETURN @foo + @bar + @foo`, func(prog *bytecode.Program) error {
			var loads []bytecode.Instruction
			for _, inst := range prog.Bytecode {
				if inst.Opcode == bytecode.OpLoadParam {
					loads = append(loads, inst)
				}
			}

			if len(loads) != 3 {
				return fmt.Errorf("unexpected number of LOADP instructions: got %d", len(loads))
			}

			got := []bytecode.Operand{
				loads[0].Operands[1],
				loads[1].Operands[1],
				loads[2].Operands[1],
			}
			want := []bytecode.Operand{1, 2, 1}

			for i := range want {
				if got[i] != want[i] {
					return fmt.Errorf("unexpected slot at LOADP #%d: got %d, want %d", i, got[i], want[i])
				}

				if got[i].IsConstant() {
					return fmt.Errorf("expected LOADP source operand to be slot-encoded, got constant %s", got[i])
				}
			}

			return nil
		}, "load param uses slot operand"),
		ProgramCheck(`
LET x = @alpha
FUNC f() => @beta
RETURN x + f()
`, func(prog *bytecode.Program) error {
			wantParams := []string{"alpha", "beta"}
			if len(prog.Params) != len(wantParams) {
				return fmt.Errorf("unexpected params count: got %d (%v), want %d", len(prog.Params), prog.Params, len(wantParams))
			}

			for i := range wantParams {
				if prog.Params[i] != wantParams[i] {
					return fmt.Errorf("unexpected param at index %d: got %q, want %q", i, prog.Params[i], wantParams[i])
				}
			}

			udf, err := findUserDefined(prog, "f")
			if err != nil {
				return err
			}

			udfEntry := udf.Entry
			if udfEntry >= len(prog.Bytecode) {
				return fmt.Errorf("invalid UDF entry: %d (bytecode len: %d)", udfEntry, len(prog.Bytecode))
			}

			inst := prog.Bytecode[udfEntry]
			if inst.Opcode != bytecode.OpLoadParam {
				return fmt.Errorf("unexpected opcode at UDF entry %d: got %s, want %s", udfEntry, inst.Opcode, bytecode.OpLoadParam)
			}

			if got := inst.Operands[1]; got != bytecode.Operand(2) {
				return fmt.Errorf("unexpected UDF LOADP slot: got %d, want %d", got, bytecode.Operand(2))
			}

			if inst.Operands[1].IsConstant() {
				return fmt.Errorf("expected UDF LOADP source operand to be slot-encoded, got constant %s", inst.Operands[1])
			}

			return nil
		}, "udf param slot matches program ordering"),
		ProgramCheck(`
FUNC a() => 1
FUNC A() => 2
RETURN a() + A()
`, func(prog *bytecode.Program) error {
			lower, err := findUserDefined(prog, "a")
			if err != nil {
				return err
			}
			upper, err := findUserDefined(prog, "A")
			if err != nil {
				return err
			}

			if lower.DisplayName != "a" {
				return fmt.Errorf("expected lowercase UDF display name, got %q", lower.DisplayName)
			}

			if upper.DisplayName != "A" {
				return fmt.Errorf("expected uppercase UDF display name, got %q", upper.DisplayName)
			}

			return nil
		}, "case-distinct udf names preserve display metadata"),
	})
}

func TestUdfNestedCaptureMetadataAcrossScopes(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`
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
`, func(prog *bytecode.Program) error {
			outer, err := findUserDefined(prog, "outer")
			if err != nil {
				return err
			}
			middle, err := findUserDefined(prog, "middle")
			if err != nil {
				return err
			}
			inner, err := findUserDefined(prog, "inner")
			if err != nil {
				return err
			}

			if outer.Params != 2 {
				return fmt.Errorf("expected outer total params/captures to be 2, got %d", outer.Params)
			}
			if middle.Params != 4 {
				return fmt.Errorf("expected middle total params/captures to be 4, got %d", middle.Params)
			}
			if inner.Params != 5 {
				return fmt.Errorf("expected inner total params/captures to be 5, got %d", inner.Params)
			}

			return nil
		}, "nested captures tracked across scopes"),
	}, compiler.O0, compiler.O1)
}

func TestUdfNestedCompileStatePropagatesMetadata(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`
LET base = 10
FUNC outer(a) (
  FUNC middle(b) (
    FUNC inner(c) => TEST_FN(@foo, base + a + b + c)
    RETURN inner(1)
  )
  RETURN middle(2)
)
RETURN outer(3)
`, func(prog *bytecode.Program) error {
			if err := hostArity(prog.Functions.Host, "TEST_FN", 2); err != nil {
				return err
			}

			return paramSet(prog.Params, "foo")
		}, "nested udf compile restores metadata after inner state swap"),
	}, compiler.O0, compiler.O1)
}

func TestUdfNestedDirectReturnStillLowersToTailCall(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`
LET base = 1
FUNC outer(a) (
  FUNC target(x) => x + 1
  FUNC forward(x) => target(x + base + a)
  RETURN forward(2)
)
RETURN outer(3)
`, func(prog *bytecode.Program) error {
			forward, err := findUserDefined(prog, "forward")
			if err != nil {
				return err
			}

			nextEntry := len(prog.Bytecode)
			for _, udf := range prog.Functions.UserDefined {
				if udf.Entry > forward.Entry && udf.Entry < nextEntry {
					nextEntry = udf.Entry
				}
			}

			for idx := forward.Entry; idx < nextEntry; idx++ {
				if prog.Bytecode[idx].Opcode == bytecode.OpTailCall {
					return nil
				}
			}

			return fmt.Errorf("expected tail call in forward body between %d and %d", forward.Entry, nextEntry)
		}, "nested udf direct return preserves tail-call lowering"),
	}, compiler.O0, compiler.O1)
}

func TestUdfNestedScopeDoesNotLeakToSiblingCompilation(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`
FUNC outer() (
  FUNC onlyInside() => 1
  RETURN onlyInside()
)
FUNC sibling() => onlyInside()
RETURN sibling()
`, func(prog *bytecode.Program) error {
			return hostArity(prog.Functions.Host, "onlyInside", 0)
		}, "sibling udf compilation does not reuse prior nested scope"),
	}, compiler.O0, compiler.O1)
}

func TestUdfMetadataO1(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`
FUNC used() => 1
FUNC unused() => TEST_FN(@foo)
RETURN used()
`, func(prog *bytecode.Program) error {
			if _, ok := prog.Functions.Host["TEST_FN"]; ok {
				return fmt.Errorf("expected TEST_FN metadata to be pruned at O1, got %v", prog.Functions.Host)
			}

			return paramSet(prog.Params)
		}, "unused udf metadata pruned at o1"),
		ProgramCheck(`
USE FOO AS F
FUNC f() => 1
RETURN f()
`, func(prog *bytecode.Program) error {
			if _, err := findUserDefined(prog, "f"); err != nil {
				return fmt.Errorf("expected UDF f to remain reachable at O1: %w", err)
			}

			if _, ok := prog.Functions.Host["FOO"]; ok {
				return fmt.Errorf("expected no bare FOO host metadata at O1, got %v", prog.Functions.Host)
			}

			return nil
		}, "namespace alias does not shadow udf call"),
	}, compiler.O1)
}
