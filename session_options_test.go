package ferret

import (
	"testing"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func mustNewSessionOptionsForTest(t *testing.T, setters ...SessionOption) *sessionOptions {
	t.Helper()

	opts, err := newSessionOptions(setters)
	if err != nil {
		t.Fatalf("failed to create session options: %v", err)
	}

	return opts
}

func mustBuildEnvironmentForTest(t *testing.T, opts *sessionOptions) *vm.Environment {
	t.Helper()

	env, err := vm.NewEnvironment(opts.env)
	if err != nil {
		t.Fatalf("failed to build environment: %v", err)
	}

	return env
}

func TestNewSessionOptionsIgnoresEmptySessionParams(t *testing.T) {
	t.Parallel()

	opts := mustNewSessionOptionsForTest(
		t,
		WithSessionParam("param1", 1),
		WithSessionParams(nil),
		WithSessionParams(map[string]any{}),
	)

	if len(opts.env) != 1 {
		t.Fatalf("expected environment options to remain unchanged, got %d entries", len(opts.env))
	}

	env := mustBuildEnvironmentForTest(t, opts)
	value, ok := env.Params.Get("param1")
	if !ok {
		t.Fatal("expected param1 to remain configured")
	}

	if value != runtime.NewInt(1) {
		t.Fatalf("expected param1 to remain 1, got: %v", value)
	}
}

func TestNewSessionOptionsIgnoresEmptySessionRuntimeParams(t *testing.T) {
	t.Parallel()

	opts := mustNewSessionOptionsForTest(
		t,
		WithSessionRuntimeParam("param1", runtime.NewInt(1)),
		WithSessionRuntimeParams(nil),
		WithSessionRuntimeParams(runtime.Params{}),
	)

	if len(opts.env) != 1 {
		t.Fatalf("expected environment options to remain unchanged, got %d entries", len(opts.env))
	}

	env := mustBuildEnvironmentForTest(t, opts)
	value, ok := env.Params.Get("param1")
	if !ok {
		t.Fatal("expected param1 to remain configured")
	}

	if value != runtime.NewInt(1) {
		t.Fatalf("expected param1 to remain 1, got: %v", value)
	}
}

func TestNewSessionOptionsIgnoresEmptySessionLogFields(t *testing.T) {
	t.Parallel()

	opts := mustNewSessionOptionsForTest(
		t,
		WithSessionLogFields(map[string]any{"component": "session"}),
		WithSessionLogFields(nil),
		WithSessionLogFields(map[string]any{}),
	)

	if len(opts.logger) != 1 {
		t.Fatalf("expected logger options to remain unchanged, got %d entries", len(opts.logger))
	}
}

func TestNewSessionOptionsKeepDefaultOutputContentTypeWithNoopOptions(t *testing.T) {
	t.Parallel()

	opts := mustNewSessionOptionsForTest(
		t,
		nil,
		WithSessionParams(nil),
		WithSessionParams(map[string]any{}),
		WithSessionRuntimeParams(nil),
		WithSessionRuntimeParams(runtime.Params{}),
		WithSessionLogFields(nil),
		WithSessionLogFields(map[string]any{}),
	)

	if opts.outputContentType != encodingjson.ContentType {
		t.Fatalf("expected default output content type %q, got %q", encodingjson.ContentType, opts.outputContentType)
	}

	if len(opts.env) != 0 {
		t.Fatalf("expected no environment options to be appended, got %d entries", len(opts.env))
	}

	if len(opts.logger) != 0 {
		t.Fatalf("expected no logger options to be appended, got %d entries", len(opts.logger))
	}
}
