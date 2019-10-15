package collections

import "github.com/MontFerret/ferret/pkg/runtime/core"

func RegisterLib(ns core.Namespace) error {
	return ns.RegisterFunctions(core.FunctionsMap{
		"LENGTH":  Length,
		"REVERSE": Reverse,
	})
}
