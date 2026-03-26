package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/optimize"
)

func TestConcatChainLowering(t *testing.T) {
	RunUseCases(t, compiler.O1, []spec.Spec{
		OpcodeCount(`RETURN "a" + 1 + "b" + 2 + "c" + 3`, map[bytecode.Opcode]int{
			bytecode.OpAdd:       0,
			bytecode.OpConcat:    0,
			bytecode.OpLoadConst: 1,
		}, "a1b2c3", "should fold fully constant concat chains into one constant"),

		OpcodeCount(`RETURN "a" + 1 + "b" + 2 + @x + "c" + 3`, map[bytecode.Opcode]int{
			bytecode.OpAdd:    0,
			bytecode.OpConcat: 1,
		}, "a1b2Xc3", "should keep one concat for mixed chains with merged constant runs").Env(vm.WithParam("x", runtime.NewString("X"))),

		OpcodeCount(`VAR str = ""
str += "a" + 1 + "b" + 2 + @x + "c" + 3
RETURN str`, map[bytecode.Opcode]int{
			bytecode.OpAdd:    0,
			bytecode.OpConcat: 1,
		}, "a1b2Xc3", "should route string += through concat-chain lowering").Env(vm.WithParam("x", runtime.NewString("X"))),

		OpcodeCount(`RETURN 1 + 2 + "x"`, map[bytecode.Opcode]int{
			bytecode.OpAdd:    0,
			bytecode.OpConcat: 0,
		}, "3x", "should preserve arithmetic boundaries before string concat"),
	})
}
