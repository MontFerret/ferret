package compiler

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/pkg/errors"
)

type errorListener struct {
	*antlr.DiagnosticErrorListener
}

func newErrorListener() antlr.ErrorListener {
	return &errorListener{
		antlr.NewDiagnosticErrorListener(false),
	}
}

func (d *errorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (d *errorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}

func (d *errorListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	panic(errors.Errorf("%s at %d:%d", msg, line, column))
}
