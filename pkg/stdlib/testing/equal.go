package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// Tests equality of the actual and expected values.
// @param actual {Any} Actual value.
// @param expected {Any} Expected value.
// @param message {String} Message to display on error.
// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
var Equal = base.EqualityAssertion(base.EqualOp)
