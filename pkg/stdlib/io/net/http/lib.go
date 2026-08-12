package http

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RegisterLib register `http` namespace functions.
// @namespace http
func RegisterLib(ns runtime.Namespace) {
	ns = ns.Namespace("http")
	ns.Function().A1().
		Add("get", GET).
		Add("post", POST).
		Add("put", PUT).
		Add("delete", DELETE).
		Add("do", REQUEST)
}
