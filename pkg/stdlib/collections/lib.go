package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) error {
	return ns.RegisterFunctions(runtime.
		NewFunctionsBuilder().
		Set1("COUNT_DISTINCT", CountDistinct).
		Set1("COUNT", Count).
		Set2("INCLUDES", Includes).
		Set1("REVERSE", Reverse).
		Build(),
	)
}
