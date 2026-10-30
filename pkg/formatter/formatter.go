package formatter

import (
	"io"
	goruntime "runtime"

	"github.com/antlr4-go/antlr/v4"
	"github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/source"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/formatter/internal"
	"github.com/MontFerret/ferret/v2/pkg/internal/sourceoffset"
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

// Format writes canonical FQL without resolving host functions or executing it.
// Diagnostic spans use zero-based, half-open UTF-8 byte offsets into src.
// Recovered panics return diagnostics.UnexpectedError; ordinary writer errors
// retain their identity. Any error may leave partial output, which must not be
// used as a successfully formatted program. The formatter remains reusable.
func (fmt *Formatter) Format(out io.Writer, src source.Source) (err error) {
	if src.Empty() {
		return parserd.NewEmptyQueryError(src)
	}

	errorHandler := parserd.NewErrorHandler(src, 10)
	var offsets sourceoffset.Offsets

	defer func() {
		if r := recover(); r != nil {
			var e *diagnostics.Diagnostic

			buf := make([]byte, 1024)
			n := goruntime.Stack(buf, false)
			stackTrace := string(buf[:n])

			switch x := r.(type) {
			case string:
				e = diagnostics.NewUnexpectedError(src, x+"\n"+stackTrace)
			case error:
				e = diagnostics.NewUnexpectedErrorWith(src, "unhandled panic\n"+stackTrace, x)
			default:
				e = diagnostics.NewUnexpectedError(src, "unhandled panic\n"+stackTrace)
			}

			// An internal failure must survive even if syntax diagnostics already
			// exhausted the handler's ordinary error limit.
			errorHandler.Errors().Add(e)
			err = errorHandler.Unwrap()
		}

		// Only handler-owned annotations use parser coordinates. Writer errors
		// and panic causes may carry their own locations and must remain intact.
		for _, diagnostic := range errorHandler.Errors().Errors() {
			for i := range diagnostic.Spans {
				diagnostic.Spans[i].Span = offsets.Span(diagnostic.Spans[i].Span)
			}
		}
	}()

	offsets = sourceoffset.New(src.Content())

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

	l := internal.NewVisitor(src, out, &fmt.config, offsets)
	p.Visit(l)

	if errorHandler.HasErrors() {
		return errorHandler.Unwrap()
	}

	return l.Err()
}
