package compiler

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	"github.com/MontFerret/ferret/v2/pkg/file"
	parser "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type Visitor struct {
	*fql.BaseFqlParserVisitor
	Ctx *internal.CompilerContext
}

func NewVisitor(src *file.Source, errors *parser.ErrorHandler, level optimization.Level) *Visitor {
	v := new(Visitor)
	v.BaseFqlParserVisitor = new(fql.BaseFqlParserVisitor)
	v.Ctx = internal.NewCompilerContext(src, errors, level)

	return v
}

func (v *Visitor) VisitProgram(ctx *fql.ProgramContext) interface{} {
	for _, head := range ctx.AllHead() {
		v.VisitHead(head.(*fql.HeadContext))
	}

	v.Ctx.UDFs = internal.CollectUDFs(v.Ctx, ctx)
	if v.Ctx.UDFs != nil {
		v.Ctx.UDFScope = v.Ctx.UDFs.GlobalScope
	}

	v.Ctx.StmtCompiler.Compile(ctx.Body())
	v.Ctx.UDFCompiler.CompileAll()

	return nil
}

func (v *Visitor) VisitHead(ctx *fql.HeadContext) interface{} {
	if ctx == nil {
		return nil
	}

	useExpr := ctx.UseExpression()
	if useExpr == nil {
		return nil
	}

	useCtx := useExpr.Use()
	if useCtx == nil {
		return nil
	}

	nsCtx := useCtx.NamespaceIdentifier()
	if nsCtx == nil {
		return nil
	}

	aliasTok := useCtx.GetAlias()
	if aliasTok == nil {
		return nil
	}

	namespace := nsCtx.GetText()
	namespace = strings.TrimSuffix(namespace, runtime.NamespaceSeparator)

	alias := aliasTok.GetText()
	if alias == "" || namespace == "" {
		return nil
	}

	if existing, ok := v.Ctx.UseAliases[alias]; ok {
		if existing != namespace {
			v.Ctx.Errors.Add(v.Ctx.Errors.Create(parser.NameError, ctx, fmt.Sprintf("USE alias '%s' is already defined", alias)))
		}

		return nil
	}

	v.Ctx.UseAliases[alias] = namespace

	return nil
}
