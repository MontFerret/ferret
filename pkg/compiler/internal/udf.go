package internal

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	UDFTable struct {
		Functions   []*UDFInfo
		GlobalScope *UDFScope
	}

	UDFScope struct {
		Parent    *UDFScope
		Functions map[string]*UDFInfo
	}

	UDFInfo struct {
		ID        int
		Name      string
		Params    []string
		Captures  []string
		Decl      fql.IFunctionDeclarationContext
		Scope     *UDFScope
		BodyScope *UDFScope
		Entry     int
		Registers int
	}
)

func NewUDFScope(parent *UDFScope) *UDFScope {
	return &UDFScope{
		Parent:    parent,
		Functions: make(map[string]*UDFInfo),
	}
}

func NewUDFTable() *UDFTable {
	return &UDFTable{
		Functions: make([]*UDFInfo, 0),
	}
}

func (t *UDFTable) Metadata() []bytecode.UDF {
	if t == nil || len(t.Functions) == 0 {
		return nil
	}

	out := make([]bytecode.UDF, 0, len(t.Functions))
	for _, fn := range t.Functions {
		if fn == nil {
			continue
		}

		out = append(out, bytecode.UDF{
			Name:      fn.Name,
			Entry:     fn.Entry,
			Registers: fn.Registers,
			Params:    len(fn.Params) + len(fn.Captures),
		})
	}

	return out
}

func (t *UDFTable) Resolve(name string, scope *UDFScope) (*UDFInfo, bool) {
	if scope == nil {
		return nil, false
	}

	for s := scope; s != nil; s = s.Parent {
		if fn, ok := s.Functions[name]; ok {
			return fn, true
		}
	}

	return nil, false
}

func CollectUDFs(ctx *CompilerContext, program *fql.ProgramContext) *UDFTable {
	table := NewUDFTable()
	table.GlobalScope = NewUDFScope(nil)

	if program == nil || program.Body() == nil {
		return table
	}

	body, ok := program.Body().(*fql.BodyContext)
	if !ok {
		return table
	}

	top := collectScopeFunctionsFromBody(ctx, table, table.GlobalScope, body)

	for _, fn := range top {
		collectNestedFunctions(ctx, table, fn)
	}

	if ctx != nil && ctx.OptimizationLevel > optimization.LevelNone {
		pruneUnusedUDFs(ctx, table, body)
	}

	analyzeCaptures(ctx, table, body)

	return table
}

func collectScopeFunctionsFromBody(
	ctx *CompilerContext,
	table *UDFTable,
	scope *UDFScope,
	body *fql.BodyContext,
) []*UDFInfo {
	if body == nil {
		return nil
	}

	out := make([]*UDFInfo, 0)

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

func collectNestedFunctions(ctx *CompilerContext, table *UDFTable, fn *UDFInfo) {
	if fn == nil || fn.Decl == nil {
		return
	}

	body := fn.Decl.FunctionBody()
	if body == nil {
		return
	}

	out := make([]*UDFInfo, 0)

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
	table *UDFTable,
	scope *UDFScope,
	decl *fql.FunctionDeclarationContext,
) *UDFInfo {
	if decl == nil {
		return nil
	}

	name := strings.ToUpper(decl.FunctionName().GetText())

	if _, exists := scope.Functions[name]; exists {
		ctx.Errors.Add(ctx.Errors.Create(parserd.NameError, decl, fmt.Sprintf("Function '%s' is already defined", name)))
		return nil
	}

	params := collectFunctionParams(ctx, decl)

	fn := &UDFInfo{
		ID:        len(table.Functions),
		Name:      name,
		Params:    params,
		Decl:      decl,
		Scope:     scope,
		BodyScope: NewUDFScope(scope),
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
	scopes []map[string]struct{}
}

func (e *varEnv) push() {
	e.scopes = append(e.scopes, make(map[string]struct{}))
}

func (e *varEnv) pop() {
	if len(e.scopes) > 0 {
		e.scopes = e.scopes[:len(e.scopes)-1]
	}
}

func (e *varEnv) add(name string) {
	if len(e.scopes) == 0 {
		return
	}

	e.scopes[len(e.scopes)-1][name] = struct{}{}
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

func analyzeCaptures(ctx *CompilerContext, table *UDFTable, body *fql.BodyContext) {
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
			name := variableDeclarationName(stmt.VariableDeclaration())
			if name != "" {
				env.add(name)
			}
		case stmt.FunctionDeclaration() != nil:
			decl := stmt.FunctionDeclaration().(*fql.FunctionDeclarationContext)
			name := strings.ToUpper(decl.FunctionName().GetText())
			if fn, ok := table.Resolve(name, table.GlobalScope); ok {
				analyzeFunctionCaptures(ctx, fn, env)
			}
		}
	}
}

func pruneUnusedUDFs(ctx *CompilerContext, table *UDFTable, body *fql.BodyContext) {
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

	filtered := make([]*UDFInfo, 0, len(reachable))

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

func computeReachableUDFs(ctx *CompilerContext, table *UDFTable, body *fql.BodyContext) map[*UDFInfo]struct{} {
	reachable := make(map[*UDFInfo]struct{})
	if ctx == nil || table == nil || body == nil {
		return reachable
	}

	roots := collectCallsInTopLevel(ctx, table, body, table.GlobalScope)
	if len(roots) == 0 {
		return reachable
	}

	callCache := make(map[*UDFInfo][]*UDFInfo)
	stack := make([]*UDFInfo, 0, len(roots))
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

func collectCallsInTopLevel(ctx *CompilerContext, table *UDFTable, body *fql.BodyContext, scope *UDFScope) []*UDFInfo {
	if ctx == nil || table == nil || body == nil {
		return nil
	}

	out := make(map[*UDFInfo]struct{})

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

func collectCallsInFunction(ctx *CompilerContext, table *UDFTable, fn *UDFInfo) []*UDFInfo {
	if ctx == nil || table == nil || fn == nil || fn.Decl == nil {
		return nil
	}

	out := make(map[*UDFInfo]struct{})
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

func collectCallsInExpression(ctx *CompilerContext, table *UDFTable, node antlr.Tree, scope *UDFScope, out map[*UDFInfo]struct{}) {
	if ctx == nil || table == nil || node == nil || out == nil {
		return
	}

	var calls []*fql.FunctionCallContext
	findFunctionCallRefs(node, &calls)

	for _, call := range calls {
		if call == nil {
			continue
		}

		name := getFunctionName(call, ctx.UseAliases)
		nameStr := name.String()
		if strings.Contains(nameStr, runtime.NamespaceSeparator) {
			continue
		}

		if fn, ok := table.Resolve(nameStr, scope); ok {
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

func udfSetToSlice(set map[*UDFInfo]struct{}) []*UDFInfo {
	if len(set) == 0 {
		return nil
	}

	out := make([]*UDFInfo, 0, len(set))
	for fn := range set {
		out = append(out, fn)
	}

	return out
}

func analyzeFunctionCaptures(ctx *CompilerContext, fn *UDFInfo, env *varEnv) {
	if fn == nil || env == nil || fn.Decl == nil {
		return
	}

	env.push()
	for _, p := range fn.Params {
		env.add(p)
	}

	captureSet := make(map[string]struct{})
	captureOrder := make([]string, 0)

	body := fn.Decl.FunctionBody()
	if body != nil {
		if arrow := body.FunctionArrow(); arrow != nil {
			if expr := arrow.Expression(); expr != nil {
				collectAndCaptureVars(ctx, expr, env, captureSet, &captureOrder)
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
					}
					name := variableDeclarationName(decl)
					if name != "" {
						env.add(name)
					}
				case stmt.FunctionDeclaration() != nil:
					decl := stmt.FunctionDeclaration().(*fql.FunctionDeclarationContext)
					name := strings.ToUpper(decl.FunctionName().GetText())
					if nested, ok := fn.BodyScope.Functions[name]; ok {
						analyzeFunctionCaptures(ctx, nested, env)
						for _, cap := range nested.Captures {
							if env.currentHas(cap) {
								continue
							}
							if _, exists := captureSet[cap]; !exists {
								captureSet[cap] = struct{}{}
								captureOrder = append(captureOrder, cap)
							}
							env.add(cap)
						}
					}
				case stmt.FunctionCallExpression() != nil:
					collectAndCaptureVars(ctx, stmt.FunctionCallExpression(), env, captureSet, &captureOrder)
				case stmt.WaitForExpression() != nil:
					collectAndCaptureVars(ctx, stmt.WaitForExpression(), env, captureSet, &captureOrder)
				case stmt.DispatchExpression() != nil:
					collectAndCaptureVars(ctx, stmt.DispatchExpression(), env, captureSet, &captureOrder)
				case stmt.ExpressionStatement() != nil:
					collectAndCaptureVars(ctx, stmt.ExpressionStatement(), env, captureSet, &captureOrder)
				}
			}

			if block.FunctionReturn() != nil {
				collectAndCaptureVars(ctx, block.FunctionReturn(), env, captureSet, &captureOrder)
			}
		}
	}

	fn.Captures = captureOrder

	env.pop()
}

func variableDeclarationName(ctx fql.IVariableDeclarationContext) string {
	if ctx == nil {
		return ""
	}

	if id := ctx.Identifier(); id != nil {
		return id.GetText()
	}

	if id := ctx.SafeReservedWord(); id != nil {
		return id.GetText()
	}

	return ""
}

func collectAndCaptureVars(
	ctx *CompilerContext,
	node antlr.Tree,
	env *varEnv,
	captureSet map[string]struct{},
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

		if env.resolve(name) {
			if _, exists := captureSet[name]; !exists {
				captureSet[name] = struct{}{}
				*captureOrder = append(*captureOrder, name)
			}
			env.add(name)
			continue
		}
	}
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
