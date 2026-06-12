package compiler

import (
	goruntime "runtime"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"

	"github.com/antlr4-go/antlr/v4"

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
func (c *Compiler) Compile(src *source.Source) (program *bytecode.Program, err error) {
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

	tokenHistory := parserd.NewTokenHistory(64)
	p := parser.New(src.Content(), func(stream antlr.TokenStream) antlr.TokenStream {
		return parserd.NewTrackingTokenStream(stream, tokenHistory)
	})

	// Remove all default error listeners
	p.RemoveErrorListeners()
	// Add custom error listener
	p.AddErrorListener(parserd.NewErrorListener(src, errorHandler, tokenHistory))
	p.Program()

	if errorHandler.HasErrors() {
		return nil, errorHandler.Unwrap()
	}

	level := c.opts.Level
	if c.opts.DebugInfo {
		level = optimization.LevelNone
	}

	l := NewVisitor(src, errorHandler, level)
	l.Session.Program.DebugInfo = c.opts.DebugInfo
	p.Visit(l)

	if errorHandler.HasErrors() {
		return nil, errorHandler.Unwrap()
	}

	var udfs []bytecode.UDF

	if l.Session.Program.UDFs != nil {
		udfs = l.Session.Program.UDFs.Metadata()
	}

	registers := l.Session.Function.Registers.Size()

	for _, udf := range udfs {
		if udf.Registers > registers {
			registers = udf.Registers
		}
	}

	program = &bytecode.Program{
		ISAVersion: bytecode.Version,
		Functions: bytecode.Functions{
			Host:        l.Session.Program.HostFunctions.All(),
			UserDefined: udfs,
		},
		Metadata: bytecode.Metadata{
			CompilerVersion:        Version,
			OptimizationLevel:      int(level),
			AggregatePlans:         l.Session.Program.AggregatePlans(),
			AggregateSelectorSlots: l.Session.Program.Emitter.AggregateSelectorSlots(),
			CallArgumentSpans:      l.Session.Program.Emitter.CallArgumentSpans(),
			MatchFailTargets:       l.Session.Program.Emitter.MatchFailTargets(),
			DebugSpans:             l.Session.Program.Emitter.Spans(),
			DebugPoints:            l.Session.Program.DebugPoints,
			Labels:                 l.Session.Program.Emitter.Labels(),
		},
		Source:     src,
		Bytecode:   l.Session.Program.Emitter.Bytecode(),
		Constants:  l.Session.Function.Symbols.Constants(),
		CatchTable: l.Session.Program.CatchTable.All(),
		Registers:  registers,
		Params:     l.Session.Program.HostParams.Names(),
	}

	if err := optimization.Run(program, level); err != nil {
		return nil, err
	}

	return program, err
}
