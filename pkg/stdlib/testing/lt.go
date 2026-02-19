package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// LT asserts that an actual value is lesser than an expected one.
// @param {Any} actual - Actual value.
// @param {Any} expected - Expected value.
// @param {String} [message] - Message to display on error.
var Lt = base.EqualityAssertion(base.LessOp)
