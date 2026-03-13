package ferret

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestSessionRunReturnsBeforeHookError(t *testing.T) {
	t.Parallel()

	beforeErr := errors.New("before run failed")
	afterCalled := false

	eng := mustNewEngine(t,
		WithBeforeRunHook(func(ctx context.Context) (context.Context, error) {
			return ctx, beforeErr
		}),
		WithAfterRunHook(func(context.Context, error) error {
			afterCalled = true
			return nil
		}),
	)
	plan := mustCompilePlan(t, eng, coverageValidQuery)
	session := mustNewSession(t, plan)

	result, err := session.Run(context.Background())
	if err == nil {
		t.Fatal("expected session run to fail on before run hook error")
	}

	if result != nil {
		t.Fatal("expected nil result when before run hook fails")
	}

	if !errors.Is(err, beforeErr) {
		t.Fatalf("expected before run hook error to be preserved, got: %v", err)
	}

	if !strings.Contains(err.Error(), "before run hooks") {
		t.Fatalf("expected before run hooks label, got: %v", err)
	}

	if afterCalled {
		t.Fatal("after run hooks must not run when before run hook fails")
	}
}

func TestSessionRunReturnsAfterHookErrorOnSuccess(t *testing.T) {
	t.Parallel()

	afterErr := errors.New("after run failed")
	var seenRunErr error

	eng := mustNewEngine(t, WithAfterRunHook(func(_ context.Context, err error) error {
		seenRunErr = err
		return afterErr
	}))
	plan := mustCompilePlan(t, eng, coverageValidQuery)
	session := mustNewSession(t, plan)

	result, err := session.Run(context.Background())
	if err == nil {
		t.Fatal("expected session run to fail when after run hook fails")
	}

	if result != nil {
		t.Fatal("expected nil result when after run hook fails")
	}

	if seenRunErr != nil {
		t.Fatalf("expected after hook to receive nil run error, got: %v", seenRunErr)
	}

	if !errors.Is(err, afterErr) {
		t.Fatalf("expected returned error to include after hook error, got: %v", err)
	}

	if !strings.Contains(err.Error(), "after run hooks") {
		t.Fatalf("expected after run hooks label, got: %v", err)
	}
}

func TestSessionRunReturnsVMError(t *testing.T) {
	t.Parallel()

	vmErr := errors.New("vm run failed")
	eng := mustNewEngine(t, WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
		fns.A0().Add("SESSION_FAIL_FN", func(context.Context) (runtime.Value, error) {
			return runtime.None, vmErr
		})
	}))
	plan := mustCompilePlan(t, eng, "RETURN SESSION_FAIL_FN()")
	session := mustNewSession(t, plan)

	result, err := session.Run(context.Background())
	if err == nil {
		t.Fatal("expected session run to fail when VM returns error")
	}

	if result != nil {
		t.Fatal("expected nil result when VM returns error")
	}

	if !errors.Is(err, vmErr) {
		t.Fatalf("expected VM error to be preserved, got: %v", err)
	}

	if strings.Contains(err.Error(), "after run hooks") {
		t.Fatalf("did not expect after run hooks label when no after hooks fail, got: %v", err)
	}
}

func TestSessionRunJoinsVMAndAfterHookErrors(t *testing.T) {
	t.Parallel()

	vmErr := errors.New("vm run failed")
	afterErr := errors.New("after run failed")
	var seenRunErr error

	eng := mustNewEngine(t,
		WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
			fns.A0().Add("SESSION_FAIL_FN_JOIN", func(context.Context) (runtime.Value, error) {
				return runtime.None, vmErr
			})
		}),
		WithAfterRunHook(func(_ context.Context, err error) error {
			seenRunErr = err
			return afterErr
		}),
	)
	plan := mustCompilePlan(t, eng, "RETURN SESSION_FAIL_FN_JOIN()")
	session := mustNewSession(t, plan)

	result, err := session.Run(context.Background())
	if err == nil {
		t.Fatal("expected session run to fail when VM and after hook fail")
	}

	if result != nil {
		t.Fatal("expected nil result when VM and after hook fail")
	}

	if seenRunErr == nil {
		t.Fatal("expected after hook to receive VM run error")
	}

	if !errors.Is(err, afterErr) {
		t.Fatalf("expected returned error to include after hook error, got: %v", err)
	}

	if !errors.Is(err, vmErr) {
		t.Fatalf("expected returned error to include VM error, got: %v", err)
	}

	if !strings.Contains(err.Error(), "after run hooks") {
		t.Fatalf("expected joined error to include after run hooks label, got: %v", err)
	}
}

func TestSessionClose(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t)
	plan := mustCompilePlan(t, eng, coverageValidQuery)
	session := mustNewSession(t, plan)

	if err := session.Close(); err != nil {
		t.Fatalf("expected session close without hooks to succeed, got: %v", err)
	}

	closeErr := errors.New("session close failed")
	eng = mustNewEngine(t, WithSessionCloseHook(func() error {
		return closeErr
	}))
	plan = mustCompilePlan(t, eng, coverageValidQuery)
	session = mustNewSession(t, plan)

	err := session.Close()
	if err == nil {
		t.Fatal("expected session close to fail when close hook fails")
	}

	if !errors.Is(err, closeErr) {
		t.Fatalf("expected close error to include session close hook error, got: %v", err)
	}

	if !strings.Contains(err.Error(), "close hooks") {
		t.Fatalf("expected close hooks label, got: %v", err)
	}
}

func TestSessionCloseReturnsBorrowedVMToPool(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithMaxIdleVMsPerPlan(1))
	plan := mustCompilePlan(t, eng, coverageValidQuery)
	first := mustNewSession(t, plan)
	firstVM := first.vm

	if err := first.Close(); err != nil {
		t.Fatalf("expected first session close to succeed, got: %v", err)
	}

	second := mustNewSession(t, plan)
	defer func() {
		_ = second.Close()
	}()

	if second.vm != firstVM {
		t.Fatal("expected second session to reuse the pooled VM from the first session")
	}
}

func TestSessionCloseAfterPlanCloseReleasesLimiter(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithMaxActiveSessions(1), WithMaxIdleVMsPerPlan(1))
	plan := mustCompilePlan(t, eng, coverageValidQuery)
	session := mustNewSession(t, plan)

	if err := plan.Close(); err != nil {
		t.Fatalf("expected plan close to succeed with active session, got: %v", err)
	}

	if err := session.Close(); err != nil {
		t.Fatalf("expected session close after plan close to succeed, got: %v", err)
	}

	nextPlan := mustCompilePlan(t, eng, coverageValidQuery)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	nextSession, err := nextPlan.NewSession(ctx)
	if err != nil {
		t.Fatalf("expected limiter permit to be released after closing the orphaned session, got: %v", err)
	}

	defer func() {
		_ = nextSession.Close()
		_ = nextPlan.Close()
	}()
}
