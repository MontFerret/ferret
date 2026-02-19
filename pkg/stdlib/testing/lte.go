package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// LTE asserts that an actual value is lesser than or equal to an expected one.
// @param {Any} actual - Actual value.
// @param {Any} expected - Expected value.
// @param {String} [message] - Message to display on error.
var Lte = base.EqualityAssertion(base.LessOrEqualOp)
