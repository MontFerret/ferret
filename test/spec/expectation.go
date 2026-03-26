package spec

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

type Expectation struct {
	Value     any
	Assertion assert.Assertion
}

func NewExpectation(fn assert.Assertion, val ...any) Expectation {
	exp := Expectation{
		Assertion: fn,
	}

	if len(val) > 0 {
		exp.Value = val[0]
	}

	return exp
}

func (exp Expectation) Defined() bool {
	return exp.Assertion != nil
}

func (exp Expectation) Assert(t *testing.T, actual any) {
	t.Helper()

	exp.Assertion(t, actual, exp.Value)
}

func (exp Expectation) Merge(other Expectation) Expectation {
	out := exp

	if other.Assertion != nil {
		out.Assertion = other.Assertion
	}

	if other.Value != nil {
		out.Value = other.Value
	}

	return out
}
