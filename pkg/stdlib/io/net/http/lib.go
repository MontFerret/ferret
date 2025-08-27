package http

import (
	"github.com/MontFerret/ferret/pkg/runtime"
)

// RegisterLib register `HTTP` namespace functions.
// @namespace HTTP
func RegisterLib(ns runtime.Namespace) error {
	return ns.
		Namespace("HTTP").
		RegisterFunctions(
			runtime.NewFunctionsFromMap(map[string]runtime.Function{
				"GET":    GET,
				"POST":   POST,
				"PUT":    PUT,
				"DELETE": DELETE,
				"DO":     REQUEST,
			}))
}
