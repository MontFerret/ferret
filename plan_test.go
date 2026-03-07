package ferret

import (
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestPlanNewSessionReturnsEnvironmentError(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t)
	plan := mustCompilePlan(t, eng, coverageValidQuery)

	_, err := plan.NewSession(
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
