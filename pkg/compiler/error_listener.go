package compiler

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/pkg/parser/fql"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/file"
	"github.com/MontFerret/ferret/pkg/parser"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
)

type (
	errorListener struct {
		*antlr.DiagnosticErrorListener
		src     *file.Source
		handler *core.ErrorHandler
		history *parser.TokenHistory
	}

	errorPattern struct {
		Name    string
		MatchFn func(tokens []antlr.Token) (matched bool, info map[string]string)
		Explain func(info map[string]string) (msg, hint string, span file.Span)
	}
)

func newErrorListener(src *file.Source, handler *core.ErrorHandler, history *parser.TokenHistory) antlr.ErrorListener {
	return &errorListener{
		DiagnosticErrorListener: antlr.NewDiagnosticErrorListener(false),
		src:                     src,
		handler:                 handler,
		history:                 history,
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

func (d *errorListener) extraneousError(err *CompilationError) (matched bool) {
	if !strings.Contains(err.Message, "extraneous input") {
		return false
	}

	last := d.history.Last()

	if last == nil {
		return false
	}

	span := core.SpanFromTokenSafe(last.Token(), d.src)
	err.Spans = []core.ErrorSpan{
		core.NewMainErrorSpan(span, "query must end with a value"),
	}

	err.Message = "Expected a RETURN or FOR clause at end of query"
	err.Hint = "All queries must return a value. Add a RETURN statement to complete the query."

	return true
}

func (d *errorListener) noViableAltError(err *CompilationError) bool {
	if !strings.Contains(err.Message, "viable alternative at input") {
		return false
	}

	if d.history.Size() < 2 {
		return false
	}

	// most recent (offending)
	last := d.history.Last()

	// CASE: RETURN [missing value]
	if isToken(last, "RETURN") && isKeyword(last.Token()) {
		span := core.SpanFromTokenSafe(last.Token(), d.src)

		err.Message = fmt.Sprintf("Expected expression after '%s'", last)
		err.Hint = "Did you forget to provide a value to return?"
		err.Spans = []core.ErrorSpan{
			core.NewMainErrorSpan(span, "missing return value"),
		}
		return true
	}

	// CASE: LET x = [missing value]
	//if strtoken(last.Token()) == "LET" && isIdentifier(tokens[n-2]) && t1.GetText() == "=" {
	//	varName := tokens[n-2].GetText()
	//	span := core.SpanFromTokenSafe(tokens[n-1], d.src)
	//
	//	err.Message = fmt.Sprintf("Expected expression after '=' for variable '%s'", varName)
	//	err.Hint = "Did you forget to provide a value?"
	//	err.Spans = []core.ErrorSpan{
	//		core.NewMainErrorSpan(span, "missing value"),
	//	}
	//	return true
	//}

	return false
}

func isIdentifier(token antlr.Token) bool {
	if token == nil {
		return false
	}

	tt := token.GetTokenType()

	return tt == fql.FqlLexerIdentifier || tt == fql.FqlLexerIgnoreIdentifier
}

func isKeyword(token antlr.Token) bool {
	if token == nil {
		return false
	}

	ttype := token.GetTokenType()

	// 0 is usually invalid; <EOF> is -1
	if ttype <= 0 || ttype >= len(fql.FqlLexerLexerStaticData.LiteralNames) {
		return false
	}

	lit := fql.FqlLexerLexerStaticData.LiteralNames[ttype]

	return strings.HasPrefix(lit, "'") && strings.HasSuffix(lit, "'")
}

func isToken(node *parser.TokenNode, expected string) bool {
	return strings.ToUpper(node.String()) == expected
}
