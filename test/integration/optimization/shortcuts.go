package optimization_test

import (
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/base/assert"
)

type (
	UseCase struct {
		spec.Spec
		Execution Execution
	}

	Execution struct {
		Expected  any
		Assertion assert.Assertion
		Options   []vm.EnvironmentOption
		Run       bool
	}
)

func NewCase(expression string, expected any, assertion assert.Assertion, exec Execution, desc ...string) UseCase {
	return UseCase{
		Spec:      spec.NewSpec(expression, expected, assertion, desc...),
		Execution: exec,
	}
}

func Skip(uc UseCase) UseCase {
	uc.Spec = spec.Skip(uc.Spec)

	return uc
}

func Debug(uc UseCase) UseCase {
	uc.Spec = spec.Debug(uc.Spec)

	return uc
}

func Options(uc UseCase, opts ...vm.EnvironmentOption) UseCase {
	uc.Execution.Options = append(uc.Execution.Options, opts...)

	return uc
}
