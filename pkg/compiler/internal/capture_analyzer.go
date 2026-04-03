package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type CaptureAnalyzer struct {
	ctx      *CompilationSession
	bindings *BindingCompiler
}

func NewCaptureAnalyzer(ctx *CompilationSession) *CaptureAnalyzer {
	return &CaptureAnalyzer{ctx: ctx}
}

func (c *CaptureAnalyzer) bind(bindings *BindingCompiler) {
	if c == nil {
		return
	}

	c.bindings = bindings
}

func (c *CaptureAnalyzer) AnalyzeProgram(body *fql.BodyContext) {
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
			binding := c.bindings.captureBindingForDeclaration(stmt.VariableDeclaration())
			if binding.Name != "" {
				env.addBinding(binding)
			}
		case stmt.FunctionDeclaration() != nil:
			decl := stmt.FunctionDeclaration().(*fql.FunctionDeclarationContext)
			name := decl.FunctionName().GetText()
			if fn, ok := c.ctx.UDFs.Resolve(name, c.ctx.UDFs.GlobalScope); ok {
				c.analyzeFunction(fn, env)
			}
		}
	}
}

func (c *CaptureAnalyzer) analyzeFunction(fn *core.UDFInfo, env *udfCaptureEnv) {
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
				c.collectVars(expr, env, captureSet, &captureOrder)
				c.collectAssignments(expr, env, captureSet, &captureOrder)
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
						c.collectVars(decl.Expression(), env, captureSet, &captureOrder)
						c.collectAssignments(decl.Expression(), env, captureSet, &captureOrder)
					}

					binding := c.bindings.captureBindingForDeclaration(decl)
					if binding.Name != "" {
						env.addBinding(binding)
					}
				case stmt.AssignmentStatement() != nil:
					c.collectVars(stmt.AssignmentStatement(), env, captureSet, &captureOrder)
					c.collectAssignments(stmt.AssignmentStatement(), env, captureSet, &captureOrder)
				case stmt.FunctionDeclaration() != nil:
					decl := stmt.FunctionDeclaration().(*fql.FunctionDeclarationContext)
					name := decl.FunctionName().GetText()
					if nested, ok := fn.BodyScope.Functions[name]; ok {
						c.analyzeFunction(nested, env)
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
					c.collectVars(stmt.FunctionCallExpression(), env, captureSet, &captureOrder)
					c.collectAssignments(stmt.FunctionCallExpression(), env, captureSet, &captureOrder)
				case stmt.WaitForExpression() != nil:
					c.collectVars(stmt.WaitForExpression(), env, captureSet, &captureOrder)
					c.collectAssignments(stmt.WaitForExpression(), env, captureSet, &captureOrder)
				case stmt.DispatchExpression() != nil:
					c.collectVars(stmt.DispatchExpression(), env, captureSet, &captureOrder)
					c.collectAssignments(stmt.DispatchExpression(), env, captureSet, &captureOrder)
				case stmt.ExpressionStatement() != nil:
					c.collectVars(stmt.ExpressionStatement(), env, captureSet, &captureOrder)
					c.collectAssignments(stmt.ExpressionStatement(), env, captureSet, &captureOrder)
				}
			}

			if block.FunctionReturn() != nil {
				c.collectVars(block.FunctionReturn(), env, captureSet, &captureOrder)
				c.collectAssignments(block.FunctionReturn(), env, captureSet, &captureOrder)
			}
		}
	}

	fn.Captures = orderedUDFCaptures(captureSet, captureOrder)

	env.pop()
}

func (c *CaptureAnalyzer) collectVars(
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

func (c *CaptureAnalyzer) collectAssignments(
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
				c.bindings.PromoteDeclaration(binding.Decl)
			}
		}

		addUDFCapture(captureSet, captureOrder, name, storage)
	}
}
