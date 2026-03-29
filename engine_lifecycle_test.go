package ferret

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type testModule struct {
	registerFn func(boot Bootstrap) error
	name       string
}

func (m testModule) Name() string {
	if m.name == "" {
		return "test-module"
	}

	return m.name
}

func (m testModule) Register(boot Bootstrap) error {
	if m.registerFn == nil {
		return nil
	}

	return m.registerFn(boot)
}

func testFn0(context.Context) (runtime.Value, error) {
	return runtime.None, nil
}

func testVarFn(context.Context, ...runtime.Value) (runtime.Value, error) {
	return runtime.None, nil
}

func TestNewRunsCloseHooksWhenHostBuildFails(t *testing.T) {
	t.Parallel()

	var (
		moduleRegistered bool
		closeHookCalled  bool
	)

	mod := testModule{
		registerFn: func(boot Bootstrap) error {
			moduleRegistered = true
			boot.Hooks().Engine().OnClose(func() error {
				closeHookCalled = true
				return nil
			})

			boot.Host().Library().Function().A0().Add("LIFECYCLE_DUPLICATE_FN", testFn0)
			boot.Host().Library().Function().A0().Add("LIFECYCLE_DUPLICATE_FN", testFn0)

			return nil
		},
	}

	_, err := New(WithModules(mod))
	if err == nil {
		t.Fatal("expected New to fail when host build fails")
	}

	if !moduleRegistered {
		t.Fatal("expected module registration to run before failure")
	}

	if !closeHookCalled {
		t.Fatal("expected engine close hooks to run on host build failure")
	}
}

func TestNewReturnsJoinedErrorWhenInitAndCloseHooksFail(t *testing.T) {
	t.Parallel()

	initErr := errors.New("init boom")
	closeErr := errors.New("close boom")

	_, err := New(
		WithEngineInitHook(func() error {
			return initErr
		}),
		WithEngineCloseHook(func() error {
			return closeErr
		}),
	)
	if err == nil {
		t.Fatal("expected New to fail")
	}

	if !errors.Is(err, initErr) {
		t.Fatalf("expected joined error to include init error, got: %v", err)
	}

	if !errors.Is(err, closeErr) {
		t.Fatalf("expected joined error to include close error, got: %v", err)
	}

	if !strings.Contains(err.Error(), "init hooks") {
		t.Fatalf("expected error to include init hooks label, got: %v", err)
	}

	if !strings.Contains(err.Error(), "close hooks") {
		t.Fatalf("expected error to include close hooks label, got: %v", err)
	}
}

func TestRunClosesPlanWhenSessionCreationFails(t *testing.T) {
	t.Parallel()

	planClosed := false

	eng, err := New(
		WithPlanCloseHook(func() error {
			planClosed = true
			return nil
		}),
	)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	_, err = eng.Run(
		context.Background(),
		source.NewAnonymous("RETURN 1"),
		WithEnvironmentOptions(
			vm.WithFunction("SESSION_DUP", testVarFn),
			vm.WithFunction("SESSION_DUP", testVarFn),
		),
	)
	if err == nil {
		t.Fatal("expected Run to fail during session creation")
	}

	if !planClosed {
		t.Fatal("expected plan close hook to run when session creation fails")
	}
}

func TestRunLogsDeferredCleanupErrorsWithoutChangingRunResult(t *testing.T) {
	t.Parallel()

	sessionCloseErr := errors.New("session close failed")
	planCloseErr := errors.New("plan close failed")
	logOutput := bytes.NewBuffer(nil)

	eng, err := New(
		WithLog(logOutput),
		WithLogLevel(runtime.ErrorLevel),
		WithSessionCloseHook(func() error {
			return sessionCloseErr
		}),
		WithPlanCloseHook(func() error {
			return planCloseErr
		}),
	)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	result, err := eng.Run(context.Background(), source.NewAnonymous("RETURN 1"))
	if err != nil {
		t.Fatalf("expected run result error to be unchanged by cleanup failures, got: %v", err)
	}

	if got := strings.TrimSpace(string(result.Content)); got != "1" {
		t.Fatalf("expected run result to stay successful, got: %s", got)
	}

	logs := logOutput.String()
	if !strings.Contains(logs, `"phase":"session"`) {
		t.Fatalf("expected cleanup logs to include session phase, got: %s", logs)
	}

	if !strings.Contains(logs, `"phase":"plan"`) {
		t.Fatalf("expected cleanup logs to include plan phase, got: %s", logs)
	}

	if !strings.Contains(logs, `"operation":"close"`) {
		t.Fatalf("expected cleanup logs to include close operation, got: %s", logs)
	}

	if !strings.Contains(logs, sessionCloseErr.Error()) {
		t.Fatalf("expected cleanup logs to include session close error, got: %s", logs)
	}

	if !strings.Contains(logs, planCloseErr.Error()) {
		t.Fatalf("expected cleanup logs to include plan close error, got: %s", logs)
	}
}
