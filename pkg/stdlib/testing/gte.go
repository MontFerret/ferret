package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// GTE asserts that an actual value is greater than or equal to an expected one.
// @param {Any} actual - Actual value.
// @param {Any} expected - Expected value.
// @param {String} [message] - Message to display on error.
var Gte = base.EqualityAssertion(base.GreaterOrEqualOp)
