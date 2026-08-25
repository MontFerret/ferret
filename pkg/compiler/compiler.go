package compiler

import (
	"github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

const Version = "2.0.0"

// Compiler translates FQL source code into bytecode programs.
//
// A Compiler is immutable after construction and safe for concurrent use.
// Multiple goroutines can call Compile on the same Compiler instance.
type Compiler struct {
	config config
}

// New creates a compiler with optional configuration. It returns any validation
// failures reported while applying the options.
//
// The returned compiler is immutable and can be shared safely across goroutines.
func New(setters ...Option) (*Compiler, error) {
	cfg, err := options.ApplyTo(defaultConfig(), setters...)
	if err != nil {
		return nil, err
	}

	return &Compiler{
		config: cfg,
	}, nil
}

// Compile parses and compiles a source into a bytecode program.
//
// Compile is safe for concurrent use by multiple goroutines.
func (c *Compiler) Compile(src *source.Source) (program *bytecode.Program, err error) {
	return c.compile(src, c.config.Level)
}

// CompileWithOptimizationLevel compiles source using level for this invocation
// without mutating the Compiler's immutable configuration.
func (c *Compiler) CompileWithOptimizationLevel(
	src *source.Source,
	level OptimizationLevel,
) (program *bytecode.Program, err error) {
	if err := validateOptimizationLevel(level); err != nil {
		return nil, err
	}

	return c.compile(src, level)
}

func (c *Compiler) compile(
	src *source.Source,
	level OptimizationLevel,
) (program *bytecode.Program, err error) {
	if src.Empty() {
		return nil, parserd.NewEmptyQueryError(src)
	}

	errorHandler := parserd.NewErrorHandler(src, 10)

	defer func() {
		if recovered := recover(); recovered != nil {
			addRecoveredAnalysisDiagnostic(src, errorHandler, recovered)

			program = nil
			err = errorHandler.Unwrap()
		}
	}()

	if c.config.DebugInfo {
		level = optimization.LevelNone
	}

	visitor := runFrontend(src, errorHandler, level, c.config.DebugInfo, nil, nil)

	if errorHandler.HasErrors() {
		return nil, errorHandler.Unwrap()
	}

	return buildProgram(visitor, src, level)
}

// Analyze parses and semantically analyzes source without constructing a bytecode program.
// It is safe for concurrent use by multiple goroutines. When source diagnostics
// exist, Analyze returns both a non-nil partial snapshot and a non-nil error.
func (c *Compiler) Analyze(src *source.Source) (analysis *Analysis, err error) {
	errorHandler := parserd.NewErrorHandler(src, 10)
	recorder := internal.NewSemanticRecorder(src)
	var syntaxTokens []SyntaxToken

	defer func() {
		if recovered := recover(); recovered != nil {
			addRecoveredAnalysisDiagnostic(src, errorHandler, recovered)

			recorder.Sort()
			analysis = buildAnalysis(src, recorder.Snapshot(), errorHandler, syntaxTokens)
			err = analysisError(analysis)
		}
	}()

	if src.Empty() {
		errorHandler.Add(parserd.NewEmptyQueryError(src))

		analysis = buildAnalysis(src, recorder.Snapshot(), errorHandler, syntaxTokens)

		return analysis, analysisError(analysis)
	}

	visitor := runFrontend(src, errorHandler, optimization.LevelNone, false, recorder, &syntaxTokens)
	if visitor != nil {
		recorder.Sort()
	}

	analysis = buildAnalysis(src, recorder.Snapshot(), errorHandler, syntaxTokens)
	if errorHandler.HasErrors() {
		return analysis, analysisError(analysis)
	}

	return analysis, nil
}
