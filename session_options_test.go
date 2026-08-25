package ferret

import (
	"errors"
	"slices"
	"testing"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/api"
	"github.com/MontFerret/ferret/v2/pkg/debugger"
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

	return &opts
}

func mustBuildEnvironmentForTest(t *testing.T, opts *sessionOptions) *vm.Environment {
	t.Helper()

	env, err := vm.NewEnvironment(opts.env)
	if err != nil {
		t.Fatalf("failed to build environment: %v", err)
	}

	return env
}

func TestSessionSimpleOptionsApplyValidValues(t *testing.T) {
	t.Parallel()

	format := DebugFormatOptions{MaxDepth: 4, MaxItems: 12, MaxBytes: 2048}
	opts := mustNewSessionOptionsForTest(
		t,
		WithDebugFormat(format),
		WithOutputContentType("  application/custom \n"),
	)

	if opts.debugFormat != format {
		t.Fatalf("debug format = %+v, want %+v", opts.debugFormat, format)
	}

	if opts.outputContentType != "application/custom" {
		t.Fatalf("output content type = %q", opts.outputContentType)
	}
}

func TestSessionSimpleOptionsReturnStructuredValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		option SessionOption
		name   string
		field  string
		value  string
		reason string
	}{
		{
			name:   "debug format",
			field:  "debug format",
			value:  "{0 8 1024}",
			reason: "debug format limits must be positive",
			option: WithDebugFormat(DebugFormatOptions{
				MaxDepth: 0,
				MaxItems: 8,
				MaxBytes: 1024,
			}),
		},
		{
			name:   "output content type",
			field:  "output content type",
			value:  `" \t\n "`,
			reason: "must not be blank",
			option: WithOutputContentType(" \t\n "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := newSessionOptions([]SessionOption{tt.option})
			if err == nil {
				t.Fatal("expected validation error")
			}

			var validationErr gooptions.ValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("expected options.ValidationError, got %T: %v", err, err)
			}

			if validationErr.Field != tt.field {
				t.Fatalf("validation field = %q, want %q", validationErr.Field, tt.field)
			}

			if validationErr.Value != tt.value {
				t.Fatalf("validation value = %q, want %q", validationErr.Value, tt.value)
			}

			if validationErr.Reason == nil || validationErr.Reason.Error() != tt.reason {
				t.Fatalf("validation reason = %v, want %q", validationErr.Reason, tt.reason)
			}
		})
	}
}

func TestInvalidSessionBuilderOptionDoesNotMutateConfig(t *testing.T) {
	t.Parallel()

	config := defaultSessionOptions()
	format := config.debugFormat
	config.outputContentType = "existing"

	for _, option := range []SessionOption{
		WithDebugFormat(DebugFormatOptions{MaxDepth: 0, MaxItems: 8, MaxBytes: 1024}),
		WithOutputContentType(" \t "),
	} {
		target := newSessionOptionTarget(1)
		if err := option(target); err != nil {
			t.Fatalf("failed to unwrap Ferret session option: %v", err)
		}

		if _, err := gooptions.ApplyTo(config, target.setters...); err == nil {
			t.Fatal("expected invalid option to fail")
		}
	}

	if config.debugFormat != format {
		t.Fatalf("invalid debug format mutated the config to %+v", config.debugFormat)
	}

	if config.outputContentType != "existing" {
		t.Fatalf("invalid output content type mutated the config to %q", config.outputContentType)
	}
}

func TestNewSessionOptionsAppliesAllOptionsAndJoinsFailures(t *testing.T) {
	t.Parallel()

	firstErr := errors.New("first session option failed")
	secondErr := errors.New("second session option failed")
	var calls []string

	_, err := newSessionOptions([]SessionOption{
		func(api.SessionOptions) error {
			calls = append(calls, "first")

			return firstErr
		},
		nil,
		func(api.SessionOptions) error {
			calls = append(calls, "middle")

			return nil
		},
		WithOutputContentType(" \t "),
		func(api.SessionOptions) error {
			calls = append(calls, "second")

			return secondErr
		},
	})
	if err == nil {
		t.Fatal("expected joined session option failures")
	}

	if !slices.Equal(calls, []string{"first", "middle", "second"}) {
		t.Fatalf("session option calls = %v", calls)
	}

	if !errors.Is(err, firstErr) || !errors.Is(err, secondErr) {
		t.Fatalf("expected both sentinel failures, got %v", err)
	}

	var validationErr gooptions.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected options.ValidationError, got %T: %v", err, err)
	}

	if validationErr.Field != "output content type" {
		t.Fatalf("validation field = %q", validationErr.Field)
	}

	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("expected joined error, got %T", err)
	}

	failures := joined.Unwrap()
	if len(failures) != 3 || failures[0] != firstErr || failures[2] != secondErr {
		t.Fatalf("unexpected joined failure order: %v", failures)
	}
}

func TestSessionParamsPreserveConversionErrorIdentity(t *testing.T) {
	t.Parallel()

	_, err := newSessionOptions([]SessionOption{
		WithSessionParams(map[string]any{"unsupported": make(chan int)}),
	})
	if err == nil {
		t.Fatal("expected unsupported session param to fail")
	}

	if !errors.Is(err, runtime.ErrInvalidType) {
		t.Fatalf("expected runtime.ErrInvalidType, got %v", err)
	}
}

func TestSessionOptionFailureDoesNotAcquirePlanCapacity(t *testing.T) {
	t.Parallel()

	engine := mustNewEngine(t, WithMaxActiveSessions(1), WithMaxVMsPerPlan(1))
	defer func() { _ = engine.Close() }()

	plan := mustCompilePlan(t, engine, coverageValidQuery)
	defer func() { _ = plan.Close() }()

	failed, err := plan.NewSession(t.Context(), WithOutputContentType(" \t "))
	if failed != nil {
		_ = failed.Close()

		t.Fatal("expected invalid option not to return a session")
	}

	if err == nil {
		t.Fatal("expected invalid session option to fail")
	}

	session, err := plan.NewSession(t.Context())
	if err != nil {
		t.Fatalf("expected capacity after option failure, got %v", err)
	}

	if err := session.Close(); err != nil {
		t.Fatalf("close session: %v", err)
	}
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

	if opts.debugFormat != debugger.DefaultFormatOptions() {
		t.Fatalf("expected default debug format, got %+v", opts.debugFormat)
	}

	if len(opts.env) != 0 {
		t.Fatalf("expected no environment options to be appended, got %d entries", len(opts.env))
	}

	if len(opts.logger) != 0 {
		t.Fatalf("expected no logger options to be appended, got %d entries", len(opts.logger))
	}
}
