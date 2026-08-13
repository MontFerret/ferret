package utils

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().Add("wait", Wait)
	ns.Function().Var().Add("print", Print)
}
