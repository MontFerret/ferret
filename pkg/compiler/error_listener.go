package compiler

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
)

type errorListener struct {
	*antlr.DiagnosticErrorListener

	handler *core.ErrorHandler
}

func newErrorListener(handler *core.ErrorHandler) antlr.ErrorListener {
	return &errorListener{
		DiagnosticErrorListener: antlr.NewDiagnosticErrorListener(false),
		handler:                 handler,
	}
}

func (d *errorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (d *errorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}

func (d *errorListener) SyntaxError(_ antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	d.handler.Add(core.NewSyntaxError(msg, line, column, offendingSymbol))
}
