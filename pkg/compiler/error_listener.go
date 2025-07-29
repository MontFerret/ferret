package compiler

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/file"
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
	message, hint := core.ExplainSyntaxError(msg)

	d.handler.Add(&core.CompilationError{
		Message:  message,
		Kind:     core.SyntaxError,
		Location: d.findErrorLocation(offendingSymbol, e),
		Hint:     hint,
		Cause:    nil,
	})
}

func (d *errorListener) findErrorLocation(offendingSymbol interface{}, e antlr.RecognitionException) *file.Location {
	line := 0
	column := 0
	start := 0
	end := 0

	if token, ok := offendingSymbol.(antlr.Token); ok {
		line = token.GetLine() - 1
		column = token.GetColumn()
		start = token.GetStart()
		end = token.GetStop()
	}

	if line < 0 {
		line = 0
	}

	if column < 0 {
		column = 0
	}

	if start < 0 {
		start = 0
	}

	if end < 0 {
		end = 0
	}

	return file.NewLocation(line, column, start, end)
}
