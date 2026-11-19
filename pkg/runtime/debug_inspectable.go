package runtime

import (
	"context"

	"github.com/MontFerret/api/debugger"
)

type (
	// DebugInfo contains optional presentation hints for development tools.
	// TypeName is a semantic name. Display is a shallow summary. Each field
	// is optional. Tools may bound both fields to their presentation budget.
	DebugInfo struct {
		TypeName string
		Display  string
	}

	// DebugChildView describes a stable named or indexed sequence. A known count
	// is exact and nonnegative. Set Known with Count zero for an absent kind.
	// Readable grants bounded arbitrary-offset access.
	DebugChildView struct {
		Count    int64
		Known    bool
		Readable bool
	}

	// DebugChildrenInfo describes the independent child views of a paused value.
	DebugChildrenInfo struct{ Named, Indexed DebugChildView }

	// DebugChildrenProvider opts a host value into immediate child inspection.
	// Metadata must be cheap and must not read children. Both hooks must be
	// deterministic and side-effect free. Reachable backing storage, including
	// borrowed slices, must remain stable for the duration of the stop.
	// DebugChildAt receives a named or indexed filter and an absolute position
	// within that view. It observes ctx, returns found=false at the end, and
	// preserves exact UTF-8 keys. Indexed names are supplied by the debugger.
	// Neither hook may perform I/O, execute queries, or consume iterators.
	DebugChildrenProvider interface {
		DebugChildrenInfo() DebugChildrenInfo
		DebugChildAt(ctx context.Context, kind debugger.VariableFilter, position int64) (name string, value Value, found bool, err error)
	}

	// DebugInspectable allows a runtime value to provide optional presentation
	// hints for development tools. Implementations must be cheap, deterministic,
	// side-effect free, and must not read children, consume lazy values or
	// perform external work. The debugger can bound its own formatting work
	// only after the hook returns. The hook grants no access to children.
	DebugInspectable interface {
		DebugInfo() DebugInfo
	}
)

// CanExpand reports whether the view can have readable children.
func (v DebugChildView) CanExpand() bool { return v.Readable && (!v.Known || v.Count > 0) }
