package spec

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

type (
	Expectation struct {
		Value     any
		Assertion assert.Assertion
	}

	Outcomes struct {
		Result Expectation
		Error  Expectation
	}
)

func (exp Expectation) Defined() bool {
	return exp.Assertion != nil
}

func (exp Expectation) Assert(t *testing.T, actual any) {
	t.Helper()

	exp.Assertion(t, actual, exp.Value)
}
