package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/compile"
	. "github.com/MontFerret/ferret/v2/test/spec/optimize"
)

func TestConstantFolding(t *testing.T) {
	RunUseCases(t, compiler.O1, []spec.Spec{
		OpcodeCount("RETURN ``", map[bytecode.Opcode]int{
			bytecode.OpConcat:    0,
			bytecode.OpLoadConst: 1,
		}, "", "should fold empty template literal into a single empty string constant"),

		OpcodeCount("RETURN `hello`", map[bytecode.Opcode]int{
			bytecode.OpConcat:    0,
			bytecode.OpLoadConst: 1,
		}, "hello", "should fold literal-only template into a single string constant"),

		OpcodeCount("RETURN `use \\`backtick\\``", map[bytecode.Opcode]int{
			bytecode.OpConcat:    0,
			bytecode.OpLoadConst: 1,
		}, "use `backtick`", "should fold escaped backtick in template literal"),

		OpcodeCount("RETURN `${NONE}`", map[bytecode.Opcode]int{
			bytecode.OpConcat:    0,
			bytecode.OpLoadConst: 1,
		}, "", "should fold NONE interpolation into empty string"),

		OpcodeCount("RETURN `foo-${1}-bar-${true}`", map[bytecode.Opcode]int{
			bytecode.OpConcat:    0,
			bytecode.OpAdd:       0,
			bytecode.OpLoadConst: 1,
		}, "foo-1-bar-true", "should fold fully constant template literal into a single string"),

		OpcodeCount("RETURN `${1s}`", map[bytecode.Opcode]int{
			bytecode.OpConcat:    0,
			bytecode.OpLoadConst: 1,
		}, "1s", "should fold Duration interpolation through String conversion"),

		OpcodeCount("LET x = \"X\" RETURN `a-${1}-b-${x}-c-${true}-d`", map[bytecode.Opcode]int{
			bytecode.OpConcat:    0,
			bytecode.OpAdd:       0,
			bytecode.OpLoadConst: 1,
		}, "a-1-b-X-c-true-d", "should fold constant expressions in template literal into single chunks"),

		OpcodeCount("RETURN `${@foo}`", map[bytecode.Opcode]int{
			bytecode.OpConcat: 1,
		}, "bar", "should not fold template literal with param interpolation").Env(vm.WithParam("foo", runtime.NewString("bar"))),

		OpcodeCount("RETURN `${@a}${@b}`", map[bytecode.Opcode]int{
			bytecode.OpConcat: 1,
		}, "xy", "should not fold template literal with adjacent param interpolations").Env(vm.WithParam("a", runtime.NewString("x")), vm.WithParam("b", runtime.NewString("y"))),

		OpcodeCount("RETURN `pre-${@foo}`", map[bytecode.Opcode]int{
			bytecode.OpConcat: 1,
		}, "pre-bar", "should not fold template literal with prefix literal and param interpolation").Env(vm.WithParam("foo", runtime.NewString("bar"))),

		OpcodeCount("RETURN `pre-${@foo}-post`", map[bytecode.Opcode]int{
			bytecode.OpConcat: 1,
		}, "pre-bar-post", "should not fold template literal with suffix literal and param interpolation").Env(vm.WithParam("foo", runtime.NewString("bar"))),

		OpcodeCount("RETURN `${@foo}-${1 + 2}`", map[bytecode.Opcode]int{
			bytecode.OpConcat: 1,
		}, "bar-3", "should keep concat with params but fold constant subexpressions").Env(vm.WithParam("foo", runtime.NewString("bar"))),

		OpcodeCount("RETURN `cost=\\${1}`", map[bytecode.Opcode]int{
			bytecode.OpConcat:    0,
			bytecode.OpLoadConst: 1,
		}, "cost=${1}", "escaped interpolation marker constant folds"),

		OpcodeCount(`RETURN 1s == "tomorrow"`, map[bytecode.Opcode]int{
			bytecode.OpEq: 0,
		}, false, "strict cross-type Duration equality folds to false"),

		OpcodeCount(`RETURN 1s != "tomorrow"`, map[bytecode.Opcode]int{
			bytecode.OpNe: 0,
		}, true, "strict cross-type Duration inequality folds to true"),

		OpcodeCount(`RETURN 1s == "1s"`, map[bytecode.Opcode]int{
			bytecode.OpEq: 0,
		}, false, "valid Duration string remains distinct without explicit conversion"),

		OpcodeCount(`RETURN 5.5 % 2`, map[bytecode.Opcode]int{
			bytecode.OpMod: 0,
		}, 1.5, "Float modulo folds to a Float remainder"),

		OpcodeErr(`RETURN "10" - 2`, compile.OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpSub},
		}, runtime.ErrInvalidOperation, "numeric strings require explicit conversion and remain runtime errors"),

		OpcodeErr(`RETURN true + 1`, compile.OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpAddConst},
		}, runtime.ErrInvalidOperation, "Boolean arithmetic remains a runtime error"),

		OpcodeErr(`RETURN 5s * "2"`, compile.OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpMul},
		}, runtime.ErrInvalidOperation, "numeric-string Duration multiplication remains a runtime error"),

		OpcodeErr(`RETURN "2" * 5s`, compile.OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpMul},
		}, runtime.ErrInvalidOperation, "reverse numeric-string Duration multiplication remains a runtime error"),

		OpcodeCount(`RETURN 5s * TO_NUMBER("2")`, map[bytecode.Opcode]int{
			bytecode.OpMul: 1,
		}, "10s", "explicit numeric conversion preserves Duration multiplication"),
	})
}
