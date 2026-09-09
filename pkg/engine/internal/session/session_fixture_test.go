package session

import (
	"context"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type sessionFixture struct {
	host    *host.Host
	hooks   *host.SessionHooks
	program *bytecode.Program
	pool    *vm.Pool
	limiter *Limiter
	closed  chan struct{}
}

func newSessionFixture(t *testing.T, capacity int) *sessionFixture {
	t.Helper()

	resources := resource.NewManager()
	t.Cleanup(func() { _ = resources.Close() })
	hooks := host.NewHooks()
	boot, err := host.NewBootstrap(host.Config{
		Library: runtime.NewLibrary(), Params: runtime.Params{"value": runtime.Int(0)},
		Encoding: encoding.NewRegistry(encodingjson.Default), FSRoot: t.TempDir(),
	}, hooks, resources)
	if err != nil {
		t.Fatal(err)
	}

	h, err := boot.Build()
	if err != nil {
		t.Fatal(err)
	}

	compiler, err := compiler.New(compiler.WithDebugInfo())
	if err != nil {
		t.Fatal(err)
	}

	program, err := compiler.Compile(t.Context(), source.NewAnonymous("RETURN @value + 1"))
	if err != nil {
		t.Fatal(err)
	}

	pool := vm.NewPoolWithLimits(program, 1, 1)
	t.Cleanup(func() { _ = pool.Close() })

	return &sessionFixture{
		host: h, hooks: hooks.SessionHooks, program: program, pool: pool,
		limiter: NewLimiter(capacity), closed: make(chan struct{}),
	}
}

func (f *sessionFixture) execution(t *testing.T, config Config) *Execution {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	execution, err := NewExecution(ctx, config, f.host, f.hooks, f.limiter, f.pool, f.closed)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = execution.Close() })

	return execution
}
