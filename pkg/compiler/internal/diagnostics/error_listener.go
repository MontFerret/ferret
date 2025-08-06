package diagnostics

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/file"
)

type ErrorListener struct {
	*antlr.DiagnosticErrorListener
	src     *file.Source
	handler *ErrorHandler
	history *TokenHistory
}

func NewErrorListener(src *file.Source, handler *ErrorHandler, history *TokenHistory) antlr.ErrorListener {
	return &ErrorListener{
		DiagnosticErrorListener: antlr.NewDiagnosticErrorListener(false),
		src:                     src,
		handler:                 handler,
		history:                 history,
	}
}

func (d *ErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (d *ErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}

func (d *ErrorListener) SyntaxError(_ antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	var offending antlr.Token

	// Get offending token
	if tok, ok := offendingSymbol.(antlr.Token); ok {
		offending = tok
	}

	d.handler.Add(d.parseError(msg, offending))
}

func (d *ErrorListener) parseError(msg string, offending antlr.Token) *CompilationError {
	span := spanFromTokenSafe(offending, d.src)

	err := &CompilationError{
		Source:  d.src,
		Kind:    SyntaxError,
		Message: "Syntax error: " + msg,
		Hint:    "Check your syntax. Did you forget to write something?",
		Spans: []ErrorSpan{
			{Span: span, Main: true},
		},
	}

	for _, handler := range []func(*CompilationError) bool{
		d.extraneousError,
		d.noViableAltError,
	} {
		if handler(err) {
			break
		}
	}

	return err
}

func (d *ErrorListener) extraneousError(err *CompilationError) (matched bool) {
	if !strings.Contains(err.Message, "extraneous input") {
		return false
	}

	last := d.history.Last()

	if last == nil {
		return false
	}

	span := spanFromTokenSafe(last.Token(), d.src)
	err.Spans = []ErrorSpan{
		NewMainErrorSpan(span, "query must end with a value"),
	}

	err.Message = "Expected a RETURN or FOR clause at end of query"
	err.Hint = "All queries must return a value. Add a RETURN statement to complete the query."

	return true
}

func (d *ErrorListener) noViableAltError(err *CompilationError) bool {
	if !strings.Contains(err.Message, "viable alternative at input") {
		return false
	}

	if d.history.Size() < 2 {
		return false
	}

	return AnalyzeSyntaxError(d.src, err, d.history.Last())
}
