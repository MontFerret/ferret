package ferret

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/module"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	sharedoptions "github.com/ziflex/go-options"
)

func TestWithOptimizationLevelConfiguresEngineCompiler(t *testing.T) {
	tests := []struct {
		name    string
		setters []Option
		level   OptimizationLevel
	}{
		{name: "default", level: OptimizationFull},
		{name: "none", setters: []Option{WithOptimizationLevel(OptimizationNone)}, level: OptimizationNone},
		{name: "basic", setters: []Option{WithOptimizationLevel(OptimizationBasic)}, level: OptimizationBasic},
		{name: "full", setters: []Option{WithOptimizationLevel(OptimizationFull)}, level: OptimizationFull},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := New(tt.setters...)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			t.Cleanup(func() {
				if err := engine.Close(); err != nil {
					t.Errorf("Engine.Close() error = %v", err)
				}
			})

			plan, err := engine.Compile(t.Context(), NewAnonymousSource("RETURN 1"))
			if err != nil {
				t.Fatalf("Compile() error = %v", err)
			}
			t.Cleanup(func() {
				if err := plan.Close(); err != nil {
					t.Errorf("Plan.Close() error = %v", err)
				}
			})

			if got := plan.prog.Metadata.OptimizationLevel; got != int(tt.level) {
				t.Fatalf("optimization level = %d, want %d", got, tt.level)
			}

			session, err := plan.NewSession(t.Context())
			if err != nil {
				t.Fatalf("NewSession() error = %v", err)
			}
			t.Cleanup(func() {
				if err := session.Close(); err != nil {
					t.Errorf("Session.Close() error = %v", err)
				}
			})

			output, err := session.Run(t.Context())
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if got := string(output.Content); got != "1" {
				t.Fatalf("Run() output = %q, want %q", got, "1")
			}
		})
	}
}

func TestCompileDebugUsesOptimizationNone(t *testing.T) {
	engine, err := New(WithOptimizationLevel(OptimizationFull))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() {
		if err := engine.Close(); err != nil {
			t.Errorf("Engine.Close() error = %v", err)
		}
	})

	plan, err := engine.CompileDebug(t.Context(), NewAnonymousSource("RETURN 1"))
	if err != nil {
		t.Fatalf("CompileDebug() error = %v", err)
	}
	t.Cleanup(func() {
		if err := plan.Close(); err != nil {
			t.Errorf("Plan.Close() error = %v", err)
		}
	})

	if got := plan.prog.Metadata.OptimizationLevel; got != int(OptimizationNone) {
		t.Fatalf("optimization level = %d, want %d", got, OptimizationNone)
	}
}

func TestNewReturnsOptimizationLevelErrorsBeforeBootstrap(t *testing.T) {
	client := &recordingHTTPClient{}
	moduleRegistered := false
	initHookCalled := false
	closeHookCalled := false

	engine, err := New(
		WithOptimizationLevel(OptimizationLevel(-1)),
		WithNetworkOptions(ferretnet.WithHTTPClient(client)),
		WithModules(testModule{registerFn: func(module.Bootstrap) error {
			moduleRegistered = true

			return nil
		}}),
		WithEngineInitHook(func() error {
			initHookCalled = true

			return nil
		}),
		WithEngineCloseHook(func() error {
			closeHookCalled = true

			return nil
		}),
	)
	if err == nil {
		t.Fatal("New() error = nil, want optimization level error")
	}
	if engine != nil {
		t.Fatal("engine != nil, want nil")
	}

	var got sharedoptions.ValidationError
	if !errors.As(err, &got) {
		t.Fatalf("New() error = %T, want sharedoptions.ValidationError", err)
	}
	if got.Field != "optimization level" || got.Value != "unknown" {
		t.Fatalf("New() validation error = %+v, want optimization level/unknown", got)
	}
	if moduleRegistered {
		t.Fatal("module registered before option validation")
	}
	if initHookCalled {
		t.Fatal("init hook called before option validation")
	}
	if closeHookCalled {
		t.Fatal("close hook called before bootstrap")
	}
	if got := client.idleCloseCount(); got != 1 {
		t.Fatalf("owned network idle closes = %d, want 1", got)
	}
}

func TestNewDoesNotCloseInjectedNetworkAfterOptimizationLevelError(t *testing.T) {
	client := &recordingHTTPClient{}
	network := mustNewTestNetwork(t, ferretnet.WithHTTPClient(client))

	engine, err := New(
		WithOptimizationLevel(OptimizationLevel(-1)),
		WithNetwork(network),
	)
	if err == nil {
		t.Fatal("New() error = nil, want optimization level error")
	}
	if engine != nil {
		t.Fatal("engine != nil, want nil")
	}
	if got := client.idleCloseCount(); got != 0 {
		t.Fatalf("injected network idle closes = %d, want 0", got)
	}
}
