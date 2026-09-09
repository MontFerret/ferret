package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// acquisition retains rollback until it constructs the final execution owner.
// It is local to one constructor and never escapes to native engine code.
type acquisition struct {
	logger      logging.Logger
	filesystem  fs.FileSystem
	environment *vm.Environment
	resources   *resource.Manager
	limiter     *Limiter
	acquired    bool
	transferred bool
}

func (a *acquisition) prepare(ctx context.Context, config Config, h *host.Host, planClosed <-chan struct{}) error {
	logger, err := logging.NewFrom(h.Logger, config.Logger...)
	if err != nil {
		return fmt.Errorf("logger: %w", err)
	}

	a.logger = logger

	if err := a.limiter.Acquire(ctx, planClosed); err != nil {
		return err
	}

	a.acquired = true
	a.resources = resource.NewManager()

	if err := a.resources.Borrow(resource.FileSystem); err != nil {
		return err
	}

	if err := a.resources.Borrow(resource.Network); err != nil {
		return err
	}

	a.filesystem = h.FileSystem

	if config.FSRoot != "" {
		a.filesystem, err = fs.New(fs.WithRoot(config.FSRoot), fs.WithReadOnly(h.FSReadOnly))
		if err != nil {
			return fmt.Errorf("filesystem: %w", err)
		}

		if err := a.resources.Own(resource.FileSystem, a.filesystem.Close); err != nil {
			return err
		}
	}

	a.environment, err = vm.ExtendEnvironment(&vm.Environment{
		Functions: h.Functions,
		Params:    h.Params,
	}, config.Environment)

	return err
}

func (a *acquisition) newExecution(config Config, h *host.Host, hooks *host.SessionHooks, pool *vm.Pool) (*Execution, error) {
	instance, err := pool.Acquire()
	if err != nil {
		if errors.Is(err, vm.ErrPoolClosed) {
			return nil, runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
		}

		return nil, err
	}

	execution := &Execution{
		vm:                instance,
		environment:       a.environment,
		logger:            a.logger,
		filesystem:        a.filesystem,
		network:           h.Network,
		encoding:          h.Encoding,
		outputContentType: config.OutputContentType,
		hooks:             hooks,
		release:           NewPermitRelease(a.limiter, pool),
		resources:         a.resources,
	}
	a.transferred = true

	return execution, nil
}

func (a *acquisition) newDebugSession(config Config, h *host.Host, hooks *host.SessionHooks, program *bytecode.Program) (*debugger.Session, error) {
	instance, err := vm.New(program)
	if err != nil {
		return nil, err
	}

	execution, err := vm.NewDebugExecution(instance, a.environment)
	if err != nil {
		_ = instance.Close()

		return nil, err
	}

	session, err := debugger.NewSession(debugger.Config{
		Execution: execution,
		Values:    vm.NewDebugValueAccess(),
		Services: NewDebugServices(DebugServicesConfig{
			Hooks:             hooks,
			Encoding:          h.Encoding,
			OutputContentType: config.OutputContentType,
			Logger:            a.logger,
			FileSystem:        a.filesystem,
			Network:           h.Network,
		}, a.limiter, a.resources),
		Source:      program.Source,
		DebugPoints: program.Metadata.DebugPoints,
		Params:      program.Params,
		Format:      config.DebugFormat,
	})
	if err != nil {
		_ = execution.Close()

		return nil, err
	}

	a.transferred = true

	return session, nil
}

func (a *acquisition) rollback(err *error) {
	if !a.acquired || a.transferred {
		return
	}

	a.acquired = false

	// Keep permit release independent of resource cleanup, including its panics.
	defer a.limiter.Release()

	if a.resources != nil {
		if closeErr := a.resources.Close(); closeErr != nil {
			*err = errors.Join(*err, closeErr)
		}
	}
}
