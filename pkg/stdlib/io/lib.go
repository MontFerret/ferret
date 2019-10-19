package io

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/io/fs"
)

// RegisterLib register `IO` namespace functions.
func RegisterLib(ns core.Namespace) error {
	io := ns.Namespace("IO")

	err := fs.RegisterLib(io)
	if err != nil {
		return core.Error(err, "register `FS`")
	}

	return nil
}
