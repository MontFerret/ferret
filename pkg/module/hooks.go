package module

import "context"

type (
	// HookRegistrar provides access to lifecycle hook registrars for the engine,
	// plan, and session stages created from a bootstrapped engine.
	HookRegistrar interface {
		// Engine returns the registrar for engine initialization and shutdown
		// hooks.
		Engine() EngineHookRegistrar
		// Plan returns the registrar for plan compilation and shutdown hooks.
		Plan() PlanHookRegistrar
		// Session returns the registrar for session run and shutdown hooks.
		Session() SessionHookRegistrar
	}

	// EngineHookRegistrar registers hooks for engine initialization and shutdown.
	EngineHookRegistrar interface {
		// OnInit registers a hook executed in FIFO order during engine
		// initialization.
		// A nil hook is ignored.
		OnInit(hook EngineInitHook)
		// OnClose registers a hook executed in LIFO order when engine
		// construction fails or the engine is closed.
		// A nil hook is ignored.
		OnClose(hook EngineCloseHook)
	}

	// PlanHookRegistrar registers hooks for plan compilation and shutdown.
	PlanHookRegistrar interface {
		// BeforeCompile registers a hook executed in FIFO order before plan
		// compilation starts.
		// A nil hook is ignored.
		BeforeCompile(hook BeforeCompileHook)
		// AfterCompile registers a hook executed in LIFO order after each
		// compilation attempt.
		// It receives the compilation error (if any), and a nil hook is ignored.
		AfterCompile(hook AfterCompileHook)
		// OnClose registers a hook executed in LIFO order when a compiled plan is
		// closed.
		// A nil hook is ignored.
		OnClose(hook PlanCloseHook)
	}

	// SessionHookRegistrar registers hooks for session execution and shutdown.
	SessionHookRegistrar interface {
		// BeforeRun registers a hook executed in FIFO order before each session
		// run.
		// Hooks can replace the context passed to subsequent hooks and VM execution.
		// A nil hook is ignored.
		BeforeRun(hook BeforeRunHook)
		// AfterRun registers a hook executed in LIFO order after each session run
		// attempt.
		// It receives the run error (if any), and a nil hook is ignored.
		AfterRun(hook AfterRunHook)
		// OnClose registers a hook executed in LIFO order when a session is
		// closed.
		// A nil hook is ignored.
		OnClose(hook SessionCloseHook)
	}
)

type (
	// EngineInitHook runs during engine initialization.
	// Returning an error stops initialization immediately.
	EngineInitHook func() error

	// EngineCloseHook runs during engine shutdown.
	// Close hooks are executed in LIFO order and their errors are aggregated.
	EngineCloseHook func() error

	// BeforeCompileHook runs before compilation starts.
	// Hooks run in FIFO order and stop on the first error.
	BeforeCompileHook func(ctx context.Context) error

	// AfterCompileHook runs after each compilation attempt.
	// Hooks run in LIFO order and receive the compilation error (if any).
	AfterCompileHook func(ctx context.Context, err error) error

	// PlanCloseHook runs when a plan is closed.
	// Close hooks are executed in LIFO order and their errors are aggregated.
	PlanCloseHook func() error

	// BeforeRunHook runs before each session run.
	// It can return a derived context for subsequent hooks and VM execution.
	BeforeRunHook func(ctx context.Context) (context.Context, error)

	// AfterRunHook runs after each session run attempt.
	// Hooks run in LIFO order, receive the run error (if any), and aggregate hook errors.
	AfterRunHook func(ctx context.Context, err error) error

	// SessionCloseHook runs when a session is closed.
	// Close hooks are executed in LIFO order and their errors are aggregated.
	SessionCloseHook func() error
)
