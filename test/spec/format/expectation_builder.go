package format

import "github.com/MontFerret/ferret/v2/test/spec/assert"

type ExpectationBuilder struct {
	spec *Spec
}

func NewExpectationBuilder(spec *Spec) ExpectationBuilder {
	return ExpectationBuilder{
		spec: spec,
	}
}

func (b ExpectationBuilder) Output(fn assert.Assertion, val ...any) Spec {
	spec := b.spec

	if len(val) > 0 {
		spec.Output.Result.Value = val[0]
	}

	spec.Output.Result.Assertion = fn

	return *spec
}

func (b ExpectationBuilder) OutputError(fn assert.Assertion, val ...any) Spec {
	spec := b.spec

	if len(val) > 0 {
		spec.Output.Error.Value = val[0]
	}

	spec.Output.Error.Assertion = fn

	return *spec
}
