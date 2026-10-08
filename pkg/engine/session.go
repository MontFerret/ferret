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
type (
	Session struct {
		execution *enginesession.Execution
		hooks     *host.SessionHooks
		closed    atomic.Bool
	}

	// runOutcome contains only detached representation and observed errors.
	// It never owns the caller's session or the materialized VM result.
	runOutcome struct {
		content      *encoding.Content
		operationErr error
		cleanupErr   error
	}
)

// Run executes and encodes eagerly, returning a caller-owned, one-shot handle.
// A nil Run error means execution was admitted; Consume or Collect observes
// terminal execution, encoding, hook, and cleanup errors, with available content.
// Keep c alive through consumption. Settle the output before reusing the session.
func (s *Session) Run(c context.Context) (encoding.Output, error) {
	outcome, err := s.execute(c)
	if err != nil {
		return nil, err
	}

	return newOutput(c, outcome), nil
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

func (s *Session) execute(c context.Context) (runOutcome, error) {
	if c == nil {
		return runOutcome{}, runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	if err := c.Err(); err != nil {
		return runOutcome{}, err
	}

	if s.closed.Load() {
		return runOutcome{}, runtime.Error(runtime.ErrInvalidOperation, "session is closed")
	}

	// Before-run hooks can replace the context used for the rest of execution.
	ctx, err := s.hooks.RunBeforeRun(c)
	if err != nil {
		return runOutcome{}, errors.Join(fmt.Errorf("before run hooks: %w", err), c.Err())
	}

	err = c.Err()
	if ctx == nil {
		ctx = c
		err = errors.Join(err, runtime.Error(runtime.ErrInvalidArgument, "before run hooks returned nil context"))
	} else if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
		err = errors.Join(err, ctxErr)
	}

	if err != nil {
		hookErr := s.hooks.RunAfterRun(ctx, err)
		if hookErr != nil {
			hookErr = fmt.Errorf("after run hooks: %w", hookErr)
		}

		return runOutcome{}, errors.Join(err, hookErr)
	}

	// VM entry is the admission boundary. From here failures belong to output.
	out, err := s.execution.Run(ctx)

	// Hooks and encoders may panic after the VM has transferred result ownership.
	// Retain rollback until materialization closes that result normally.
	defer func() {
		if out != nil {
			_ = out.Close()
		}
	}()

	// Every after-run hook receives the same primary execution error.
	hookErr := s.hooks.RunAfterRun(ctx, err)
	if hookErr != nil {
		hookErr = fmt.Errorf("after run hooks: %w", hookErr)
	}

	if err != nil {
		return runOutcome{operationErr: joinOutputErrors(err, hookErr)}, nil
	}

	// Before-run replacement contexts remain execution-only. Encoding belongs to
	// the original invocation lifetime, just like the resulting output handle.
	content, encodeErr, closeErr := s.execution.MaterializeAndClose(c, out)
	out = nil

	return runOutcome{content: content, operationErr: joinOutputErrors(hookErr, encodeErr), cleanupErr: closeErr}, nil
}
