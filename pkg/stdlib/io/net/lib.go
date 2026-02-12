package net

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/io/net/http"
)

// RegisterLib register `NET` namespace functions.
// @namespace NET
func RegisterLib(ns runtime.Namespace) error {
	io := ns.Namespace("NET")

	if err := http.RegisterLib(io); err != nil {
		return runtime.Error(err, "register `HTTP`")
	}

	return nil
}
