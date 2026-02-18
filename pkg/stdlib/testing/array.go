package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// ARRAY asserts that value is a array type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var Array = base.TypeAssertion(runtime.TypeArray)
