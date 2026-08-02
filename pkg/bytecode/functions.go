package bytecode

// Functions groups host and user-defined function metadata required for execution.
type Functions struct {
	Host        []HostFunction `json:"host,omitempty"`
	UserDefined []UDF          `json:"userDefined,omitempty"`
}
