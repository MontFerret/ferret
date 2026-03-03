package bytecode

// Functions groups host and user-defined function metadata required for execution.
type Functions struct {
	Host        map[string]int `json:"host,omitempty"`
	UserDefined []UDF          `json:"userDefined,omitempty"`
}
