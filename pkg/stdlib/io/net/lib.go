package net

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/io/net/http"
)

// RegisterLib register `net` namespace functions.
// @namespace net
func RegisterLib(ns runtime.Namespace) {
	net := ns.Namespace("net")

	http.RegisterLib(net)
}
