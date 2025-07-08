package compiler

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal"
	"github.com/MontFerret/ferret/pkg/parser/fql"
)

type Visitor struct {
	*fql.BaseFqlParserVisitor

	Ctx *internal.CompilerContext
	Err error
	Src string
}

func NewVisitor(src string) *Visitor {
	v := new(Visitor)
	v.BaseFqlParserVisitor = new(fql.BaseFqlParserVisitor)
	v.Ctx = internal.NewCompilerContext()
	v.Src = src

	return v
}

func (v *Visitor) VisitProgram(ctx *fql.ProgramContext) interface{} {
	for _, head := range ctx.AllHead() {
		v.VisitHead(head.(*fql.HeadContext))
	}

	v.Ctx.StmtCompiler.Compile(ctx.Body())

	return nil
}

func (v *Visitor) VisitHead(_ *fql.HeadContext) interface{} {
	return nil
}
