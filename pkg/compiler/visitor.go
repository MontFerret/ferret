package compiler

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	parser "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type Visitor struct {
	*fql.BaseFqlParserVisitor
	Session  *internal.CompilationSession
	Frontend *internal.CompilationFrontend
}

func NewVisitor(src *source.Source, errors *parser.ErrorHandler, level optimization.Level) *Visitor {
	v := new(Visitor)
	v.BaseFqlParserVisitor = new(fql.BaseFqlParserVisitor)
	v.Session = internal.NewCompilationSession(src, errors, level)
	v.Frontend = internal.NewCompilationFrontend(v.Session)

	return v
}

func (v *Visitor) VisitProgram(ctx *fql.ProgramContext) interface{} {
	for _, head := range ctx.AllHead() {
		v.VisitHead(head.(*fql.HeadContext))
	}

	v.Frontend.UDFCatalog.BuildCatalog(ctx)
	if ctx != nil {
		if body, ok := ctx.Body().(*fql.BodyContext); ok {
			v.Frontend.CaptureAnalyzer.AnalyzeProgram(body)
		}
	}
	v.Frontend.Statements.Compile(ctx.Body())
	v.Frontend.UDFs.CompileAll()

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

	if existing, ok := v.Session.Program.UseAliases[alias]; ok {
		if existing != namespace {
			v.Session.Program.Errors.Add(v.Session.Program.Errors.Create(parser.NameError, ctx, fmt.Sprintf("USE alias '%s' is already defined", alias)))
		}

		return nil
	}

	v.Session.Program.UseAliases[alias] = namespace

	return nil
}
