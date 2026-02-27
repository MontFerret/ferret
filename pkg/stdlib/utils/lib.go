package utils

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().Add("WAIT", Wait)
	ns.Function().Var().Add("PRINT", Print)
}
