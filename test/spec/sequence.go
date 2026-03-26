package spec

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

type (
	Sequence struct {
		Env   []vm.EnvironmentOption
		VM    []vm.Option
		Steps []SequenceStep
		Base  BaseSpec
	}

	SequenceStep struct {
		Result     Expectation
		Error      Expectation
		Panic      Expectation
		EnvFactory func() (*vm.Environment, error)
		Name       string
		Env        []vm.EnvironmentOption
	}
)

func NewSequence(expression string, desc ...string) Sequence {
	return Sequence{
		Base: BaseSpec{
			Input:       NewExpressionInput(expression),
			Description: strings.Join(desc, " "),
		},
	}
}

func ResultStep(name string, fn assert.Assertion, val ...any) SequenceStep {
	return SequenceStep{
		Name:   name,
		Result: NewExpectation(fn, val...),
	}
}

func ErrorStep(name string, fn assert.Assertion, val ...any) SequenceStep {
	return SequenceStep{
		Name:  name,
		Error: NewExpectation(fn, val...),
	}
}

func PanicStep(name string, fn assert.Assertion, val ...any) SequenceStep {
	return SequenceStep{
		Name:  name,
		Panic: NewExpectation(fn, val...),
	}
}
