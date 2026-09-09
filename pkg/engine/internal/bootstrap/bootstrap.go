// Package bootstrap orchestrates native engine construction and failure rollback.
package bootstrap

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/session"
	"github.com/MontFerret/ferret/v2/pkg/module"
)

type (
	// Config supplies applied configuration and the resources already acquired
	// by engine options. Hooks and Resources must be non-nil.
	Config struct {
		Hooks             *host.Hooks
		Resources         *resource.Manager
		Modules           []module.Module
		Host              host.Config
		OptimizationLevel compiler.OptimizationLevel
		MaxActiveSessions int
	}

	// Result contains initialized dependencies for the native Engine to own.
	// A successful build transfers Resources to the caller without closing them.
	Result struct {
		Compiler      *compiler.Compiler
		DebugCompiler *compiler.Compiler
		Host          *host.Host
		Hooks         *host.Hooks
		Resources     *resource.Manager
		Limiter       *session.Limiter
	}
)

// Build takes responsibility for Config.Resources until successful return.
// Returned errors close eligible hooks and owned resources; borrowed resources
// remain untouched. Construction panics propagate without error rollback.
func Build(config Config) (result Result, resultErr error) {
	// Close hooks become eligible only after host services are constructed.
	var rollbackHooks *host.EngineHooks

	defer func() {
		resultErr = rollback(resultErr, rollbackHooks, config.Resources)
	}()

	compilerInstance, err := compiler.New(compiler.WithOptimizationLevel(config.OptimizationLevel))
	if err != nil {
		return Result{}, fmt.Errorf("compiler: %w", err)
	}

	debugCompiler, err := compiler.New(compiler.WithDebugInfo())
	if err != nil {
		return Result{}, fmt.Errorf("debug compiler: %w", err)
	}

	boot, err := host.NewBootstrap(config.Host, config.Hooks, config.Resources)
	if err != nil {
		return Result{}, fmt.Errorf("bootstrap: %w", err)
	}

	rollbackHooks = config.Hooks.EngineHooks

	for _, m := range config.Modules {
		if err := m.Register(boot); err != nil {
			return Result{}, err
		}
	}

	h, err := boot.Build()
	if err != nil {
		return Result{}, err
	}

	hooks := config.Hooks.Clone()
	rollbackHooks = hooks.EngineHooks

	// Run initialization after host finalization using the fixed hook snapshot.
	if err := hooks.EngineHooks.RunInit(); err != nil {
		return Result{}, fmt.Errorf("init hooks: %w", err)
	}

	return Result{
		Compiler:      compilerInstance,
		DebugCompiler: debugCompiler,
		Host:          h,
		Hooks:         hooks,
		Resources:     config.Resources,
		Limiter:       session.NewLimiter(config.MaxActiveSessions),
	}, nil
}
