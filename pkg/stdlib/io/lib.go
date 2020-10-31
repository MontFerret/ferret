package io

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/io/fs"
	"github.com/MontFerret/ferret/pkg/stdlib/io/net"
)

// RegisterLib register `IO` namespace functions.
// @namespace IO
func RegisterLib(ns core.Namespace) error {
	io := ns.Namespace("IO")

	if err := fs.RegisterLib(io); err != nil {
		return core.Error(err, "register `FS`")
	}

	if err := net.RegisterLib(io); err != nil {
		return core.Error(err, "register `NET`")
	}

	return nil
}
