package compiler

import (
	"errors"
	"testing"

	"github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
)

func mustNewCompiler(t testing.TB, setters ...Option) *Compiler {
	t.Helper()

	compilerInstance, err := New(setters...)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	return compilerInstance
}

func TestCompilerOptions(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		compilerInstance := mustNewCompiler(t)

		if compilerInstance.config.Level != O1 {
			t.Fatalf("optimization level = %v, want %v", compilerInstance.config.Level, O1)
		}
		if compilerInstance.config.DebugInfo {
			t.Fatal("debug info = true, want false")
		}
	})

	t.Run("optimization level", func(t *testing.T) {
		tests := []struct {
			name  string
			level OptimizationLevel
		}{
			{name: "none", level: optimization.LevelNone},
			{name: "basic", level: optimization.LevelBasic},
			{name: "full", level: optimization.LevelFull},
			{name: "aggressive", level: optimization.LevelAggressive},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				compilerInstance := mustNewCompiler(t, WithOptimizationLevel(tt.level))

				if compilerInstance.config.Level != tt.level {
					t.Fatalf("optimization level = %v, want %v", compilerInstance.config.Level, tt.level)
				}
			})
		}
	})

	t.Run("debug info", func(t *testing.T) {
		compilerInstance := mustNewCompiler(t, WithDebugInfo())

		if !compilerInstance.config.DebugInfo {
			t.Fatal("debug info = false, want true")
		}
		if compilerInstance.config.Level != O0 {
			t.Fatalf("optimization level = %v, want %v", compilerInstance.config.Level, O0)
		}
	})

	t.Run("debug info after optimization", func(t *testing.T) {
		compilerInstance := mustNewCompiler(
			t,
			WithOptimizationLevel(O1),
			WithDebugInfo(),
		)

		if !compilerInstance.config.DebugInfo {
			t.Fatal("debug info = false, want true")
		}
		if compilerInstance.config.Level != O0 {
			t.Fatalf("optimization level = %v, want %v", compilerInstance.config.Level, O0)
		}
	})

	t.Run("optimization after debug info", func(t *testing.T) {
		compilerInstance := mustNewCompiler(
			t,
			WithDebugInfo(),
			WithOptimizationLevel(O1),
		)

		if !compilerInstance.config.DebugInfo {
			t.Fatal("debug info = false, want true")
		}
		if compilerInstance.config.Level != O1 {
			t.Fatalf("optimization level = %v, want %v", compilerInstance.config.Level, O1)
		}
	})

	t.Run("repeated option", func(t *testing.T) {
		compilerInstance := mustNewCompiler(
			t,
			WithOptimizationLevel(O0),
			WithOptimizationLevel(O1),
		)

		if compilerInstance.config.Level != O1 {
			t.Fatalf("optimization level = %v, want %v", compilerInstance.config.Level, O1)
		}
	})

	t.Run("nil option", func(t *testing.T) {
		compilerInstance, err := New(nil)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		if compilerInstance == nil {
			t.Fatal("compiler = nil, want non-nil")
		}
		if compilerInstance.config.Level != O1 {
			t.Fatalf("optimization level = %v, want %v", compilerInstance.config.Level, O1)
		}
	})
}

func TestWithOptimizationLevelRejectsUnsupportedLevels(t *testing.T) {
	tests := []struct {
		name      string
		wantValue string
		level     OptimizationLevel
	}{
		{name: "below minimum", level: optimization.LevelNone - 1, wantValue: "-1"},
		{name: "above maximum", level: optimization.LevelAggressive + 1, wantValue: "4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compilerInstance, err := New(WithOptimizationLevel(tt.level))
			if err == nil {
				t.Fatal("New() error = nil, want validation error")
			}
			if compilerInstance != nil {
				t.Fatal("compiler != nil, want nil")
			}

			want := options.ValidationError{
				Field:  "optimization level",
				Value:  tt.wantValue,
				Reason: "must be one of [0 1 2 3]",
			}
			var got options.ValidationError
			if !errors.As(err, &got) {
				t.Fatalf("New() error = %T, want options.ValidationError", err)
			}
			if got != want {
				t.Fatalf("New() validation error = %+v, want %+v", got, want)
			}
		})
	}
}

func TestNewReturnsOptionValidationError(t *testing.T) {
	want := options.ValidationError{
		Field:  "compiler",
		Value:  "invalid",
		Reason: "test validation failure",
	}
	invalid := func(_ *config, report options.Report) {
		report(want)
	}

	compilerInstance, err := New(invalid)
	if err == nil {
		t.Fatal("New() error = nil, want validation error")
	}
	if compilerInstance != nil {
		t.Fatal("compiler != nil, want nil")
	}

	var got options.ValidationError
	if !errors.As(err, &got) {
		t.Fatalf("New() error = %T, want options.ValidationError", err)
	}
	if got != want {
		t.Fatalf("New() validation error = %+v, want %+v", got, want)
	}
}
