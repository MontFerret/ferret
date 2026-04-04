package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

// HostParamTable stores the host-level query parameters (@name references)
// referenced anywhere in a compilation unit, including inside UDF bodies.
// It assigns each distinct name a stable 1-based slot in first-seen order.
//
// The table is owned by ProgramContext and lives for the entire compilation
// of one source file - host params are a program-wide concern, not per-function.
type HostParamTable struct {
	slots map[string]int
	order []string
}

func NewHostParamTable() *HostParamTable {
	return &HostParamTable{
		slots: make(map[string]int),
		order: make([]string, 0),
	}
}

// Bind assigns (or returns the existing) stable 1-based slot for name.
// The returned operand is the slot used by EmitLoadParam instructions.
func (t *HostParamTable) Bind(name string) bytecode.Operand {
	if slot, exists := t.slots[name]; exists {
		return bytecode.Operand(slot + 1)
	}

	slot := len(t.order)
	t.slots[name] = slot
	t.order = append(t.order, name)

	return bytecode.Operand(slot + 1)
}

// Names returns a defensive copy of the parameter names in first-seen order.
// The returned slice is what bytecode.Program.Params is populated from.
func (t *HostParamTable) Names() []string {
	out := make([]string, len(t.order))
	copy(out, t.order)

	return out
}

// HostFunctionTable tracks every host function referenced in a compilation
// unit, keeping the maximum arity observed per distinct name. Arity merging
// is required because the same host function may be called with different
// argument counts in different places (top-level and inside UDFs); the host
// environment needs to know the maximum.
//
// Owned by ProgramContext, shared across the main body and all UDF compiles.
type HostFunctionTable struct {
	arities map[string]int
}

func NewHostFunctionTable() *HostFunctionTable {
	return &HostFunctionTable{
		arities: make(map[string]int),
	}
}

// Bind records that name was called with args arguments. If name was
// previously bound with a larger arity, the larger value is kept.
func (t *HostFunctionTable) Bind(name string, args int) {
	if curr, exists := t.arities[name]; exists && curr > args {
		return
	}

	t.arities[name] = args
}

// All returns a defensive copy of name -> max arity. The returned map is
// what bytecode.Program.Functions.Host is populated from.
func (t *HostFunctionTable) All() map[string]int {
	out := make(map[string]int, len(t.arities))

	for k, v := range t.arities {
		out[k] = v
	}

	return out
}
