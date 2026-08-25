package compiler

import (
	"errors"
	"testing"

	"github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	"github.com/MontFerret/ferret/v2/pkg/source"
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
		wantError string
		level     OptimizationLevel
	}{
		{
			name:      "below minimum",
			level:     optimization.LevelNone - 1,
			wantValue: "-1",
			wantError: "optimization level: must be one of [0 1 2 3]: value=-1",
		},
		{
			name:      "above maximum",
			level:     optimization.LevelAggressive + 1,
			wantValue: "4",
			wantError: "optimization level: must be one of [0 1 2 3]: value=4",
		},
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
			if err.Error() != tt.wantError {
				t.Fatalf("New() error = %q, want %q", err, tt.wantError)
			}

			var got options.ValidationError
			if !errors.As(err, &got) {
				t.Fatalf("New() error = %T, want options.ValidationError", err)
			}
			if got.Field != "optimization level" || got.Value != tt.wantValue {
				t.Fatalf("New() validation error = %+v, want named value %q", got, tt.wantValue)
			}

			var nested options.ValidationError
			if errors.As(got.Reason, &nested) {
				t.Fatalf("New() validation reason = %+v, want flat validation error", nested)
			}
			if got.Reason == nil || got.Reason.Error() != "must be one of [0 1 2 3]" {
				t.Fatalf("New() validation reason = %v, want supported-level error", got.Reason)
			}
		})
	}
}

func TestCompileWithOptimizationLevelOverridesOneCompilation(t *testing.T) {
	t.Parallel()

	compilerInstance := mustNewCompiler(t, WithOptimizationLevel(O1))
	overridden, err := compilerInstance.CompileWithOptimizationLevel(
		source.NewAnonymous("LET value = 1 RETURN value"),
		O0,
	)
	if err != nil {
		t.Fatalf("CompileWithOptimizationLevel() error = %v", err)
	}

	if got := overridden.Metadata.OptimizationLevel; got != int(O0) {
		t.Fatalf("override optimization level = O%d, want O%d", got, O0)
	}

	configured, err := compilerInstance.Compile(source.NewAnonymous("RETURN 1"))
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	if got := configured.Metadata.OptimizationLevel; got != int(O1) {
		t.Fatalf("configured optimization level = O%d, want O%d", got, O1)
	}
}

func TestCompileWithOptimizationLevelRejectsUnsupportedLevel(t *testing.T) {
	t.Parallel()

	compilerInstance := mustNewCompiler(t)
	program, err := compilerInstance.CompileWithOptimizationLevel(
		source.NewAnonymous("RETURN 1"),
		optimization.LevelAggressive+1,
	)
	if err == nil {
		t.Fatal("CompileWithOptimizationLevel() error = nil, want validation error")
	}

	if program != nil {
		t.Fatal("CompileWithOptimizationLevel() program != nil, want nil")
	}
}

func TestNewReturnsOptionValidationError(t *testing.T) {
	reason := errors.New("test validation failure")
	want := options.ValidationError{
		Field:  "compiler",
		Value:  "invalid",
		Reason: reason,
	}
	invalid := func(_ *config) error {
		return want
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
	if got.Field != want.Field || got.Value != want.Value || got.Reason != reason {
		t.Fatalf("New() validation error = %+v, want %+v", got, want)
	}
	if !errors.Is(err, reason) {
		t.Fatalf("New() error = %v, want reason %v", err, reason)
	}
}
