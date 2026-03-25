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
		Opcode("RETURN ``", compile.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpLoadConst: 1,
			},
		}, "", "should fold empty template literal into a single empty string constant"),

		Opcode("RETURN `hello`", compile.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpLoadConst: 1,
			},
		}, "hello", "should fold literal-only template into a single string constant"),

		Opcode("RETURN `use \\`backtick\\``", compile.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpLoadConst: 1,
			},
		}, "use `backtick`", "should fold escaped backtick in template literal"),

		Opcode("RETURN `${NONE}`", compile.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpLoadConst: 1,
			},
		}, "", "should fold NONE interpolation into empty string"),

		Opcode("RETURN `foo-${1}-bar-${true}`", compile.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpAdd:       0,
				bytecode.OpLoadConst: 1,
			},
		}, "foo-1-bar-true", "should fold fully constant template literal into a single string"),

		Opcode("LET x = \"X\" RETURN `a-${1}-b-${x}-c-${true}-d`", compile.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpAdd:       0,
				bytecode.OpLoadConst: 1,
			},
		}, "a-1-b-X-c-true-d", "should fold constant expressions in template literal into single chunks"),

		Opcode("RETURN `${@foo}`", compile.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat: 1,
			},
		}, "bar", "should not fold template literal with param interpolation").Env(vm.WithParam("foo", runtime.NewString("bar"))),

		Opcode("RETURN `${@a}${@b}`", compile.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat: 1,
			},
		}, "xy", "should not fold template literal with adjacent param interpolations").Env(vm.WithParam("a", runtime.NewString("x")), vm.WithParam("b", runtime.NewString("y"))),

		Opcode("RETURN `pre-${@foo}`", compile.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat: 1,
			},
		}, "pre-bar", "should not fold template literal with prefix literal and param interpolation").Env(vm.WithParam("foo", runtime.NewString("bar"))),

		Opcode("RETURN `pre-${@foo}-post`", compile.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat: 1,
			},
		}, "pre-bar-post", "should not fold template literal with suffix literal and param interpolation").Env(vm.WithParam("foo", runtime.NewString("bar"))),

		Opcode("RETURN `${@foo}-${1 + 2}`", compile.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat: 1,
			},
		}, "bar-3", "should keep concat with params but fold constant subexpressions").Env(vm.WithParam("foo", runtime.NewString("bar"))),

		Opcode("RETURN `cost=\\${1}`", compile.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpLoadConst: 1,
			},
		}, "cost=${1}", "escaped interpolation marker constant folds"),
	})
}
