// Package compat provides a Ferret v1-compatible public API surface, built on top of
// the Ferret v2 engine. It is intended to ease migration of existing v1 code with
// minimal changes.
//
// # Migration notes
//
//   - Import paths change from github.com/MontFerret/ferret/… to
//     github.com/MontFerret/ferret/v2/compat/…
//   - New() now panics on engine initialisation errors (same contract as v1 New).
//     Use NewE() if you need the error.
//   - Functions registered via Instance.Functions() after New() are not reflected in
//     subsequent Compile/Exec calls (v2 builds functions immutably at engine creation).
//   - The drivers package (CDP, HTTP) has moved to a separate repository.
//     See the Ferret v2 migration guide for details.
package compat

import (
	"context"

	ferret "github.com/MontFerret/ferret/v2"
	compatruntime "github.com/MontFerret/ferret/v2/compat/runtime"
	"github.com/MontFerret/ferret/v2/compat/runtime/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Instance is the v1-compatible entry point for compiling and executing FQL queries.
type Instance struct {
	engine  *ferret.Engine
	library runtime.Library
}

// New creates a new Instance with the provided options.
// It panics if the underlying v2 Engine cannot be initialised.
// Use NewE if you need to handle the error.
func New(setters ...Option) *Instance {
	inst, err := NewE(setters...)
	if err != nil {
		panic(err)
	}

	return inst
}

// NewE creates a new Instance, returning an error if initialisation fails.
func NewE(setters ...Option) (*Instance, error) {
	o, err := newInstanceOptions(setters)
	if err != nil {
		return nil, err
	}

	lib := runtime.NewLibrary()

	engineOpts := []ferret.Option{
		ferret.WithNamespace(lib),
	}

	if o.noStdlib {
		engineOpts = append(engineOpts, ferret.WithoutStdlib())
	}

	engineOpts = append(engineOpts, o.engineOpts...)

	eng, err := ferret.New(engineOpts...)
	if err != nil {
		return nil, err
	}

	return &Instance{
		engine:  eng,
		library: lib,
	}, nil
}

// Functions returns a v1-compatible Namespace backed by the engine's function library.
//
// WARNING: Functions registered after New() are NOT automatically picked up by
// subsequent Compile/Exec calls. Pre-register all custom functions before constructing
// the Instance.
func (inst *Instance) Functions() core.Namespace {
	return core.WrapNamespace(inst.library)
}

// Compile compiles the FQL query string into a Program.
func (inst *Instance) Compile(query string) (*compatruntime.Program, error) {
	return compatruntime.CompileFromSource(context.Background(), inst.engine, query)
}

// MustCompile compiles the query and panics on error.
func (inst *Instance) MustCompile(query string) *compatruntime.Program {
	prog, err := inst.Compile(query)
	if err != nil {
		panic(err)
	}

	return prog
}

// Exec compiles and immediately executes the FQL query, returning the JSON result.
func (inst *Instance) Exec(ctx context.Context, query string, opts ...compatruntime.Option) ([]byte, error) {
	src := source.NewAnonymousSource(query)

	out, err := inst.engine.Run(ctx, src, compatruntime.ToSessionOptions(opts)...)
	if err != nil {
		return nil, err
	}

	return out.Content, nil
}

// MustExec compiles and executes the query, panicking on error.
func (inst *Instance) MustExec(ctx context.Context, query string, opts ...compatruntime.Option) []byte {
	out, err := inst.Exec(ctx, query, opts...)
	if err != nil {
		panic(err)
	}

	return out
}

// Run executes a previously compiled Program.
func (inst *Instance) Run(ctx context.Context, prog *compatruntime.Program, opts ...compatruntime.Option) ([]byte, error) {
	return prog.Run(ctx, opts...)
}

// MustRun executes a previously compiled Program, panicking on error.
func (inst *Instance) MustRun(ctx context.Context, prog *compatruntime.Program, opts ...compatruntime.Option) []byte {
	out, err := inst.Run(ctx, prog, opts...)
	if err != nil {
		panic(err)
	}

	return out
}
