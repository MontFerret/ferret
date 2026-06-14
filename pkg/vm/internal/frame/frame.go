package frame

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

// CallFrame captures the caller state needed to resume execution after a UDF call.
type CallFrame struct {
	OwnedResources   mem.OwnedResources
	Aliases          mem.AliasTracker
	CellHandles      []mem.CellHandle
	FnName           string
	CallerRegisters  []runtime.Value
	FnID             int
	CallSitePC       int
	ReturnPC         int
	ReturnDest       bytecode.Operand
	HasCallSite      bool
	RecoveryBoundary bool
}

// structuralCallSitePC returns the original caller source PC even after a
// tail call replaces the diagnostic call-site metadata.
func (f *CallFrame) structuralCallSitePC() int {
	// runCore advances PC before entering the UDF, while call descriptors map
	// the source location to the instruction immediately before the call.
	return f.ReturnPC - 2
}
