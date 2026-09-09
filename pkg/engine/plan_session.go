package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	enginesession "github.com/MontFerret/ferret/v2/pkg/engine/internal/session"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type (
	planSessionSetup struct {
		requiresDebugInfo bool
	}

	planSessionDependencies struct {
		logger     logging.Logger
		hooks      *host.SessionHooks
		filesystem fs.FileSystem
		program    *bytecode.Program
		host       *host.Host
		limiter    *enginesession.Limiter
		pool       *vm.Pool
		resources  *resource.Manager
		options    sessionConfig
	}
)

func newPlanSession[T interface{ Close() error }](
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
		return session, runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	if err := ctx.Err(); err != nil {
		return session, err
	}

	select {
	case <-plan.closed:
		return session, runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
	default:
	}

	program := plan.prog
	h := plan.host
	hooks := plan.sessionHooks
	limiter := plan.limiter
	pool := plan.pool

	if setup.requiresDebugInfo && len(program.Metadata.DebugPoints) == 0 {
		return session, runtime.Error(runtime.ErrInvalidOperation, "plan was not compiled for debugging")
	}

	options, err := newSessionConfig(setters)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
			err = errors.Join(err, ctxErr)
		}

		return session, err
	}

	if err := ctx.Err(); err != nil {
		return session, err
	}

	logger, err := logging.NewFrom(h.Logger, options.logger...)
	if err != nil {
		return session, fmt.Errorf("logger: %w", err)
	}

	if err = limiter.Acquire(ctx, plan.closed); err != nil {
		return session, err
	}

	transferred := false
	defer func() {
		// Construction errors and panics retain ownership here. Successful
		// construction transfers the permit release to the returned session.
		// Keep this defer separate so resource cleanup cannot skip permit release.
		if !transferred {
			limiter.Release()
		}
	}()

	resources := resource.NewManager()
	defer func() {
		if !transferred {
			if closeErr := resources.Close(); closeErr != nil {
				err = errors.Join(err, closeErr)
			}
		}
	}()

	if err := resources.Borrow(resource.FileSystem); err != nil {
		return session, err
	}

	if err := resources.Borrow(resource.Network); err != nil {
		return session, err
	}

	filesystem := h.FileSystem
	if options.fsRoot != "" {
		filesystem, err = fs.New(fs.WithRoot(options.fsRoot), fs.WithReadOnly(h.FSReadOnly))
		if err != nil {
			return session, fmt.Errorf("filesystem: %w", err)
		}

		if err := resources.Own(resource.FileSystem, filesystem.Close); err != nil {
			return session, err
		}
	}

	session, err = build(planSessionDependencies{
		program:    program,
		host:       h,
		hooks:      hooks,
		limiter:    limiter,
		pool:       pool,
		filesystem: filesystem,
		options:    options,
		logger:     logger,
		resources:  resources,
	})

	if err == nil {
		// The constructed session now owns the permit/host resources. Only request
		// cancellation rolls it back; parent close does not revoke ownership.
		transferred = true

		if err = ctx.Err(); err != nil {
			if closeErr := session.Close(); closeErr != nil {
				err = errors.Join(err, closeErr)
			}

			var zero T

			return zero, err
		}
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
		fs:                dependencies.filesystem,
		network:           dependencies.host.Network,
		encoding:          dependencies.host.Encoding,
		outputContentType: dependencies.options.outputContentType,
		hooks:             dependencies.hooks,
		release:           enginesession.NewPermitRelease(dependencies.limiter, dependencies.pool),
		resources:         dependencies.resources,
	}, nil
}

func buildDebugSession(dependencies planSessionDependencies) (*debugger.Session, error) {
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
		Services: enginesession.NewDebugServices(enginesession.DebugServicesConfig{
			Hooks:             dependencies.hooks,
			Encoding:          dependencies.host.Encoding,
			OutputContentType: dependencies.options.outputContentType,
			Logger:            dependencies.logger,
			FileSystem:        dependencies.filesystem,
			Network:           dependencies.host.Network,
		}, dependencies.limiter, dependencies.resources),
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

func newPlanSessionEnvironment(h *host.Host, options sessionConfig) (*vm.Environment, error) {
	return vm.ExtendEnvironment(&vm.Environment{
		Functions: h.Functions,
		Params:    h.Params,
	}, options.env)
}
