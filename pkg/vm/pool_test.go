package vm

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func newPoolTestProgram() *bytecode.Program {
	return newTestProgram(
		1,
		nil,
		bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)
}

func TestPoolAcquireAllowsMoreActiveVMsThanIdleCap(t *testing.T) {
	t.Parallel()

	pool := NewPool(newPoolTestProgram(), 1)

	first, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected first acquire to succeed, got: %v", err)
	}

	second, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected second acquire to succeed beyond idle capacity, got: %v", err)
	}

	if first == second {
		t.Fatal("expected distinct VMs while no idle VM is available")
	}
}

func TestPoolReleaseRetainsOnlyIdleCap(t *testing.T) {
	t.Parallel()

	pool := NewPool(newPoolTestProgram(), 1)
	first, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected first acquire to succeed, got: %v", err)
	}

	second, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected second acquire to succeed, got: %v", err)
	}

	pool.Release(first)
	pool.Release(second)

	reusedA, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected first reuse acquire to succeed, got: %v", err)
	}

	reusedB, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected second reuse acquire to succeed, got: %v", err)
	}

	matches := 0
	if reusedA == first || reusedA == second {
		matches++
	}

	if reusedB == first || reusedB == second {
		matches++
	}

	if matches != 1 {
		t.Fatalf("expected exactly one borrowed VM to be retained, got %d", matches)
	}
}

func TestPoolZeroIdleCapacityDisablesReuse(t *testing.T) {
	t.Parallel()

	pool := NewPool(newPoolTestProgram(), 0)
	first, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected acquire to succeed, got: %v", err)
	}

	pool.Release(first)

	second, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected second acquire to succeed, got: %v", err)
	}

	if first == second {
		t.Fatal("expected zero idle capacity to disable VM reuse")
	}
}

func TestPoolAcquireRespectsTotalCapacity(t *testing.T) {
	t.Parallel()

	pool := NewPoolWithLimits(newPoolTestProgram(), 1, 1)
	first, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected first acquire to succeed, got: %v", err)
	}
	defer pool.Release(first)

	_, err = pool.Acquire()
	if !errors.Is(err, ErrPoolExhausted) {
		t.Fatalf("expected acquire at total cap to fail with ErrPoolExhausted, got: %v", err)
	}
}

func TestPoolReleaseIgnoresDuplicateReleaseOfIdleVM(t *testing.T) {
	t.Parallel()

	pool := NewPoolWithLimits(newPoolTestProgram(), 1, 1)
	instance, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected acquire to succeed, got: %v", err)
	}

	closer := newTrackingCloser("retained")
	instance.state.owned.Track(closer)

	pool.Release(instance)
	pool.Release(instance)

	if got := closer.closed; got != 0 {
		t.Fatalf("expected duplicate release not to close retained VM, got %d closes", got)
	}

	reused, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected retained VM to be reused, got: %v", err)
	}

	if reused != instance {
		t.Fatal("expected duplicate release not to replace the retained VM")
	}

	_, err = pool.Acquire()
	if !errors.Is(err, ErrPoolExhausted) {
		t.Fatalf("expected duplicate release not to free extra capacity, got: %v", err)
	}

	pool.Release(reused)
}

func TestPoolReleaseOfManuallyClosedVMIgnoresDuplicates(t *testing.T) {
	t.Parallel()

	pool := NewPoolWithLimits(newPoolTestProgram(), 1, 1)
	instance, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected acquire to succeed, got: %v", err)
	}

	closer := newTrackingCloser("manual-close")
	instance.state.owned.Track(closer)

	if err := instance.Close(); err != nil {
		t.Fatalf("expected manual close to succeed, got: %v", err)
	}

	if got := closer.closed; got != 1 {
		t.Fatalf("expected manual close to close VM once, got %d", got)
	}

	pool.Release(instance)
	pool.Release(instance)

	if got := closer.closed; got != 1 {
		t.Fatalf("expected duplicate release not to re-close VM, got %d", got)
	}

	replacement, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected capacity to be freed after releasing closed VM, got: %v", err)
	}

	_, err = pool.Acquire()
	if !errors.Is(err, ErrPoolExhausted) {
		t.Fatalf("expected duplicate release not to free extra capacity, got: %v", err)
	}

	pool.Release(replacement)
}

func TestPoolDuplicateReleaseOfClosedVMDoesNotBypassTotalCapacity(t *testing.T) {
	t.Parallel()

	pool := NewPoolWithLimits(newPoolTestProgram(), 0, 1)
	first, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected first acquire to succeed, got: %v", err)
	}

	closer := newTrackingCloser("dropped")
	first.state.owned.Track(closer)

	pool.Release(first)

	if got := closer.closed; got != 1 {
		t.Fatalf("expected zero-idle release to close VM once, got %d", got)
	}

	second, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected capacity to be available after first release, got: %v", err)
	}

	pool.Release(first)

	if got := closer.closed; got != 1 {
		t.Fatalf("expected duplicate release not to re-close VM, got %d", got)
	}

	_, err = pool.Acquire()
	if !errors.Is(err, ErrPoolExhausted) {
		t.Fatalf("expected duplicate release not to bypass total capacity, got: %v", err)
	}

	pool.Release(second)
}

func TestPoolReleaseWhenIdleIsFullClosesDroppedVMAndFreesCapacity(t *testing.T) {
	t.Parallel()

	pool := NewPoolWithLimits(newPoolTestProgram(), 1, 2)
	first, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected first acquire to succeed, got: %v", err)
	}

	second, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected second acquire to succeed, got: %v", err)
	}

	closer := newTrackingCloser("dropped")
	second.state.owned.Track(closer)

	pool.Release(first)
	pool.Release(second)

	if got, want := closer.closed, 1; got != want {
		t.Fatalf("expected dropped VM to be closed once, got %d", got)
	}

	reused, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected idle VM to be reused, got: %v", err)
	}

	fresh, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected freed total capacity to allow another acquire, got: %v", err)
	}

	if fresh == second {
		t.Fatal("expected dropped VM not to be reused after close")
	}

	pool.Release(reused)
	pool.Release(fresh)
}

func TestPoolCloseClosesRetainedIdleVMs(t *testing.T) {
	t.Parallel()

	pool := NewPool(newPoolTestProgram(), 1)
	instance, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected acquire to succeed, got: %v", err)
	}

	closer := newTrackingCloser("idle")
	instance.state.owned.Track(closer)
	pool.Release(instance)

	if err := pool.Close(); err != nil {
		t.Fatalf("expected pool close to succeed, got: %v", err)
	}

	if got, want := closer.closed, 1; got != want {
		t.Fatalf("expected idle VM close on pool shutdown, got %d", got)
	}
}

func TestPoolReleaseAfterCloseClosesBorrowedVM(t *testing.T) {
	t.Parallel()

	pool := NewPool(newPoolTestProgram(), 1)
	instance, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected acquire to succeed, got: %v", err)
	}

	closer := newTrackingCloser("borrowed")
	instance.state.owned.Track(closer)

	if err := pool.Close(); err != nil {
		t.Fatalf("expected pool close to succeed, got: %v", err)
	}

	pool.Release(instance)

	if got, want := closer.closed, 1; got != want {
		t.Fatalf("expected borrowed VM to close when returned after pool shutdown, got %d", got)
	}
}

func TestPoolCloseRejectsAcquireAndIgnoresRelease(t *testing.T) {
	t.Parallel()

	pool := NewPool(newPoolTestProgram(), 1)
	instance, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected acquire to succeed, got: %v", err)
	}

	if err := pool.Close(); err != nil {
		t.Fatalf("expected pool close to succeed, got: %v", err)
	}

	pool.Release(instance)

	_, err = pool.Acquire()
	if !errors.Is(err, ErrPoolClosed) {
		t.Fatalf("expected acquire on closed pool to fail with ErrPoolClosed, got: %v", err)
	}
}
