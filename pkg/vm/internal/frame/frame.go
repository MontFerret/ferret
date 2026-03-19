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
	CellIDs          []uint64
	FnName           string
	CallerRegisters  []runtime.Value
	FnID             int
	CallSitePC       int
	ReturnPC         int
	ReturnDest       bytecode.Operand
	HasCallSite      bool
	RecoveryBoundary bool
}
