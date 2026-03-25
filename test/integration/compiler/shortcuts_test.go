package compiler_test

import (
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
