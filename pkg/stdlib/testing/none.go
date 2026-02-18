package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// NONE asserts that value is none.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var None = base.TypeAssertion(runtime.TypeNone)
