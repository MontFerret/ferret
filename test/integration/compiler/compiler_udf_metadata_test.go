package compiler_test

import (
	"fmt"
	"slices"
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

func hostSignature(host []bytecode.HostFunction, name string, argCount int) error {
	for _, fn := range host {
		if fn.Name == name && fn.ArgCount == argCount {
			return nil
		}
	}

	return fmt.Errorf("expected host function %q with %d arguments in %v", name, argCount, host)
}

func hasHostName(host []bytecode.HostFunction, name string) bool {
	for _, fn := range host {
		if fn.Name == name {
			return true
		}
	}

	return false
}

func expectHostSignatures(expected ...bytecode.HostFunction) func(*bytecode.Program) error {
	return func(prog *bytecode.Program) error {
		if !slices.Equal(prog.Functions.Host, expected) {
			return fmt.Errorf("unexpected host signatures: got %v, want %v", prog.Functions.Host, expected)
		}

		return nil
	}
}

func TestHostFunctionSignatureBindings(t *testing.T) {
	foo := []bytecode.HostFunction{
		{Name: "foo", ArgCount: 1},
		{Name: "foo", ArgCount: 2},
	}
	namespaced := []bytecode.HostFunction{
		{Name: "ns::foo", ArgCount: 1},
		{Name: "ns::foo", ArgCount: 2},
	}

	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(
			`RETURN [FOO(1), FOO(1, 2), FOO(3)]`,
			expectHostSignatures(foo...),
			"top-level host overload bindings preserve first-seen order and reuse matching signatures",
		),
		ProgramCheck(`
FUNC call() => [FOO(1), FOO(1, 2), FOO(3)]
RETURN call()
`, expectHostSignatures(foo...), "udf host overload bindings preserve first-seen order and reuse matching signatures"),
		ProgramCheck(
			`RETURN [NS::FOO(1), ns::foo(1, 2), Ns::FoO(3)]`,
			expectHostSignatures(namespaced...),
			"namespaced host overload bindings preserve first-seen order and reuse matching signatures",
		),
		ProgramCheck(
			`RETURN [FOO(1)?, FOO(1, 2)?, FOO(3)?]`,
			expectHostSignatures(foo...),
			"protected host overload bindings preserve first-seen order and reuse matching signatures",
		),
	}, compiler.O0, compiler.O1)
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

			return hostSignature(prog.Functions.Host, "test_fn", 1)
		}, "udf host call included in metadata"),
		ProgramCheck(`
FUNC f() => TEST_FN(1, 2)
RETURN TEST_FN(1)
`, func(prog *bytecode.Program) error {
			if len(prog.Functions.Host) != 2 {
				return fmt.Errorf("expected exactly 2 host signatures, got %d (%v)", len(prog.Functions.Host), prog.Functions.Host)
			}

			if err := hostSignature(prog.Functions.Host, "test_fn", 1); err != nil {
				return err
			}

			return hostSignature(prog.Functions.Host, "test_fn", 2)
		}, "host overloads remain distinct across scopes"),
		ProgramCheck(`
FUNC outer() {
  FUNC inner(x) => TEST_FN(x)
  RETURN inner(1)
}
RETURN outer()
`, func(prog *bytecode.Program) error {
			if len(prog.Functions.Host) != 1 {
				return fmt.Errorf("expected exactly 1 host function, got %d (%v)", len(prog.Functions.Host), prog.Functions.Host)
			}

			return hostSignature(prog.Functions.Host, "test_fn", 1)
		}, "nested host call included in metadata"),
		ProgramCheck(`
FUNC f() => @foo
RETURN f()
`, func(prog *bytecode.Program) error {
			return paramSet(prog.Params, "foo")
		}, "udf param included in program params"),
		ProgramCheck(`
FUNC outer() {
  FUNC inner() => @foo
  RETURN inner()
}
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
			if len(prog.Functions.Host) != 3 {
				return fmt.Errorf("expected exactly 3 host signatures, got %d (%v)", len(prog.Functions.Host), prog.Functions.Host)
			}

			for _, argCount := range []int{1, 2, 3} {
				if err := hostSignature(prog.Functions.Host, "test_fn", argCount); err != nil {
					return err
				}
			}

			return nil
		}, "host overloads remain distinct across multiple udfs and top level"),
		ProgramCheck(`
FUNC used() => 1
FUNC unused() => TEST_FN(@foo)
RETURN used()
`, func(prog *bytecode.Program) error {
			if err := hostSignature(prog.Functions.Host, "test_fn", 1); err != nil {
				return err
			}

			return paramSet(prog.Params, "foo")
		}, "unused udf metadata kept at o0"),
		ProgramCheck(`
USE FOO::TEST_FN AS FN
FUNC f() => FN(1)
RETURN FN()
`, func(prog *bytecode.Program) error {
			if len(prog.Functions.Host) != 2 {
				return fmt.Errorf("expected exactly 2 host signatures, got %d (%v)", len(prog.Functions.Host), prog.Functions.Host)
			}

			if err := hostSignature(prog.Functions.Host, "foo::test_fn", 1); err != nil {
				return err
			}

			return hostSignature(prog.Functions.Host, "foo::test_fn", 0)
		}, "function alias preserves host metadata"),
		ProgramCheck(`
LET upper = Foo()
LET lower = foo()
RETURN [upper, lower]
`, func(prog *bytecode.Program) error {
			if len(prog.Functions.Host) != 1 {
				return fmt.Errorf("expected one canonical host function, got %d (%v)", len(prog.Functions.Host), prog.Functions.Host)
			}

			return hostSignature(prog.Functions.Host, "foo", 0)
		}, "host call casing shares one canonical binding"),
		ProgramCheck(`
USE Foo::Test_FN AS Fn
RETURN Fn()
`, func(prog *bytecode.Program) error {
			if len(prog.Functions.Host) != 1 {
				return fmt.Errorf("expected exactly 1 host function, got %d (%v)", len(prog.Functions.Host), prog.Functions.Host)
			}

			return hostSignature(prog.Functions.Host, "foo::test_fn", 0)
		}, "function alias canonicalizes its qualified host target"),
		ProgramCheck(`
USE Foo AS F
RETURN f::Test_FN()
`, func(prog *bytecode.Program) error {
			if err := hostSignature(prog.Functions.Host, "f::test_fn", 0); err != nil {
				return err
			}

			if hasHostName(prog.Functions.Host, "foo::test_fn") {
				return fmt.Errorf("expected no exact-case alias rewrite on mismatch, got %v", prog.Functions.Host)
			}

			return nil
		}, "namespace alias mismatch does not rewrite metadata"),
		ProgramCheck(`
USE FOO AS F
FUNC f() => F::TEST_FN()
RETURN f()
`, func(prog *bytecode.Program) error {
			if err := hostSignature(prog.Functions.Host, "foo::test_fn", 0); err != nil {
				return err
			}

			if hasHostName(prog.Functions.Host, "foo") {
				return fmt.Errorf("expected no bare FOO host metadata, got %v", prog.Functions.Host)
			}

			return nil
		}, "namespace alias canonicalizes its fully qualified host target"),
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
FUNC outer(a) {
  LET outerLocal = 10
  FUNC middle(b) {
    FUNC inner(c) => global + a + outerLocal + b + c
    RETURN inner(1)
  }
  RETURN middle(2)
}
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

func TestUdfTransitiveCaptureMetadata(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`
LET base = 1
FUNC first(value) => second(value)
FUNC second(value) => third(value)
FUNC third(value) => base + value
RETURN first(1)
`, func(prog *bytecode.Program) error {
			for _, name := range []string{"first", "second", "third"} {
				fn, err := findUserDefined(prog, name)
				if err != nil {
					return err
				}

				if fn.Params != 2 {
					return fmt.Errorf("expected %s to have one argument and one capture, got %d total parameters", name, fn.Params)
				}
			}

			return nil
		}, "transitive captures are included in udf metadata"),
	}, compiler.O0, compiler.O1)
}

func TestUdfNestedCompileStatePropagatesMetadata(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`
LET base = 10
FUNC outer(a) {
  FUNC middle(b) {
    FUNC inner(c) => TEST_FN(@foo, base + a + b + c)
    RETURN inner(1)
  }
  RETURN middle(2)
}
RETURN outer(3)
`, func(prog *bytecode.Program) error {
			if err := hostSignature(prog.Functions.Host, "test_fn", 2); err != nil {
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
FUNC outer(a) {
  FUNC target(x) => x + 1
  FUNC forward(x) => target(x + base + a)
  RETURN forward(2)
}
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
FUNC outer() {
  FUNC onlyInside() => 1
  RETURN onlyInside()
}
FUNC sibling() => onlyInside()
RETURN sibling()
`, func(prog *bytecode.Program) error {
			return hostSignature(prog.Functions.Host, "onlyinside", 0)
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
			if hasHostName(prog.Functions.Host, "test_fn") {
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

			if hasHostName(prog.Functions.Host, "foo") {
				return fmt.Errorf("expected no bare FOO host metadata at O1, got %v", prog.Functions.Host)
			}

			return nil
		}, "namespace alias does not shadow udf call"),
	}, compiler.O1)
}
