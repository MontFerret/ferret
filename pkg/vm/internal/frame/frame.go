package frame

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// CallFrame captures the caller state needed to resume execution after a UDF call.
type CallFrame struct {
	Registers  []runtime.Value
	FnID       int
	ReturnPC   int
	ReturnDest bytecode.Operand
	Protected  bool
}
