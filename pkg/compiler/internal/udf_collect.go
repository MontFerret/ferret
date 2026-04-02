package internal

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type UDFCatalogBuilder struct {
	ctx   *CompilationSession
	front *CompilationFrontend
}

func NewUDFCatalogBuilder(ctx *CompilationSession) *UDFCatalogBuilder {
	return &UDFCatalogBuilder{ctx: ctx}
}

func (c *UDFCatalogBuilder) BuildCatalog(program *fql.ProgramContext) {
	if c == nil || c.ctx == nil {
		return
	}

	table := core.NewUDFTable()
	table.GlobalScope = core.NewUDFScope(nil)

	c.ctx.UDFs = table
	c.ctx.UDFScope = table.GlobalScope

	if program == nil || program.Body() == nil {
		return
	}

	body, ok := program.Body().(*fql.BodyContext)
	if !ok {
		return
	}

	top := c.collectScopeFunctionsFromBody(body, table.GlobalScope)
	for _, fn := range top {
		c.collectNestedFunctions(fn)
	}

	if c.ctx.OptimizationLevel > optimization.LevelNone {
		c.pruneUnusedFunctions(body)
	}

}

func (c *UDFCatalogBuilder) collectScopeFunctionsFromBody(body *fql.BodyContext, scope *core.UDFScope) []*core.UDFInfo {
	if c == nil || c.ctx == nil || c.ctx.UDFs == nil || body == nil {
		return nil
	}

	out := make([]*core.UDFInfo, 0)

	for _, stmt := range body.AllBodyStatement() {
		if stmt == nil {
			continue
		}

		if decl := stmt.FunctionDeclaration(); decl != nil {
			if fn := c.registerFunction(scope, decl.(*fql.FunctionDeclarationContext)); fn != nil {
				out = append(out, fn)
			}
		}
	}

	return out
}

func (c *UDFCatalogBuilder) collectNestedFunctions(fn *core.UDFInfo) {
	if c == nil || c.ctx == nil || c.ctx.UDFs == nil || fn == nil || fn.Decl == nil {
		return
	}

	body := fn.Decl.FunctionBody()
	if body == nil {
		return
	}

	out := make([]*core.UDFInfo, 0)

	if block := body.FunctionBlock(); block != nil {
		for _, stmt := range block.AllFunctionStatement() {
			if stmt == nil {
				continue
			}

			if decl := stmt.FunctionDeclaration(); decl != nil {
				if nested := c.registerFunction(fn.BodyScope, decl.(*fql.FunctionDeclarationContext)); nested != nil {
					out = append(out, nested)
				}
			}
		}
	}

	for _, nested := range out {
		c.collectNestedFunctions(nested)
	}
}

func (c *UDFCatalogBuilder) registerFunction(scope *core.UDFScope, decl *fql.FunctionDeclarationContext) *core.UDFInfo {
	if c == nil || c.ctx == nil || c.ctx.UDFs == nil || scope == nil || decl == nil {
		return nil
	}

	displayName := decl.FunctionName().GetText()
	name := displayName

	if _, exists := scope.Functions[name]; exists {
		c.ctx.Errors.Add(c.ctx.Errors.Create(parserd.NameError, decl, fmt.Sprintf("Function '%s' is already defined", displayName)))
		return nil
	}

	fn := &core.UDFInfo{
		ID:          len(c.ctx.UDFs.Functions),
		Name:        name,
		DisplayName: displayName,
		Params:      c.collectFunctionParams(decl),
		Decl:        decl,
		Scope:       scope,
		BodyScope:   core.NewUDFScope(scope),
	}

	scope.Functions[name] = fn
	c.ctx.UDFs.Functions = append(c.ctx.UDFs.Functions, fn)

	return fn
}

func (c *UDFCatalogBuilder) collectFunctionParams(decl *fql.FunctionDeclarationContext) []string {
	if c == nil || c.ctx == nil || decl == nil {
		return nil
	}

	params := make([]string, 0)
	seen := make(map[string]struct{})

	list := decl.FunctionParameterList()
	if list == nil {
		return params
	}

	for _, param := range list.AllFunctionParameter() {
		if param == nil || param.Identifier() == nil {
			continue
		}

		name := param.Identifier().GetText()
		if _, exists := seen[name]; exists {
			c.ctx.Errors.Add(c.ctx.Errors.Create(parserd.NameError, param.(antlr.ParserRuleContext), fmt.Sprintf("Parameter '%s' is already defined", name)))
			continue
		}

		seen[name] = struct{}{}
		params = append(params, name)
	}

	return params
}

func (c *UDFCatalogBuilder) pruneUnusedFunctions(body *fql.BodyContext) {
	if c == nil || c.ctx == nil || c.ctx.UDFs == nil || body == nil {
		return
	}

	reachable := c.computeReachableFunctions(body)
	if len(reachable) == 0 {
		for _, fn := range c.ctx.UDFs.Functions {
			if fn != nil && fn.Scope != nil {
				delete(fn.Scope.Functions, fn.Name)
			}
		}
		c.ctx.UDFs.Functions = nil
		return
	}

	filtered := make([]*core.UDFInfo, 0, len(reachable))
	for _, fn := range c.ctx.UDFs.Functions {
		if fn == nil {
			continue
		}

		if _, ok := reachable[fn]; ok {
			filtered = append(filtered, fn)
			continue
		}

		if fn.Scope != nil {
			delete(fn.Scope.Functions, fn.Name)
		}
	}

	for i, fn := range filtered {
		fn.ID = i
	}

	c.ctx.UDFs.Functions = filtered
}

func (c *UDFCatalogBuilder) computeReachableFunctions(body *fql.BodyContext) map[*core.UDFInfo]struct{} {
	reachable := make(map[*core.UDFInfo]struct{})
	if c == nil || c.ctx == nil || c.ctx.UDFs == nil || body == nil || c.ctx.UDFs.GlobalScope == nil {
		return reachable
	}

	roots := c.collectCallsInTopLevel(body, c.ctx.UDFs.GlobalScope)
	if len(roots) == 0 {
		return reachable
	}

	callCache := make(map[*core.UDFInfo][]*core.UDFInfo)
	stack := append(make([]*core.UDFInfo, 0, len(roots)), roots...)

	for len(stack) > 0 {
		idx := len(stack) - 1
		fn := stack[idx]
		stack = stack[:idx]

		if fn == nil {
			continue
		}

		if _, ok := reachable[fn]; ok {
			continue
		}

		reachable[fn] = struct{}{}

		callees, ok := callCache[fn]
		if !ok {
			callees = c.collectCallsInFunction(fn)
			callCache[fn] = callees
		}

		for _, callee := range callees {
			if callee == nil {
				continue
			}
			if _, ok := reachable[callee]; !ok {
				stack = append(stack, callee)
			}
		}
	}

	return reachable
}

func (c *UDFCatalogBuilder) collectCallsInTopLevel(body *fql.BodyContext, scope *core.UDFScope) []*core.UDFInfo {
	if c == nil || body == nil {
		return nil
	}

	out := make(map[*core.UDFInfo]struct{})

	for _, stmt := range body.AllBodyStatement() {
		if stmt == nil {
			continue
		}

		st, ok := stmt.(*fql.BodyStatementContext)
		if !ok || st.FunctionDeclaration() != nil {
			continue
		}

		c.collectCallsInExpression(st, scope, out)
	}

	if expr := body.BodyExpression(); expr != nil {
		c.collectCallsInExpression(expr, scope, out)
	}

	return udfSetToSlice(out)
}

func (c *UDFCatalogBuilder) collectCallsInFunction(fn *core.UDFInfo) []*core.UDFInfo {
	if c == nil || fn == nil || fn.Decl == nil {
		return nil
	}

	out := make(map[*core.UDFInfo]struct{})
	scope := fn.BodyScope
	body := fn.Decl.FunctionBody()
	if body == nil {
		return nil
	}

	if arrow := body.FunctionArrow(); arrow != nil {
		c.collectCallsInExpression(arrow.Expression(), scope, out)
		return udfSetToSlice(out)
	}

	if block := body.FunctionBlock(); block != nil {
		for _, stmt := range block.AllFunctionStatement() {
			if stmt == nil {
				continue
			}

			st, ok := stmt.(*fql.FunctionStatementContext)
			if !ok || st.FunctionDeclaration() != nil {
				continue
			}

			c.collectCallsInExpression(st, scope, out)
		}

		if ret := block.FunctionReturn(); ret != nil {
			c.collectCallsInExpression(ret, scope, out)
		}
	}

	return udfSetToSlice(out)
}

func (c *UDFCatalogBuilder) collectCallsInExpression(
	node antlr.Tree,
	scope *core.UDFScope,
	out map[*core.UDFInfo]struct{},
) {
	if c == nil || node == nil || out == nil {
		return
	}

	var calls []*fql.FunctionCallContext
	findFunctionCallRefs(node, &calls)

	for _, call := range calls {
		if call == nil {
			continue
		}

		if fn, ok := c.resolveCallInScope(call, scope); ok {
			out[fn] = struct{}{}
		}
	}
}

func (c *UDFCatalogBuilder) resolveCallInScope(ctx fql.IFunctionCallContext, scope *core.UDFScope) (*core.UDFInfo, bool) {
	if c == nil || c.ctx == nil || c.ctx.UDFs == nil || ctx == nil || scope == nil {
		return nil, false
	}

	name, ok := c.front.Calls.ResolveLocalFunctionName(ctx)
	if !ok {
		return nil, false
	}

	return c.ctx.UDFs.Resolve(name, scope)
}

func variableName(ctx *fql.VariableContext) (string, antlr.Token) {
	if ctx == nil {
		return "", nil
	}

	if id := ctx.Identifier(); id != nil {
		return id.GetText(), id.GetSymbol()
	}

	if id := ctx.SafeReservedWord(); id != nil {
		if prc, ok := id.(antlr.ParserRuleContext); ok {
			return id.GetText(), prc.GetStart()
		}

		return id.GetText(), nil
	}

	return "", nil
}

func findVariableRefs(node antlr.Tree, out *[]*fql.VariableContext) {
	if node == nil || out == nil {
		return
	}

	if v, ok := node.(*fql.VariableContext); ok {
		*out = append(*out, v)
	}

	for i := 0; i < node.GetChildCount(); i++ {
		findVariableRefs(node.GetChild(i), out)
	}
}

func findAssignmentRefs(node antlr.Tree, out *[]*fql.AssignmentStatementContext) {
	if node == nil || out == nil {
		return
	}

	if stmt, ok := node.(*fql.AssignmentStatementContext); ok {
		*out = append(*out, stmt)
	}

	for i := 0; i < node.GetChildCount(); i++ {
		findAssignmentRefs(node.GetChild(i), out)
	}
}

func findFunctionCallRefs(node antlr.Tree, out *[]*fql.FunctionCallContext) {
	if node == nil || out == nil {
		return
	}

	if _, ok := node.(*fql.FunctionDeclarationContext); ok {
		return
	}

	if call, ok := node.(*fql.FunctionCallContext); ok {
		*out = append(*out, call)
	}

	for i := 0; i < node.GetChildCount(); i++ {
		findFunctionCallRefs(node.GetChild(i), out)
	}
}

func addUDFCapture(captures map[string]core.UDFCapture, order *[]string, name string, storage core.BindingStorage) {
	if captures == nil || order == nil || name == "" {
		return
	}

	capture, exists := captures[name]
	if !exists {
		captures[name] = core.UDFCapture{
			Name:    name,
			Mutable: storage == core.BindingStorageCell,
			Storage: storage,
		}
		*order = append(*order, name)
		return
	}

	if storage == core.BindingStorageCell && capture.Storage != core.BindingStorageCell {
		capture.Storage = core.BindingStorageCell
		capture.Mutable = true
		captures[name] = capture
	}
}

func orderedUDFCaptures(captures map[string]core.UDFCapture, order []string) []core.UDFCapture {
	if len(order) == 0 {
		return nil
	}

	out := make([]core.UDFCapture, 0, len(order))
	for _, name := range order {
		capture, ok := captures[name]
		if !ok {
			continue
		}
		out = append(out, capture)
	}

	return out
}

func udfSetToSlice(set map[*core.UDFInfo]struct{}) []*core.UDFInfo {
	if len(set) == 0 {
		return nil
	}

	out := make([]*core.UDFInfo, 0, len(set))
	for fn := range set {
		out = append(out, fn)
	}

	return out
}
