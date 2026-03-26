package compiler_test

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
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
