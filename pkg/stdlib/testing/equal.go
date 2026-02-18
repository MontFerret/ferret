package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// EQUAL asserts equality of actual and expected values.
// @param {Any} actual - Actual value.
// @param {Any} expected - Expected value.
// @param {String} [message] - Message to display on error.
var Equal = base.EqualityAssertion(base.EqualOp)
