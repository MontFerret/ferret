package debugger

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

type debugEvalErrorListener struct {
	*antlr.DefaultErrorListener
	errs []error
}

func newDebugEvalErrorListener() *debugEvalErrorListener {
	return &debugEvalErrorListener{DefaultErrorListener: antlr.NewDefaultErrorListener()}
}

func (l *debugEvalErrorListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	l.errs = append(l.errs, fmt.Errorf("line %d:%d %s", line, column+1, msg))
}
