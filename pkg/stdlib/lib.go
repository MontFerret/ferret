package stdlib

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// New creates a new standard library.
// It registers all available functions and namespaces to the root namespace and returns it.
func New() runtime.Namespace {
	ns := runtime.NewLibrary()

	RegisterLib(ns)

	return ns
}

func RegisterLib(ns runtime.Namespace) {
	if err := Full().Register(ns); err != nil {
		panic(err)
	}
}
