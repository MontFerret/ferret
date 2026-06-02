package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	// ForwardBindingIndex records source-level binding declarations before
	// lowering so unresolved uses can point at declarations that appear later.
	ForwardBindingIndex struct {
		declarations []forwardBindingDeclaration
	}

	forwardBindingDeclaration struct {
		Ctx   antlr.ParserRuleContext
		Name  string
		Scope forwardBindingScope
		Span  source.Span
	}

	forwardBindingScope struct {
		Span  source.Span
		Depth int
	}
)

// NewForwardBindingIndex creates an empty declaration index for one compile.
func NewForwardBindingIndex(_ *CompilationSession) *ForwardBindingIndex {
	return &ForwardBindingIndex{}
}

// BuildProgram indexes declarations by lexical scope without changing compiler state.
func (i *ForwardBindingIndex) BuildProgram(program *fql.ProgramContext) {
	if i == nil {
		return
	}

	i.declarations = i.declarations[:0]

	if program == nil || program.Body() == nil {
		return
	}

	body, ok := program.Body().(*fql.BodyContext)
	if !ok || body == nil {
		return
	}

	i.buildBody(body, i.scopeFromContext(body, 0))
}

// Lookup finds the nearest later declaration visible from the unresolved token.
func (i *ForwardBindingIndex) Lookup(token antlr.Token, name string) (antlr.ParserRuleContext, bool) {
	if i == nil || token == nil || name == "" || name == core.IgnorePseudoVariable || name == core.PseudoVariable {
		return nil, false
	}

	useStart := token.GetStart()
	var best *forwardBindingDeclaration

	for idx := range i.declarations {
		decl := &i.declarations[idx]
		if decl.Name != name || decl.Ctx == nil || decl.Span.Start <= useStart || !i.scopeContains(decl.Scope, useStart) {
			continue
		}

		if best == nil || decl.Scope.Depth > best.Scope.Depth || decl.Scope.Depth == best.Scope.Depth && decl.Span.Start < best.Span.Start {
			best = decl
		}
	}

	if best == nil {
		return nil, false
	}

	return best.Ctx, true
}

func (i *ForwardBindingIndex) buildBody(body *fql.BodyContext, scope forwardBindingScope) {
	if i == nil || body == nil {
		return
	}

	for _, stmt := range body.AllBodyStatement() {
		i.buildBodyStatement(stmt, scope)
	}

	i.buildBodyExpression(body.BodyExpression(), scope)
}

func (i *ForwardBindingIndex) buildBodyStatement(ctx fql.IBodyStatementContext, scope forwardBindingScope) {
	if i == nil || ctx == nil {
		return
	}

	switch {
	case ctx.VariableDeclaration() != nil:
		i.buildVariableDeclaration(ctx.VariableDeclaration(), scope)
	case ctx.FunctionDeclaration() != nil:
		i.buildFunctionDeclaration(ctx.FunctionDeclaration(), scope)
	default:
		i.collectNestedScopes(ctx, scope)
	}
}

func (i *ForwardBindingIndex) buildBodyExpression(ctx fql.IBodyExpressionContext, scope forwardBindingScope) {
	if i == nil || ctx == nil {
		return
	}

	if fe := ctx.ForExpression(); fe != nil {
		i.buildForExpression(fe, scope)
		return
	}

	if ret := ctx.ReturnExpression(); ret != nil {
		i.collectNestedScopes(ret.Expression(), scope)
	}
}

func (i *ForwardBindingIndex) buildFunctionDeclaration(ctx fql.IFunctionDeclarationContext, parent forwardBindingScope) {
	if i == nil || ctx == nil || ctx.FunctionBody() == nil {
		return
	}

	body := ctx.FunctionBody()
	scope := i.scopeFromNode(body, parent.Depth+1)

	if arrow := body.FunctionArrow(); arrow != nil {
		i.collectNestedScopes(arrow.Expression(), scope)
		return
	}

	block := body.FunctionBlock()
	if block == nil {
		return
	}

	for _, stmt := range block.AllFunctionStatement() {
		i.buildFunctionStatement(stmt, scope)
	}

	if ret := block.FunctionReturn(); ret != nil {
		i.collectNestedScopes(ret.Expression(), scope)
	}
}

func (i *ForwardBindingIndex) buildFunctionStatement(ctx fql.IFunctionStatementContext, scope forwardBindingScope) {
	if i == nil || ctx == nil {
		return
	}

	switch {
	case ctx.VariableDeclaration() != nil:
		i.buildVariableDeclaration(ctx.VariableDeclaration(), scope)
	case ctx.FunctionDeclaration() != nil:
		i.buildFunctionDeclaration(ctx.FunctionDeclaration(), scope)
	default:
		i.collectNestedScopes(ctx, scope)
	}
}

func (i *ForwardBindingIndex) buildForExpression(ctx fql.IForExpressionContext, parent forwardBindingScope) {
	if i == nil || ctx == nil {
		return
	}

	scope := i.loopScope(ctx, parent)

	if src := ctx.ForExpressionSource(); src != nil {
		i.collectNestedScopes(src, parent)
	} else {
		i.collectNestedScopes(ctx.Expression(), scope)
	}

	if val := ctx.GetValueVariable(); val != nil {
		i.recordDeclaration(textOfLoopVariable(val), i.ruleContext(val), scope)
	}

	if counter := ctx.GetCounterVariable(); counter != nil {
		i.recordDeclaration(textOfBindingIdentifier(counter), i.ruleContext(counter), scope)
	}

	for _, body := range ctx.AllForExpressionBody() {
		i.buildForExpressionBody(body, scope)
	}

	i.buildForExpressionReturn(ctx.ForExpressionReturn(), scope)
}

func (i *ForwardBindingIndex) buildForExpressionBody(ctx fql.IForExpressionBodyContext, scope forwardBindingScope) {
	if i == nil || ctx == nil {
		return
	}

	if stmt := ctx.ForExpressionStatement(); stmt != nil {
		i.buildForExpressionStatement(stmt, scope)
		return
	}

	if clause := ctx.ForExpressionClause(); clause != nil {
		i.buildForExpressionClause(clause, scope)
	}
}

func (i *ForwardBindingIndex) buildForExpressionStatement(ctx fql.IForExpressionStatementContext, scope forwardBindingScope) {
	if i == nil || ctx == nil {
		return
	}

	if decl := ctx.VariableDeclaration(); decl != nil {
		i.buildVariableDeclaration(decl, scope)
		return
	}

	i.collectNestedScopes(ctx, scope)
}

func (i *ForwardBindingIndex) buildForExpressionClause(ctx fql.IForExpressionClauseContext, scope forwardBindingScope) {
	if i == nil || ctx == nil {
		return
	}

	if collect := ctx.CollectClause(); collect != nil {
		i.buildCollectClause(collect, scope)
		return
	}

	i.collectNestedScopes(ctx, scope)
}

func (i *ForwardBindingIndex) buildForExpressionReturn(ctx fql.IForExpressionReturnContext, scope forwardBindingScope) {
	if i == nil || ctx == nil {
		return
	}

	if ret := ctx.ReturnExpression(); ret != nil {
		i.collectNestedScopes(ret.Expression(), scope)
		return
	}

	if nested := ctx.ForExpression(); nested != nil {
		i.buildForExpression(nested, scope)
	}
}

func (i *ForwardBindingIndex) buildCollectClause(ctx fql.ICollectClauseContext, scope forwardBindingScope) {
	if i == nil || ctx == nil {
		return
	}

	i.collectNestedScopes(ctx, scope)

	if grouping := ctx.CollectGrouping(); grouping != nil {
		for _, selector := range grouping.AllCollectSelector() {
			i.recordCollectSelector(selector, scope)
		}
	}

	if aggregator := ctx.CollectAggregator(); aggregator != nil {
		for _, selector := range aggregator.AllCollectAggregateSelector() {
			i.recordBindingIdentifier(selector.BindingIdentifier(), scope)
		}
	}

	if projection := ctx.CollectGroupProjection(); projection != nil {
		if selector := projection.CollectSelector(); selector != nil {
			i.recordCollectSelector(selector, scope)
		} else {
			i.recordBindingIdentifier(projection.BindingIdentifier(), scope)
		}
	}

	if counter := ctx.CollectCounter(); counter != nil {
		i.recordBindingIdentifier(counter.BindingIdentifier(), scope)
	}
}

func (i *ForwardBindingIndex) buildVariableDeclaration(ctx fql.IVariableDeclarationContext, scope forwardBindingScope) {
	if i == nil || ctx == nil {
		return
	}

	i.collectNestedScopes(ctx.Expression(), scope)
	i.recordDeclaration(bindingDeclarationName(ctx), i.declarationContext(ctx), scope)
}

func (i *ForwardBindingIndex) buildWaitForTriggerClause(ctx fql.IWaitForTriggerClauseContext, parent forwardBindingScope) {
	if i == nil || ctx == nil {
		return
	}

	scope := i.scopeFromNode(ctx, parent.Depth+1)

	if inline := ctx.WaitForTriggerInlineStatement(); inline != nil {
		i.buildWaitForTriggerInlineStatement(inline, scope)
		return
	}

	for _, stmt := range ctx.AllWaitForTriggerStatement() {
		i.buildWaitForTriggerStatement(stmt, scope)
	}
}

func (i *ForwardBindingIndex) buildWaitForTriggerStatement(ctx fql.IWaitForTriggerStatementContext, scope forwardBindingScope) {
	if i == nil || ctx == nil {
		return
	}

	if decl := ctx.VariableDeclaration(); decl != nil {
		i.buildVariableDeclaration(decl, scope)
		return
	}

	i.collectNestedScopes(ctx, scope)
}

func (i *ForwardBindingIndex) buildWaitForTriggerInlineStatement(ctx fql.IWaitForTriggerInlineStatementContext, scope forwardBindingScope) {
	if i == nil || ctx == nil {
		return
	}

	if decl := ctx.VariableDeclaration(); decl != nil {
		i.buildVariableDeclaration(decl, scope)
		return
	}

	i.collectNestedScopes(ctx, scope)
}

func (i *ForwardBindingIndex) collectNestedScopes(node antlr.Tree, scope forwardBindingScope) {
	if i == nil || node == nil {
		return
	}

	switch ctx := node.(type) {
	case *fql.FunctionDeclarationContext:
		i.buildFunctionDeclaration(ctx, scope)
		return
	case *fql.ForExpressionContext:
		i.buildForExpression(ctx, scope)
		return
	case *fql.WaitForTriggerClauseContext:
		i.buildWaitForTriggerClause(ctx, scope)
		return
	case *fql.VariableDeclarationContext:
		i.buildVariableDeclaration(ctx, scope)
		return
	}

	for childIdx := 0; childIdx < node.GetChildCount(); childIdx++ {
		i.collectNestedScopes(node.GetChild(childIdx), scope)
	}
}

func (i *ForwardBindingIndex) recordCollectSelector(selector fql.ICollectSelectorContext, scope forwardBindingScope) {
	if i == nil || selector == nil {
		return
	}

	i.recordBindingIdentifier(selector.BindingIdentifier(), scope)
}

func (i *ForwardBindingIndex) recordBindingIdentifier(id fql.IBindingIdentifierContext, scope forwardBindingScope) {
	if i == nil || id == nil {
		return
	}

	i.recordDeclaration(textOfBindingIdentifier(id), i.ruleContext(id), scope)
}

func (i *ForwardBindingIndex) recordDeclaration(name string, ctx antlr.ParserRuleContext, scope forwardBindingScope) {
	if i == nil || name == "" || name == core.IgnorePseudoVariable || ctx == nil || !i.scopeValid(scope) {
		return
	}

	i.declarations = append(i.declarations, forwardBindingDeclaration{
		Name:  name,
		Span:  parserd.SpanFromRuleContext(ctx),
		Ctx:   ctx,
		Scope: scope,
	})
}

func (i *ForwardBindingIndex) declarationContext(ctx fql.IVariableDeclarationContext) antlr.ParserRuleContext {
	if i == nil || ctx == nil {
		return nil
	}

	if id := ctx.BindingIdentifier(); id != nil {
		return i.ruleContext(id)
	}

	return i.ruleContext(ctx)
}

func (i *ForwardBindingIndex) loopScope(ctx fql.IForExpressionContext, parent forwardBindingScope) forwardBindingScope {
	scope := i.scopeFromNode(ctx, parent.Depth+1)
	if ctx == nil || ctx.ForExpressionSource() == nil {
		return scope
	}

	srcCtx := i.ruleContext(ctx.ForExpressionSource())
	if srcCtx == nil {
		return scope
	}

	srcSpan := parserd.SpanFromRuleContext(srcCtx)
	if srcSpan.End > scope.Span.Start && srcSpan.End < scope.Span.End {
		scope.Span.Start = srcSpan.End
	}

	return scope
}

func (i *ForwardBindingIndex) scopeFromNode(node any, depth int) forwardBindingScope {
	return i.scopeFromContext(i.ruleContext(node), depth)
}

func (i *ForwardBindingIndex) scopeFromContext(ctx antlr.ParserRuleContext, depth int) forwardBindingScope {
	if ctx == nil {
		return forwardBindingScope{Depth: depth}
	}

	return forwardBindingScope{
		Span:  parserd.SpanFromRuleContext(ctx),
		Depth: depth,
	}
}

func (i *ForwardBindingIndex) ruleContext(node any) antlr.ParserRuleContext {
	if ctx, ok := node.(antlr.ParserRuleContext); ok {
		return ctx
	}

	return nil
}

func (i *ForwardBindingIndex) scopeContains(scope forwardBindingScope, pos int) bool {
	return i.scopeValid(scope) && pos >= scope.Span.Start && pos < scope.Span.End
}

func (i *ForwardBindingIndex) scopeValid(scope forwardBindingScope) bool {
	return scope.Span.End > scope.Span.Start
}
