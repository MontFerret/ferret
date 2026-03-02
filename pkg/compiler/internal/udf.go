package internal

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
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

	for _, stmt := range body.AllFunctionStatement() {
		if stmt == nil {
			continue
		}

		if decl := stmt.FunctionDeclaration(); decl != nil {
			if nested := registerFunction(ctx, table, fn.BodyScope, decl.(*fql.FunctionDeclarationContext)); nested != nil {
				out = append(out, nested)
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
		for _, stmt := range body.AllFunctionStatement() {
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
	}

	if body != nil && body.FunctionReturn() != nil {
		collectAndCaptureVars(ctx, body.FunctionReturn(), env, captureSet, &captureOrder)
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
