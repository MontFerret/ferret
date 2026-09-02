package ferret

import (
	"fmt"
	"io"
	"strings"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	encodingmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/module"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

type (
	config struct {
		library           runtime.Library
		network           ferretnet.Network
		managedNetworks   []ferretnet.Network
		hooks             *hookRegistry
		encoding          *encoding.Registry
		params            runtime.Params
		programLoader     *artifact.Loader
		fsRoot            string
		stdlib            stdlib.Set
		modules           []module.Module
		logger            []logging.Option
		compiler          []compiler.Option
		maxActiveSessions int
		maxIdleVMsPerPlan int
		maxVMsPerPlan     int
		hostNetwork       bool
		fsReadOnly        bool
	}

	// Option configures an Engine during construction.
	Option = gooptions.Option[config]
)

type encodingCodecAlias struct {
	encoding.Codec
	contentType string
}

const (
	defaultMaxActiveSessions = 0 // 0 means no limit on active sessions.
	defaultVMPoolSize        = 8
	defaultMaxVMsPerPlan     = 0 // 0 means no limit on total VMs per plan.
)

func (c encodingCodecAlias) ContentType() string {
	return c.contentType
}

func defaultOptions() config {
	return config{
		library:           runtime.NewLibrary(),
		params:            make(map[string]runtime.Value),
		encoding:          encoding.NewRegistry(encodingjson.Default, encodingmsgpack.Default),
		programLoader:     artifact.NewDefaultLoader(),
		hooks:             newHookRegistry(),
		maxActiveSessions: defaultMaxActiveSessions,
		maxIdleVMsPerPlan: defaultVMPoolSize,
		maxVMsPerPlan:     defaultMaxVMsPerPlan,
		stdlib:            stdlib.Full(),
	}
}

func newOptions(setters []Option) (config, error) {
	opts, err := gooptions.ApplyTo(defaultOptions(), setters...)
	if err == nil {
		if registerErr := opts.stdlib.Register(opts.library); registerErr != nil {
			err = fmt.Errorf("stdlib: %w", registerErr)
		}
	}

	managedNetworks := opts.managedNetworks
	if err == nil && !opts.hostNetwork && len(managedNetworks) > 0 {
		// Managed networks are recorded in option order, so the last entry is
		// the selected network when the final network option is managed.
		managedNetworks = managedNetworks[:len(managedNetworks)-1]
	}

	for _, network := range managedNetworks {
		ferretnet.CloseIdleNetworkConnections(network)
	}

	// Tracking references must not outlive option finalization. On success,
	// only the selected network is transferred to Engine ownership.
	opts.managedNetworks = nil

	if err != nil {
		return config{}, err
	}

	return opts, nil
}

// WithParams applies custom parameters to the options by merging them with existing ones, initializing if necessary.
// If a parameter already exists, it will be overwritten.
// All host values will be converted to a runtime.Value.
func WithParams(params map[string]any) Option {
	return func(opts *config) error {
		if len(params) == 0 {
			return nil
		}

		if opts.params == nil {
			opts.params = runtime.NewParams()
		}

		merged, err := opts.params.Merge(params)

		if err != nil {
			return err
		}

		opts.params = merged

		return nil
	}
}

// WithRuntimeParams configures runtime parameters by merging the provided params with existing ones in options.
// If a parameter already exists, it will be overwritten.
func WithRuntimeParams(params runtime.Params) Option {
	return func(opts *config) error {
		if len(params) == 0 {
			return nil
		}

		if opts.params == nil {
			opts.params = runtime.NewParams()
		}

		opts.params = opts.params.MergeParams(params)

		return nil
	}
}

// WithParam returns an Option that sets a parameter with the specified name and value in the options configuration.
// The name cannot be empty, and the value cannot be nil. It ensures the parameter value is correctly parsed and stored.
func WithParam(name string, value any) Option {
	return func(opts *config) error {
		if name == "" {
			return fmt.Errorf("param name cannot be empty")
		}

		if value == nil {
			return fmt.Errorf("param value cannot be nil")
		}

		if opts.params == nil {
			opts.params = runtime.NewParams()
		}

		parsed, err := runtime.ValueOf(value)

		if err != nil {
			return fmt.Errorf("invalid param value: %w", err)
		}

		opts.params.SetValue(name, parsed)

		return nil
	}
}

// WithRuntimeParam returns an Option that sets a runtime parameter with the specified name and value in the options configuration.
// The name cannot be empty, and the value cannot be nil.
func WithRuntimeParam(name string, value runtime.Value) Option {
	return func(opts *config) error {
		if name == "" {
			return fmt.Errorf("param name cannot be empty")
		}

		if value == nil {
			return fmt.Errorf("param value cannot be nil")
		}

		if opts.params == nil {
			opts.params = runtime.NewParams()
		}

		opts.params.SetValue(name, value)

		return nil
	}
}

// WithNamespace merges the functions from the provided runtime.Namespace into the engine's function library.
func WithNamespace(ns runtime.Namespace) Option {
	return func(opts *config) error {
		if ns == nil {
			return fmt.Errorf("namespace cannot be nil")
		}

		opts.library.Function().From(ns.Function())

		return nil
	}
}

// WithFunctionsRegistrar creates an Option that invokes the provided registrar with the engine's runtime.Namespace if the registrar is not nil.
// Registered host-function names and namespace segments are canonicalized to lowercase and resolve case-insensitively in FQL.
func WithFunctionsRegistrar(setter func(ns runtime.Namespace)) Option {
	return func(env *config) error {
		if setter == nil {
			return fmt.Errorf("functions registrar cannot be nil")
		}

		setter(env.library)

		return nil
	}
}

// WithFunctions merges the provided *runtime.Functions into the engine's function library.
func WithFunctions(funcs *runtime.Functions) Option {
	return func(opts *config) error {
		if funcs == nil {
			return fmt.Errorf("functions cannot be nil")
		}

		opts.library.Function().From(runtime.NewFunctionsBuilderFrom(funcs))

		return nil
	}
}

// WithLog sets the writer for logging output.
// The writer can be any io.Writer, such as os.Stdout or a file.
func WithLog(writer io.Writer) Option {
	return func(opts *config) error {
		if writer == nil {
			return fmt.Errorf("log writer cannot be nil")
		}

		opts.logger = append(opts.logger, logging.WithWriter(writer))

		return nil
	}
}

// WithLogLevel sets the logging level for the engine.
// The logging level determines the severity of log messages that will be recorded.
func WithLogLevel(lvl logging.LogLevel) Option {
	return func(opts *config) error {
		if lvl < logging.TraceLevel || lvl > logging.Disabled {
			return fmt.Errorf("invalid log level: %v", lvl)
		}

		opts.logger = append(opts.logger, logging.WithLevel(lvl))

		return nil
	}
}

// WithLogFields sets the fields to be included in log entries.
// These fields can provide additional context for debugging and monitoring purposes.
func WithLogFields(fields map[string]any) Option {
	return func(opts *config) error {
		if len(fields) == 0 {
			return nil
		}

		opts.logger = append(opts.logger, logging.WithFields(fields))

		return nil
	}
}

// WithEncodingRegistry sets a custom encoding registry for query execution.
func WithEncodingRegistry(registry *encoding.Registry) Option {
	return gooptions.New(func(opts *config, registry *encoding.Registry) {
		opts.encoding = registry
	}).
		Value(registry).
		Named("encoding registry").
		Validators(gooptions.NotNilPtr[encoding.Registry]()).
		Build()
}

// WithProgramLoader sets a custom artifact loader for Engine.Load.
func WithProgramLoader(loader *artifact.Loader) Option {
	return gooptions.New(func(opts *config, loader *artifact.Loader) {
		opts.programLoader = loader
	}).
		Value(loader).
		Named("program loader").
		Validators(gooptions.NotNilPtr[artifact.Loader]()).
		Build()
}

// WithoutStdlib disables the standard library, so no built-in functions are registered by default.
func WithoutStdlib() Option {
	return WithStdlib(stdlib.Empty())
}

// WithStdlib configures which standard library groups are registered by default.
func WithStdlib(set stdlib.Set) Option {
	return gooptions.New(func(opts *config, set stdlib.Set) {
		opts.stdlib = set
	}).
		Value(set).
		Named("stdlib").
		Build()
}

// WithModules creates an Option that appends the provided modules to the options if not empty.
func WithModules(mods ...module.Module) Option {
	return func(env *config) error {
		if len(mods) == 0 {
			return nil
		}

		if env.modules == nil {
			env.modules = make([]module.Module, 0, len(mods))
		}

		for _, m := range mods {
			if m == nil {
				return fmt.Errorf("module cannot be nil")
			}

			env.modules = append(env.modules, m)
		}

		return nil
	}
}

// WithEncodingCodec registers or overrides a codec for the given content type.
func WithEncodingCodec(contentType string, codec encoding.Codec) Option {
	return func(opts *config) error {
		if codec == nil {
			return encoding.ErrNilCodec
		}

		if opts.encoding == nil {
			opts.encoding = encoding.NewRegistry()
		}

		return opts.encoding.Register(encodingCodecAlias{
			Codec:       codec,
			contentType: contentType,
		})
	}
}

// WithCompilerOptions creates an Option that appends the provided compiler options to the options if not empty.
func WithCompilerOptions(opts ...compiler.Option) Option {
	return func(o *config) error {
		if len(opts) == 0 {
			return nil
		}

		if o.compiler == nil {
			o.compiler = make([]compiler.Option, 0, len(opts))
		}

		for _, opt := range opts {
			if opt == nil {
				return fmt.Errorf("compiler option cannot be nil")
			}

			o.compiler = append(o.compiler, opt)
		}

		return nil
	}
}

// WithEngineInitHook returns an Option that registers a hook to execute during engine initialization.
// It returns an error if hook is nil.
func WithEngineInitHook(hook module.EngineInitHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("engine init hook is nil")
		}

		opts.hooks.engine.OnInit(hook)

		return nil
	}
}

// WithEngineCloseHook returns an Option that registers a hook to execute when the engine is closed.
// It returns an error if hook is nil.
func WithEngineCloseHook(hook module.EngineCloseHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("engine close hook is nil")
		}

		opts.hooks.engine.OnClose(hook)

		return nil
	}
}

// WithBeforeCompileHook returns an Option that registers a hook to execute before each compilation attempt.
// It returns an error if hook is nil.
func WithBeforeCompileHook(hook module.BeforeCompileHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("before compile hook is nil")
		}

		opts.hooks.plan.BeforeCompile(hook)

		return nil
	}
}

// WithAfterCompileHook returns an Option that registers a hook to execute after each compilation attempt.
// The hook receives the compilation error (if any). It returns an error if hook is nil.
func WithAfterCompileHook(hook module.AfterCompileHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("after compile hook is nil")
		}

		opts.hooks.plan.AfterCompile(hook)

		return nil
	}
}

// WithPlanCloseHook returns an Option that registers a hook to execute when a plan is closed.
// It returns an error if hook is nil.
func WithPlanCloseHook(hook module.PlanCloseHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("plan close hook is nil")
		}

		opts.hooks.plan.OnClose(hook)

		return nil
	}
}

// WithBeforeRunHook returns an Option that registers a hook to execute before each session run.
// The hook can replace the context used by subsequent hooks and VM execution.
// It returns an error if hook is nil.
func WithBeforeRunHook(hook module.BeforeRunHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("before run hook is nil")
		}

		opts.hooks.session.BeforeRun(hook)

		return nil
	}
}

// WithAfterRunHook returns an Option that registers a hook to execute after each session run attempt.
// The hook receives the run error (if any). It returns an error if hook is nil.
func WithAfterRunHook(hook module.AfterRunHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("after run hook is nil")
		}

		opts.hooks.session.AfterRun(hook)

		return nil
	}
}

// WithSessionCloseHook returns an Option that registers a hook to execute when a session is closed.
// It returns an error if hook is nil.
func WithSessionCloseHook(hook module.SessionCloseHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("session close hook is nil")
		}

		opts.hooks.session.OnClose(hook)

		return nil
	}
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
	return gooptions.New(func(opts *config, n int) {
		opts.maxActiveSessions = n
	}).
		Value(n).
		Named("max active sessions").
		Validators(gooptions.NonNegative[int]()).
		Build()
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
	return gooptions.New(func(opts *config, n int) {
		opts.maxIdleVMsPerPlan = n
	}).
		Value(n).
		Named("max idle VMs per plan").
		Validators(gooptions.NonNegative[int]()).
		Build()
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
	return gooptions.New(func(opts *config, n int) {
		opts.maxVMsPerPlan = n
	}).
		Value(n).
		Named("max VMs per plan").
		Validators(gooptions.NonNegative[int]()).
		Build()
}

// WithFSRoot sets the root directory for the engine's file system.
func WithFSRoot(root string) Option {
	return gooptions.New(func(opts *config, root string) {
		opts.fsRoot = strings.TrimSpace(root)
	}).
		Value(root).
		Named("fs root").
		Validators(gooptions.NotBlank[string]()).
		Build()
}

// WithFSReadOnly sets the engine's file system to read-only mode.
func WithFSReadOnly() Option {
	return gooptions.New(func(opts *config, readOnly bool) {
		opts.fsReadOnly = readOnly
	}).
		Value(true).
		Build()
}

// WithNetwork sets the engine network service used by derived executions.
// If a network is provided, the engine will use it directly and will not manage its lifecycle.
// The host application is responsible for closing the network when it is no longer needed.
func WithNetwork(network ferretnet.Network) Option {
	return func(opts *config) error {
		if network == nil {
			return fmt.Errorf("network cannot be nil")
		}

		opts.network = network
		opts.hostNetwork = true

		return nil
	}
}

// WithNetworkOptions creates an Option that constructs a new network service using the provided Ferret network options.
// If no options are provided, the engine will use the default network service.
func WithNetworkOptions(setters ...ferretnet.Option) Option {
	return func(opts *config) error {
		if len(setters) == 0 {
			return nil
		}

		net, err := ferretnet.New(setters...)

		if err != nil {
			return fmt.Errorf("create network: %w", err)
		}

		opts.network = net
		opts.managedNetworks = append(opts.managedNetworks, net)
		opts.hostNetwork = false

		return nil
	}
}
