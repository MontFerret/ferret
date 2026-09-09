package ferret

import (
	"io"

	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/engine"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

// Option configures an Engine during construction.
type Option = engine.Option

// WithOptimizationLevel configures optimization for normal query compilation.
// OptimizationFull is used by default; debug compilation always uses OptimizationNone.
func WithOptimizationLevel(level OptimizationLevel) Option {
	return engine.WithOptimizationLevel(level)
}

// WithParams applies custom parameters to the options by merging them with existing ones, initializing if necessary.
// If a parameter already exists, it will be overwritten.
// All host values will be converted to a Value.
func WithParams(params map[string]any) Option {
	return engine.WithParams(params)
}

// WithRuntimeParams configures runtime parameters by merging the provided Params with existing ones in options.
// If a parameter already exists, it will be overwritten.
func WithRuntimeParams(params Params) Option {
	return engine.WithRuntimeParams(params)
}

// WithParam returns an Option that sets a parameter with the specified name and value in the options configuration.
// The name cannot be empty, and the value cannot be nil. It ensures the parameter value is correctly parsed and stored.
func WithParam(name string, value any) Option {
	return engine.WithParam(name, value)
}

// WithRuntimeParam returns an Option that sets a runtime parameter with the specified name and pre-converted Value.
// The name cannot be empty, and the value cannot be nil.
func WithRuntimeParam(name string, value Value) Option {
	return engine.WithRuntimeParam(name, value)
}

// WithNamespace merges the functions from the provided runtime.Namespace into the engine's function library.
func WithNamespace(ns runtime.Namespace) Option {
	return engine.WithNamespace(ns)
}

// WithFunctionsRegistrar creates an Option that invokes the provided registrar with the engine's runtime.Namespace if the registrar is not nil.
// Registered host-function names and namespace segments are canonicalized to lowercase and resolve case-insensitively in FQL.
func WithFunctionsRegistrar(setter func(ns runtime.Namespace)) Option {
	return engine.WithFunctionsRegistrar(setter)
}

// WithFunctions merges the provided *runtime.Functions into the engine's function library.
func WithFunctions(funcs *runtime.Functions) Option {
	return engine.WithFunctions(funcs)
}

// WithLog sets the writer for logging output.
// The writer can be any io.Writer, such as os.Stdout or a file.
func WithLog(writer io.Writer) Option {
	return engine.WithLog(writer)
}

// WithLogLevel sets the logging level for the engine.
// The logging level determines the severity of log messages that will be recorded.
func WithLogLevel(lvl LogLevel) Option {
	return engine.WithLogLevel(lvl)
}

// WithLogFields sets the fields to be included in log entries.
// These fields can provide additional context for debugging and monitoring purposes.
func WithLogFields(fields map[string]any) Option {
	return engine.WithLogFields(fields)
}

// WithEncodingRegistry sets a custom encoding registry for query execution.
func WithEncodingRegistry(registry *encoding.Registry) Option {
	return engine.WithEncodingRegistry(registry)
}

// WithProgramLoader sets a custom artifact loader for Engine.Load.
func WithProgramLoader(loader *artifact.Loader) Option {
	return engine.WithProgramLoader(loader)
}

// WithoutStdlib disables the standard library, so no built-in functions are registered by default.
func WithoutStdlib() Option {
	return engine.WithoutStdlib()
}

// WithStdlib configures which standard library groups are registered by default.
func WithStdlib(set stdlib.Set) Option {
	return engine.WithStdlib(set)
}

// WithModules creates an Option that appends the provided Module values if not empty.
func WithModules(mods ...Module) Option {
	return engine.WithModules(mods...)
}

// WithEncodingCodec registers or overrides a codec for the given content type.
func WithEncodingCodec(contentType string, codec encoding.Codec) Option {
	return engine.WithEncodingCodec(contentType, codec)
}

// WithEngineInitHook returns an Option that registers a hook to execute during engine initialization.
// It returns an error if hook is nil.
func WithEngineInitHook(hook EngineInitHook) Option {
	return engine.WithEngineInitHook(hook)
}

// WithEngineCloseHook returns an Option that registers a hook to execute when the engine is closed.
// It returns an error if hook is nil.
func WithEngineCloseHook(hook EngineCloseHook) Option {
	return engine.WithEngineCloseHook(hook)
}

// WithBeforeCompileHook returns an Option that registers a hook to execute before each compilation attempt.
// It returns an error if hook is nil.
func WithBeforeCompileHook(hook BeforeCompileHook) Option {
	return engine.WithBeforeCompileHook(hook)
}

// WithAfterCompileHook returns an Option that registers a hook to execute after each compilation attempt.
// The hook receives the compilation error (if any). It returns an error if hook is nil.
func WithAfterCompileHook(hook AfterCompileHook) Option {
	return engine.WithAfterCompileHook(hook)
}

// WithPlanCloseHook returns an Option that registers a hook to execute when a plan is closed.
// It returns an error if hook is nil.
func WithPlanCloseHook(hook PlanCloseHook) Option {
	return engine.WithPlanCloseHook(hook)
}

// WithBeforeRunHook returns an Option that registers a hook to execute before each session run.
// The hook can replace the context used by subsequent hooks and VM execution.
// It returns an error if hook is nil.
func WithBeforeRunHook(hook BeforeRunHook) Option {
	return engine.WithBeforeRunHook(hook)
}

// WithAfterRunHook returns an Option that registers a hook to execute after each session run attempt.
// The hook receives the run error (if any). It returns an error if hook is nil.
func WithAfterRunHook(hook AfterRunHook) Option {
	return engine.WithAfterRunHook(hook)
}

// WithSessionCloseHook returns an Option that registers a hook to execute when a session is closed.
// It returns an error if hook is nil.
func WithSessionCloseHook(hook SessionCloseHook) Option {
	return engine.WithSessionCloseHook(hook)
}

// WithMaxActiveSessions sets an engine-wide limit on concurrently active sessions.
//
// This limit applies to Session objects created from any plan compiled by the
// engine. When the limit is reached, Plan.NewSession blocks until another
// session is closed or the provided context is canceled.
//
// Use this when you want to put a global cap on query execution concurrency and
// the host-side resources that come with it, such as CPU, memory, network
// traffic, or downstream service pressure.
//
// This is different from WithMaxIdleVMsPerPlan and WithMaxVMsPerPlan:
// WithMaxActiveSessions controls how many sessions may be running or checked
// out at once across the engine, while the VM options control how each
// individual plan manages its VM pool.
//
// A value of 0 disables the limit.
func WithMaxActiveSessions(n int) Option {
	return engine.WithMaxActiveSessions(n)
}

// WithMaxIdleVMsPerPlan sets how many closed-session VMs each plan keeps warm
// for reuse after they become idle.
//
// This is a retention setting, not a concurrency limit. It only controls how
// many unused VMs remain cached in a plan's pool after sessions are closed.
// When the idle cache is full, additional returned VMs are closed instead of
// retained.
//
// Use this when the same compiled plan is executed repeatedly and you want to
// trade some steady-state memory for faster session creation by reusing already
// initialized VMs.
//
// This is different from WithMaxVMsPerPlan:
// WithMaxIdleVMsPerPlan controls how many unused VMs stay cached, while
// WithMaxVMsPerPlan controls the maximum total number of VMs the plan may own
// at all, including both idle and currently borrowed VMs.
//
// A value of 0 disables idle retention for the plan.
func WithMaxIdleVMsPerPlan(n int) Option {
	return engine.WithMaxIdleVMsPerPlan(n)
}

// WithMaxVMsPerPlan sets a hard per-plan limit on the total number of VMs the
// plan's pool may own at one time.
//
// The total includes both idle VMs kept in the pool and VMs currently borrowed
// by active sessions created from that plan. When the limit is reached and no
// idle VM is available to reuse, session creation fails with vm.ErrPoolExhausted.
//
// Use this when you need a strict upper bound on the memory or resource
// footprint of a single hot plan, even if that plan is under heavy concurrent
// load.
//
// This is different from WithMaxActiveSessions:
// WithMaxVMsPerPlan limits VM ownership for one plan, while
// WithMaxActiveSessions limits active session concurrency across the entire
// engine.
//
// This is also different from WithMaxIdleVMsPerPlan:
// WithMaxVMsPerPlan is a hard cap, while WithMaxIdleVMsPerPlan only decides how
// many unused VMs are retained after demand drops.
//
// A value of 0 means the plan may create as many VMs as needed, subject only to
// other limits such as WithMaxActiveSessions.
func WithMaxVMsPerPlan(n int) Option {
	return engine.WithMaxVMsPerPlan(n)
}

// WithFSRoot sets the root directory for the engine's file system.
func WithFSRoot(root string) Option {
	return engine.WithFSRoot(root)
}

// WithFSReadOnly sets the engine's file system to read-only mode.
func WithFSReadOnly() Option {
	return engine.WithFSReadOnly()
}

// WithNetwork sets the engine network service used by derived executions.
// If a network is provided, the engine will use it directly and will not manage its lifecycle.
// The host application is responsible for closing the network when it is no longer needed.
func WithNetwork(network ferretnet.Network) Option {
	return engine.WithNetwork(network)
}

// WithNetworkOptions creates an Option that constructs a new network service using the provided Ferret network options.
// The engine owns the resulting service and cleans it up when replaced, on
// construction failure, or at shutdown. If no options are provided, it is a no-op.
func WithNetworkOptions(setters ...ferretnet.Option) Option {
	return engine.WithNetworkOptions(setters...)
}
