package collections

import "github.com/MontFerret/ferret/pkg/runtime/core"

func RegisterLib(ns core.Namespace) error {
	return ns.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
			"INCLUDES": Includes,
			"LENGTH":   Length,
			"REVERSE":  Reverse,
		}))
}
