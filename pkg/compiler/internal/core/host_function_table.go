package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/internal/hostfunction"
)

type (
	hostFunctionKey struct {
		Name     string
		ArgCount int
	}

	// HostFunctionTable assigns stable binding IDs to host call signatures.
	// It is owned by ProgramContext and shared by the main body and all UDFs.
	HostFunctionTable struct {
		bindings map[hostFunctionKey]int
		order    []bytecode.HostFunction
	}
)

func NewHostFunctionTable() *HostFunctionTable {
	return &HostFunctionTable{
		bindings: make(map[hostFunctionKey]int),
		order:    make([]bytecode.HostFunction, 0),
	}
}

// Bind returns the stable binding ID for a qualified name and call argument count.
func (t *HostFunctionTable) Bind(name string, argCount int) int {
	name = hostfunction.CanonicalName(name)
	key := hostFunctionKey{Name: name, ArgCount: argCount}

	if id, exists := t.bindings[key]; exists {
		return id
	}

	id := len(t.order)
	t.bindings[key] = id
	t.order = append(t.order, bytecode.HostFunction{Name: name, ArgCount: argCount})

	return id
}

// All returns the signatures in binding-ID order.
func (t *HostFunctionTable) All() []bytecode.HostFunction {
	return append([]bytecode.HostFunction(nil), t.order...)
}
