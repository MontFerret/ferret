package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// INT asserts that value is a int type.
// @param {Any} actual - Actual value.
// @param {String} [message] - Message to display on error.
var Int = base.TypeAssertion(runtime.TypeInt)
