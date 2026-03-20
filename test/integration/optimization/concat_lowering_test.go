package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestConcatChainLowering(t *testing.T) {
	RunUseCases(t, compiler.O1, []UseCase{
		OpcodeCase(`RETURN "a" + 1 + "b" + 2 + "c" + 3`, OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpAdd:       0,
				bytecode.OpConcat:    0,
				bytecode.OpLoadConst: 1,
			},
		}, "a1b2c3", "should fold fully constant concat chains into one constant"),

		Options(OpcodeCase(`RETURN "a" + 1 + "b" + 2 + @x + "c" + 3`, OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpAdd:    0,
				bytecode.OpConcat: 1,
			},
		}, "a1b2Xc3", "should keep one concat for mixed chains with merged constant runs"), vm.WithParam("x", runtime.NewString("X"))),

		Options(OpcodeCase(`VAR str = ""
str += "a" + 1 + "b" + 2 + @x + "c" + 3
RETURN str`, OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpAdd:    0,
				bytecode.OpConcat: 1,
			},
		}, "a1b2Xc3", "should route string += through concat-chain lowering"), vm.WithParam("x", runtime.NewString("X"))),

		OpcodeCase(`RETURN 1 + 2 + "x"`, OpcodeCount{
			Count: map[bytecode.Opcode]int{
				bytecode.OpAdd:    0,
				bytecode.OpConcat: 0,
			},
		}, "3x", "should preserve arithmetic boundaries before string concat"),
	})
}
