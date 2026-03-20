package core

import (
	"fmt"
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Functions is a mutable registry of named functions, mirroring v1's core.Functions.
type Functions struct {
	data map[string]Function
}

// NewFunctions creates a new empty Functions registry.
func NewFunctions() *Functions {
	return &Functions{data: make(map[string]Function)}
}

// NewFunctionsFromMap creates a Functions registry pre-populated with the given map.
func NewFunctionsFromMap(m map[string]Function) *Functions {
	f := NewFunctions()
	for k, v := range m {
		f.data[k] = v
	}

	return f
}

// Get returns the function registered under name, and whether it was found.
func (f *Functions) Get(name string) (Function, bool) {
	fn, ok := f.data[name]
	return fn, ok
}

// Set registers fn under name, overwriting any existing entry.
func (f *Functions) Set(name string, fn Function) {
	f.data[name] = fn
}

// Unset removes the function registered under name.
func (f *Functions) Unset(name string) {
	delete(f.data, name)
}

// Names returns a sorted list of all registered function names.
func (f *Functions) Names() []string {
	names := make([]string, 0, len(f.data))
	for k := range f.data {
		names = append(names, k)
	}

	sort.Strings(names)

	return names
}

// --- Namespace ---

// Namespace mirrors the v1 core.Namespace interface.
type Namespace interface {
	Namespace(name string) Namespace
	RegisterFunction(name string, fun Function) error
	MustRegisterFunction(name string, fun Function)
	RegisterFunctions(funcs *Functions) error
	MustRegisterFunctions(funcs *Functions)
	RegisteredFunctions() []string
	RemoveFunction(name string)
}

// namespaceAdapter wraps a v2 runtime.Namespace and exposes the v1 Namespace interface.
type namespaceAdapter struct {
	ns       runtime.Namespace
	onChange func()
}

// WrapNamespace wraps a v2 runtime.Namespace into a compat Namespace.
func WrapNamespace(ns runtime.Namespace) Namespace {
	return &namespaceAdapter{ns: ns}
}

// WrapNamespaceWithObservability wraps a v2 runtime.Namespace into a compat Namespace with an optional change hook.
func WrapNamespaceWithObservability(ns runtime.Namespace, hook func()) Namespace {
	return &namespaceAdapter{ns: ns, onChange: hook}
}

func (a *namespaceAdapter) Namespace(name string) Namespace {
	return WrapNamespaceWithObservability(a.ns.Namespace(name), a.onChange)
}

func (a *namespaceAdapter) RegisterFunction(name string, fun Function) error {
	fns := a.ns.Function()
	if fns.Has(name) {
		return fmt.Errorf("function '%s' already registered", name)
	}

	fns.Var().Add(name, UnwrapFunction(fun))
	a.emitOnChange()

	return nil
}

func (a *namespaceAdapter) MustRegisterFunction(name string, fun Function) {
	if err := a.RegisterFunction(name, fun); err != nil {
		panic(err)
	}
}

func (a *namespaceAdapter) RegisterFunctions(funcs *Functions) error {
	if funcs == nil {
		return fmt.Errorf("functions cannot be nil")
	}

	for _, name := range funcs.Names() {
		fn, _ := funcs.Get(name)
		if err := a.RegisterFunction(name, fn); err != nil {
			return err
		}
	}

	return nil
}

func (a *namespaceAdapter) MustRegisterFunctions(funcs *Functions) {
	if err := a.RegisterFunctions(funcs); err != nil {
		panic(err)
	}
}

func (a *namespaceAdapter) RegisteredFunctions() []string {
	fns := a.ns.Function()
	seen := make(map[string]struct{})
	var names []string

	for _, defs := range []interface{ List() []string }{
		fns.Var(),
		fns.A0(),
		fns.A1(),
		fns.A2(),
		fns.A3(),
		fns.A4(),
	} {
		for _, n := range defs.List() {
			if _, ok := seen[n]; !ok {
				seen[n] = struct{}{}
				names = append(names, n)
			}
		}
	}

	sort.Strings(names)

	return names
}

func (a *namespaceAdapter) RemoveFunction(name string) {
	fns := a.ns.Function()

	switch {
	case fns.Var().Has(name):
		fns.Var().Remove(name)
	case fns.A0().Has(name):
		fns.A0().Remove(name)
	case fns.A1().Has(name):
		fns.A1().Remove(name)
	case fns.A2().Has(name):
		fns.A2().Remove(name)
	case fns.A3().Has(name):
		fns.A3().Remove(name)
	case fns.A4().Has(name):
		fns.A4().Remove(name)
	}

	a.emitOnChange()
}

func (a *namespaceAdapter) emitOnChange() {
	if a.onChange != nil {
		a.onChange()
	}
}
