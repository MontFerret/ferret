package debugger

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type evalScope struct {
	locals   map[string]runtime.Value
	params   runtime.Params
	values   vm.DebugValueAccess
	maxBytes int
}

func (s evalScope) typeError(value runtime.Value, expected string) error {
	typeName := formatTypeName(s.values.Describe(value), s.maxBytes)

	return runtime.Errorf(runtime.ErrInvalidArgument, "debugger expression requires %s, got %s", expected, typeName)
}
