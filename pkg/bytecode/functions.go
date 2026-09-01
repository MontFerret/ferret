package bytecode

type (
	// FunctionID identifies a function within one compiled program.
	FunctionID int

	// Functions groups host and user-defined function metadata required for execution.
	Functions struct {
		Host        []HostFunction `json:"host,omitempty"`
		UserDefined []UDF          `json:"userDefined,omitempty"`
	}
)

// NoFunction identifies the top-level program body rather than a callable
// function table entry.
const NoFunction FunctionID = -1

// Valid reports whether id represents a valid function ID.
func (id FunctionID) Valid() bool {
	return id >= 0
}

// InRange reports whether id can reference a function table with count entries.
func (id FunctionID) InRange(count int) bool {
	return id.Valid() && int(id) < count
}
