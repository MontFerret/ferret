package vm

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestNewEnvironmentUsesDefaultsAndIgnoresNilOptions(t *testing.T) {
	t.Parallel()

	env, err := NewEnvironment([]EnvironmentOption{nil})
	if err != nil {
		t.Fatalf("expected environment construction to succeed, got: %v", err)
	}

	if env == nil {
		t.Fatal("expected environment")
	}

	if got := env.Functions.Size(); got != 0 {
		t.Fatalf("expected no default functions, got %d", got)
	}

	if got := len(env.Params); got != 0 {
		t.Fatalf("expected no default parameters, got %d", got)
	}
}

func TestNewEnvironmentPreservesOptionOrderingAndNoOps(t *testing.T) {
	t.Parallel()

	function := func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, nil
	}
	env, err := NewEnvironment([]EnvironmentOption{
		WithParams(runtime.Params{
			"first": runtime.NewInt(1),
			"last":  runtime.NewInt(2),
		}),
		WithParam("last", runtime.NewInt(3)),
		WithParams(runtime.Params{"last": runtime.NewInt(4)}),
		WithParams(nil),
		WithFunctions(nil),
		WithFunction("", function),
		WithFunction("nil", nil),
		WithNamespace(nil),
		WithFunctionsBuilder(nil),
		WithFunctionsRegistrar(nil),
		WithFunction("PING", function),
	})
	if err != nil {
		t.Fatalf("expected environment construction to succeed, got: %v", err)
	}

	if got, want := env.Params.MustGet("first"), runtime.NewInt(1); got != want {
		t.Fatalf("unexpected first parameter: got %v, want %v", got, want)
	}

	if got, want := env.Params.MustGet("last"), runtime.NewInt(4); got != want {
		t.Fatalf("unexpected last parameter: got %v, want %v", got, want)
	}

	if !env.Functions.Has("PING") {
		t.Fatal("expected valid function to be registered")
	}

	if env.Functions.Has("") || env.Functions.Has("nil") {
		t.Fatal("expected invalid function inputs to remain no-ops")
	}
}

func TestNewEnvironmentJoinsOptionErrorsAndReturnsNoEnvironment(t *testing.T) {
	t.Parallel()

	firstErr := errors.New("first option")
	secondErr := errors.New("second option")
	applied := make([]string, 0, 2)
	env, err := NewEnvironment([]EnvironmentOption{
		func(cfg *environmentConfig) error {
			applied = append(applied, "first")
			cfg.params["value"] = runtime.NewInt(1)

			return firstErr
		},
		nil,
		func(cfg *environmentConfig) error {
			applied = append(applied, "second")
			cfg.params["value"] = runtime.NewInt(2)

			return secondErr
		},
	})
	if env != nil {
		t.Fatal("expected option failure to return no environment")
	}

	if !errors.Is(err, firstErr) || !errors.Is(err, secondErr) {
		t.Fatalf("expected both option errors, got: %v", err)
	}

	if got, want := len(applied), 2; got != want {
		t.Fatalf("expected later options to execute after an error: got %d applications, want %d", got, want)
	}

	if applied[0] != "first" || applied[1] != "second" {
		t.Fatalf("unexpected option application order: %v", applied)
	}
}

func TestExtendEnvironmentPropagatesOptionErrorsWithoutChangingBase(t *testing.T) {
	t.Parallel()

	base := NewDefaultEnvironment()
	base.Params["value"] = runtime.NewInt(1)
	wantErr := errors.New("extend option")
	extended, err := ExtendEnvironment(base, []EnvironmentOption{
		func(cfg *environmentConfig) error {
			cfg.params["value"] = runtime.NewInt(2)

			return wantErr
		},
	})
	if extended != nil {
		t.Fatal("expected option failure to return no environment")
	}

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected option error, got: %v", err)
	}

	if got, want := base.Params.MustGet("value"), runtime.NewInt(1); got != want {
		t.Fatalf("expected base environment to remain unchanged: got %v, want %v", got, want)
	}
}
