package http

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RegisterLib register `HTTP` namespace functions.
// @namespace HTTP
func RegisterLib(ns runtime.Namespace) {
	ns = ns.Namespace("HTTP")
	ns.Function().A1().
		Add("GET", GET).
		Add("POST", POST).
		Add("PUT", PUT).
		Add("DELETE", DELETE).
		Add("DO", REQUEST)
}
