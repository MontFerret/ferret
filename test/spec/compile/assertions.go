package compile

import (
	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

var (
	ShouldHaveSameBytecode   = assert.NewBinaryAssertion(EqualBytecode)
	ShouldContainOpcode      = assert.NewBinaryAssertion(ContainsOpcode)
	ShouldHaveRegisters      = assert.NewBinaryAssertion(EqualRegisters)
	ShouldBeCompilationError = assert.NewBinaryAssertion(IsCompilationError)
)
