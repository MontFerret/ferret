package utils

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) error {
	ns.Functions().
		Set1("WAIT", Wait).
		Set("PRINT", Print)

	return nil
}
