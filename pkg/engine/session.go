package engine

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	enginesession "github.com/MontFerret/ferret/v2/pkg/engine/internal/session"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Session represents the execution of a compiled Ferret program.
// It owns native run admission and hooks, with operational state held by its execution.
// A Session is created from a Plan and can be run to obtain results.
//
// Session is not safe for concurrent use by multiple goroutines, except that
// Close is idempotent and safe to call multiple times, including concurrently.
// It is typically used for a single logical execution. When a Session is created
// directly via Plan.NewSession, it may be reused for multiple sequential runs as
// long as the environment and encoding registry are not modified between runs.
// Helper APIs such as Engine.Run may take ownership of the Session and close it
// after a single execution, in which case the caller must not attempt to reuse it.
type Session struct {
	execution *enginesession.Execution
	hooks     *host.SessionHooks
	closed    atomic.Bool
}

// Run executes the session with the provided context and returns encoded output.
// Successful output remains available alongside after-run hook or result-cleanup
// errors. The caller owns the encoded data; Run releases the underlying VM result.
func (s *Session) Run(c context.Context) (*encoding.Output, error) {
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
	ctx, err := s.hooks.RunBeforeRun(c)
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
		out, err = s.execution.Run(ctx)
	}

	// Successful before-run hooks must be paired even when context validation
	// prevents VM entry. Every after-run hook receives the same primary error.
	hookErr := s.hooks.RunAfterRun(ctx, err)
	if hookErr != nil {
		hookErr = fmt.Errorf("after run hooks: %w", hookErr)
	}

	if err != nil {
		if hookErr != nil {
			return nil, errors.Join(err, hookErr)
		}

		return nil, err
	}

	output, outputErr, closeErr := s.execution.MaterializeAndClose(out)

	if outputErr != nil {
		return nil, errors.Join(hookErr, outputErr, closeErr)
	}

	if hookErr != nil {
		return output, errors.Join(hookErr, closeErr)
	}

	return output, closeErr
}

// Close rejects further runs, runs close hooks, closes session-owned host
// resources, and returns the session permit and borrowed VM.
// It is idempotent and safe to call multiple times, including concurrently.
func (s *Session) Close() error {
	if s == nil {
		return nil
	}

	s.closed.Store(true)

	return s.execution.Close()
}
