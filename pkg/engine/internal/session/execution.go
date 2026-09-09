package session

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Execution owns an ordinary session's borrowed VM, host resources, and permit.
// The native Session owns run admission and hook orchestration. Execution must
// be used serially, except that concurrent Close calls retain the same result.
type Execution struct {
	logger            logging.Logger
	closeErr          error
	network           ferretnet.Network
	hooks             *host.SessionHooks
	filesystem        fs.FileSystem
	environment       *vm.Environment
	encoding          *encoding.Registry
	vm                *vm.VM
	release           PermitRelease
	resources         *resource.Manager
	outputContentType string
	closeOnce         sync.Once
}

// NewExecution acquires session resources and a pooled VM. Failure rolls back
// without close hooks; success transfers cleanup to Execution. The native
// caller must check cancellation before publishing and close on cancellation.
func NewExecution(ctx context.Context, config Config, h *host.Host, hooks *host.SessionHooks, limiter *Limiter, pool *vm.Pool, planClosed <-chan struct{}) (_ *Execution, err error) {
	a := acquisition{limiter: limiter}
	defer a.rollback(&err)

	if err := a.prepare(ctx, config, h, planClosed); err != nil {
		return nil, err
	}

	return a.newExecution(config, h, hooks, pool)
}

// Run invokes the VM with session services. The caller must settle run hooks
// before passing a successful result to MaterializeAndClose.
func (e *Execution) Run(ctx context.Context) (*vm.Result, error) {
	return e.vm.Run(e.extendContext(ctx), e.environment)
}

// MaterializeAndClose encodes and releases a successful VM result. Separate
// errors let the native Session preserve hook, encoding, and cleanup ordering.
func (e *Execution) MaterializeAndClose(result *vm.Result) (output *encoding.Output, materializeErr, closeErr error) {
	output, materializeErr = Materialize(e.encoding, e.outputContentType, result)
	closeErr = result.Close()

	return output, materializeErr, closeErr
}

// Close runs close hooks, closes owned host resources, and returns the permit
// and borrowed VM even when cleanup returns errors. Repeated and concurrent
// calls return the same cleanup result. Close must not overlap execution.
func (e *Execution) Close() error {
	e.closeOnce.Do(func() {
		instance := e.vm
		release := e.release
		resources := e.resources

		e.vm = nil
		e.release = nil
		e.filesystem = nil
		e.resources = nil

		if hookErr := e.hooks.RunClose(); hookErr != nil {
			e.closeErr = fmt.Errorf("close hooks: %w", hookErr)
		}

		if closeErr := resources.Close(); closeErr != nil {
			e.closeErr = errors.Join(e.closeErr, closeErr)
		}

		if release != nil && instance != nil {
			release(instance)
		}
	})

	return e.closeErr
}

func (e *Execution) extendContext(ctx context.Context) context.Context {
	ctx = e.logger.WithContext(ctx)
	ctx = encoding.WithRegistry(ctx, e.encoding)
	ctx = fs.WithFileSystem(ctx, e.filesystem)

	return ferretnet.WithNetwork(ctx, e.network)
}
