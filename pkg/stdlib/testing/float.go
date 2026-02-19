package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// FLOAT asserts that value is a float type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var Float = base.TypeAssertion(runtime.TypeFloat)
