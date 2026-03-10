package bytecode

// UDF describes a user-defined function compiled into bytecode.
type UDF struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Entry       int    `json:"entry"`
	Registers   int    `json:"registers"`
	Params      int    `json:"params"`
}
