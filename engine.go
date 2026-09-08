package ferret

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Engine compiles queries into reusable plans and runs them against the configured host.
// Callers own directly created plans and must close them before the engine.
type Engine struct {
	closeErr          error
	compiler          *compiler.Compiler
	debugCompiler     *compiler.Compiler
	loader            *artifact.Loader
	host              *host
	hooks             *hookRegistry
	limiter           *sessionLimiter
	closeOnce         sync.Once
	closed            atomic.Bool
	optimizationLevel compiler.OptimizationLevel
	idleCap           int
	totalCap          int
	ownsNetwork       bool
}

// New constructs an Engine from the provided options, registers all modules,
// builds the host, and runs engine init hooks. It returns an error if any
// option, module registration, or init hook fails. All non-nil options are
// applied in order, and failures from multiple options are joined before
// construction stops.
func New(setters ...Option) (*Engine, error) {
	opts, err := newOptions(setters)
	if err != nil {
		return nil, err
	}

	ownsNetwork := opts.hostNetwork == false
	compilerLevel, err := opts.optimizationLevel.compilerLevel()
	if err != nil {
		if ownsNetwork {
			ferretnet.CloseIdleNetworkConnections(opts.network)
		}

		return nil, fmt.Errorf("compiler: %w", err)
	}

	compilerInstance, err := compiler.New(compiler.WithOptimizationLevel(compilerLevel))
	if err != nil {
		if ownsNetwork {
			ferretnet.CloseIdleNetworkConnections(opts.network)
		}

		return nil, fmt.Errorf("compiler: %w", err)
	}

	debugCompiler, err := compiler.New(compiler.WithDebugInfo())
	if err != nil {
		if ownsNetwork {
			ferretnet.CloseIdleNetworkConnections(opts.network)
		}

		return nil, fmt.Errorf("debug compiler: %w", err)
	}

	boot, err := newBootstrap(&opts)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: %w", err)
	}

	for _, m := range opts.modules {
		if err := m.Register(boot); err != nil {
			return nil, closeEngineOnError(
				err,
				boot.hooks.engine,
				boot.host.FileSystem(),
				boot.host.Network(),
				ownsNetwork,
			)
		}
	}

	h, err := boot.host.Build()
	if err != nil {
		return nil, closeEngineOnError(
			err,
			boot.hooks.engine,
			boot.host.FileSystem(),
			boot.host.Network(),
			ownsNetwork,
		)
	}

	hooks := boot.hooks.clone()
	// Run init hooks after bootstrap is finalized and before returning the engine.
	if err := hooks.engine.runInitHooks(); err != nil {
		initErr := fmt.Errorf("init hooks: %w", err)

		return nil, closeEngineOnError(initErr, hooks.engine, h.fs, h.network, ownsNetwork)
	}

	return &Engine{
		compiler:          compilerInstance,
		optimizationLevel: compilerLevel,
		debugCompiler:     debugCompiler,
		loader:            opts.programLoader,
		host:              h,
		hooks:             hooks,
		limiter:           newSessionLimiter(opts.maxActiveSessions),
		idleCap:           opts.maxIdleVMsPerPlan,
		totalCap:          opts.maxVMsPerPlan,
		ownsNetwork:       ownsNetwork,
	}, nil
}

// Compile synchronously compiles source into a reusable execution plan. Omitted
// optimization inherits the engine configuration; explicit per-plan options do
// not mutate the engine. The context must be non-nil. Compilation admitted before
// Close may finish afterward; callers must settle it before releasing resources
// used by their compile hooks.
func (e *Engine) Compile(ctx context.Context, src Source, opts ...PlanOption) (*Plan, error) {
	return e.compile(ctx, src, false, opts)
}

// CompileDebug compiles a reusable plan with debug metadata and no optimization.
// Context, concurrent closure, and ownership follow Compile.
func (e *Engine) CompileDebug(ctx context.Context, src Source, opts ...PlanOption) (*Plan, error) {
	return e.compile(ctx, src, true, opts)
}

func (e *Engine) compile(ctx context.Context, src Source, debug bool, setters []PlanOption) (*Plan, error) {
	if ctx == nil {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if e.closed.Load() {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "engine is closed")
	}

	level := e.optimizationLevel

	if debug {
		level = compiler.None
	}

	opts, err := newPlanOptions(level, debug, setters)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
			err = errors.Join(err, ctxErr)
		}

		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	selected := e.compiler
	if debug {
		selected = e.debugCompiler
	} else if opts.level != e.optimizationLevel {
		selected, err = compiler.New(compiler.WithOptimizationLevel(opts.level))
		if err != nil {
			return nil, err
		}
	}

	if err := e.hooks.plan.runBeforeCompileHooks(ctx); err != nil {
		err = fmt.Errorf("before compile hooks: %w", err)
		if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
			err = errors.Join(err, ctxErr)
		}

		return nil, err
	}

	var prog *bytecode.Program

	err = ctx.Err()
	if err == nil {
		prog, err = selected.Compile(ctx, src)
	}

	if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
		err = errors.Join(err, ctxErr)
	}

	if hookErr := e.hooks.plan.runAfterCompileHooks(ctx, err); hookErr != nil {
		err = errors.Join(err, fmt.Errorf("after compile hooks: %w", hookErr))
	}

	if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
		err = errors.Join(err, ctxErr)
	}

	if err != nil {
		return nil, err
	}

	return e.newPlan(prog)
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
// Similar to Session.Run, it may return a non-nil *Output together with a non-nil error
// (for example, if execution produced output but a deferred cleanup step failed).
func (e *Engine) Run(ctx context.Context, src Source, opts ...SessionOption) (output *Output, resultErr error) {
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
		e.closeErr = closeEngine(e.hooks.engine, e.host.fs, e.host.network, e.ownsNetwork)
	})

	return e.closeErr
}

func (e *Engine) newPlan(prog *bytecode.Program) (*Plan, error) {
	return &Plan{
		prog:         prog,
		closed:       make(chan struct{}),
		host:         e.host,
		hooks:        e.hooks.plan,
		sessionHooks: e.hooks.session,
		limiter:      e.limiter,
		pool:         vm.NewPoolWithLimits(prog, e.idleCap, e.totalCap),
	}, nil
}
