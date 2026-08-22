package ferret

import (
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/module"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	sharedoptions "github.com/ziflex/go-options"
)

func reportCompilerOptionOnApplication[T any](
	base sharedoptions.Option[T],
	application int,
	validationErr sharedoptions.ValidationError,
) sharedoptions.Option[T] {
	applied := 0

	return func(config *T, report sharedoptions.Report) {
		applied++
		base(config, report)

		if applied == application {
			report(validationErr)
		}
	}
}

func TestNewReturnsCompilerOptionErrorsBeforeBootstrap(t *testing.T) {
	tests := []struct {
		name        string
		phase       string
		application int
	}{
		{
			name:        "normal compiler",
			phase:       "compiler",
			application: 1,
		},
		{
			name:        "debug compiler",
			phase:       "debug compiler",
			application: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validationErr := sharedoptions.ValidationError{
				Field:  "compiler",
				Value:  tt.name,
				Reason: "test validation failure",
			}
			invalid := reportCompilerOptionOnApplication(
				compiler.WithOptimizationLevel(compiler.O1),
				tt.application,
				validationErr,
			)
			client := &recordingHTTPClient{}
			moduleRegistered := false
			initHookCalled := false
			closeHookCalled := false

			engine, err := New(
				WithCompilerOptions(invalid),
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
				t.Fatal("New() error = nil, want compiler option error")
			}
			if engine != nil {
				t.Fatal("engine != nil, want nil")
			}
			if !strings.HasPrefix(err.Error(), tt.phase+": ") {
				t.Fatalf("New() error = %q, want %q phase", err, tt.phase)
			}

			var got sharedoptions.ValidationError
			if !errors.As(err, &got) {
				t.Fatalf("New() error = %T, want sharedoptions.ValidationError", err)
			}
			if got != validationErr {
				t.Fatalf("New() validation error = %+v, want %+v", got, validationErr)
			}
			if moduleRegistered {
				t.Fatal("module registered before compiler validation")
			}
			if initHookCalled {
				t.Fatal("init hook called before compiler validation")
			}
			if closeHookCalled {
				t.Fatal("close hook called before bootstrap")
			}
			if got := client.idleCloseCount(); got != 1 {
				t.Fatalf("owned network idle closes = %d, want 1", got)
			}
		})
	}
}

func TestNewDoesNotCloseInjectedNetworkAfterCompilerOptionError(t *testing.T) {
	validationErr := sharedoptions.ValidationError{Reason: "test validation failure"}
	invalid := reportCompilerOptionOnApplication(
		compiler.WithOptimizationLevel(compiler.O1),
		1,
		validationErr,
	)
	client := &recordingHTTPClient{}
	network := mustNewTestNetwork(t, ferretnet.WithHTTPClient(client))

	engine, err := New(
		WithCompilerOptions(invalid),
		WithNetwork(network),
	)
	if err == nil {
		t.Fatal("New() error = nil, want compiler option error")
	}
	if engine != nil {
		t.Fatal("engine != nil, want nil")
	}
	if got := client.idleCloseCount(); got != 0 {
		t.Fatalf("injected network idle closes = %d, want 0", got)
	}
}
