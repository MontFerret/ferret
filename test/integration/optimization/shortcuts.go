package optimization_test

import (
	"github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/test/integration/base"
)

type (
	UseCase struct {
		base.TestCase
		Execution Execution
	}

	Execution struct {
		Run       bool
		Expected  any
		Assertion convey.Assertion
	}
)

func NewCase(expression string, expected any, assertion convey.Assertion, exec Execution, desc ...string) UseCase {
	return UseCase{
		TestCase:  base.NewCase(expression, expected, assertion, desc...),
		Execution: exec,
	}
}

func Skip(uc UseCase) UseCase {
	uc.Skip = true

	return uc
}

var Debug = base.Debug
