package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	exprMatchCallbacks struct {
		compileExpr                   func(fql.IExpressionContext) bytecode.Operand
		emitConditionJump             func(fql.IExpressionContext, core.Label, bool)
		compileFunctionCallByNameWith func(fql.IFunctionCallContext, runtime.String, bool, core.RegisterSequence) bytecode.Operand
	}

	exprMatchCompiler struct {
		ctx       *CompilationSession
		literals  *LiteralCompiler
		facts     *TypeFacts
		callbacks exprMatchCallbacks
	}
)

func newExprMatchCompiler(ctx *CompilationSession, callbacks exprMatchCallbacks) *exprMatchCompiler {
	return &exprMatchCompiler{
		ctx:       ctx,
		callbacks: callbacks,
	}
}

func (c *exprMatchCompiler) bind(literals *LiteralCompiler, facts *TypeFacts) {
	if c == nil {
		return
	}

	c.literals = literals
	c.facts = facts
}

func (c *exprMatchCompiler) compileMatchExpression(ctx fql.IMatchExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	dst := c.ctx.Function.Registers.Allocate()
	end := c.ctx.Program.Emitter.NewLabel("match.end")

	if arms := ctx.MatchPatternArms(); arms != nil {
		scrutinee := ctx.Expression()
		if scrutinee == nil {
			return bytecode.NoopOperand
		}

		if c.tryCompileMatchConstantFold(scrutinee, arms, dst) {
			c.ctx.Function.Types.Set(dst, core.TypeAny)

			return dst
		}

		scrReg := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(scrutinee))
		c.compileMatchPatternArms(scrReg, arms, dst, end)
	} else if guards := ctx.MatchGuardArms(); guards != nil {
		c.compileMatchGuardArms(guards, dst, end)
	}

	c.ctx.Program.Emitter.MarkLabel(end)

	if dst.IsRegister() {
		c.ctx.Function.Types.Set(dst, core.TypeAny)
	}

	return dst
}

func (c *exprMatchCompiler) compileMatchPatternArms(scrReg bytecode.Operand, ctx fql.IMatchPatternArmsContext, dst bytecode.Operand, end core.Label) {
	if ctx == nil {
		return
	}

	arms := collectMatchPatternArms(ctx)
	mergeLabels, mergeGroups := collectMatchResultMerges(c.ctx, arms)
	defaultLabel, hasDefaultLabel := c.matchMergeDefaultLabel(mergeGroups)

	for idx, arm := range arms {
		c.compileMatchPatternArm(scrReg, arm, idx, mergeLabels, dst, end)
	}

	if hasDefaultLabel {
		c.compileMatchMergedResults(mergeGroups, defaultLabel, dst, end)
	}

	c.compileMatchPatternDefaultArm(ctx.MatchDefaultArm(), dst)
}

func (c *exprMatchCompiler) matchMergeDefaultLabel(groups []matchResultGroup) (core.Label, bool) {
	if len(groups) == 0 {
		return core.Label{}, false
	}

	return c.ctx.Program.Emitter.NewLabel("match.default"), true
}

func (c *exprMatchCompiler) compileMatchPatternArm(scrReg bytecode.Operand, arm fql.IMatchPatternArmContext, idx int, mergeLabels map[int]core.Label, dst bytecode.Operand, end core.Label) {
	if arm == nil {
		return
	}

	next := c.ctx.Program.Emitter.NewLabel("match.next")
	c.ctx.Function.Symbols.EnterScope()
	c.compileMatchPatternArmConditions(scrReg, arm, next)
	c.compileMatchPatternArmResult(arm, idx, mergeLabels, dst, end)
	c.ctx.Function.Symbols.ExitScope()
	c.ctx.Program.Emitter.MarkLabel(next)
}

func (c *exprMatchCompiler) compileMatchPatternArmConditions(scrReg bytecode.Operand, arm fql.IMatchPatternArmContext, next core.Label) {
	if pattern := arm.MatchPattern(); pattern != nil {
		c.compileMatchPatternValue(scrReg, pattern, next)
	}

	guard := arm.MatchPatternGuard()
	if guard == nil {
		return
	}

	expr := guard.Expression()
	if expr == nil {
		return
	}

	c.callbacks.emitConditionJump(expr, next, false)
}

func (c *exprMatchCompiler) compileMatchPatternArmResult(arm fql.IMatchPatternArmContext, idx int, mergeLabels map[int]core.Label, dst bytecode.Operand, end core.Label) {
	result := arm.Expression()
	if result == nil {
		c.ctx.Program.Emitter.EmitJump(end)

		return
	}

	if label, ok := mergeLabels[idx]; ok {
		c.ctx.Program.Emitter.EmitJump(label)

		return
	}

	out := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(result))
	if out != bytecode.NoopOperand && out != dst {
		c.ctx.Program.Emitter.EmitMove(dst, out)
	}

	c.ctx.Program.Emitter.EmitJump(end)
}

func (c *exprMatchCompiler) compileMatchMergedResults(groups []matchResultGroup, defaultLabel core.Label, dst bytecode.Operand, end core.Label) {
	c.ctx.Program.Emitter.EmitJump(defaultLabel)

	for _, group := range groups {
		c.ctx.Program.Emitter.MarkLabel(group.label)
		c.ctx.Function.Symbols.EnterScope()
		out := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(group.result))
		if out != bytecode.NoopOperand && out != dst {
			c.ctx.Program.Emitter.EmitMove(dst, out)
		}
		c.ctx.Function.Symbols.ExitScope()
		c.ctx.Program.Emitter.EmitJump(end)
	}

	c.ctx.Program.Emitter.MarkLabel(defaultLabel)
}

func (c *exprMatchCompiler) compileMatchPatternDefaultArm(def fql.IMatchDefaultArmContext, dst bytecode.Operand) {
	if def == nil {
		return
	}

	c.ctx.Function.Symbols.EnterScope()
	result := def.Expression()
	if result != nil {
		out := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(result))
		if out != bytecode.NoopOperand && out != dst {
			c.ctx.Program.Emitter.EmitMove(dst, out)
		}
	}
	c.ctx.Function.Symbols.ExitScope()
}

func (c *exprMatchCompiler) compileMatchGuardArms(ctx fql.IMatchGuardArmsContext, dst bytecode.Operand, end core.Label) {
	if ctx == nil {
		return
	}

	var arms []fql.IMatchGuardArmContext
	if list := ctx.MatchGuardArmList(); list != nil {
		arms = list.AllMatchGuardArm()
	}

	for _, arm := range arms {
		if arm == nil {
			continue
		}

		next := c.ctx.Program.Emitter.NewLabel("match.next")
		c.ctx.Function.Symbols.EnterScope()

		exprs := arm.AllExpression()
		if len(exprs) > 0 {
			c.callbacks.emitConditionJump(exprs[0], next, false)
		}

		if len(exprs) > 1 {
			out := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(exprs[1]))
			if out != bytecode.NoopOperand && out != dst {
				c.ctx.Program.Emitter.EmitMove(dst, out)
			}
		}

		c.ctx.Program.Emitter.EmitJump(end)
		c.ctx.Function.Symbols.ExitScope()
		c.ctx.Program.Emitter.MarkLabel(next)
	}

	if def := ctx.MatchDefaultArm(); def != nil {
		c.ctx.Function.Symbols.EnterScope()
		if result := def.Expression(); result != nil {
			out := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(result))
			if out != bytecode.NoopOperand && out != dst {
				c.ctx.Program.Emitter.EmitMove(dst, out)
			}
		}
		c.ctx.Function.Symbols.ExitScope()
	}
}

func (c *exprMatchCompiler) compileMatchPatternValue(valueReg bytecode.Operand, ctx fql.IMatchPatternContext, onFail core.Label) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.MatchLiteralPattern() != nil:
		litOp := c.compileMatchLiteralOperand(ctx.MatchLiteralPattern())
		if litOp == bytecode.NoopOperand {
			return
		}

		if litOp.IsConstant() {
			c.ctx.Program.Emitter.EmitJumpCompare(bytecode.OpJumpIfNeConst, valueReg, litOp, onFail)
		} else {
			c.ctx.Program.Emitter.EmitJumpCompare(bytecode.OpJumpIfNe, valueReg, litOp, onFail)
		}
	case ctx.MatchBindingPattern() != nil:
		binding := ctx.MatchBindingPattern()
		if binding == nil {
			return
		}

		var name string
		if id := binding.Identifier(); id != nil {
			name = id.GetText()
		} else if srw := binding.SafeReservedWord(); srw != nil {
			name = srw.GetText()
		}

		if name != "" {
			c.declareMatchBinding(binding, name, valueReg)
		}
	case ctx.MatchObjectPattern() != nil:
		c.compileMatchObjectPattern(valueReg, ctx.MatchObjectPattern(), onFail)
	}
}

func (c *exprMatchCompiler) compileMatchLiteralOperand(ctx fql.IMatchLiteralPatternContext) bytecode.Operand {
	return compileScalarLiteralOperand(c.ctx, c.literals, ctx)
}

func (c *exprMatchCompiler) tryCompileMatchConstantFold(scrutinee fql.IExpressionContext, armsCtx fql.IMatchPatternArmsContext, dst bytecode.Operand) bool {
	scrutineeVal, arms, ok := matchConstantFoldPreconditions(scrutinee, armsCtx)
	if !ok {
		return false
	}

	selected, ok := selectMatchConstantFoldExpression(scrutineeVal, arms, armsCtx.MatchDefaultArm())
	if !ok {
		return false
	}

	return c.emitMatchConstantFoldExpression(selected, dst)
}

func (c *exprMatchCompiler) emitMatchConstantFoldExpression(expr fql.IExpressionContext, dst bytecode.Operand) bool {
	if expr == nil {
		return false
	}

	c.ctx.Function.Symbols.EnterScope()
	out := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(expr))
	if out != bytecode.NoopOperand && out != dst {
		c.ctx.Program.Emitter.EmitMove(dst, out)
	}
	c.ctx.Function.Symbols.ExitScope()

	return true
}

func (c *exprMatchCompiler) compileMatchObjectPattern(valueReg bytecode.Operand, ctx fql.IMatchObjectPatternContext, onFail core.Label) {
	if ctx == nil {
		return
	}

	props := ctx.AllMatchObjectPatternProperty()
	if len(props) == 0 {
		keys := c.emitObjectsKeys(valueReg)
		c.ctx.Program.Emitter.EmitJumpIfNone(keys, onFail)

		return
	}

	for _, prop := range props {
		if prop == nil {
			continue
		}

		keyOp := c.compileMatchObjectPatternKey(prop.MatchObjectPatternKey())
		if keyOp == bytecode.NoopOperand {
			continue
		}

		val := c.ctx.Function.Registers.Allocate()
		if keyOp.IsConstant() {
			c.ctx.Program.Emitter.EmitMatchLoadPropertyConst(val, valueReg, keyOp, onFail)
		} else {
			c.ctx.Program.Emitter.EmitJumpCompare(bytecode.OpJumpIfMissingProperty, valueReg, keyOp, onFail)
			c.ctx.Program.Emitter.EmitABC(bytecode.OpLoadProperty, val, valueReg, keyOp)
		}

		c.compileMatchPatternValue(val, prop.MatchPattern(), onFail)
	}
}

func (c *exprMatchCompiler) compileMatchObjectPatternKey(ctx fql.IMatchObjectPatternKeyContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	if sl := ctx.StringLiteral(); sl != nil {
		if val, ok := parseStringLiteralConst(sl); ok {
			return c.ctx.Function.Symbols.AddConstant(val)
		}

		return c.literals.CompileStringLiteral(sl)
	}

	var name string

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if srw := ctx.SafeReservedWord(); srw != nil {
		name = srw.GetText()
	} else if urw := ctx.UnsafeReservedWord(); urw != nil {
		name = urw.GetText()
	}

	if name == "" {
		return bytecode.NoopOperand
	}

	return c.ctx.Function.Symbols.AddConstant(runtime.NewString(name))
}

func (c *exprMatchCompiler) emitObjectsKeys(scrReg bytecode.Operand) bytecode.Operand {
	scrReg = ensureOperandRegister(c.ctx, c.facts, scrReg)
	seq := c.ctx.Function.Registers.AllocateSequence(1)
	c.ctx.Program.Emitter.EmitMove(seq[0], scrReg)

	return c.callbacks.compileFunctionCallByNameWith(nil, runtime.NewString("KEYS"), true, seq)
}

func (c *exprMatchCompiler) declareMatchBinding(ctx antlr.ParserRuleContext, name string, valueReg bytecode.Operand) bytecode.Operand {
	valueReg = ensureOperandRegister(c.ctx, c.facts, valueReg)
	reg, ok := c.ctx.Function.Symbols.DeclareLocal(name, core.TypeAny)
	if ok {
		c.ctx.Program.Emitter.EmitMove(reg, valueReg)
		c.ctx.Function.Types.Set(reg, c.facts.OperandType(valueReg))

		return reg
	}

	if ctx != nil {
		c.ctx.Program.Errors.DuplicateMatchBinding(ctx, name)
	}

	if existing, _, found := c.ctx.Function.Symbols.Resolve(name); found {
		return existing
	}

	return valueReg
}
