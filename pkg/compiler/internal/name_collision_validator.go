package internal

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	diag "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type (
	// NameCollisionValidator reports same-scope ambiguity between a binding
	// declaration and a source-level bare function call with the same name.
	NameCollisionValidator struct {
		ctx *CompilationSession
	}

	nameCollisionUseKind string

	nameCollisionUse struct {
		ctx  antlr.ParserRuleContext
		kind nameCollisionUseKind
	}

	nameCollisionScope struct {
		bindings map[string]nameCollisionUse
		calls    map[string]nameCollisionUse
		reported map[string]struct{}
	}
)

const (
	nameCollisionBinding nameCollisionUseKind = "binding"
	nameCollisionCall    nameCollisionUseKind = "function call"
)

// NewNameCollisionValidator creates a compiler semantic validator for local
// function-call and binding-name collisions.
func NewNameCollisionValidator(ctx *CompilationSession) *NameCollisionValidator {
	return &NameCollisionValidator{ctx: ctx}
}

func (v *NameCollisionValidator) ValidateProgram(program *fql.ProgramContext) {
	if v == nil || v.ctx == nil || program == nil || program.Body() == nil {
		return
	}

	body, ok := program.Body().(*fql.BodyContext)
	if !ok || body == nil {
		return
	}

	v.validateBody(body)
}

func (v *NameCollisionValidator) validateBody(body *fql.BodyContext) {
	if body == nil {
		return
	}

	scope := v.newScope()

	for _, stmt := range body.AllBodyStatement() {
		v.collectBodyStatement(scope, stmt)
	}

	v.collectBodyExpression(scope, body.BodyExpression())
}

func (v *NameCollisionValidator) collectBodyStatement(scope *nameCollisionScope, ctx fql.IBodyStatementContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.VariableDeclaration() != nil:
		v.collectVariableDeclaration(scope, ctx.VariableDeclaration())
	case ctx.AssignmentStatement() != nil:
		v.collectCallsInNode(scope, ctx.AssignmentStatement())
	case ctx.FunctionDeclaration() != nil:
		v.validateFunctionDeclaration(ctx.FunctionDeclaration())
	case ctx.FunctionCallExpression() != nil:
		v.collectCallsInNode(scope, ctx.FunctionCallExpression())
	case ctx.WaitForExpression() != nil:
		v.collectCallsInNode(scope, ctx.WaitForExpression())
	case ctx.DispatchExpression() != nil:
		v.collectCallsInNode(scope, ctx.DispatchExpression())
	}
}

func (v *NameCollisionValidator) collectBodyExpression(scope *nameCollisionScope, ctx fql.IBodyExpressionContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.ReturnExpression() != nil:
		v.collectCallsInNode(scope, ctx.ReturnExpression().Expression())
	case ctx.ForExpression() != nil:
		v.validateForExpression(scope, ctx.ForExpression())
	}
}

func (v *NameCollisionValidator) validateFunctionDeclaration(ctx fql.IFunctionDeclarationContext) {
	if ctx == nil || ctx.FunctionBody() == nil {
		return
	}

	scope := v.newScope()
	body := ctx.FunctionBody()

	if arrow := body.FunctionArrow(); arrow != nil {
		v.collectCallsInNode(scope, arrow.Expression())
		return
	}

	block := body.FunctionBlock()
	if block == nil {
		return
	}

	for _, stmt := range block.AllFunctionStatement() {
		v.collectFunctionStatement(scope, stmt)
	}

	if ret := block.FunctionReturn(); ret != nil {
		v.collectCallsInNode(scope, ret.Expression())
	}
}

func (v *NameCollisionValidator) collectFunctionStatement(scope *nameCollisionScope, ctx fql.IFunctionStatementContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.VariableDeclaration() != nil:
		v.collectVariableDeclaration(scope, ctx.VariableDeclaration())
	case ctx.AssignmentStatement() != nil:
		v.collectCallsInNode(scope, ctx.AssignmentStatement())
	case ctx.FunctionDeclaration() != nil:
		v.validateFunctionDeclaration(ctx.FunctionDeclaration())
	case ctx.FunctionCallExpression() != nil:
		v.collectCallsInNode(scope, ctx.FunctionCallExpression())
	case ctx.WaitForExpression() != nil:
		v.collectCallsInNode(scope, ctx.WaitForExpression())
	case ctx.DispatchExpression() != nil:
		v.collectCallsInNode(scope, ctx.DispatchExpression())
	case ctx.ExpressionStatement() != nil:
		v.collectCallsInNode(scope, ctx.ExpressionStatement())
	}
}

func (v *NameCollisionValidator) validateForExpression(parent *nameCollisionScope, ctx fql.IForExpressionContext) {
	if ctx == nil {
		return
	}

	if src := ctx.ForExpressionSource(); src != nil {
		v.collectCallsInNode(parent, src)
	}

	scope := v.newScope()

	if value := ctx.GetValueVariable(); value != nil {
		v.recordBinding(scope, textOfLoopVariable(value), v.ruleContext(value))
	}

	if counter := ctx.GetCounterVariable(); counter != nil {
		v.recordBinding(scope, textOfBindingIdentifier(counter), v.ruleContext(counter))
	}

	if ctx.ForExpressionSource() == nil {
		v.collectCallsInNode(scope, ctx.Expression())
	}

	for _, body := range ctx.AllForExpressionBody() {
		v.collectForExpressionBody(scope, body)
	}

	v.collectForExpressionReturn(scope, ctx.ForExpressionReturn())
}

func (v *NameCollisionValidator) collectForExpressionBody(scope *nameCollisionScope, ctx fql.IForExpressionBodyContext) {
	if ctx == nil {
		return
	}

	if stmt := ctx.ForExpressionStatement(); stmt != nil {
		v.collectForExpressionStatement(scope, stmt)
		return
	}

	if clause := ctx.ForExpressionClause(); clause != nil {
		v.collectForExpressionClause(scope, clause)
	}
}

func (v *NameCollisionValidator) collectForExpressionStatement(scope *nameCollisionScope, ctx fql.IForExpressionStatementContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.VariableDeclaration() != nil:
		v.collectVariableDeclaration(scope, ctx.VariableDeclaration())
	case ctx.AssignmentStatement() != nil:
		v.collectCallsInNode(scope, ctx.AssignmentStatement())
	case ctx.FunctionCallExpression() != nil:
		v.collectCallsInNode(scope, ctx.FunctionCallExpression())
	}
}

func (v *NameCollisionValidator) collectForExpressionClause(scope *nameCollisionScope, ctx fql.IForExpressionClauseContext) {
	if ctx == nil {
		return
	}

	if collect := ctx.CollectClause(); collect != nil {
		v.collectCallsInNode(scope, collect)
		v.collectCollectBindings(scope, collect)
		return
	}

	v.collectCallsInNode(scope, ctx)
}

func (v *NameCollisionValidator) collectForExpressionReturn(scope *nameCollisionScope, ctx fql.IForExpressionReturnContext) {
	if ctx == nil {
		return
	}

	if ret := ctx.ReturnExpression(); ret != nil {
		v.collectCallsInNode(scope, ret.Expression())
		return
	}

	if nested := ctx.ForExpression(); nested != nil {
		v.validateForExpression(scope, nested)
	}
}

func (v *NameCollisionValidator) collectVariableDeclaration(scope *nameCollisionScope, ctx fql.IVariableDeclarationContext) {
	if ctx == nil {
		return
	}

	v.recordBinding(scope, bindingDeclarationName(ctx), v.ruleContext(ctx))
	v.collectCallsInNode(scope, ctx.Expression())
}

func (v *NameCollisionValidator) collectCollectBindings(scope *nameCollisionScope, ctx fql.ICollectClauseContext) {
	if ctx == nil {
		return
	}

	if grouping := ctx.CollectGrouping(); grouping != nil {
		for _, selector := range grouping.AllCollectSelector() {
			v.recordCollectSelectorBinding(scope, selector)
		}
	}

	if aggregator := ctx.CollectAggregator(); aggregator != nil {
		for _, selector := range aggregator.AllCollectAggregateSelector() {
			v.recordBinding(scope, textOfBindingIdentifier(selector.BindingIdentifier()), v.ruleContext(selector.BindingIdentifier()))
		}
	}

	if projection := ctx.CollectGroupProjection(); projection != nil {
		if selector := projection.CollectSelector(); selector != nil {
			v.recordCollectSelectorBinding(scope, selector)
		} else if id := projection.BindingIdentifier(); id != nil {
			v.recordBinding(scope, textOfBindingIdentifier(id), v.ruleContext(id))
		}
	}

	if counter := ctx.CollectCounter(); counter != nil {
		v.recordBinding(scope, textOfBindingIdentifier(counter.BindingIdentifier()), v.ruleContext(counter.BindingIdentifier()))
	}
}

func (v *NameCollisionValidator) recordCollectSelectorBinding(scope *nameCollisionScope, selector fql.ICollectSelectorContext) {
	if selector == nil {
		return
	}

	id := selector.BindingIdentifier()
	v.recordBinding(scope, textOfBindingIdentifier(id), v.ruleContext(id))
}

func (v *NameCollisionValidator) collectCallsInNode(scope *nameCollisionScope, node antlr.Tree) {
	if node == nil {
		return
	}

	switch ctx := node.(type) {
	case *fql.FunctionDeclarationContext:
		v.validateFunctionDeclaration(ctx)
		return
	case *fql.ForExpressionContext:
		v.validateForExpression(scope, ctx)
		return
	case *fql.FunctionCallContext:
		v.recordFunctionCall(scope, ctx)
	}

	for i := 0; i < node.GetChildCount(); i++ {
		v.collectCallsInNode(scope, node.GetChild(i))
	}
}

func (v *NameCollisionValidator) recordFunctionCall(scope *nameCollisionScope, ctx fql.IFunctionCallContext) {
	if ctx == nil || ctx.FunctionName() == nil {
		return
	}

	if ns := ctx.Namespace(); ns != nil && ns.GetText() != "" {
		return
	}

	v.recordCall(scope, ctx.FunctionName().GetText(), v.ruleContext(ctx))
}

func (v *NameCollisionValidator) recordBinding(scope *nameCollisionScope, name string, ctx antlr.ParserRuleContext) {
	if scope == nil || name == "" || name == core.IgnorePseudoVariable {
		return
	}

	use := nameCollisionUse{ctx: ctx, kind: nameCollisionBinding}
	if _, exists := scope.bindings[name]; !exists {
		scope.bindings[name] = use
	}

	if call, exists := scope.calls[name]; exists {
		v.reportCollision(scope, name, use, call)
	}
}

func (v *NameCollisionValidator) recordCall(scope *nameCollisionScope, name string, ctx antlr.ParserRuleContext) {
	if scope == nil || name == "" {
		return
	}

	use := nameCollisionUse{ctx: ctx, kind: nameCollisionCall}
	if _, exists := scope.calls[name]; !exists {
		scope.calls[name] = use
	}

	if binding, exists := scope.bindings[name]; exists {
		v.reportCollision(scope, name, use, binding)
	}
}

func (v *NameCollisionValidator) reportCollision(scope *nameCollisionScope, name string, main, secondary nameCollisionUse) {
	if v == nil || v.ctx == nil || v.ctx.Program == nil || v.ctx.Program.Errors == nil || main.ctx == nil {
		return
	}

	if _, exists := scope.reported[name]; exists {
		return
	}

	scope.reported[name] = struct{}{}

	err := v.ctx.Program.Errors.Create(
		parserd.NameError,
		main.ctx,
		fmt.Sprintf("Variable '%s' conflicts with function call '%s'", name, name),
	)
	err.Hint = fmt.Sprintf("Rename either the variable '%s' or the function call '%s' to make the reference unambiguous.", name, name)

	if secondary.ctx != nil {
		err.Spans = []diag.ErrorSpan{
			diag.NewMainErrorSpan(parserd.SpanFromRuleContext(main.ctx), fmt.Sprintf("%s uses '%s'", main.kind, name)),
			diag.NewSecondaryErrorSpan(parserd.SpanFromRuleContext(secondary.ctx), fmt.Sprintf("%s uses '%s'", secondary.kind, name)),
		}
	}

	v.ctx.Program.Errors.Add(err)
}

func (v *NameCollisionValidator) newScope() *nameCollisionScope {
	return &nameCollisionScope{
		bindings: make(map[string]nameCollisionUse),
		calls:    make(map[string]nameCollisionUse),
		reported: make(map[string]struct{}),
	}
}

func (v *NameCollisionValidator) ruleContext(node any) antlr.ParserRuleContext {
	if ctx, ok := node.(antlr.ParserRuleContext); ok {
		return ctx
	}

	return nil
}
