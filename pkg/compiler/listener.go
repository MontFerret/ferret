package compiler

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/pkg/errors"
)

type errorListener struct {
	*antlr.DefaultErrorListener
}

func (d *errorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	println(dfa.ToLexerString())
}

func (d *errorListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	panic(errors.Errorf("%s at %d:%d", msg, line, column))
}
