package ferret

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestPlanNewSessionReturnsEnvironmentError(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t)
	plan := mustCompilePlan(t, eng, coverageValidQuery)

	_, err := plan.NewSession(
		context.Background(),
		vm.WithFunction("SESSION_DUP_COVER", coverageVarFn),
		vm.WithFunction("SESSION_DUP_COVER", coverageVarFn),
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
		vm.WithFunction("SESSION_DUP_LIMIT", coverageVarFn),
		vm.WithFunction("SESSION_DUP_LIMIT", coverageVarFn),
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
