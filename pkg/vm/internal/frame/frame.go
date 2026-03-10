package frame

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// CallFrame captures the caller state needed to resume execution after a UDF call.
type CallFrame struct {
	FnName      string
	Registers   []runtime.Value
	FnID        int
	CallSitePC  int
	ReturnPC    int
	ReturnDest  bytecode.Operand
	HasCallSite bool
	Protected   bool
}
