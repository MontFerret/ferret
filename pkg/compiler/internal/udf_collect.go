package internal

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (c *UDFCompiler) CollectProgram(program *fql.ProgramContext) {
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

	c.analyzeCaptures(body)
}

func (c *UDFCompiler) ResolveCall(ctx fql.IFunctionCallContext) (*core.UDFInfo, bool) {
	if c == nil || c.ctx == nil || c.ctx.UDFs == nil || c.ctx.UDFScope == nil {
		return nil, false
	}

	return c.resolveCallInScope(ctx, c.ctx.UDFScope)
}

func (c *UDFCompiler) collectScopeFunctionsFromBody(body *fql.BodyContext, scope *core.UDFScope) []*core.UDFInfo {
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

func (c *UDFCompiler) collectNestedFunctions(fn *core.UDFInfo) {
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

func (c *UDFCompiler) registerFunction(scope *core.UDFScope, decl *fql.FunctionDeclarationContext) *core.UDFInfo {
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

func (c *UDFCompiler) collectFunctionParams(decl *fql.FunctionDeclarationContext) []string {
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

func (c *UDFCompiler) analyzeCaptures(body *fql.BodyContext) {
	if c == nil || c.ctx == nil || c.ctx.UDFs == nil || body == nil {
		return
	}

	env := &udfCaptureEnv{}
	env.push()

	for _, stmt := range body.AllBodyStatement() {
		if stmt == nil {
			continue
		}

		switch {
		case stmt.VariableDeclaration() != nil:
			binding := c.ctx.BindingCompiler.captureBindingForDeclaration(stmt.VariableDeclaration())
			if binding.Name != "" {
				env.addBinding(udfCaptureBinding{
					Decl:    binding.Decl,
					Name:    binding.Name,
					Mutable: binding.Mutable,
				})
			}
		case stmt.FunctionDeclaration() != nil:
			decl := stmt.FunctionDeclaration().(*fql.FunctionDeclarationContext)
			name := decl.FunctionName().GetText()
			if fn, ok := c.ctx.UDFs.Resolve(name, c.ctx.UDFs.GlobalScope); ok {
				c.analyzeFunctionCaptures(fn, env)
			}
		}
	}
}

func (c *UDFCompiler) analyzeFunctionCaptures(fn *core.UDFInfo, env *udfCaptureEnv) {
	if c == nil || c.ctx == nil || fn == nil || env == nil || fn.Decl == nil {
		return
	}

	env.push()
	for _, p := range fn.Params {
		env.add(p)
	}

	captureSet := make(map[string]core.UDFCapture)
	captureOrder := make([]string, 0)

	body := fn.Decl.FunctionBody()
	if body != nil {
		if arrow := body.FunctionArrow(); arrow != nil {
			if expr := arrow.Expression(); expr != nil {
				c.collectAndCaptureVars(expr, env, captureSet, &captureOrder)
				c.collectAndCaptureAssignments(expr, env, captureSet, &captureOrder)
			}
		}

		if block := body.FunctionBlock(); block != nil {
			for _, stmt := range block.AllFunctionStatement() {
				if stmt == nil {
					continue
				}

				switch {
				case stmt.VariableDeclaration() != nil:
					decl := stmt.VariableDeclaration()
					if decl != nil && decl.Expression() != nil {
						c.collectAndCaptureVars(decl.Expression(), env, captureSet, &captureOrder)
						c.collectAndCaptureAssignments(decl.Expression(), env, captureSet, &captureOrder)
					}

					binding := c.ctx.BindingCompiler.captureBindingForDeclaration(decl)
					if binding.Name != "" {
						env.addBinding(udfCaptureBinding{
							Decl:    binding.Decl,
							Name:    binding.Name,
							Mutable: binding.Mutable,
						})
					}
				case stmt.AssignmentStatement() != nil:
					c.collectAndCaptureVars(stmt.AssignmentStatement(), env, captureSet, &captureOrder)
					c.collectAndCaptureAssignments(stmt.AssignmentStatement(), env, captureSet, &captureOrder)
				case stmt.FunctionDeclaration() != nil:
					decl := stmt.FunctionDeclaration().(*fql.FunctionDeclarationContext)
					name := decl.FunctionName().GetText()
					if nested, ok := fn.BodyScope.Functions[name]; ok {
						c.analyzeFunctionCaptures(nested, env)
						for _, capture := range nested.Captures {
							if env.currentHas(capture.Name) {
								continue
							}
							if _, ok := env.resolveBinding(capture.Name); !ok {
								continue
							}
							addUDFCapture(captureSet, &captureOrder, capture.Name, capture.Storage)
						}
					}
				case stmt.FunctionCallExpression() != nil:
					c.collectAndCaptureVars(stmt.FunctionCallExpression(), env, captureSet, &captureOrder)
					c.collectAndCaptureAssignments(stmt.FunctionCallExpression(), env, captureSet, &captureOrder)
				case stmt.WaitForExpression() != nil:
					c.collectAndCaptureVars(stmt.WaitForExpression(), env, captureSet, &captureOrder)
					c.collectAndCaptureAssignments(stmt.WaitForExpression(), env, captureSet, &captureOrder)
				case stmt.DispatchExpression() != nil:
					c.collectAndCaptureVars(stmt.DispatchExpression(), env, captureSet, &captureOrder)
					c.collectAndCaptureAssignments(stmt.DispatchExpression(), env, captureSet, &captureOrder)
				case stmt.ExpressionStatement() != nil:
					c.collectAndCaptureVars(stmt.ExpressionStatement(), env, captureSet, &captureOrder)
					c.collectAndCaptureAssignments(stmt.ExpressionStatement(), env, captureSet, &captureOrder)
				}
			}

			if block.FunctionReturn() != nil {
				c.collectAndCaptureVars(block.FunctionReturn(), env, captureSet, &captureOrder)
				c.collectAndCaptureAssignments(block.FunctionReturn(), env, captureSet, &captureOrder)
			}
		}
	}

	fn.Captures = orderedUDFCaptures(captureSet, captureOrder)

	env.pop()
}

func (c *UDFCompiler) collectAndCaptureVars(
	node antlr.Tree,
	env *udfCaptureEnv,
	captureSet map[string]core.UDFCapture,
	captureOrder *[]string,
) {
	if c == nil || node == nil || env == nil || captureSet == nil || captureOrder == nil {
		return
	}

	var vars []*fql.VariableContext
	findVariableRefs(node, &vars)

	for _, v := range vars {
		name, _ := variableName(v)
		if name == "" || env.currentHas(name) {
			continue
		}

		if _, ok := env.resolveBinding(name); ok {
			addUDFCapture(captureSet, captureOrder, name, core.BindingStorageValue)
		}
	}
}

func (c *UDFCompiler) collectAndCaptureAssignments(
	node antlr.Tree,
	env *udfCaptureEnv,
	captureSet map[string]core.UDFCapture,
	captureOrder *[]string,
) {
	if c == nil || c.ctx == nil || node == nil || env == nil || captureSet == nil || captureOrder == nil {
		return
	}

	var assignments []*fql.AssignmentStatementContext
	findAssignmentRefs(node, &assignments)

	for _, stmt := range assignments {
		if stmt == nil || stmt.AssignmentTarget() == nil {
			continue
		}

		name := textOfBindingIdentifier(stmt.AssignmentTarget().BindingIdentifier())
		if name == "" || env.currentHas(name) {
			continue
		}

		binding, ok := env.resolveBinding(name)
		if !ok {
			continue
		}

		storage := core.BindingStorageValue
		if binding.Mutable {
			storage = core.BindingStorageCell
			if binding.Decl != nil {
				c.ctx.BindingCompiler.PromoteDeclaration(binding.Decl)
			}
		}

		addUDFCapture(captureSet, captureOrder, name, storage)
	}
}

func (c *UDFCompiler) pruneUnusedFunctions(body *fql.BodyContext) {
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

func (c *UDFCompiler) computeReachableFunctions(body *fql.BodyContext) map[*core.UDFInfo]struct{} {
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

func (c *UDFCompiler) collectCallsInTopLevel(body *fql.BodyContext, scope *core.UDFScope) []*core.UDFInfo {
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

func (c *UDFCompiler) collectCallsInFunction(fn *core.UDFInfo) []*core.UDFInfo {
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

func (c *UDFCompiler) collectCallsInExpression(
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

func (c *UDFCompiler) resolveCallInScope(ctx fql.IFunctionCallContext, scope *core.UDFScope) (*core.UDFInfo, bool) {
	if c == nil || c.ctx == nil || c.ctx.UDFs == nil || ctx == nil || scope == nil {
		return nil, false
	}

	name, ok := lookupUDFName(ctx, c.ctx.UseAliases)
	if !ok {
		return nil, false
	}

	return c.ctx.UDFs.Resolve(name, scope)
}

func lookupUDFName(ctx fql.IFunctionCallContext, aliases map[string]string) (string, bool) {
	if ctx == nil {
		return "", false
	}

	if ns := ctx.Namespace(); ns != nil && ns.GetText() != "" {
		return "", false
	}

	fnCtx := ctx.FunctionName()
	if fnCtx == nil {
		return "", false
	}

	name := fnCtx.GetText()
	if name == "" {
		return "", false
	}

	if len(aliases) > 0 {
		if target, ok := aliases[name]; ok && target != "" {
			if strings.Contains(target, runtime.NamespaceSeparator) {
				return "", false
			}
		}
	}

	return name, true
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
