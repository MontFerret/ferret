package internal

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func loadConstant(ctx *CompilerContext, value runtime.Value) bytecode.Operand {
	reg := ctx.Registers.Allocate()
	ctx.Emitter.EmitLoadConst(reg, ctx.Symbols.AddConstant(value))
	ctx.Types.Set(reg, valueTypeFromRuntime(value))

	return reg
}

func sortDirection(dir antlr.TerminalNode) runtime.SortDirection {
	if dir == nil {
		return runtime.SortDirectionAsc
	}

	if strings.ToLower(dir.GetText()) == "desc" {
		return runtime.SortDirectionDesc
	}

	return runtime.SortDirectionAsc
}

func getFunctionName(ctx fql.IFunctionCallContext, aliases map[string]string) runtime.String {
	var name string
	funcNS := ctx.Namespace()
	nsText := ""

	if funcNS != nil {
		nsText = funcNS.GetText()
	}

	if nsText != "" {
		ns := nsText

		if len(aliases) > 0 {
			ns = applyNamespaceAlias(ns, aliases)
		}

		name += ns
		name += ctx.FunctionName().GetText()

		return runtime.NewString(name)
	}

	fn := ctx.FunctionName().GetText()

	if len(aliases) > 0 {
		if target, ok := aliases[fn]; ok && target != "" {
			// Bare calls should only use function aliases (e.g. USE NS::FN AS ALIAS).
			// Namespace aliases (e.g. USE NS AS ALIAS) are intended for qualified
			// calls such as ALIAS::FN and must not rewrite ALIAS().
			if strings.Contains(target, runtime.NamespaceSeparator) {
				return runtime.NewString(target)
			}
		}
	}

	name += fn

	return runtime.NewString(name)
}

func getUDFName(ctx fql.IFunctionCallContext, aliases map[string]string) (string, bool) {
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
			// Bare function aliases targeting a qualified host function must bypass UDF lookup.
			if strings.Contains(target, runtime.NamespaceSeparator) {
				return "", false
			}
		}
	}

	return name, true
}

func applyNamespaceAlias(ns string, aliases map[string]string) string {
	if ns == "" || len(aliases) == 0 {
		return ns
	}

	trimmed := strings.TrimSuffix(ns, runtime.NamespaceSeparator)
	if trimmed == "" {
		return ns
	}

	parts := strings.Split(trimmed, runtime.NamespaceSeparator)
	if len(parts) == 0 {
		return ns
	}

	target, ok := aliases[parts[0]]
	if !ok {
		return ns
	}

	target = strings.TrimSuffix(target, runtime.NamespaceSeparator)
	if target == "" {
		return ns
	}

	parts = parts[1:]
	if len(parts) == 0 {
		return target + runtime.NamespaceSeparator
	}

	return target + runtime.NamespaceSeparator + strings.Join(parts, runtime.NamespaceSeparator) + runtime.NamespaceSeparator
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

func collectAndCaptureVars(
	ctx *CompilerContext,
	node antlr.Tree,
	env *varEnv,
	captureSet map[string]core.UDFCapture,
	captureOrder *[]string,
) {
	if node == nil || env == nil || captureSet == nil || captureOrder == nil {
		return
	}

	var vars []*fql.VariableContext
	findVariableRefs(node, &vars)

	for _, v := range vars {
		name, _ := variableName(v)
		if name == "" {
			continue
		}

		if env.currentHas(name) {
			continue
		}

		if _, ok := env.resolveBinding(name); ok {
			addCapture(captureSet, captureOrder, name, core.BindingStorageValue)
		}
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

func collectAndCaptureAssignments(
	ctx *CompilerContext,
	node antlr.Tree,
	env *varEnv,
	captureSet map[string]core.UDFCapture,
	captureOrder *[]string,
) {
	if node == nil || env == nil || captureSet == nil || captureOrder == nil {
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
			if ctx != nil && ctx.BindingCompiler != nil && binding.Decl != nil {
				ctx.BindingCompiler.PromoteDeclaration(binding.Decl)
			}
		}

		addCapture(captureSet, captureOrder, name, storage)
	}
}

func collectScopeFunctionsFromBody(
	ctx *CompilerContext,
	table *core.UDFTable,
	scope *core.UDFScope,
	body *fql.BodyContext,
) []*core.UDFInfo {
	if body == nil {
		return nil
	}

	out := make([]*core.UDFInfo, 0)

	for _, stmt := range body.AllBodyStatement() {
		if stmt == nil {
			continue
		}

		if decl := stmt.FunctionDeclaration(); decl != nil {
			if fn := registerFunction(ctx, table, scope, decl.(*fql.FunctionDeclarationContext)); fn != nil {
				out = append(out, fn)
			}
		}
	}

	return out
}

func collectNestedFunctions(ctx *CompilerContext, table *core.UDFTable, fn *core.UDFInfo) {
	if fn == nil || fn.Decl == nil {
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
				if nested := registerFunction(ctx, table, fn.BodyScope, decl.(*fql.FunctionDeclarationContext)); nested != nil {
					out = append(out, nested)
				}
			}
		}
	}

	for _, nested := range out {
		collectNestedFunctions(ctx, table, nested)
	}
}

func registerFunction(
	ctx *CompilerContext,
	table *core.UDFTable,
	scope *core.UDFScope,
	decl *fql.FunctionDeclarationContext,
) *core.UDFInfo {
	if decl == nil {
		return nil
	}

	displayName := decl.FunctionName().GetText()
	name := displayName

	if _, exists := scope.Functions[name]; exists {
		ctx.Errors.Add(ctx.Errors.Create(parserd.NameError, decl, fmt.Sprintf("Function '%s' is already defined", displayName)))
		return nil
	}

	params := collectFunctionParams(ctx, decl)

	fn := &core.UDFInfo{
		ID:          len(table.Functions),
		Name:        name,
		DisplayName: displayName,
		Params:      params,
		Decl:        decl,
		Scope:       scope,
		BodyScope:   core.NewUDFScope(scope),
	}

	scope.Functions[name] = fn
	table.Functions = append(table.Functions, fn)

	return fn
}

func collectFunctionParams(ctx *CompilerContext, decl *fql.FunctionDeclarationContext) []string {
	if decl == nil {
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
			ctx.Errors.Add(ctx.Errors.Create(parserd.NameError, param.(antlr.ParserRuleContext), fmt.Sprintf("Parameter '%s' is already defined", name)))
			continue
		}

		seen[name] = struct{}{}
		params = append(params, name)
	}

	return params
}

type varEnv struct {
	scopes []map[string]captureBinding
}

type captureBinding struct {
	Decl    antlr.ParserRuleContext
	Name    string
	Mutable bool
}

func (e *varEnv) push() {
	e.scopes = append(e.scopes, make(map[string]captureBinding))
}

func (e *varEnv) pop() {
	if len(e.scopes) > 0 {
		e.scopes = e.scopes[:len(e.scopes)-1]
	}
}

func (e *varEnv) add(name string) {
	e.addBinding(captureBinding{Name: name})
}

func (e *varEnv) addBinding(binding captureBinding) {
	if len(e.scopes) == 0 {
		return
	}

	e.scopes[len(e.scopes)-1][binding.Name] = binding
}

func (e *varEnv) currentHas(name string) bool {
	if len(e.scopes) == 0 {
		return false
	}

	_, ok := e.scopes[len(e.scopes)-1][name]
	return ok
}

func (e *varEnv) resolve(name string) bool {
	for i := len(e.scopes) - 1; i >= 0; i-- {
		if _, ok := e.scopes[i][name]; ok {
			return true
		}
	}

	return false
}

func (e *varEnv) resolveBinding(name string) (captureBinding, bool) {
	for i := len(e.scopes) - 1; i >= 0; i-- {
		if binding, ok := e.scopes[i][name]; ok {
			return binding, true
		}
	}

	return captureBinding{}, false
}

func addCapture(captures map[string]core.UDFCapture, order *[]string, name string, storage core.BindingStorage) {
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

func orderedCaptures(captures map[string]core.UDFCapture, order []string) []core.UDFCapture {
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

func analyzeCaptures(ctx *CompilerContext, table *core.UDFTable, body *fql.BodyContext) {
	if ctx == nil || table == nil || body == nil {
		return
	}

	env := &varEnv{}
	env.push() // global scope

	for _, stmt := range body.AllBodyStatement() {
		if stmt == nil {
			continue
		}

		switch {
		case stmt.VariableDeclaration() != nil:
			binding := ctx.BindingCompiler.captureBindingForDeclaration(stmt.VariableDeclaration())
			if binding.Name != "" {
				env.addBinding(binding)
			}
		case stmt.FunctionDeclaration() != nil:
			decl := stmt.FunctionDeclaration().(*fql.FunctionDeclarationContext)
			name := decl.FunctionName().GetText()
			if fn, ok := table.Resolve(name, table.GlobalScope); ok {
				analyzeFunctionCaptures(ctx, fn, env)
			}
		}
	}
}

func pruneUnusedUDFs(ctx *CompilerContext, table *core.UDFTable, body *fql.BodyContext) {
	if ctx == nil || table == nil || body == nil {
		return
	}

	reachable := computeReachableUDFs(ctx, table, body)
	if len(reachable) == 0 {
		for _, fn := range table.Functions {
			if fn != nil && fn.Scope != nil {
				delete(fn.Scope.Functions, fn.Name)
			}
		}
		table.Functions = nil

		return
	}

	filtered := make([]*core.UDFInfo, 0, len(reachable))

	for _, fn := range table.Functions {
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

	table.Functions = filtered
}

func computeReachableUDFs(ctx *CompilerContext, table *core.UDFTable, body *fql.BodyContext) map[*core.UDFInfo]struct{} {
	reachable := make(map[*core.UDFInfo]struct{})
	if ctx == nil || table == nil || body == nil {
		return reachable
	}

	roots := collectCallsInTopLevel(ctx, table, body, table.GlobalScope)
	if len(roots) == 0 {
		return reachable
	}

	callCache := make(map[*core.UDFInfo][]*core.UDFInfo)
	stack := make([]*core.UDFInfo, 0, len(roots))
	stack = append(stack, roots...)

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
			callees = collectCallsInFunction(ctx, table, fn)
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

func collectCallsInTopLevel(ctx *CompilerContext, table *core.UDFTable, body *fql.BodyContext, scope *core.UDFScope) []*core.UDFInfo {
	if ctx == nil || table == nil || body == nil {
		return nil
	}

	out := make(map[*core.UDFInfo]struct{})

	for _, stmt := range body.AllBodyStatement() {
		if stmt == nil {
			continue
		}

		st, ok := stmt.(*fql.BodyStatementContext)
		if !ok {
			continue
		}

		if st.FunctionDeclaration() != nil {
			continue
		}

		collectCallsInExpression(ctx, table, st, scope, out)
	}

	if expr := body.BodyExpression(); expr != nil {
		collectCallsInExpression(ctx, table, expr, scope, out)
	}

	return udfSetToSlice(out)
}

func collectCallsInFunction(ctx *CompilerContext, table *core.UDFTable, fn *core.UDFInfo) []*core.UDFInfo {
	if ctx == nil || table == nil || fn == nil || fn.Decl == nil {
		return nil
	}

	out := make(map[*core.UDFInfo]struct{})
	scope := fn.BodyScope

	body := fn.Decl.FunctionBody()
	if body == nil {
		return nil
	}

	if arrow := body.FunctionArrow(); arrow != nil {
		collectCallsInExpression(ctx, table, arrow.Expression(), scope, out)
		return udfSetToSlice(out)
	}

	if block := body.FunctionBlock(); block != nil {
		for _, stmt := range block.AllFunctionStatement() {
			if stmt == nil {
				continue
			}

			st, ok := stmt.(*fql.FunctionStatementContext)
			if !ok {
				continue
			}

			if st.FunctionDeclaration() != nil {
				continue
			}

			collectCallsInExpression(ctx, table, st, scope, out)
		}

		if ret := block.FunctionReturn(); ret != nil {
			collectCallsInExpression(ctx, table, ret, scope, out)
		}
	}

	return udfSetToSlice(out)
}

func collectCallsInExpression(ctx *CompilerContext, table *core.UDFTable, node antlr.Tree, scope *core.UDFScope, out map[*core.UDFInfo]struct{}) {
	if ctx == nil || table == nil || node == nil || out == nil {
		return
	}

	var calls []*fql.FunctionCallContext
	findFunctionCallRefs(node, &calls)

	for _, call := range calls {
		if call == nil {
			continue
		}

		name, ok := getUDFName(call, ctx.UseAliases)
		if !ok {
			continue
		}

		if fn, ok := table.Resolve(name, scope); ok {
			out[fn] = struct{}{}
		}
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

func analyzeFunctionCaptures(ctx *CompilerContext, fn *core.UDFInfo, env *varEnv) {
	if fn == nil || env == nil || fn.Decl == nil {
		return
	}

	env.push()
	for _, p := range fn.Params {
		env.addBinding(captureBinding{Name: p})
	}

	captureSet := make(map[string]core.UDFCapture)
	captureOrder := make([]string, 0)

	body := fn.Decl.FunctionBody()
	if body != nil {
		if arrow := body.FunctionArrow(); arrow != nil {
			if expr := arrow.Expression(); expr != nil {
				collectAndCaptureVars(ctx, expr, env, captureSet, &captureOrder)
				collectAndCaptureAssignments(ctx, expr, env, captureSet, &captureOrder)
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
						collectAndCaptureVars(ctx, decl.Expression(), env, captureSet, &captureOrder)
						collectAndCaptureAssignments(ctx, decl.Expression(), env, captureSet, &captureOrder)
					}
					binding := ctx.BindingCompiler.captureBindingForDeclaration(decl)
					if binding.Name != "" {
						env.addBinding(binding)
					}
				case stmt.AssignmentStatement() != nil:
					collectAndCaptureVars(ctx, stmt.AssignmentStatement(), env, captureSet, &captureOrder)
					collectAndCaptureAssignments(ctx, stmt.AssignmentStatement(), env, captureSet, &captureOrder)
				case stmt.FunctionDeclaration() != nil:
					decl := stmt.FunctionDeclaration().(*fql.FunctionDeclarationContext)
					name := decl.FunctionName().GetText()
					if nested, ok := fn.BodyScope.Functions[name]; ok {
						analyzeFunctionCaptures(ctx, nested, env)
						for _, cap := range nested.Captures {
							if env.currentHas(cap.Name) {
								continue
							}
							if _, ok := env.resolveBinding(cap.Name); !ok {
								continue
							}
							addCapture(captureSet, &captureOrder, cap.Name, cap.Storage)
						}
					}
				case stmt.FunctionCallExpression() != nil:
					collectAndCaptureVars(ctx, stmt.FunctionCallExpression(), env, captureSet, &captureOrder)
					collectAndCaptureAssignments(ctx, stmt.FunctionCallExpression(), env, captureSet, &captureOrder)
				case stmt.WaitForExpression() != nil:
					collectAndCaptureVars(ctx, stmt.WaitForExpression(), env, captureSet, &captureOrder)
					collectAndCaptureAssignments(ctx, stmt.WaitForExpression(), env, captureSet, &captureOrder)
				case stmt.DispatchExpression() != nil:
					collectAndCaptureVars(ctx, stmt.DispatchExpression(), env, captureSet, &captureOrder)
					collectAndCaptureAssignments(ctx, stmt.DispatchExpression(), env, captureSet, &captureOrder)
				case stmt.ExpressionStatement() != nil:
					collectAndCaptureVars(ctx, stmt.ExpressionStatement(), env, captureSet, &captureOrder)
					collectAndCaptureAssignments(ctx, stmt.ExpressionStatement(), env, captureSet, &captureOrder)
				}
			}

			if block.FunctionReturn() != nil {
				collectAndCaptureVars(ctx, block.FunctionReturn(), env, captureSet, &captureOrder)
				collectAndCaptureAssignments(ctx, block.FunctionReturn(), env, captureSet, &captureOrder)
			}
		}
	}

	fn.Captures = orderedCaptures(captureSet, captureOrder)

	env.pop()
}
