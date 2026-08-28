package formatter

import (
	"io"
	goruntime "runtime"

	"github.com/antlr4-go/antlr/v4"
	"github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/source"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/formatter/internal"
	"github.com/MontFerret/ferret/v2/pkg/parser"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

type (
	// Option configures a Formatter during construction.
	Option = options.Option[internal.Config]

	// Formatter produces canonical FQL source formatting.
	//
	// A Formatter is immutable after construction and safe for concurrent use.
	Formatter struct {
		config internal.Config
	}
)

// New creates a formatter with optional configuration. It returns any
// validation failures reported while applying the options.
//
// The returned formatter is immutable and can be shared safely across goroutines.
func New(setters ...Option) (*Formatter, error) {
	cfg, err := options.ApplyTo(internal.DefaultConfig(), setters...)
	if err != nil {
		return nil, err
	}

	return &Formatter{
		config: cfg,
	}, nil
}

func (fmt *Formatter) Format(out io.Writer, src source.Source) error {
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
		return errorHandler.Unwrap()
	}

	l := internal.NewVisitor(src, out, &fmt.config)
	p.Visit(l)

	if errorHandler.HasErrors() {
		return errorHandler.Unwrap()
	}

	return l.Err()
}
