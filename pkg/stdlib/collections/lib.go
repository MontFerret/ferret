package collections

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().
		Add("COUNT_DISTINCT", CountDistinct).
		Add("COUNT", Count).
		Add("REVERSE", Reverse)

	ns.Function().A2().
		Add("INCLUDES", Includes)
}
