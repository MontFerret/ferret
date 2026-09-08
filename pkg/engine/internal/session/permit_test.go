package session

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestPermitReleaseIsConcurrentSafeAndIdempotent(t *testing.T) {
	t.Parallel()

	compiler, err := compiler.New()
	if err != nil {
		t.Fatal(err)
	}

	program, err := compiler.Compile(t.Context(), source.NewAnonymous("RETURN 1"))
	if err != nil {
		t.Fatal(err)
	}

	pool := vm.NewPoolWithLimits(program, 1, 1)
	defer func() { _ = pool.Close() }()
	limiter := NewLimiter(1)
	closed := make(chan struct{})

	instance, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected direct pool acquire to succeed, got: %v", err)
	}

	if err := limiter.Acquire(context.Background(), closed); err != nil {
		t.Fatalf("expected limiter acquire to succeed, got: %v", err)
	}

	release := NewPermitRelease(limiter, pool)

	const callers = 8

	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < callers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			<-start
			release(instance)
		}()
	}

	close(start)
	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := limiter.Acquire(ctx, closed); err != nil {
		t.Fatalf("expected concurrent release to free the limiter slot, got: %v", err)
	}

	release(instance)
	if len(limiter.ch) != 1 {
		t.Fatal("repeated release consumed the next session's permit")
	}

	limiter.Release()

	reused, err := pool.Acquire()
	if err != nil {
		t.Fatalf("expected concurrent release to return the borrowed VM, got: %v", err)
	}

	if reused != instance {
		t.Fatal("expected concurrent release to retain the original VM in the pool")
	}

	_, err = pool.Acquire()
	if !errors.Is(err, vm.ErrPoolExhausted) {
		t.Fatalf("expected concurrent release not to free extra pool capacity, got: %v", err)
	}

	pool.Release(reused)
}
