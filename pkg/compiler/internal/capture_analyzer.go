package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type CaptureAnalyzer struct {
	ctx      *CompilationSession
	bindings *BindingCompiler
	calls    *CallResolver
	states   map[*core.UDFInfo]*udfCaptureState
}

func NewCaptureAnalyzer(ctx *CompilationSession) *CaptureAnalyzer {
	return &CaptureAnalyzer{ctx: ctx}
}

func (c *CaptureAnalyzer) bind(bindings *BindingCompiler, calls *CallResolver) {
	if c == nil {
		return
	}

	c.bindings = bindings
	c.calls = calls
}

func (c *CaptureAnalyzer) AnalyzeProgram(body *fql.BodyContext) {
	if c == nil || c.ctx == nil || c.ctx.Program.UDFs == nil || body == nil {
		return
	}

	env := &udfCaptureEnv{}
	env.push()
	c.states = make(map[*core.UDFInfo]*udfCaptureState, len(c.ctx.Program.UDFs.Functions))

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
			if fn, ok := c.ctx.Program.UDFs.Resolve(name, c.ctx.Program.UDFs.GlobalScope); ok {
				c.analyzeFunction(fn, env)
			}
		}
	}

	c.collectCallGraph()
	c.propagateCallCaptures()
	c.publishCaptures()
}

func (c *CaptureAnalyzer) analyzeFunction(fn *core.UDFInfo, env *udfCaptureEnv) {
	if c == nil || c.ctx == nil || fn == nil || env == nil || fn.Decl == nil {
		return
	}

	state := &udfCaptureState{
		captures: make(map[core.BindingID]core.UDFCapture),
		order:    make([]core.BindingID, 0),
		outer:    env.bindingsByID(),
		owned:    make(map[core.BindingID]captureBindingInfo),
	}
	c.states[fn] = state

	env.push()
	for _, param := range fn.Params {
		binding := captureBindingInfo{ID: param.ID, Name: param.Name}
		env.addBinding(binding)
		state.owned[binding.ID] = binding
	}

	body := fn.Decl.FunctionBody()
	if body != nil {
		if arrow := body.FunctionArrow(); arrow != nil {
			if expr := arrow.Expression(); expr != nil {
				c.collectVars(expr, env, state)
				c.collectAssignments(expr, env, state)
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
						c.collectVars(decl.Expression(), env, state)
						c.collectAssignments(decl.Expression(), env, state)
					}

					binding := c.bindings.captureBindingForDeclaration(decl)
					if binding.Name != "" {
						env.addBinding(binding)
						state.owned[binding.ID] = binding
					}
				case stmt.AssignmentStatement() != nil:
					c.collectVars(stmt.AssignmentStatement(), env, state)
					c.collectAssignments(stmt.AssignmentStatement(), env, state)
				case stmt.FunctionDeclaration() != nil:
					decl := stmt.FunctionDeclaration().(*fql.FunctionDeclarationContext)
					name := decl.FunctionName().GetText()

					if nested, ok := fn.BodyScope.Functions[name]; ok {
						c.analyzeFunction(nested, env)
						for _, capture := range nested.Captures {
							if _, ok := state.owned[capture.ID]; ok {
								continue
							}

							if _, ok := state.outer[capture.ID]; !ok {
								continue
							}

							capture.Visible = false
							addUDFCapture(state.captures, &state.order, capture)
						}
					}
				case stmt.FunctionCallExpression() != nil:
					c.collectVars(stmt.FunctionCallExpression(), env, state)
					c.collectAssignments(stmt.FunctionCallExpression(), env, state)
				case stmt.WaitForExpression() != nil:
					c.collectVars(stmt.WaitForExpression(), env, state)
					c.collectAssignments(stmt.WaitForExpression(), env, state)
				case stmt.DispatchExpression() != nil:
					c.collectVars(stmt.DispatchExpression(), env, state)
					c.collectAssignments(stmt.DispatchExpression(), env, state)
				case stmt.ExpressionStatement() != nil:
					c.collectVars(stmt.ExpressionStatement(), env, state)
					c.collectAssignments(stmt.ExpressionStatement(), env, state)
				}
			}

			if block.FunctionReturn() != nil {
				c.collectVars(block.FunctionReturn(), env, state)
				c.collectAssignments(block.FunctionReturn(), env, state)
			}
		}
	}

	fn.Captures = orderedUDFCaptures(state.captures, state.order)

	env.pop()
}

func (c *CaptureAnalyzer) collectVars(
	node antlr.Tree,
	env *udfCaptureEnv,
	state *udfCaptureState,
) {
	if c == nil || node == nil || env == nil || state == nil {
		return
	}

	var vars []*fql.VariableContext
	findVariableRefs(node, &vars)

	for _, v := range vars {
		name, _ := variableName(v)
		if name == "" || env.currentHas(name) {
			continue
		}

		if binding, ok := env.resolveBinding(name); ok {
			addUDFCapture(state.captures, &state.order, core.UDFCapture{
				ID:      binding.ID,
				Name:    binding.Name,
				Storage: core.BindingStorageValue,
				Visible: true,
			})
		}
	}
}

func (c *CaptureAnalyzer) collectAssignments(
	node antlr.Tree,
	env *udfCaptureEnv,
	state *udfCaptureState,
) {
	if c == nil || c.ctx == nil || node == nil || env == nil || state == nil {
		return
	}

	var assignments []*fql.AssignmentStatementContext
	findAssignmentRefs(node, &assignments)

	for _, stmt := range assignments {
		if stmt == nil || stmt.AssignmentTarget() == nil {
			continue
		}

		target, ok := newAssignmentTarget(stmt.AssignmentTarget())
		if !ok {
			continue
		}

		name := target.Root
		if name == "" || env.currentHas(name) {
			continue
		}

		binding, ok := env.resolveBinding(name)
		if !ok {
			continue
		}

		if len(target.Segments) > 0 {
			addUDFCapture(state.captures, &state.order, core.UDFCapture{
				ID:      binding.ID,
				Name:    binding.Name,
				Storage: core.BindingStorageValue,
				Visible: true,
			})
			continue
		}

		storage := core.BindingStorageValue

		if binding.Mutable {
			storage = core.BindingStorageCell

			if binding.Decl != nil {
				c.bindings.PromoteDeclaration(binding.Decl)
			}
		}

		addUDFCapture(state.captures, &state.order, core.UDFCapture{
			ID:      binding.ID,
			Name:    binding.Name,
			Storage: storage,
			Visible: true,
		})
	}
}

func (c *CaptureAnalyzer) collectCallGraph() {
	if c == nil || c.ctx == nil || c.ctx.Program.UDFs == nil || c.calls == nil {
		return
	}

	for _, fn := range c.ctx.Program.UDFs.Functions {
		state := c.states[fn]
		if fn == nil || state == nil || fn.Decl == nil || fn.Decl.FunctionBody() == nil {
			continue
		}

		var calls []*fql.FunctionCallContext
		findFunctionCallRefs(fn.Decl.FunctionBody(), &calls)
		seen := make(map[*core.UDFInfo]struct{})

		for _, call := range calls {
			callee, ok := c.calls.ResolveUDFInScope(call, fn.BodyScope)
			if !ok || callee == nil {
				continue
			}

			if _, ok := seen[callee]; ok {
				continue
			}

			seen[callee] = struct{}{}
			state.callees = append(state.callees, callee)
		}
	}
}

// propagateCallCaptures computes the transitive closure of captures required at
// UDF call sites. Forwarded captures remain tied to their lexical binding ID.
func (c *CaptureAnalyzer) propagateCallCaptures() {
	if c == nil || c.ctx == nil || c.ctx.Program.UDFs == nil {
		return
	}

	for {
		changed := false

		for _, fn := range c.ctx.Program.UDFs.Functions {
			state := c.states[fn]
			if state == nil {
				continue
			}

			for _, callee := range state.callees {
				calleeState := c.states[callee]
				if calleeState == nil {
					continue
				}

				for _, id := range calleeState.order {
					capture, ok := calleeState.captures[id]
					if !ok {
						continue
					}

					if _, ok := state.owned[id]; ok {
						continue
					}

					if _, ok := state.outer[id]; !ok {
						continue
					}

					capture.Visible = false

					if addUDFCapture(state.captures, &state.order, capture) {
						changed = true
					}
				}
			}
		}

		if !changed {
			return
		}
	}
}

func (c *CaptureAnalyzer) publishCaptures() {
	if c == nil || c.ctx == nil || c.ctx.Program.UDFs == nil {
		return
	}

	for _, fn := range c.ctx.Program.UDFs.Functions {
		state := c.states[fn]
		if fn == nil || state == nil {
			continue
		}

		fn.Captures = orderedUDFCaptures(state.captures, state.order)
	}
}
