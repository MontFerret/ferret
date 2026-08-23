package internal

import (
	"io"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	configuration interface {
		PrintWidth() uint64
		TabWidth() uint64
		SingleQuote() bool
		BracketSpacing() bool
		FormatKeyword(string) string
	}

	context struct {
		config configuration
		p      *printer
		src    *source.Source
	}

	engine struct {
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
		bindings   *bindingPatternFormatter
	}
)

func newEngine(src *source.Source, out io.Writer, config configuration) *engine {
	ctx := &context{
		config: config,
		p:      newPrinter(out, config.TabWidth()),
		src:    src,
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
	e.bindings = &bindingPatternFormatter{engine: e}

	return e
}

func (e *engine) Err() error {
	return e.p.Err()
}

func (e *engine) writeKeyword(val string) {
	e.p.write(e.config.FormatKeyword(val))
}

func (e *engine) inlineFits(inline string) bool {
	return e.inlineFitsWith(e.p, inline)
}

func (e *engine) inlineFitsWith(p *printer, inline string) bool {
	if p == nil {
		return len(inline) <= int(e.config.PrintWidth())
	}

	column := p.currentColumn()

	if p.atLineStart {
		column += int(p.tabWidth) * p.indent
	}

	return column+len(inline) <= int(e.config.PrintWidth())
}

func (e *engine) renderInline(fn func(p *printer)) (string, bool) {
	var b strings.Builder

	p := newPrinter(&b, e.config.TabWidth())
	p.forceSingleLine = true
	fn(p)

	return b.String(), !p.sawHardNewline
}
