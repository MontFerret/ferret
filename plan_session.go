package ferret

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type (
	planSessionSetup struct {
		requiresDebugInfo bool
	}

	planSessionDependencies struct {
		program *bytecode.Program
		host    *host
		hooks   sessionHooks
		limiter *sessionLimiter
		pool    *vm.Pool
		options *sessionOptions
		logger  logging.Logger
	}
)

func newPlanSession[T any](
	plan *Plan,
	ctx context.Context,
	setters []SessionOption,
	setup planSessionSetup,
	build func(planSessionDependencies) (T, error),
) (session T, err error) {
	if plan == nil {
		return session, runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	plan.mu.RLock()
	if plan.closed {
		plan.mu.RUnlock()
		return session, runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
	}

	program := plan.prog
	h := plan.host
	hooks := plan.sessionHooks
	limiter := plan.limiter
	pool := plan.pool
	plan.mu.RUnlock()

	if setup.requiresDebugInfo && len(program.Metadata.DebugPoints) == 0 {
		return session, runtime.Error(runtime.ErrInvalidOperation, "plan was not compiled for debugging")
	}

	options, err := newSessionOptions(setters)
	if err != nil {
		return session, err
	}

	if err = limiter.Acquire(ctx); err != nil {
		return session, err
	}

	releaseOnReturn := true
	defer func() {
		// Construction errors and panics retain ownership here. Successful
		// construction transfers the permit release to the returned session.
		if releaseOnReturn {
			limiter.Release()
		}
	}()

	session, err = build(planSessionDependencies{
		program: program,
		host:    h,
		hooks:   hooks,
		limiter: limiter,
		pool:    pool,
		options: options,
		logger:  logging.NewFrom(h.logger, options.logger...),
	})
	if err == nil {
		releaseOnReturn = false
	}

	return session, err
}

func buildSession(dependencies planSessionDependencies) (*Session, error) {
	environment, err := newPlanSessionEnvironment(dependencies.host, dependencies.options)
	if err != nil {
		return nil, err
	}

	instance, err := dependencies.pool.Acquire()
	if err != nil {
		if errors.Is(err, vm.ErrPoolClosed) {
			return nil, runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
		}

		return nil, err
	}

	return &Session{
		vm:                instance,
		env:               environment,
		logger:            dependencies.logger,
		fs:                dependencies.host.fs,
		network:           dependencies.host.network,
		encoding:          dependencies.host.encoding,
		outputContentType: dependencies.options.outputContentType,
		hooks:             dependencies.hooks,
		release:           newSessionPermitRelease(dependencies.limiter, dependencies.pool),
	}, nil
}

func buildDebugSession(dependencies planSessionDependencies) (*DebugSession, error) {
	environment, err := newPlanSessionEnvironment(dependencies.host, dependencies.options)
	if err != nil {
		return nil, err
	}

	instance, err := vm.New(dependencies.program)
	if err != nil {
		return nil, err
	}

	execution, err := vm.NewDebugExecution(instance, environment)
	if err != nil {
		_ = instance.Close()
		return nil, err
	}

	session, err := debugger.NewSession(debugger.Config{
		Execution: execution,
		Values:    vm.NewDebugValueAccess(),
		Services: &debugSessionServices{
			hooks:             dependencies.hooks,
			releasePermit:     newSessionPermitRelease(dependencies.limiter, nil),
			encoding:          dependencies.host.encoding,
			outputContentType: dependencies.options.outputContentType,
			logger:            dependencies.logger,
			fs:                dependencies.host.fs,
			network:           dependencies.host.network,
		},
		Source:      dependencies.program.Source,
		DebugPoints: dependencies.program.Metadata.DebugPoints,
		Params:      dependencies.program.Params,
		Format:      dependencies.options.debugFormat,
	})
	if err != nil {
		_ = execution.Close()
		return nil, err
	}

	return session, nil
}

func newPlanSessionEnvironment(h *host, options *sessionOptions) (*vm.Environment, error) {
	return vm.ExtendEnvironment(&vm.Environment{
		Functions: h.functions,
		Params:    h.params,
	}, options.env)
}
