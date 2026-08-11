package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// Tests whether the actual value is greater than the expected value.
// @param actual {Any} Actual value.
// @param expected {Any} Expected value.
// @param message {String} Message to display on error.
// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
var Gt = base.EqualityAssertion(base.GreaterOp)
