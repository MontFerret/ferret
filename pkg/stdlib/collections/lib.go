package collections

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RegisterLib registers canonical collection functions and deprecated global aliases.
// @namespace collections
func RegisterLib(ns runtime.Namespace) {
	canonical := ns.Namespace("collections")

	canonical.Function().A1().
		Add("count_distinct", CountDistinct).
		Add("count", Count).
		Add("reverse", Reverse)

	canonical.Function().A2().
		Add("includes", Includes)

	registerLegacy(ns)
}
