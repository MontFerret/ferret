package optimization_test

import (
	"github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"

	"github.com/MontFerret/ferret/v2/pkg/vm"

	"github.com/MontFerret/ferret/v2/test/integration/base"
)

type (
	UseCase struct {
		base.TestCase
		Execution Execution
	}

	Execution struct {
		Expected  any
		Assertion convey.Assertion
		Options   []vm.EnvironmentOption
		Run       bool
	}

	OpcodeExpectation interface {
		OpcodeExistence | OpcodeCount
	}

	OpcodeExistence struct {
		Exists    []bytecode.Opcode
		NotExists []bytecode.Opcode
	}

	OpcodeCount struct {
		Count map[bytecode.Opcode]int
	}
)

func NewCase(expression string, expected any, assertion convey.Assertion, exec Execution, desc ...string) UseCase {
	return UseCase{
		TestCase:  base.NewCase(expression, expected, assertion, desc...),
		Execution: exec,
	}
}

func Skip(uc UseCase) UseCase {
	uc.TestCase = base.Skip(uc.TestCase)

	return uc
}

func Debug(uc UseCase) UseCase {
	uc.TestCase = base.Debug(uc.TestCase)

	return uc
}

func Options(uc UseCase, opts ...vm.EnvironmentOption) UseCase {
	uc.Execution.Options = append(uc.Execution.Options, opts...)

	return uc
}
