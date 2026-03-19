package compiler

import (
	goruntime "runtime"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/file"

	"github.com/MontFerret/ferret/v2/pkg/parser"
)

const Version = "2.0.0"

// Compiler translates FQL source code into bytecode programs.
//
// A Compiler is immutable after construction and safe for concurrent use.
// Multiple goroutines can call Compile on the same Compiler instance.
type Compiler struct {
	opts *options
}

// New creates a compiler with optional configuration.
//
// The returned compiler is immutable and can be shared safely across goroutines.
func New(setters ...Option) *Compiler {
	c := &Compiler{
		opts: &options{
			Level: optimization.LevelBasic,
		},
	}

	for _, setter := range setters {
		setter(c.opts)
	}

	return c
}

// Compile parses and compiles a source into a bytecode program.
//
// Compile is safe for concurrent use by multiple goroutines.
func (c *Compiler) Compile(src *file.Source) (program *bytecode.Program, err error) {
	if src.Empty() {
		return nil, parserd.NewEmptyQueryError(src)
	}

	errorHandler := parserd.NewErrorHandler(src, 10)

	defer func() {
		if r := recover(); r != nil {
			var e *diagnostics.Diagnostic

			buf := make([]byte, 1024)
			n := goruntime.Stack(buf, false)
			stackTrace := string(buf[:n])

			// Find out exactly what the error was and add the e
			switch x := r.(type) {
			case string:
				e = diagnostics.NewUnexpectedError(src, x+"\n"+stackTrace)
			case error:
				e = diagnostics.NewUnexpectedErrorWith(src, "unhandled panic\n"+stackTrace, x)
			default:
				e = diagnostics.NewUnexpectedError(src, "unhandled panic\n"+stackTrace)
			}

			errorHandler.Add(e)

			program = nil
			err = errorHandler.Unwrap()
		}
	}()

	l := NewVisitor(src, errorHandler, c.opts.Level)
	tokenHistory := parserd.NewTokenHistory(10)
	p := parser.New(src.Content(), func(stream antlr.TokenStream) antlr.TokenStream {
		return parserd.NewTrackingTokenStream(stream, tokenHistory)
	})
	// Remove all default error listeners
	p.RemoveErrorListeners()
	// Add custom error listener
	p.AddErrorListener(parserd.NewErrorListener(src, l.Ctx.Errors, tokenHistory))
	p.Visit(l)

	if l.Ctx.Errors.HasErrors() {
		return nil, l.Ctx.Errors.Unwrap()
	}

	var udfs []bytecode.UDF
	if l.Ctx.UDFs != nil {
		udfs = l.Ctx.UDFs.Metadata()
	}

	registers := l.Ctx.Registers.Size()
	for _, udf := range udfs {
		if udf.Registers > registers {
			registers = udf.Registers
		}
	}

	program = &bytecode.Program{
		ISAVersion: bytecode.Version,
		Functions: bytecode.Functions{
			Host:        l.Ctx.Symbols.Functions(),
			UserDefined: udfs,
		},
		Metadata: bytecode.Metadata{
			CompilerVersion:        Version,
			OptimizationLevel:      int(c.opts.Level),
			AggregatePlans:         l.Ctx.AggregatePlans(),
			AggregateSelectorSlots: l.Ctx.Emitter.AggregateSelectorSlots(),
			DebugSpans:             l.Ctx.Emitter.Spans(),
			Labels:                 l.Ctx.Emitter.Labels(),
		},
		Source:     src,
		Bytecode:   l.Ctx.Emitter.Bytecode(),
		Constants:  l.Ctx.Symbols.Constants(),
		CatchTable: l.Ctx.CatchTable.All(),
		Registers:  registers,
		Params:     l.Ctx.Symbols.Params(),
	}

	if err := optimization.Run(program, c.opts.Level); err != nil {
		return nil, err
	}

	return program, err
}
