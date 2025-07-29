package compiler

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal"
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/file"
	"github.com/MontFerret/ferret/pkg/parser/fql"
)

type Visitor struct {
	*fql.BaseFqlParserVisitor
	Ctx *internal.CompilerContext
}

func NewVisitor(src *file.Source, errors *core.ErrorHandler) *Visitor {
	v := new(Visitor)
	v.BaseFqlParserVisitor = new(fql.BaseFqlParserVisitor)
	v.Ctx = internal.NewCompilerContext(src, errors)

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
