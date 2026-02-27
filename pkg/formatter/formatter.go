package formatter

import (
	"io"
	goruntime "runtime"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/formatter/internal"
	"github.com/MontFerret/ferret/v2/pkg/parser"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

type Formatter struct {
	opts *internal.Options
}

func New(setters ...Option) *Formatter {
	opts := internal.DefaultOptions()

	for _, setter := range setters {
		setter(opts)
	}

	return &Formatter{
		opts: opts,
	}
}

func (fmt *Formatter) Format(out io.Writer, src *file.Source) error {
	if src.Empty() {
		return parserd.NewEmptyQueryError(src)
	}

	errorHandler := parserd.NewErrorHandler(src, 10)

	defer func() {
		if r := recover(); r != nil {
			var e *diagnostics.Diagnostic

			buf := make([]byte, 1024)
			n := goruntime.Stack(buf, false)
			stackTrace := string(buf[:n])

			// Find out exactly what the error was and add the error
			switch x := r.(type) {
			case string:
				e = diagnostics.NewUnexpectedError(src, x+"\n"+stackTrace)
			case error:
				e = diagnostics.NewUnexpectedErrorWith(src, "unhandled panic\n"+stackTrace, x)
			default:
				e = diagnostics.NewUnexpectedError(src, "unhandled panic\n"+stackTrace)
			}

			errorHandler.Add(e)
		}
	}()

	tokenHistory := parserd.NewTokenHistory(10)
	p := parser.New(src.Content(), func(stream antlr.TokenStream) antlr.TokenStream {
		return parserd.NewTrackingTokenStream(stream, tokenHistory)
	})
	// Remove all default error listeners
	p.RemoveErrorListeners()
	// Add custom error listener
	p.AddErrorListener(parserd.NewErrorListener(src, errorHandler, tokenHistory))
	l := internal.NewVisitor(src, out, fmt.opts)
	p.Visit(l)

	if errorHandler.HasErrors() {
		return errorHandler.Unwrap()
	}

	return l.Err()
}
