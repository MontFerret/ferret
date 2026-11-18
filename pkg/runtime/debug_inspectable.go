package runtime

type (
	// DebugInfo contains optional presentation hints for development tools.
	// TypeName is a semantic name. Display is a shallow summary. Each field
	// is optional. Tools may bound both fields to their presentation budget.
	DebugInfo struct {
		TypeName string
		Display  string
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
