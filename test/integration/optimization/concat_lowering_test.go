package optimization_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/compile"
	. "github.com/MontFerret/ferret/v2/test/spec/optimize"
)

func TestConcatChainLowering(t *testing.T) {
	RunUseCases(t, compiler.OptimizationFull, []spec.Spec{
		OpcodeCount(`RETURN "a" + 1 + "b" + 2 + "c" + 3`, map[bytecode.Opcode]int{
			bytecode.OpAdd:       0,
			bytecode.OpConcat:    0,
			bytecode.OpLoadConst: 1,
		}, "a1b2c3", "should fold fully constant concat chains into one constant"),

		OpcodeCount(`RETURN "a" + 1 + "b" + 2 + @x + "c" + 3`, map[bytecode.Opcode]int{
			bytecode.OpAdd:    0,
			bytecode.OpConcat: 1,
		}, "a1b2Xc3", "should concatenate unknown values once a String anchors the expression").Env(vm.WithParam("x", runtime.NewString("X"))),

		OpcodeCount(`VAR str = ""
str += "a" + 1 + "b" + 2 + @x + "c" + 3
		RETURN str`, map[bytecode.Opcode]int{
			bytecode.OpAdd:    0,
			bytecode.OpConcat: 1,
		}, "a1b2Xc3", "should concatenate unknown values in the String expression and assignment").Env(vm.WithParam("x", runtime.NewString("X"))),

		OpcodeCount(`RETURN 1 + 2 + "x"`, map[bytecode.Opcode]int{
			bytecode.OpAdd:    0,
			bytecode.OpConcat: 0,
		}, "3x", "should preserve arithmetic boundaries before string concat"),
	})
}

func TestStringAnchoredTemporalConcatenation(t *testing.T) {
	RunUseCases(t, compiler.OptimizationFull, []spec.Spec{
		Opcode(`RETURN TYPENAME(NOW() + "5m")`, compile.OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpAddConst},
		}, "String", "String concatenation takes precedence over DateTime arithmetic"),
		Opcode(`RETURN @value + "500ms"`, compile.OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpAddConst},
		}, "5s500ms", "String concatenation takes precedence over Duration arithmetic").Env(vm.WithParam("value", runtime.NewDuration(5*time.Second))),
	})
}
