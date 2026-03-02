package bytecode

// UDF describes a user-defined function compiled into bytecode.
type UDF struct {
	Name      string `json:"name"`
	Entry     int    `json:"entry"`
	Registers int    `json:"registers"`
	Params    int    `json:"params"`
}
