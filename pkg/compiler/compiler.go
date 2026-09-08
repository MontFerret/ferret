package compiler

import (
	"context"
	"errors"

	"github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
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

// Compile synchronously compiles source and checks cancellation between
// parsing, lowering, and program construction. Individual phases are not preempted.
// The context must be non-nil. The compiler remains safe for concurrent use.
func (c *Compiler) Compile(ctx context.Context, src source.Source) (program *bytecode.Program, err error) {
	if ctx == nil {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	defer func() {
		if ctxErr := ctx.Err(); ctxErr != nil {
			if err == nil {
				err = ctxErr
			} else if !errors.Is(err, ctxErr) {
				err = errors.Join(err, ctxErr)
			}

			program = nil
		}
	}()

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

		// The handler and lowering keep ANTLR offsets until compilation settles.
		// Convert the owned diagnostics in place to preserve their error tree.
		if errorHandler.HasErrors() {
			offsets := sourceByteOffsets(src)
			for _, diagnostic := range errorHandler.Errors().Errors() {
				for i := range diagnostic.Spans {
					diagnostic.Spans[i].Span = sourceByteSpan(offsets, diagnostic.Spans[i].Span)
				}
			}
		}
	}()

	level := c.config.Level
	if c.config.DebugInfo {
		level = optimization.None
	}

	visitor := runFrontend(ctx, src, errorHandler, level, c.config.DebugInfo, nil, nil)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if errorHandler.HasErrors() {
		return nil, errorHandler.Unwrap()
	}

	return buildProgram(visitor, src, level)
}

// Analyze parses and semantically analyzes source without constructing a bytecode program.
// It is safe for concurrent use by multiple goroutines. When source diagnostics
// exist, Analyze returns both a non-nil partial snapshot and a non-nil error.
func (c *Compiler) Analyze(src source.Source) (analysis *Analysis, err error) {
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

	visitor := runFrontend(context.Background(), src, errorHandler, optimization.None, false, recorder, &syntaxTokens)
	if visitor != nil {
		recorder.Sort()
	}

	analysis = buildAnalysis(src, recorder.Snapshot(), errorHandler, syntaxTokens)
	if errorHandler.HasErrors() {
		return analysis, analysisError(analysis)
	}

	return analysis, nil
}
