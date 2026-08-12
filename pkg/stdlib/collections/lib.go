package collections

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().
		Add("count_distinct", CountDistinct).
		Add("count", Count).
		Add("reverse", Reverse)

	ns.Function().A2().
		Add("includes", Includes)
}
