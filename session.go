package ferret

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type (
	sessionPermitRelease func(*vm.VM)

	// Session represents the execution of a compiled Ferret program.
	// It holds the state of the execution, including the virtual machine, environment, and encoding registry.
	// A Session is created from a Plan and can be run to obtain results.
	//
	// Session is not safe for concurrent use by multiple goroutines, except that
	// Close is idempotent and safe to call multiple times, including concurrently.
	// It is typically used for a single logical execution. When a Session is created
	// directly via Plan.NewSession, it may be reused for multiple sequential runs as
	// long as the environment and encoding registry are not modified between runs.
	// Helper APIs such as Engine.Run may take ownership of the Session and close it
	// after a single execution, in which case the caller must not attempt to reuse it.
	Session struct {
		logger            logging.Logger
		closeErr          error
		network           ferretnet.Network
		hooks             sessionHooks
		fs                fs.FileSystem
		env               *vm.Environment
		encoding          *encoding.Registry
		vm                *vm.VM
		release           sessionPermitRelease
		resources         *resource.Manager
		outputContentType string
		closeOnce         sync.Once
		closed            atomic.Bool
	}
)

// Run executes the session with the provided context and returns encoded output.
// Successful output remains available alongside after-run hook or result-cleanup
// errors. The caller owns the encoded data; Run releases the underlying VM result.
func (s *Session) Run(c context.Context) (*Output, error) {
	if c == nil {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	if err := c.Err(); err != nil {
		return nil, err
	}

	if s.closed.Load() {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "session is closed")
	}

	// Before-run hooks can replace the context used for the rest of execution.
	ctx, err := s.hooks.runBeforeRunHooks(c)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("before run hooks: %w", err), c.Err())
	}

	err = c.Err()
	if ctx == nil {
		ctx = c
		err = errors.Join(err, runtime.Error(runtime.ErrInvalidArgument, "before run hooks returned nil context"))
	} else if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
		err = errors.Join(err, ctxErr)
	}

	var out *vm.Result

	if err == nil {
		out, err = s.vm.Run(s.extendContext(ctx), s.env)
	}

	// Successful before-run hooks must be paired even when context validation
	// prevents VM entry. Every after-run hook receives the same primary error.
	hookErr := s.hooks.runAfterRunHooks(ctx, err)
	if hookErr != nil {
		hookErr = fmt.Errorf("after run hooks: %w", hookErr)
	}

	if err != nil {
		if hookErr != nil {
			return nil, errors.Join(err, hookErr)
		}

		return nil, err
	}

	output, outputErr := newOutput(s.encoding, s.outputContentType, out)
	closeErr := out.Close()

	if outputErr != nil {
		return nil, errors.Join(hookErr, outputErr, closeErr)
	}

	if hookErr != nil {
		return output, errors.Join(hookErr, closeErr)
	}

	return output, closeErr
}

// Close releases the session's borrowed VM, runs close hooks, and closes any
// filesystem created specifically for the session.
// It is idempotent and safe to call multiple times, including concurrently.
func (s *Session) Close() error {
	if s == nil {
		return nil
	}

	s.closeOnce.Do(func() {
		s.closed.Store(true)

		instance := s.vm
		release := s.release
		resources := s.resources

		s.vm = nil
		s.release = nil
		s.fs = nil
		s.resources = nil

		if hookErr := s.hooks.runCloseHooks(); hookErr != nil {
			s.closeErr = fmt.Errorf("close hooks: %w", hookErr)
		}

		if closeErr := resources.Close(); closeErr != nil {
			s.closeErr = errors.Join(s.closeErr, closeErr)
		}

		if release != nil && instance != nil {
			// Returning the borrowed VM is best-effort cleanup and must still happen when
			// close hooks fail so the pool/limiter do not leak capacity.
			release(instance)
		}
	})

	return s.closeErr
}

func (s *Session) extendContext(ctx context.Context) context.Context {
	ctx = s.logger.WithContext(ctx)
	ctx = encoding.WithRegistry(ctx, s.encoding)
	ctx = fs.WithFileSystem(ctx, s.fs)
	ctx = ferretnet.WithNetwork(ctx, s.network)

	return ctx
}
