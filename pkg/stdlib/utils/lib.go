package utils

import (
	"github.com/MontFerret/ferret/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) error {
	return ns.RegisterFunctions(
		runtime.NewFunctionsFromMap(map[string]runtime.Function{
			"WAIT":  Wait,
			"PRINT": Print,
		}),
	)
}
