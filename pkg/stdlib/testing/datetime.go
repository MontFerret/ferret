package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// DATETIME asserts that value is a datetime type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var DateTime = base.TypeAssertion(runtime.TypeDateTime)
