package diagnostics

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"

	"github.com/MontFerret/ferret/v2/pkg/file"
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

func (d *ErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	if recognizer == nil || dfa == nil || d.DiagnosticErrorListener == nil {
		return
	}

	if ctx := recognizer.GetParserRuleContext(); ctx != nil {
		for _, rule := range recognizer.GetRuleInvocationStack(ctx) {
			if rule == "expressionAtom" {
				return
			}
		}
	}

	d.DiagnosticErrorListener.ReportAmbiguity(recognizer, dfa, startIndex, stopIndex, exact, ambigAlts, configs)
}

func (d *ErrorListener) SyntaxError(_ antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	var offending antlr.Token

	// Get offending token
	if tok, ok := offendingSymbol.(antlr.Token); ok {
		offending = tok
	}

	if !d.handler.HasErrorOnLine(line) {
		if err := d.parseError(msg, offending); err != nil {
			d.handler.Add(err)
		}
	}
}

func (d *ErrorListener) parseError(msg string, offending antlr.Token) *diagnostics.Diagnostic {
	span := spanFromTokenSafe(offending, d.src)

	err := &diagnostics.Diagnostic{
		Kind:    SyntaxError,
		Source:  d.src,
		Message: "Syntax error: " + msg,
		Hint:    "Check your syntax. Did you forget to write something?",
		Spans: []diagnostics.ErrorSpan{
			{Span: span, Main: true},
		},
	}

	AnalyzeSyntaxError(d.src, err, d.history.Last())

	return err
}
