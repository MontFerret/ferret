package http

import "github.com/MontFerret/ferret/pkg/runtime/core"

// RegisterLib register `HTTP` namespace functions.
// @namespace HTTP
func RegisterLib(ns core.Namespace) error {
	return ns.
		Namespace("HTTP").
		RegisterFunctions(
			core.NewFunctionsFromMap(map[string]core.Function{
				"GET":    GET,
				"POST":   POST,
				"PUT":    PUT,
				"DELETE": DELETE,
				"DO":     REQUEST,
			}))
}
