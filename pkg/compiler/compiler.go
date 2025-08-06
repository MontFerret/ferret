package compiler

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/diagnostics"
	goruntime "runtime"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/file"

	"github.com/MontFerret/ferret/pkg/vm"

	"github.com/MontFerret/ferret/pkg/parser"
)

type Compiler struct{}

func New(setters ...Option) *Compiler {
	c := &Compiler{}
	opts := &Options{}

	for _, setter := range setters {
		setter(opts)
	}

	return c
}

func (c *Compiler) Compile(src *file.Source) (program *vm.Program, err error) {
	if src.Empty() {
		return nil, diagnostics.NewEmptyQueryErr(src)
	}

	errorHandler := diagnostics.NewErrorHandler(src, 10)

	defer func() {
		if r := recover(); r != nil {
			var e *CompilationError

			buf := make([]byte, 1024)
			n := goruntime.Stack(buf, false)
			stackTrace := string(buf[:n])

			// Find out exactly what the error was and add the e
			switch x := r.(type) {
			case string:
				e = diagnostics.NewInternalErr(src, x+"\n"+stackTrace)
			case error:
				e = diagnostics.NewInternalErrWith(src, "unknown panic\n"+stackTrace, x)
			default:
				e = diagnostics.NewInternalErr(src, "unknown panic\n"+stackTrace)
			}

			errorHandler.Add(e)

			program = nil
			err = errorHandler.Unwrap()
		}
	}()

	l := NewVisitor(src, errorHandler)
	tokenHistory := diagnostics.NewTokenHistory(10)
	p := parser.New(src.Content(), func(stream antlr.TokenStream) antlr.TokenStream {
		return diagnostics.NewTrackingTokenStream(stream, tokenHistory)
	})
	p.AddErrorListener(diagnostics.NewErrorListener(src, l.Ctx.Errors, tokenHistory))
	p.Visit(l)

	if l.Ctx.Errors.HasErrors() {
		return nil, l.Ctx.Errors.Unwrap()
	}

	program = &vm.Program{}
	program.Bytecode = l.Ctx.Emitter.Bytecode()
	program.Constants = l.Ctx.Symbols.Constants()
	program.CatchTable = l.Ctx.CatchTable.All()
	program.Registers = l.Ctx.Registers.Size()
	program.Params = l.Ctx.Symbols.Params()
	program.Functions = l.Ctx.Symbols.Functions()
	program.Labels = l.Ctx.Emitter.Labels()

	return program, err
}

func (c *Compiler) MustCompile(src *file.Source) *vm.Program {
	program, err := c.Compile(src)

	if err != nil {
		panic(err)
	}

	return program
}
