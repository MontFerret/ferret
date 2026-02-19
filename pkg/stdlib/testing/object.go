package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// OBJECT asserts that value is a object type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var Object = base.TypeAssertion(runtime.TypeObject)
