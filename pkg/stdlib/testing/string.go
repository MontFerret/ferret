package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// STRING asserts that value is a string type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var String = base.TypeAssertion(runtime.TypeString)
