package utils

import (
	"github.com/MontFerret/ferret/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) error {
	ns.Functions().
		Set1("WAIT", Wait).
		Set("PRINT", Print)

	return nil
}
