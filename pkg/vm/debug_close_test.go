package vm

import (
	"errors"
	"testing"
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
