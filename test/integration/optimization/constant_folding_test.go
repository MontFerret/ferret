package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestConstantFolding(t *testing.T) {
	RunUseCases(t, compiler.O1, []UseCase{
		OpcodeCase("RETURN `foo-${1}-bar-${true}`", OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpAdd:       0,
				bytecode.OpLoadConst: 1,
			},
		}, "foo-1-bar-true", "should fold fully constant template literal into a single string"),

		OpcodeCase("LET x = \"X\" RETURN `a-${1}-b-${x}-c-${true}-d`", OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpAdd:       0,
				bytecode.OpLoadConst: 1,
			},
		}, "a-1-b-X-c-true-d", "should fold constant expressions in template literal into single chunks"),

		Options(OpcodeCase("RETURN `${@foo}`", OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat: 1,
			},
		}, "bar", "should not fold template literal with param interpolation"), vm.WithParam("foo", runtime.NewString("bar"))),

		Options(OpcodeCase("RETURN `${@a}${@b}`", OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat: 1,
			},
		}, "xy", "should not fold template literal with adjacent param interpolations"), vm.WithParam("a", runtime.NewString("x")), vm.WithParam("b", runtime.NewString("y"))),

		Options(OpcodeCase("RETURN `pre-${@foo}`", OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat: 1,
			},
		}, "pre-bar", "should not fold template literal with prefix literal and param interpolation"), vm.WithParam("foo", runtime.NewString("bar"))),

		Options(OpcodeCase("RETURN `cost=\\${1}`", OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpConcat:    0,
				bytecode.OpLoadConst: 1,
			},
		}, "cost=${1}", "escaped interpolation marker constant folds")),
	})
}
