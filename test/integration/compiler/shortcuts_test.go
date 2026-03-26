package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

type (
	BC = []bytecode.Instruction
	E  = assert.ExpectedError
	ME = assert.ExpectedMultiError
)

var (
	I = bytecode.NewInstruction
	C = bytecode.NewConstant
	R = bytecode.NewRegister
)

func RunSpecsWith(t *testing.T, name string, c *compiler.Compiler, specs []spec.Spec) {
	t.Helper()

	runner := &spec.Runner{
		Name:     name,
		Compiler: c,
	}

	runner.Run(t, specs)
}

func RunSpecs(t *testing.T, useCases []spec.Spec) {
	t.Helper()
	RunSpecsWith(t, "Compiler/O0", compiler.New(compiler.WithOptimizationLevel(compiler.O0)), useCases)
}

func RunSpecsLevels(t *testing.T, specs []spec.Spec, levels ...compiler.OptimizationLevel) {
	t.Helper()

	if len(levels) == 0 {
		levels = []compiler.OptimizationLevel{compiler.O0}
	}

	for _, level := range levels {
		RunSpecsWith(
			t,
			fmt.Sprintf("Compiler/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			specs,
		)
	}
}

func ProgramCheck(expression string, fn func(*bytecode.Program) error, desc ...string) spec.Spec {
	return spec.NewSpec(expression, desc...).Expect().Compile(assert.NewUnaryAssertion(func(actual any) error {
		prog, ok := actual.(*bytecode.Program)
		if !ok {
			return fmt.Errorf("expected *bytecode.Program, got %T", actual)
		}

		return fn(prog)
	}))
}
