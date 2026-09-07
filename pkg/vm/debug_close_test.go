package vm

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestVMCloseReturnsAndCachesRetainedCleanupError(t *testing.T) {
	closeErr := errors.New("retained cleanup failed")
	closer := newFailingCloser(closeErr)
	instance := mustNewVM(t, sourcePointTestProgram())
	instance.state.deferred.AddCloser(closer)

	if err := instance.Close(); !errors.Is(err, closeErr) {
		t.Fatalf("expected retained cleanup error, got %v", err)
	}
	if err := instance.Close(); !errors.Is(err, closeErr) {
		t.Fatalf("expected cached retained cleanup error, got %v", err)
	}
	if got := closer.count(); got != 1 {
		t.Fatalf("expected one cleanup call, got %d", got)
	}
}

func TestDebugExecutionTerminationRetainsOnlyCleanupErrors(t *testing.T) {
	for _, mode := range []string{"cancel", "deadline", "source_point"} {
		for _, failCleanup := range []bool{false, true} {
			t.Run(mode+map[bool]string{false: "/clean", true: "/failed"}[failCleanup], func(t *testing.T) {
				var cleanupErr error
				if failCleanup {
					cleanupErr = errors.New("termination cleanup failed")
				}

				closer := newFailingCloser(cleanupErr)
				instance := mustNewVM(t, sourcePointTestProgram())
				execution, err := NewDebugExecution(instance, nil)
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = execution.Close() })
				if _, err := execution.Start(t.Context()); err != nil {
					t.Fatal(err)
				}

				instance.state.deferred.AddCloser(closer)
				ctx := t.Context()
				var cause error
				switch mode {
				case "cancel":
					var cancel context.CancelFunc
					ctx, cancel = context.WithCancel(ctx)
					cancel()
					cause = context.Canceled
				case "deadline":
					var cancel context.CancelFunc
					ctx, cancel = context.WithDeadline(ctx, time.Unix(1, 0))
					defer cancel()
					cause = context.DeadlineExceeded
				case "source_point":
					instance.sourcePointObserver = &recordingSourcePointObserver{action: sourcePointTerminate}
				}

				event, err := execution.Resume(ctx, DebugResumeContinue, nil)
				if err != nil || event == nil || event.Reason != DebugStopTerminated {
					t.Fatalf("termination: event=%+v err=%v", event, err)
				}

				for _, want := range []error{cause, cleanupErr} {
					if want != nil && !errors.Is(event.Error, want) {
						t.Fatalf("termination lost %v: %v", want, event.Error)
					}
				}

				if cause == nil && cleanupErr == nil && event.Error != nil {
					t.Fatalf("unexpected termination error: %v", event.Error)
				}

				if closer.count() != 1 {
					t.Fatal("termination did not clean up before returning")
				}

				for range 2 {
					closeErr := execution.Close()
					if !errors.Is(closeErr, cleanupErr) || (cause != nil && errors.Is(closeErr, cause)) {
						t.Fatalf("close must retain only cleanup errors: %v", closeErr)
					}
				}

				if closer.count() != 1 {
					t.Fatal("close repeated termination cleanup")
				}
			})
		}
	}
}

func TestDebugExecutionCloseReturnsAndCachesVMCloseError(t *testing.T) {
	closeErr := errors.New("debug cleanup failed")
	closer := newFailingCloser(closeErr)
	instance := mustNewVM(t, sourcePointTestProgram())
	execution, err := NewDebugExecution(instance, nil)
	if err != nil {
		t.Fatal(err)
	}
	instance.state.deferred.AddCloser(closer)

	if err := execution.Close(); !errors.Is(err, closeErr) {
		t.Fatalf("expected VM close error, got %v", err)
	}
	if err := execution.Close(); !errors.Is(err, closeErr) {
		t.Fatalf("expected cached VM close error, got %v", err)
	}
	if got := closer.count(); got != 1 {
		t.Fatalf("expected one cleanup call, got %d", got)
	}
}
