package compiler_test

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/integration/base"
)

type BC = []bytecode.Instruction
type UseCase = base.TestCase
type E = base.ExpectedError
type ME = base.ExpectedMultiError

var I = bytecode.NewInstruction
var C = bytecode.NewConstant
var R = bytecode.NewRegister

var NewCase = base.NewCase
var Skip = base.Skip
