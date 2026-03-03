package frame

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// CallFrame captures the caller state needed to resume execution after a UDF call.
type CallFrame struct {
	// FnID identifies the current UDF for diagnostics and tail-call updates.
	FnID int
	// ReturnPC is the program counter to resume at after returning.
	ReturnPC int
	// ReturnDest is the destination register for the return value.
	ReturnDest bytecode.Operand
	// Registers holds the caller register window to restore on return.
	Registers []runtime.Value
	// Protected marks frames that act as unwind targets for protected calls.
	Protected bool
}
