package compiler_test

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/base/assert"
)

type BC = []bytecode.Instruction
type UseCase = spec.Spec
type E = assert.ExpectedError
type ME = assert.ExpectedMultiError

var I = bytecode.NewInstruction
var C = bytecode.NewConstant
var R = bytecode.NewRegister

var NewCase = spec.NewSpec
var Skip = spec.Skip
