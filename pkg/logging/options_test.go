package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/ziflex/go-options"
)

func mustApplyConfig(t testing.TB, setters ...Option) config {
	t.Helper()

	cfg, err := applyOptions(setters...)
	if err != nil {
		t.Fatalf("applyOptions() error = %v", err)
	}

	return cfg
}

func mustNewLogger(t testing.TB, setters ...Option) Logger {
	t.Helper()

	logger, err := New(setters...)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	return logger
}

func mustNewLoggerFrom(t testing.TB, base Logger, setters ...Option) Logger {
	t.Helper()

	logger, err := NewFrom(base, setters...)
	if err != nil {
		t.Fatalf("NewFrom() error = %v", err)
	}

	return logger
}

func decodeLogEntry(t testing.TB, output *bytes.Buffer) map[string]any {
	t.Helper()

	entry := make(map[string]any)
	if err := json.Unmarshal(bytes.TrimSpace(output.Bytes()), &entry); err != nil {
		t.Fatalf("decode log entry: %v", err)
	}

	return entry
}

func TestLoggingOptions(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg := mustApplyConfig(t)

		if cfg.writer != io.Discard {
			t.Fatalf("writer = %T, want io.Discard", cfg.writer)
		}

		if cfg.fields != nil {
			t.Fatalf("fields = %v, want nil", cfg.fields)
		}

		if cfg.level != ErrorLevel {
			t.Fatalf("level = %v, want %v", cfg.level, ErrorLevel)
		}

		if cfg.hasLevel {
			t.Fatal("hasLevel = true, want false")
		}
	})

	t.Run("options apply in order", func(t *testing.T) {
		var firstWriter bytes.Buffer
		var secondWriter bytes.Buffer
		cfg := mustApplyConfig(
			t,
			WithWriter(&firstWriter),
			nil,
			WithWriter(&secondWriter),
			WithWriter(nil),
			WithField("shared", "first"),
			WithFields(nil),
			WithFields(map[string]any{}),
			WithFields(map[string]any{
				"batch":  "value",
				"shared": "second",
			}),
			WithField("shared", "last"),
			WithLevel(InfoLevel),
			WithLevel(DebugLevel),
		)

		if cfg.writer != &secondWriter {
			t.Fatalf("writer = %T, want second writer", cfg.writer)
		}

		if len(cfg.fields) != 2 {
			t.Fatalf("fields = %v, want two fields", cfg.fields)
		}

		if cfg.fields["batch"] != "value" {
			t.Fatalf("batch field = %v, want %q", cfg.fields["batch"], "value")
		}

		if cfg.fields["shared"] != "last" {
			t.Fatalf("shared field = %v, want %q", cfg.fields["shared"], "last")
		}

		if cfg.level != DebugLevel {
			t.Fatalf("level = %v, want %v", cfg.level, DebugLevel)
		}

		if !cfg.hasLevel {
			t.Fatal("hasLevel = false, want true")
		}
	})
}

func TestWithLevelValidation(t *testing.T) {
	valid := []LogLevel{
		TraceLevel,
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
		FatalLevel,
		PanicLevel,
		Disabled,
	}

	for _, level := range valid {
		t.Run(level.String(), func(t *testing.T) {
			cfg, err := applyOptions(WithLevel(level))
			if err != nil {
				t.Fatalf("applyOptions() error = %v", err)
			}

			if cfg.level != level {
				t.Fatalf("level = %v, want %v", cfg.level, level)
			}

			if !cfg.hasLevel {
				t.Fatal("hasLevel = false, want true")
			}
		})
	}

	invalid := []struct {
		name      string
		wantValue string
		level     LogLevel
	}{
		{
			name:      "below minimum",
			level:     TraceLevel - 1,
			wantValue: "-2",
		},
		{
			name:      "no level",
			level:     NoLevel,
			wantValue: "",
		},
		{
			name:      "above maximum",
			level:     Disabled + 1,
			wantValue: "8",
		},
	}

	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := applyOptions(WithLevel(tt.level))
			if err == nil {
				t.Fatal("applyOptions() error = nil, want validation error")
			}

			if cfg.level != ErrorLevel || cfg.hasLevel {
				t.Fatalf("config = %+v, want default level without explicit level", cfg)
			}

			var got options.ValidationError
			if !errors.As(err, &got) {
				t.Fatalf("applyOptions() error = %T, want options.ValidationError", err)
			}

			if got.Field != "" || got.Value != tt.wantValue {
				t.Fatalf("validation error = %+v, want unnamed value %q", got, tt.wantValue)
			}

			var reason options.ValidationError
			if !errors.As(got.Reason, &reason) {
				t.Fatalf("validation reason = %T, want options.ValidationError", got.Reason)
			}

			if reason.Field != "" || reason.Value != tt.wantValue {
				t.Fatalf("nested validation error = %+v, want unnamed value %q", reason, tt.wantValue)
			}

			const wantReason = "must be one of [trace debug info warn error fatal panic disabled]"
			if reason.Reason == nil || reason.Reason.Error() != wantReason {
				t.Fatalf("nested validation reason = %v, want %q", reason.Reason, wantReason)
			}
		})
	}
}

func TestNewAppliesLoggingOptions(t *testing.T) {
	var output bytes.Buffer
	logger := mustNewLogger(
		t,
		WithWriter(&output),
		WithField("single", "one"),
		WithFields(map[string]any{"batch": "two"}),
		WithLevel(DebugLevel),
	)

	logger.Debug().Msg("configured")

	entry := decodeLogEntry(t, &output)
	if entry["level"] != "debug" {
		t.Fatalf("level = %v, want %q", entry["level"], "debug")
	}

	if entry["message"] != "configured" {
		t.Fatalf("message = %v, want %q", entry["message"], "configured")
	}

	if entry["single"] != "one" {
		t.Fatalf("single field = %v, want %q", entry["single"], "one")
	}

	if entry["batch"] != "two" {
		t.Fatalf("batch field = %v, want %q", entry["batch"], "two")
	}
}

func TestNewFromPreservesBaseLoggerBehavior(t *testing.T) {
	var baseOutput bytes.Buffer
	var ignoredOutput bytes.Buffer
	base := mustNewLogger(
		t,
		WithWriter(&baseOutput),
		WithField("base", "value"),
		WithLevel(WarnLevel),
	)
	unchanged, err := NewFrom(base)
	if err != nil {
		t.Fatalf("NewFrom() error = %v", err)
	}

	if !reflect.DeepEqual(unchanged, base) {
		t.Fatal("NewFrom() without options changed the base logger")
	}

	derived := mustNewLoggerFrom(
		t,
		base,
		WithWriter(&ignoredOutput),
		WithField("derived", "value"),
	)

	if LogLevel(derived.GetLevel()) != WarnLevel {
		t.Fatalf("level = %v, want %v", derived.GetLevel(), WarnLevel)
	}

	derived.Warn().Msg("derived")

	if ignoredOutput.Len() != 0 {
		t.Fatalf("ignored writer output = %q, want empty", ignoredOutput.Bytes())
	}

	entry := decodeLogEntry(t, &baseOutput)
	if entry["base"] != "value" {
		t.Fatalf("base field = %v, want %q", entry["base"], "value")
	}

	if entry["derived"] != "value" {
		t.Fatalf("derived field = %v, want %q", entry["derived"], "value")
	}

	overridden := mustNewLoggerFrom(t, base, WithLevel(DebugLevel))
	if LogLevel(overridden.GetLevel()) != DebugLevel {
		t.Fatalf("overridden level = %v, want %v", overridden.GetLevel(), DebugLevel)
	}
}

func TestWithInstallsDerivedLoggerInContext(t *testing.T) {
	var output bytes.Buffer
	base := mustNewLogger(
		t,
		WithWriter(&output),
		WithLevel(InfoLevel),
	)
	ctx := base.WithContext(context.Background())
	ctx, err := With(ctx, WithField("context", "value"))
	if err != nil {
		t.Fatalf("With() error = %v", err)
	}

	logger := From(ctx)
	logger.Info().Msg("contextual")

	entry := decodeLogEntry(t, &output)
	if entry["context"] != "value" {
		t.Fatalf("context field = %v, want %q", entry["context"], "value")
	}
}

func TestLoggingConstructorsReturnOptionError(t *testing.T) {
	want := errors.New("option failed")
	invalid := func(_ *config) error {
		return want
	}

	t.Run("New", func(t *testing.T) {
		logger, err := New(invalid)
		if !errors.Is(err, want) {
			t.Fatalf("New() error = %v, want error containing %v", err, want)
		}

		if !reflect.DeepEqual(logger, Logger{}) {
			t.Fatalf("New() logger = %+v, want zero logger", logger)
		}
	})

	t.Run("NewFrom", func(t *testing.T) {
		logger, err := NewFrom(Logger{}, invalid)
		if !errors.Is(err, want) {
			t.Fatalf("NewFrom() error = %v, want error containing %v", err, want)
		}

		if !reflect.DeepEqual(logger, Logger{}) {
			t.Fatalf("NewFrom() logger = %+v, want zero logger", logger)
		}
	})

	t.Run("With", func(t *testing.T) {
		ctx, err := With(context.Background(), invalid)
		if !errors.Is(err, want) {
			t.Fatalf("With() error = %v, want error containing %v", err, want)
		}

		if ctx != nil {
			t.Fatalf("With() context = %v, want nil", ctx)
		}
	})
}
