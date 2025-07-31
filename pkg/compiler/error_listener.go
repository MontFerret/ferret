package compiler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/file"
	"github.com/MontFerret/ferret/pkg/parser"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
)

type (
	errorListener struct {
		*antlr.DiagnosticErrorListener
		src        *file.Source
		handler    *core.ErrorHandler
		lastTokens *parser.TokenHistory
	}

	errorPattern struct {
		Name    string
		MatchFn func(tokens []antlr.Token) (matched bool, info map[string]string)
		Explain func(info map[string]string) (msg, hint string, span file.Span)
	}
)

func newErrorListener(src *file.Source, handler *core.ErrorHandler, lastTokens *parser.TokenHistory) antlr.ErrorListener {
	return &errorListener{
		DiagnosticErrorListener: antlr.NewDiagnosticErrorListener(false),
		src:                     src,
		handler:                 handler,
		lastTokens:              lastTokens,
	}
}

func (d *errorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (d *errorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}

func (d *errorListener) SyntaxError(_ antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	var offending antlr.Token

	// Get offending token
	if tok, ok := offendingSymbol.(antlr.Token); ok {
		offending = tok
	}

	d.handler.Add(d.parseError(msg, offending))
}

func (d *errorListener) parseError(msg string, offending antlr.Token) *CompilationError {
	span := core.SpanFromTokenSafe(offending, d.src)

	err := &CompilationError{
		Kind:    SyntaxError,
		Message: "Syntax error: " + msg,
		Hint:    "Check your syntax. Did you forget to write something?",
		Spans: []core.ErrorSpan{
			{Span: span, Main: true},
		},
	}

	for _, handler := range []func(*CompilationError, antlr.Token) bool{
		d.extraneousError,
		d.noViableAltError,
	} {
		if handler(err, offending) {
			break
		}
	}

	return err
}

func (d *errorListener) extraneousError(err *CompilationError, offending antlr.Token) (matched bool) {
	if !strings.Contains(err.Message, "extraneous input") {
		return false
	}

	span := core.SpanFromTokenSafe(offending, d.src)
	err.Spans = []core.ErrorSpan{
		core.NewMainErrorSpan(span, "query must end with a value"),
	}

	err.Message = "Expected a RETURN or FOR clause at end of query"
	err.Hint = "All queries must return a value. Add a RETURN statement to complete the query."

	return true
}

func (d *errorListener) noViableAltError(err *CompilationError, offending antlr.Token) bool {
	recognizer := regexp.MustCompile("no viable alternative at input '(\\w+).+'")

	matches := recognizer.FindAllStringSubmatch(err.Message, -1)

	if len(matches) == 0 {
		return false
	}

	last := d.lastTokens.Last()
	keyword := matches[0][1]
	start := file.SkipWhitespaceForward(d.src.Content(), last.GetStop()+1)
	span := file.Span{
		Start: start,
		End:   start + len(keyword),
	}

	switch strings.ToLower(keyword) {
	case "return":
		err.Message = fmt.Sprintf("Expected expression after '%s'", keyword)
		err.Hint = fmt.Sprintf("Did you forget to provide a value after '%s'?", keyword)

		// Replace span with RETURN token’s span
		err.Spans = []core.ErrorSpan{
			core.NewMainErrorSpan(span, "missing return value"),
		}

		return true
	}

	return false
}
