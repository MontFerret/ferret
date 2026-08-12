package io

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/io/fs"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/io/net"
)

// RegisterLib register `io` namespace functions.
// @namespace io
func RegisterLib(ns runtime.Namespace) {
	io := ns.Namespace("io")

	fs.RegisterLib(io)
	net.RegisterLib(io)
}
