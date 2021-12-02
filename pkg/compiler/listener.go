package compiler

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/pkg/errors"
)

type errorListener struct {
	*antlr.DiagnosticErrorListener
}

func newErrorListener() *errorListener {
	return &errorListener{
		antlr.NewDiagnosticErrorListener(false),
	}
}

func (d *errorListener) ReportAttemptingFullContext(_ antlr.Parser, _ *antlr.DFA, _, _ int, _ *antlr.BitSet, _ antlr.ATNConfigSet) {
}

func (d *errorListener) ReportContextSensitivity(_ antlr.Parser, _ *antlr.DFA, _, _, _ int, _ antlr.ATNConfigSet) {
}

func (d *errorListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	panic(errors.Errorf("%s at %d:%d", msg, line, column))
}
