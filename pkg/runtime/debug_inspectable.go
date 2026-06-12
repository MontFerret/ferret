package runtime

// DebugInfo contains optional presentation hints for development tools.
// Empty fields leave the tool's existing presentation unchanged.
type DebugInfo struct {
	TypeName string
	Display  string
}

// DebugInspectable allows a runtime value to provide optional presentation
// hints for development tools. Implementations must be cheap, deterministic,
// side-effect free, and must not consume lazy values or perform external work.
type DebugInspectable interface {
	DebugInfo() DebugInfo
}
