package engine

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/bootstrap"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	enginesession "github.com/MontFerret/ferret/v2/pkg/engine/internal/session"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

// Engine compiles queries into reusable plans and runs them against the configured host.
// Callers own directly created plans and must close them before the engine.
type Engine struct {
	closeErr          error
	compiler          *compiler.Compiler
	debugCompiler     *compiler.Compiler
	loader            *artifact.Loader
	host              *host.Host
	hooks             *host.Hooks
	resources         *resource.Manager
	limiter           *enginesession.Limiter
	closeOnce         sync.Once
	closed            atomic.Bool
	optimizationLevel compiler.OptimizationLevel
	idleCap           int
	totalCap          int
}

// New constructs an Engine from the provided options, registers all modules,
// builds the host, and runs engine init hooks. It returns an error if any
// option, module registration, or init hook fails. All non-nil options are
// applied in order, and failures from multiple options are joined before
// construction stops.
func New(setters ...Option) (*Engine, error) {
	opts, err := newConfig(setters)
	if err != nil {
		return nil, err
	}

	compilerLevel, err := opts.optimizationLevel.compilerLevel()
	if err != nil {
		err = fmt.Errorf("compiler: %w", err)
		if closeErr := opts.resources.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close engine: %w", closeErr))
		}

		return nil, err
	}

	state, err := bootstrap.Build(bootstrap.Config{
		Host: host.Config{
			Library:    opts.library,
			Params:     opts.params,
			Encoding:   opts.encoding,
			Logger:     opts.logger,
			FSRoot:     opts.fsRoot,
			Network:    opts.network,
			FSReadOnly: opts.fsReadOnly,
		},
		Hooks:             opts.hooks,
		Resources:         opts.resources,
		Modules:           opts.modules,
		OptimizationLevel: compilerLevel,
		MaxActiveSessions: opts.maxActiveSessions,
	})
	if err != nil {
		return nil, err
	}

	return &Engine{
		compiler:          state.Compiler,
		optimizationLevel: compilerLevel,
		debugCompiler:     state.DebugCompiler,
		loader:            opts.programLoader,
		host:              state.Host,
		hooks:             state.Hooks,
		resources:         state.Resources,
		limiter:           state.Limiter,
		idleCap:           opts.maxIdleVMsPerPlan,
		totalCap:          opts.maxVMsPerPlan,
	}, nil
}

// Compile synchronously compiles source into a reusable execution plan. Omitted
// optimization inherits the engine configuration; explicit per-plan options do
// not mutate the engine. The context must be non-nil. Compilation admitted before
// Close may finish afterward; callers must settle it before releasing resources
// used by their compile hooks.
func (e *Engine) Compile(ctx context.Context, src source.Source, opts ...PlanOption) (*Plan, error) {
	return e.compile(ctx, src, false, opts)
}

// CompileDebug compiles a reusable plan with debug metadata and no optimization.
// Context, concurrent closure, and ownership follow Compile.
func (e *Engine) CompileDebug(ctx context.Context, src source.Source, opts ...PlanOption) (*Plan, error) {
	return e.compile(ctx, src, true, opts)
}

// Load decodes a serialized program artifact and wraps it in a reusable plan.
// An admitted load may finish after Close; the caller owns the returned plan.
func (e *Engine) Load(data []byte) (*Plan, error) {
	if e.closed.Load() {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "engine is closed")
	}

	prog, err := e.loader.Load(data)
	if err != nil {
		return nil, err
	}

	return e.newPlan(prog)
}

// Run compiles source, executes it in a fresh session, and returns encoded output and an error.
// Similar to Session.Run, it may return a non-nil *encoding.Output together with a non-nil error
// (for example, if execution produced output but an after-run hook or cleanup failed).
func (e *Engine) Run(ctx context.Context, src source.Source, opts ...SessionOption) (output *encoding.Output, resultErr error) {
	plan, err := e.Compile(ctx, src)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := plan.Close(); closeErr != nil {
			resultErr = errors.Join(resultErr, closeErr)
		}
	}()

	session, err := plan.NewSession(ctx, opts...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := session.Close(); closeErr != nil {
			resultErr = errors.Join(resultErr, closeErr)
		}
	}()

	return session.Run(ctx)
}

// Close runs the engine close hooks and releases engine-scoped resources,
// including the configured rooted filesystem and owned network idle connections.
// It rejects new compilation and loading without waiting for admitted operations
// or closing descendants. Callers must settle outstanding work and close children
// first. Concurrent/repeated calls retain the cleanup result. An engine close
// hook must not recursively call Close on the same engine.
func (e *Engine) Close() error {
	e.closeOnce.Do(func() {
		e.closed.Store(true)
		e.closeErr = e.close()
	})

	return e.closeErr
}
