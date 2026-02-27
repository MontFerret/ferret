package io

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/io/fs"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/io/net"
)

// RegisterLib register `IO` namespace functions.
// @namespace IO
func RegisterLib(ns runtime.Namespace) {
	io := ns.Namespace("IO")

	fs.RegisterLib(io)
	net.RegisterLib(io)
}
