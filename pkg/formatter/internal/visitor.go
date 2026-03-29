package internal

import (
	"io"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type Visitor struct {
	*fql.BaseFqlParserVisitor
	engine *engine
}

func NewVisitor(src *source.Source, out io.Writer, opts *Options) *Visitor {
	return &Visitor{
		BaseFqlParserVisitor: new(fql.BaseFqlParserVisitor),
		engine:               newEngine(src, out, opts),
	}
}

func (v *Visitor) Err() error {
	return v.engine.Err()
}

func (v *Visitor) VisitProgram(ctx *fql.ProgramContext) interface{} {
	v.engine.program.formatProgram(ctx)

	return nil
}
