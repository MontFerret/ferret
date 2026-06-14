package debugger

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type evalScope struct {
	locals map[string]runtime.Value
	params runtime.Params
	values vm.DebugValueAccess
}
