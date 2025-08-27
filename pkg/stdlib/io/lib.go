package io

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/io/fs"
	"github.com/MontFerret/ferret/pkg/stdlib/io/net"
)

// RegisterLib register `IO` namespace functions.
// @namespace IO
func RegisterLib(ns runtime.Namespace) error {
	io := ns.Namespace("IO")

	if err := fs.RegisterLib(io); err != nil {
		return runtime.Error(err, "register `FS`")
	}

	if err := net.RegisterLib(io); err != nil {
		return runtime.Error(err, "register `NET`")
	}

	return nil
}
