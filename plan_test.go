package ferret

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestPlanNewSessionReturnsEnvironmentError(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t)
	plan := mustCompilePlan(t, eng, coverageValidQuery)

	_, err := plan.NewSession(
		context.Background(),
		WithEnvironmentOptions(
			vm.WithFunction("SESSION_DUP_COVER", coverageVarFn),
			vm.WithFunction("SESSION_DUP_COVER", coverageVarFn),
		),
	)
	if err == nil {
		t.Fatal("expected plan.NewSession to fail on conflicting session options")
	}

	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected duplicate function error, got: %v", err)
	}
}

func TestPlanClose(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t)
	plan := mustCompilePlan(t, eng, coverageValidQuery)

	if err := plan.Close(); err != nil {
		t.Fatalf("expected plan close without hooks to succeed, got: %v", err)
	}

	closeErr := errors.New("plan close failed")
	eng = mustNewEngine(t, WithPlanCloseHook(func() error {
		return closeErr
	}))
	plan = mustCompilePlan(t, eng, coverageValidQuery)

	err := plan.Close()
	if err == nil {
		t.Fatal("expected plan close to fail when close hook fails")
	}

	if !errors.Is(err, closeErr) {
		t.Fatalf("expected close error to include plan close hook error, got: %v", err)
	}

	if !strings.Contains(err.Error(), "close hooks") {
		t.Fatalf("expected close hooks label, got: %v", err)
	}
}

func TestPlanNewSessionReleasesLimiterOnEnvironmentError(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithMaxActiveSessions(1))
	plan := mustCompilePlan(t, eng, coverageValidQuery)

	_, err := plan.NewSession(
		context.Background(),
		WithEnvironmentOptions(
			vm.WithFunction("SESSION_DUP_LIMIT", coverageVarFn),
			vm.WithFunction("SESSION_DUP_LIMIT", coverageVarFn),
		),
	)
	if err == nil {
		t.Fatal("expected plan.NewSession to fail on conflicting session options")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	session, err := plan.NewSession(ctx)
	if err != nil {
		t.Fatalf("expected limiter slot to be released after failed session creation, got: %v", err)
	}

	defer func() {
		_ = session.Close()
	}()
}

func TestPlanCloseIsIdempotentAndRejectsNewSessions(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t)
	plan := mustCompilePlan(t, eng, coverageValidQuery)

	if err := plan.Close(); err != nil {
		t.Fatalf("expected initial plan close to succeed, got: %v", err)
	}

	if err := plan.Close(); err != nil {
		t.Fatalf("expected repeated plan close to be idempotent, got: %v", err)
	}

	_, err := plan.NewSession(context.Background())
	if err == nil {
		t.Fatal("expected plan.NewSession to fail after plan close")
	}

	if !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("expected invalid operation after plan close, got: %v", err)
	}

	if !strings.Contains(err.Error(), "plan is closed") {
		t.Fatalf("expected closed-plan message, got: %v", err)
	}
}

func TestPlanNewSessionReturnsPoolExhaustedAtPerPlanVMLimit(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithMaxActiveSessions(2), WithMaxVMsPerPlan(1))
	plan := mustCompilePlan(t, eng, coverageValidQuery)
	first := mustNewSession(t, plan)
	defer func() {
		_ = first.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := plan.NewSession(ctx)
	if !errors.Is(err, vm.ErrPoolExhausted) {
		t.Fatalf("expected pool exhaustion at per-plan VM limit, got: %v", err)
	}
}

func TestPlanNewSessionReleasesLimiterOnPoolExhaustion(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithMaxActiveSessions(1), WithMaxVMsPerPlan(1), WithMaxIdleVMsPerPlan(1))
	plan := mustCompilePlan(t, eng, coverageValidQuery)

	borrowed, err := plan.pool.Acquire()
	if err != nil {
		t.Fatalf("expected direct pool acquire to succeed, got: %v", err)
	}

	_, err = plan.NewSession(context.Background())
	if !errors.Is(err, vm.ErrPoolExhausted) {
		t.Fatalf("expected pool exhaustion during session creation, got: %v", err)
	}

	plan.pool.Release(borrowed)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	session, err := plan.NewSession(ctx)
	if err != nil {
		t.Fatalf("expected limiter permit to be released after pool exhaustion, got: %v", err)
	}

	defer func() {
		_ = session.Close()
	}()
}

func TestSessionPermitReleaseIsConcurrentSafeAndIdempotent(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithMaxActiveSessions(1), WithMaxIdleVMsPerPlan(1), WithMaxVMsPerPlan(1))
	plan := mustCompilePlan(t, eng, coverageValidQuery)
	defer func() {
		_ = plan.Close()
	}()

	instance, err := plan.pool.Acquire()
	if err != nil {
		t.Fatalf("expected direct pool acquire to succeed, got: %v", err)
	}

	if err := plan.limiter.Acquire(context.Background()); err != nil {
		t.Fatalf("expected limiter acquire to succeed, got: %v", err)
	}

	release := newSessionPermitRelease(plan.limiter, plan.pool)

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

	if err := plan.limiter.Acquire(ctx); err != nil {
		t.Fatalf("expected concurrent release to free the limiter slot, got: %v", err)
	}

	plan.limiter.Release()

	reused, err := plan.pool.Acquire()
	if err != nil {
		t.Fatalf("expected concurrent release to return the borrowed VM, got: %v", err)
	}

	if reused != instance {
		t.Fatal("expected concurrent release to retain the original VM in the pool")
	}

	_, err = plan.pool.Acquire()
	if !errors.Is(err, vm.ErrPoolExhausted) {
		t.Fatalf("expected concurrent release not to free extra pool capacity, got: %v", err)
	}

	plan.pool.Release(reused)
}

func TestPlanNewDebugSessionReleasesLimiterOnEnvironmentError(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithMaxActiveSessions(1))
	plan, err := eng.CompileDebug(context.Background(), source.NewAnonymous(coverageValidQuery))
	if err != nil {
		t.Fatalf("failed to compile debug plan: %v", err)
	}

	_, err = plan.NewDebugSession(
		context.Background(),
		WithEnvironmentOptions(
			vm.WithFunction("DEBUG_SESSION_DUP_LIMIT", coverageVarFn),
			vm.WithFunction("DEBUG_SESSION_DUP_LIMIT", coverageVarFn),
		),
	)
	if err == nil {
		t.Fatal("expected plan.NewDebugSession to fail on conflicting session options")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	session, err := plan.NewDebugSession(ctx)
	if err != nil {
		t.Fatalf("expected limiter permit to be released after failed debug session creation, got: %v", err)
	}
	defer func() {
		_ = session.Close()
	}()
}

func TestPlanNewDebugSessionMetadataRejectionDoesNotAcquireLimiter(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithMaxActiveSessions(1))
	plan := mustCompilePlan(t, eng, coverageValidQuery)

	_, err := plan.NewDebugSession(context.Background())
	if err == nil {
		t.Fatal("expected non-debug plan to be rejected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	session, err := plan.NewSession(ctx)
	if err != nil {
		t.Fatalf("expected debug metadata rejection not to consume a limiter permit, got: %v", err)
	}
	defer func() {
		_ = session.Close()
	}()
}

func TestNewPlanSessionReleasesLimiterOnBuilderPanic(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithMaxActiveSessions(1))
	plan := mustCompilePlan(t, eng, coverageValidQuery)

	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("expected session builder panic")
			}
		}()

		_, _ = newPlanSession(
			plan,
			context.Background(),
			nil,
			planSessionSetup{},
			func(planSessionDependencies) (struct{}, error) {
				panic("session builder failed")
			},
		)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	session, err := plan.NewSession(ctx)
	if err != nil {
		t.Fatalf("expected builder panic to release limiter permit, got: %v", err)
	}
	defer func() {
		_ = session.Close()
	}()
}

func TestDebugSessionCloseReleasesLimiterOnce(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithMaxActiveSessions(2))
	plan, err := eng.CompileDebug(context.Background(), source.NewAnonymous(coverageValidQuery))
	if err != nil {
		t.Fatalf("failed to compile debug plan: %v", err)
	}

	debugSession, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatalf("failed to create debug session: %v", err)
	}

	firstSession, err := plan.NewSession(context.Background())
	if err != nil {
		t.Fatalf("failed to create first normal session: %v", err)
	}
	defer func() {
		_ = firstSession.Close()
	}()

	if err := debugSession.Close(); err != nil {
		t.Fatalf("expected debug session close to succeed, got: %v", err)
	}
	if err := debugSession.Close(); err != nil {
		t.Fatalf("expected repeated debug session close to succeed, got: %v", err)
	}

	secondSession, err := plan.NewSession(context.Background())
	if err != nil {
		t.Fatalf("expected debug session close to release one limiter permit, got: %v", err)
	}
	defer func() {
		_ = secondSession.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, err = plan.NewSession(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected repeated debug session close not to release an extra permit, got: %v", err)
	}
}
