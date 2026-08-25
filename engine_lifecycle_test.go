package ferret

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/module"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type testModule struct {
	registerFn func(boot module.Bootstrap) error
	name       string
}

func (m testModule) Name() string {
	if m.name == "" {
		return "test-module"
	}

	return m.name
}

func (m testModule) Register(boot module.Bootstrap) error {
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
	client := &recordingHTTPClient{}

	mod := testModule{
		registerFn: func(boot module.Bootstrap) error {
			moduleRegistered = true
			internal, ok := boot.(*bootstrap)
			if !ok {
				t.Fatalf("expected internal bootstrap, got %T", boot)
			}
			internal.host.network = mustNewTestNetwork(t, ferretnet.WithHTTPClient(client))

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

	if got := client.idleCloseCount(); got != 1 {
		t.Fatalf("expected host-build-failure cleanup, got %d calls", got)
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

func TestEngineCloseClosesOwnedNetworkIdleConnections(t *testing.T) {
	t.Parallel()

	client := &recordingHTTPClient{}
	eng := mustNewEngine(t)
	eng.host.network = mustNewTestNetwork(t, ferretnet.WithHTTPClient(client))

	if err := eng.Close(); err != nil {
		t.Fatalf("close engine: %v", err)
	}

	if got := client.idleCloseCount(); got != 1 {
		t.Fatalf("expected one idle-connection cleanup, got %d", got)
	}
}

func TestEngineCloseCleansOwnedNetworkAfterHookFailure(t *testing.T) {
	t.Parallel()

	hookErr := errors.New("close hook failed")
	client := &recordingHTTPClient{}
	eng := mustNewEngine(t, WithEngineCloseHook(func() error {
		return hookErr
	}))
	eng.host.network = mustNewTestNetwork(t, ferretnet.WithHTTPClient(client))

	err := eng.Close()
	if !errors.Is(err, hookErr) {
		t.Fatalf("expected close hook error, got %v", err)
	}

	if got := client.idleCloseCount(); got != 1 {
		t.Fatalf("expected cleanup after hook failure, got %d calls", got)
	}
}

func TestEngineCloseDoesNotCleanInjectedNetwork(t *testing.T) {
	t.Parallel()

	client := &recordingHTTPClient{}
	network := mustNewTestNetwork(t, ferretnet.WithHTTPClient(client))
	eng := mustNewEngine(t, WithNetwork(network))

	if err := eng.Close(); err != nil {
		t.Fatalf("close engine: %v", err)
	}

	if got := client.idleCloseCount(); got != 0 {
		t.Fatalf("expected injected network ownership to remain with caller, got %d cleanup calls", got)
	}
}

func TestEngineNetworkOwnershipFollowsLastOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		managedLast bool
	}{
		{name: "network options last", managedLast: true},
		{name: "injected network last", managedLast: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			managedClient := &recordingHTTPClient{}
			injectedClient := &recordingHTTPClient{}
			injectedNetwork := mustNewTestNetwork(t, ferretnet.WithHTTPClient(injectedClient))

			var setters []Option
			if tt.managedLast {
				setters = []Option{
					WithNetwork(injectedNetwork),
					WithNetworkOptions(ferretnet.WithHTTPClient(managedClient)),
				}
			} else {
				setters = []Option{
					WithNetworkOptions(ferretnet.WithHTTPClient(managedClient)),
					WithNetwork(injectedNetwork),
				}
			}

			eng := mustNewEngine(t, setters...)
			if tt.managedLast {
				if got := eng.host.network.HTTP(); got != managedClient {
					t.Fatalf("expected network options client, got %T", got)
				}
			} else if eng.host.network != injectedNetwork {
				t.Fatalf("expected injected network, got %T", eng.host.network)
			}

			if eng.ownsNetwork != tt.managedLast {
				t.Fatalf("expected ownsNetwork to be %t", tt.managedLast)
			}

			wantManagedCloses := 1
			if tt.managedLast {
				wantManagedCloses = 0
			}

			if got := managedClient.idleCloseCount(); got != wantManagedCloses {
				t.Fatalf("expected %d managed network cleanup calls after construction, got %d", wantManagedCloses, got)
			}

			if got := injectedClient.idleCloseCount(); got != 0 {
				t.Fatalf("expected injected network to remain caller-owned after construction, got %d cleanup calls", got)
			}

			if err := eng.Close(); err != nil {
				t.Fatalf("close engine: %v", err)
			}

			if tt.managedLast {
				wantManagedCloses = 1
			}

			if got := managedClient.idleCloseCount(); got != wantManagedCloses {
				t.Fatalf("expected %d managed network cleanup calls, got %d", wantManagedCloses, got)
			}

			if got := injectedClient.idleCloseCount(); got != 0 {
				t.Fatalf("expected injected network to remain caller-owned, got %d cleanup calls", got)
			}
		})
	}
}

func TestEngineClosesSupersededManagedNetworks(t *testing.T) {
	t.Parallel()

	firstClient := &recordingHTTPClient{}
	secondClient := &recordingHTTPClient{}
	eng := mustNewEngine(
		t,
		WithNetworkOptions(ferretnet.WithHTTPClient(firstClient)),
		WithNetworkOptions(ferretnet.WithHTTPClient(secondClient)),
	)

	if !eng.ownsNetwork {
		t.Fatal("expected final managed network to be engine-owned")
	}

	if got := eng.host.network.HTTP(); got != secondClient {
		t.Fatalf("expected second managed network to be selected, got %T", got)
	}

	if got := firstClient.idleCloseCount(); got != 1 {
		t.Fatalf("expected superseded managed network to close once after construction, got %d", got)
	}

	if got := secondClient.idleCloseCount(); got != 0 {
		t.Fatalf("expected selected managed network to remain active after construction, got %d closes", got)
	}

	if err := eng.Close(); err != nil {
		t.Fatalf("close engine: %v", err)
	}

	if got := firstClient.idleCloseCount(); got != 1 {
		t.Fatalf("expected superseded managed network to remain closed exactly once, got %d", got)
	}

	if got := secondClient.idleCloseCount(); got != 1 {
		t.Fatalf("expected engine shutdown to close selected managed network once, got %d", got)
	}
}

func TestNewCleansNetworkCreatedFromOptionsOnInitFailure(t *testing.T) {
	t.Parallel()

	initErr := errors.New("init failed")
	client := &recordingHTTPClient{}

	_, err := New(
		WithNetworkOptions(ferretnet.WithHTTPClient(client)),
		WithEngineInitHook(func() error {
			return initErr
		}),
	)
	if !errors.Is(err, initErr) {
		t.Fatalf("expected init error, got %v", err)
	}

	if got := client.idleCloseCount(); got != 1 {
		t.Fatalf("expected construction-failure cleanup, got %d calls", got)
	}
}

func TestNewRollsBackAllConstructedNetworksOnOptionFailure(t *testing.T) {
	t.Parallel()

	firstManagedClient := &recordingHTTPClient{}
	secondManagedClient := &recordingHTTPClient{}
	injectedClient := &recordingHTTPClient{}
	injectedNetwork := mustNewTestNetwork(t, ferretnet.WithHTTPClient(injectedClient))

	engine, err := New(
		WithNetworkOptions(ferretnet.WithHTTPClient(firstManagedClient)),
		WithNetworkOptions(ferretnet.WithHTTPClient(secondManagedClient)),
		WithEngineParams(map[string]any{"unsupported": make(chan int)}),
		WithNetwork(injectedNetwork),
	)
	if engine != nil {
		_ = engine.Close()

		t.Fatal("expected option failure not to return an engine")
	}

	if !errors.Is(err, runtime.ErrInvalidType) {
		t.Fatalf("expected runtime.ErrInvalidType, got %v", err)
	}

	if got := firstManagedClient.idleCloseCount(); got != 1 {
		t.Fatalf("expected first constructed network to close once, got %d", got)
	}

	if got := secondManagedClient.idleCloseCount(); got != 1 {
		t.Fatalf("expected later constructed network to close once, got %d", got)
	}

	if got := injectedClient.idleCloseCount(); got != 0 {
		t.Fatalf("expected injected network to remain caller-owned, got %d closes", got)
	}
}

func TestNewRollsBackConstructedNetworkOnStdlibRegistrationFailure(t *testing.T) {
	t.Parallel()

	client := &recordingHTTPClient{}
	engine, err := New(
		WithNetworkOptions(ferretnet.WithHTTPClient(client)),
		WithStdlib(stdlib.Only(stdlib.Group("unknown"))),
	)
	if engine != nil {
		_ = engine.Close()

		t.Fatal("expected stdlib registration failure not to return an engine")
	}

	if err == nil || !strings.Contains(err.Error(), "stdlib: invalid stdlib group(s): unknown") {
		t.Fatalf("expected stdlib registration failure, got %v", err)
	}

	if got := client.idleCloseCount(); got != 1 {
		t.Fatalf("expected constructed network to close once, got %d", got)
	}
}

func TestNewCleansOwnedNetworkOnRegistrationFailure(t *testing.T) {
	t.Parallel()

	registerErr := errors.New("register failed")
	client := &recordingHTTPClient{}
	mod := testModule{
		registerFn: func(boot module.Bootstrap) error {
			internal, ok := boot.(*bootstrap)
			if !ok {
				t.Fatalf("expected internal bootstrap, got %T", boot)
			}

			internal.host.network = mustNewTestNetwork(t, ferretnet.WithHTTPClient(client))

			return registerErr
		},
	}

	_, err := New(WithModules(mod))
	if !errors.Is(err, registerErr) {
		t.Fatalf("expected registration error, got %v", err)
	}

	if got := client.idleCloseCount(); got != 1 {
		t.Fatalf("expected construction-failure cleanup, got %d calls", got)
	}
}

func TestNewCleansOwnedNetworkOnInitFailure(t *testing.T) {
	t.Parallel()

	initErr := errors.New("init failed")
	client := &recordingHTTPClient{}
	mod := testModule{
		registerFn: func(boot module.Bootstrap) error {
			internal, ok := boot.(*bootstrap)
			if !ok {
				t.Fatalf("expected internal bootstrap, got %T", boot)
			}

			internal.host.network = mustNewTestNetwork(t, ferretnet.WithHTTPClient(client))

			return nil
		},
	}

	_, err := New(
		WithModules(mod),
		WithEngineInitHook(func() error {
			return initErr
		}),
	)
	if !errors.Is(err, initErr) {
		t.Fatalf("expected init error, got %v", err)
	}

	if got := client.idleCloseCount(); got != 1 {
		t.Fatalf("expected init-failure cleanup, got %d calls", got)
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
		NewAnonymousSource("RETURN 1"),
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
		WithLogLevel(logging.ErrorLevel),
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

	result, err := eng.Run(context.Background(), NewAnonymousSource("RETURN 1"))
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
