package compiler_test

import (
	"github.com/MontFerret/ferret/pkg/vm"
	"github.com/MontFerret/ferret/test/integration/base"
)

type BC = []vm.Instruction
type UseCase = base.TestCase
type E = base.ExpectedError
type ME = base.ExpectedMultiError

var I = vm.NewInstruction
var C = vm.NewConstant
var R = vm.NewRegister

var NewCase = base.NewCase
var Skip = base.Skip
