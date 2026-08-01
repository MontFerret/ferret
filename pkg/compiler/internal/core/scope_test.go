package core_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
)

func TestScopeProjectionPreservesHiddenBindingIdentity(t *testing.T) {
	registers := core.NewRegisterAllocator()
	symbols := core.NewSymbolTable(registers, nil)
	symbols.EnterScope()

	if _, ok := symbols.DeclareLocalWithOptions("value", core.TypeInt, core.BindingOptions{
		ID:     core.BindingID(1),
		Hidden: true,
	}); !ok {
		t.Fatal("expected hidden binding declaration")
	}
	if _, ok := symbols.DeclareLocalWithOptions("value", core.TypeString, core.BindingOptions{
		ID: core.BindingID(2),
	}); !ok {
		t.Fatal("expected visible binding declaration")
	}

	projection := core.NewScopeProjection(
		registers,
		core.NewEmitter(),
		symbols,
		core.NewTypeTracker(),
		symbols.ProjectionVariables(),
	)

	symbols.ExitScope()
	symbols.EnterScope()
	projection.RestoreFromArray(bytecode.NewRegister(20))

	hidden, ok := symbols.ResolveBindingByID(core.BindingID(1))
	if !ok || !hidden.Hidden {
		t.Fatalf("expected restored hidden binding, got %#v", hidden)
	}
	visible, ok := symbols.ResolveBindingByID(core.BindingID(2))
	if !ok || visible.Hidden {
		t.Fatalf("expected restored visible binding, got %#v", visible)
	}
	byName, ok := symbols.ResolveBinding("value")
	if !ok || byName.ID != core.BindingID(2) {
		t.Fatalf("expected source lookup to ignore hidden binding, got %#v", byName)
	}
}
