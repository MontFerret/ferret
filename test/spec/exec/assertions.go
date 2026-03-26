package exec

import "github.com/MontFerret/ferret/v2/test/spec/assert"

var (
	ShouldBeRuntimeError = assert.NewBinaryAssertion(IsRuntimeError)
)
