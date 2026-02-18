package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// BINARY asserts that value is a binary type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var Binary = base.TypeAssertion(runtime.TypeBinary)
