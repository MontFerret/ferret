package internal

import (
	"io"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/file"
)

type engine struct {
	*context
	trivia     *triviaEmitter
	list       *listFormatter
	program    *programFormatter
	statement  *statementFormatter
	expression *expressionFormatter
	literal    *literalFormatter
	member     *memberFormatter
	clause     *clauseFormatter
	values     *valueFormatter
}

type context struct {
	opts *Options
	p    *printer
	src  *file.Source
}

func newEngine(src *file.Source, out io.Writer, opts *Options) *engine {
	ctx := &context{
		opts: opts,
		p:    newPrinter(out, opts),
		src:  src,
	}

	e := &engine{context: ctx}
	e.trivia = &triviaEmitter{engine: e}
	e.list = &listFormatter{engine: e}
	e.program = &programFormatter{engine: e}
	e.statement = &statementFormatter{engine: e}
	e.expression = &expressionFormatter{engine: e}
	e.literal = &literalFormatter{engine: e}
	e.member = &memberFormatter{engine: e}
	e.clause = &clauseFormatter{engine: e}
	e.values = &valueFormatter{engine: e}

	return e
}

func (e *engine) Err() error {
	return e.p.Err()
}

func (e *engine) writeKeyword(val string) {
	e.p.write(applyCase(e.opts.caseMode, val))
}

func (e *engine) renderInline(fn func(p *printer)) (string, bool) {
	var b strings.Builder

	p := newPrinter(&b, e.opts)
	p.forceSingleLine = true
	fn(p)

	return b.String(), !p.sawHardNewline
}
