// Package compiler provides a v1-compatible Compiler for the Ferret compatibility layer.
// It mirrors the github.com/MontFerret/ferret/pkg/compiler package from Ferret v1.
package compiler

import (
	"context"
	"sync"

	ferret "github.com/MontFerret/ferret/v2"
	compatruntime "github.com/MontFerret/ferret/v2/compat/runtime"
	"github.com/MontFerret/ferret/v2/compat/runtime/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Compiler mirrors the v1 compiler.Compiler.
// It compiles FQL query strings into Programs and exposes a mutable function namespace.
//
// The underlying v2 Engine is created lazily on the first Compile call and is
// re-created whenever the function namespace is modified after a successful compile.
type Compiler struct {
	library  runtime.Library
	noStdlib bool
	mu       sync.Mutex
	engine   *ferret.Engine
	dirty    bool
}

// New creates a new Compiler with stdlib enabled.
func New() *Compiler {
	return &Compiler{
		library: runtime.NewLibrary(),
		dirty:   true,
	}
}

// Compile compiles the FQL query string and returns a Program.
func (c *Compiler) Compile(query string) (*compatruntime.Program, error) {
	eng, err := c.getEngine()
	if err != nil {
		return nil, err
	}

	return compatruntime.CompileFromSource(context.Background(), eng, query)
}

// MustCompile compiles the query and panics on error.
func (c *Compiler) MustCompile(query string) *compatruntime.Program {
	prog, err := c.Compile(query)
	if err != nil {
		panic(err)
	}

	return prog
}

// getEngine returns the cached engine, rebuilding it if needed.
func (c *Compiler) getEngine() (*ferret.Engine, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.dirty && c.engine != nil {
		return c.engine, nil
	}

	opts := []ferret.Option{ferret.WithNamespace(c.library)}

	if c.noStdlib {
		opts = append(opts, ferret.WithoutStdlib())
	}

	eng, err := ferret.New(opts...)
	if err != nil {
		return nil, err
	}

	c.engine = eng
	c.dirty = false

	return eng, nil
}

// markDirty invalidates the cached engine so it will be rebuilt on next Compile.
func (c *Compiler) markDirty() {
	c.mu.Lock()
	c.dirty = true
	c.mu.Unlock()
}

// --- Namespace methods ---

// Namespace returns a child namespace with the given name.
func (c *Compiler) Namespace(name string) core.Namespace {
	return core.WrapNamespace(c.library.Namespace(name))
}

// RegisterFunction registers a function in the root namespace.
func (c *Compiler) RegisterFunction(name string, fun core.Function) error {
	fns := c.library.Function()
	if fns.Has(name) {
		return nil // already registered — match v1 silent-overwrite semantics
	}

	fns.Var().Add(name, core.UnwrapFunction(fun))
	c.markDirty()

	return nil
}

// MustRegisterFunction registers a function, panicking on error.
func (c *Compiler) MustRegisterFunction(name string, fun core.Function) {
	if err := c.RegisterFunction(name, fun); err != nil {
		panic(err)
	}
}

// RegisterFunctions registers all functions from a Functions registry.
func (c *Compiler) RegisterFunctions(funcs *core.Functions) error {
	if funcs == nil {
		return nil
	}

	fns := c.library.Function()

	for _, name := range funcs.Names() {
		fn, _ := funcs.Get(name)
		fns.Var().Add(name, core.UnwrapFunction(fn))
	}

	c.markDirty()

	return nil
}

// MustRegisterFunctions registers all functions, panicking on error.
func (c *Compiler) MustRegisterFunctions(funcs *core.Functions) {
	if err := c.RegisterFunctions(funcs); err != nil {
		panic(err)
	}
}

// RemoveFunction removes a function from the root namespace.
func (c *Compiler) RemoveFunction(name string) {
	fns := c.library.Function()

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

	c.markDirty()
}

// RegisteredFunctions returns all function names registered in the root namespace.
func (c *Compiler) RegisteredFunctions() []string {
	return core.WrapNamespace(c.library).RegisteredFunctions()
}
