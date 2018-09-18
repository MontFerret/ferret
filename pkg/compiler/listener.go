package compiler

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/pkg/errors"
)

type errorListener struct {
	*antlr.DefaultErrorListener
}

func (d *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	panic(errors.Errorf("%s at %d:%d", msg, line, column))
}
