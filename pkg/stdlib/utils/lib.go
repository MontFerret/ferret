package utils

import "github.com/MontFerret/ferret/pkg/runtime/core"

func RegisterLib(ns core.Namespace) error {
	return ns.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
			"WAIT":  Wait,
			"PRINT": Print,
		}),
	)
}
