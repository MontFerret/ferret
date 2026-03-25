package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/base/compilation"
)

func TestConstantFolding(t *testing.T) {
	RunUseCases(t, compiler.O1, []UseCase{
		OpcodeCase("RETURN ``", compilation.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpLoadConst: 1,
			},
		}, "", "should fold empty template literal into a single empty string constant"),

		OpcodeCase("RETURN `hello`", compilation.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpLoadConst: 1,
			},
		}, "hello", "should fold literal-only template into a single string constant"),

		OpcodeCase("RETURN `use \\`backtick\\``", compilation.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpLoadConst: 1,
			},
		}, "use `backtick`", "should fold escaped backtick in template literal"),

		OpcodeCase("RETURN `${NONE}`", compilation.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpLoadConst: 1,
			},
		}, "", "should fold NONE interpolation into empty string"),

		OpcodeCase("RETURN `foo-${1}-bar-${true}`", compilation.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpAdd:       0,
				bytecode.OpLoadConst: 1,
			},
		}, "foo-1-bar-true", "should fold fully constant template literal into a single string"),

		OpcodeCase("LET x = \"X\" RETURN `a-${1}-b-${x}-c-${true}-d`", compilation.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpAdd:       0,
				bytecode.OpLoadConst: 1,
			},
		}, "a-1-b-X-c-true-d", "should fold constant expressions in template literal into single chunks"),

		Options(OpcodeCase("RETURN `${@foo}`", compilation.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat: 1,
			},
		}, "bar", "should not fold template literal with param interpolation"), vm.WithParam("foo", runtime.NewString("bar"))),

		Options(OpcodeCase("RETURN `${@a}${@b}`", compilation.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat: 1,
			},
		}, "xy", "should not fold template literal with adjacent param interpolations"), vm.WithParam("a", runtime.NewString("x")), vm.WithParam("b", runtime.NewString("y"))),

		Options(OpcodeCase("RETURN `pre-${@foo}`", compilation.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat: 1,
			},
		}, "pre-bar", "should not fold template literal with prefix literal and param interpolation"), vm.WithParam("foo", runtime.NewString("bar"))),

		Options(OpcodeCase("RETURN `pre-${@foo}-post`", compilation.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat: 1,
			},
		}, "pre-bar-post", "should not fold template literal with suffix literal and param interpolation"), vm.WithParam("foo", runtime.NewString("bar"))),

		Options(OpcodeCase("RETURN `${@foo}-${1 + 2}`", compilation.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat: 1,
			},
		}, "bar-3", "should keep concat with params but fold constant subexpressions"), vm.WithParam("foo", runtime.NewString("bar"))),

		Options(OpcodeCase("RETURN `cost=\\${1}`", compilation.OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpLoadConst: 1,
			},
		}, "cost=${1}", "escaped interpolation marker constant folds")),
	})
}
