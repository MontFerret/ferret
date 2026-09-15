package debugger

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// sessionLifecycle owns execution and embedding resources. commandMu protects
// execution, inspection, hooks, and cleanup. lifecycleMu orders cancellation,
// native terminal commitment, and breakpoint publication; it never waits for
// commandMu or the breakpoint writer gate.
type (
	sessionLifecycle struct {
		services        SessionServices
		execution       vm.DebugExecution
		closeErr        error
		runCtx          context.Context
		executionCtx    context.Context
		runCancel       context.CancelCauseFunc
		activeCancel    context.CancelCauseFunc
		terminalDone    chan struct{}
		terminalState   string
		activeCancelID  uint64
		closeOnce       sync.Once
		commandMu       sync.Mutex
		lifecycleMu     sync.Mutex
		closed          atomic.Bool
		started         atomic.Bool
		afterRun        bool
		executionClosed bool
	}

	// resumeCommand keeps cancellation bookkeeping on the command stack.
	resumeCommand struct {
		ctx      context.Context
		cleanup  func()
		cancelID uint64
	}
)

func newSessionLifecycle(execution vm.DebugExecution, services SessionServices) sessionLifecycle {
	return sessionLifecycle{
		execution:    execution,
		services:     services,
		terminalDone: make(chan struct{}),
	}
}

func (s *sessionLifecycle) prepareStart(ctx context.Context) error {
	if err := checkContext(ctx); err != nil {
		return err
	}

	if err := s.ensureOpen(); err != nil {
		return err
	}

	if s.started.Load() {
		return &StateError{Operation: "start", State: "started"}
	}

	runCtx, err := s.services.BeforeRun(ctx)
	if err != nil {
		return errors.Join(fmt.Errorf("before run hooks: %w", err), ctx.Err())
	}

	err = checkContext(runCtx)
	if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
		err = errors.Join(err, ctxErr)
	}

	if err != nil {
		if runCtx == nil {
			runCtx = ctx
		}

		// Settle this attempt without consuming a later successful Start's hooks
		// or making Close repeat them. Execution has not been entered yet.
		if hookErr := s.services.AfterRun(runCtx, err); hookErr != nil {
			err = errors.Join(err, fmt.Errorf("after run hooks: %w", hookErr))
		}

		return err
	}

	executionCtx, runCancel := context.WithCancelCause(runCtx)

	s.started.Store(true)
	s.runCtx = runCtx
	s.executionCtx = executionCtx
	s.setRunCancel(runCancel)

	return nil
}

func (s *sessionLifecycle) start() (*vm.DebugExecutionEvent, error) {
	event, err := s.execution.Start(s.services.ExtendContext(s.executionCtx))
	if err != nil {
		s.recordExecutionTerminal()
	}

	return event, err
}

func (s *sessionLifecycle) prepareResume(ctx context.Context) (resumeCommand, error) {
	if err := s.ensureOpen(); err != nil {
		return resumeCommand{}, err
	}

	if !s.started.Load() {
		return resumeCommand{}, &StateError{Operation: "resume", State: "new"}
	}

	switch s.execution.Status() {
	case vm.DebugExecutionCompleted:
		return resumeCommand{}, &StateError{Operation: "resume", State: "completed"}
	case vm.DebugExecutionTerminated:
		return resumeCommand{}, &StateError{Operation: "resume", State: "terminated"}
	}

	if ctx == s.runCtx || ctx == s.executionCtx {
		return resumeCommand{ctx: s.executionCtx}, nil
	}

	ctx, cleanup, cancel := newResumeContext(s.runCtx, ctx)
	activeID := s.setActiveCancel(cancel)

	return resumeCommand{ctx: ctx, cleanup: cleanup, cancelID: activeID}, nil
}

func (s *sessionLifecycle) finishResume(command resumeCommand) {
	if command.cleanup != nil {
		s.clearActiveCancel(command.cancelID)
		command.cleanup()
	}
}

func (s *sessionLifecycle) resume(ctx context.Context, mode vm.DebugResumeMode, predicate vm.DebugBreakpointPredicate) (*vm.DebugExecutionEvent, error) {
	event, err := s.execution.Resume(s.services.ExtendContext(ctx), mode, predicate)
	if err != nil {
		s.recordExecutionTerminal()
	}

	return event, err
}

func (s *sessionLifecycle) pause() error {
	if err := s.ensureOpen(); err != nil {
		return err
	}

	if !s.started.Load() {
		return &StateError{Operation: "pause", State: "new"}
	}

	s.execution.RequestPause()

	return nil
}

func (s *sessionLifecycle) complete(result *vm.Result) (*encoding.Output, error) {
	output, outputErr := s.services.Materialize(result)
	closeErr := result.Close()
	hookErr := s.runAfterHooks(nil)

	return output, errors.Join(outputErr, closeErr, hookErr)
}

// terminateStop drains a stop overtaken by Close before publishing termination.
// Keep only cleanup failures in closeErr so repeated Close does not return the
// execution's cancellation or runtime failure.
func (s *sessionLifecycle) terminateStop(runErr error) error {
	err := context.Cause(s.executionCtx)
	if err == nil {
		err = context.Canceled
	}

	if runErr != nil {
		err = errors.Join(err, runErr)
	}

	s.executionClosed = true
	if closeErr := s.execution.Close(); closeErr != nil {
		s.closeErr = errors.Join(s.closeErr, closeErr)
		err = errors.Join(err, closeErr)
	}

	return s.settleRunError(err)
}

func (s *sessionLifecycle) settleRunError(runErr error) error {
	if hookErr := s.runAfterHooks(runErr); hookErr != nil {
		return errors.Join(runErr, hookErr)
	}

	return runErr
}

func (s *sessionLifecycle) runAfterHooks(runErr error) error {
	if s.afterRun {
		return nil
	}

	s.afterRun = true

	ctx := s.runCtx
	if s.executionCtx != nil {
		ctx = s.executionCtx
	}

	if hookErr := s.services.AfterRun(ctx, runErr); hookErr != nil {
		return fmt.Errorf("after run hooks: %w", hookErr)
	}

	return nil
}

func (s *sessionLifecycle) ensureOpen() error {
	if s == nil || s.closed.Load() {
		return &StateError{Operation: "use", State: "closed"}
	}

	return nil
}

func (s *sessionLifecycle) lockCommand() error {
	if s == nil || s.closed.Load() {
		return &StateError{Operation: "use", State: "closed"}
	}

	s.commandMu.Lock()
	if s.closed.Load() {
		s.commandMu.Unlock()
		return &StateError{Operation: "use", State: "closed"}
	}

	return nil
}

func (s *sessionLifecycle) setRunCancel(cancel context.CancelCauseFunc) {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()

	s.runCancel = cancel
	if s.closed.Load() && s.runCancel != nil {
		s.runCancel(context.Canceled)
	}
}

func (s *sessionLifecycle) setActiveCancel(cancel context.CancelCauseFunc) uint64 {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()

	s.activeCancelID++
	s.activeCancel = cancel
	if s.closed.Load() && s.activeCancel != nil {
		s.activeCancel(context.Canceled)
	}

	return s.activeCancelID
}

func (s *sessionLifecycle) clearActiveCancel(id uint64) {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()

	if s.activeCancelID == id {
		s.activeCancel = nil
	}
}

func (s *sessionLifecycle) requestTermination() {
	s.lifecycleMu.Lock()

	if s.runCancel != nil {
		s.runCancel(context.Canceled)
	}

	if s.activeCancel != nil {
		s.activeCancel(context.Canceled)
	}

	s.lifecycleMu.Unlock()

	if s.execution != nil {
		s.execution.RequestPause()
	}
}

func (s *sessionLifecycle) unlockCommand() {
	s.commandMu.Unlock()
}

// beginClose rejects new work and interrupts execution before waiting for the
// active command. The caller clears inspection state before finishClose.
func (s *sessionLifecycle) beginClose() {
	s.lifecycleMu.Lock()
	s.closed.Store(true)
	s.finishTerminalLocked("closed")
	s.lifecycleMu.Unlock()
	s.requestTermination()
	s.commandMu.Lock()
}

func (s *sessionLifecycle) finishClose() error {
	s.closeOnce.Do(func() {
		if s.started.Load() {
			s.closeErr = errors.Join(s.closeErr, s.runAfterHooks(context.Canceled))
		}

		if s.execution != nil && !s.executionClosed {
			s.executionClosed = true
			s.closeErr = errors.Join(s.closeErr, s.execution.Close())
		}

		if s.services != nil {
			s.closeErr = errors.Join(s.closeErr, s.services.Close())
		}
	})

	return s.closeErr
}

func (s *sessionLifecycle) finishTerminal(state string) {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()
	s.finishTerminalLocked(state)
}

func (s *sessionLifecycle) finishTerminalLocked(state string) {
	if s.terminalState == "" {
		s.terminalState = state
		if s.terminalDone != nil {
			close(s.terminalDone)
		}
	}
}

// Called only after execution returns; Status may acquire the VM run lock.
func (s *sessionLifecycle) recordExecutionTerminal() {
	switch s.execution.Status() {
	case vm.DebugExecutionCompleted:
		s.finishTerminal("completed")
	case vm.DebugExecutionTerminated, vm.DebugExecutionClosed:
		s.finishTerminal("terminated")
	}
}

func (s *sessionLifecycle) checkBreakpointContext(ctx context.Context, replacement bool) error {
	if err := checkContext(ctx); err != nil {
		return err
	}

	if err := s.ensureOpen(); err != nil {
		return err
	}

	if replacement {
		select {
		case <-s.terminalDone:
			s.lifecycleMu.Lock()
			state := s.terminalState
			s.lifecycleMu.Unlock()

			return &StateError{Operation: "replace breakpoints", State: state}
		default:
		}
	}

	return nil
}

// Success leaves lifecycleMu held until the caller publishes its snapshot.
func (s *sessionLifecycle) lockPublication(ctx context.Context, replacement bool) error {
	s.lifecycleMu.Lock()

	if err := checkContext(ctx); err != nil {
		s.lifecycleMu.Unlock()

		return err
	}

	if err := s.ensureOpen(); err != nil {
		s.lifecycleMu.Unlock()

		return err
	}

	if replacement && s.terminalState != "" {
		state := s.terminalState
		s.lifecycleMu.Unlock()

		return &StateError{Operation: "replace breakpoints", State: state}
	}

	return nil
}

func (s *sessionLifecycle) unlockPublication() {
	s.lifecycleMu.Unlock()
}
