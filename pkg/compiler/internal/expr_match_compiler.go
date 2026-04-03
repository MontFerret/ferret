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

	dst := c.ctx.Registers.Allocate()
	end := c.ctx.Emitter.NewLabel("match.end")

	if arms := ctx.MatchPatternArms(); arms != nil {
		scrutinee := ctx.Expression()
		if scrutinee == nil {
			return bytecode.NoopOperand
		}

		if c.tryCompileMatchConstantFold(scrutinee, arms, dst) {
			c.ctx.Types.Set(dst, core.TypeAny)

			return dst
		}

		scrReg := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(scrutinee))
		c.compileMatchPatternArms(scrReg, arms, dst, end)
	} else if guards := ctx.MatchGuardArms(); guards != nil {
		c.compileMatchGuardArms(guards, dst, end)
	}

	c.ctx.Emitter.MarkLabel(end)

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeAny)
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

	return c.ctx.Emitter.NewLabel("match.default"), true
}

func (c *exprMatchCompiler) compileMatchPatternArm(scrReg bytecode.Operand, arm fql.IMatchPatternArmContext, idx int, mergeLabels map[int]core.Label, dst bytecode.Operand, end core.Label) {
	if arm == nil {
		return
	}

	next := c.ctx.Emitter.NewLabel("match.next")
	c.ctx.Symbols.EnterScope()
	c.compileMatchPatternArmConditions(scrReg, arm, next)
	c.compileMatchPatternArmResult(arm, idx, mergeLabels, dst, end)
	c.ctx.Symbols.ExitScope()
	c.ctx.Emitter.MarkLabel(next)
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
		c.ctx.Emitter.EmitJump(end)

		return
	}

	if label, ok := mergeLabels[idx]; ok {
		c.ctx.Emitter.EmitJump(label)

		return
	}

	out := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(result))
	if out != bytecode.NoopOperand && out != dst {
		c.ctx.Emitter.EmitMove(dst, out)
	}

	c.ctx.Emitter.EmitJump(end)
}

func (c *exprMatchCompiler) compileMatchMergedResults(groups []matchResultGroup, defaultLabel core.Label, dst bytecode.Operand, end core.Label) {
	c.ctx.Emitter.EmitJump(defaultLabel)

	for _, group := range groups {
		c.ctx.Emitter.MarkLabel(group.label)
		c.ctx.Symbols.EnterScope()
		out := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(group.result))
		if out != bytecode.NoopOperand && out != dst {
			c.ctx.Emitter.EmitMove(dst, out)
		}
		c.ctx.Symbols.ExitScope()
		c.ctx.Emitter.EmitJump(end)
	}

	c.ctx.Emitter.MarkLabel(defaultLabel)
}

func (c *exprMatchCompiler) compileMatchPatternDefaultArm(def fql.IMatchDefaultArmContext, dst bytecode.Operand) {
	if def == nil {
		return
	}

	c.ctx.Symbols.EnterScope()
	result := def.Expression()
	if result != nil {
		out := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(result))
		if out != bytecode.NoopOperand && out != dst {
			c.ctx.Emitter.EmitMove(dst, out)
		}
	}
	c.ctx.Symbols.ExitScope()
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

		next := c.ctx.Emitter.NewLabel("match.next")
		c.ctx.Symbols.EnterScope()

		exprs := arm.AllExpression()
		if len(exprs) > 0 {
			c.callbacks.emitConditionJump(exprs[0], next, false)
		}

		if len(exprs) > 1 {
			out := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(exprs[1]))
			if out != bytecode.NoopOperand && out != dst {
				c.ctx.Emitter.EmitMove(dst, out)
			}
		}

		c.ctx.Emitter.EmitJump(end)
		c.ctx.Symbols.ExitScope()
		c.ctx.Emitter.MarkLabel(next)
	}

	if def := ctx.MatchDefaultArm(); def != nil {
		c.ctx.Symbols.EnterScope()
		if result := def.Expression(); result != nil {
			out := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(result))
			if out != bytecode.NoopOperand && out != dst {
				c.ctx.Emitter.EmitMove(dst, out)
			}
		}
		c.ctx.Symbols.ExitScope()
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
			c.ctx.Emitter.EmitJumpCompare(bytecode.OpJumpIfNeConst, valueReg, litOp, onFail)
		} else {
			c.ctx.Emitter.EmitJumpCompare(bytecode.OpJumpIfNe, valueReg, litOp, onFail)
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

	c.ctx.Symbols.EnterScope()
	out := ensureOperandRegister(c.ctx, c.facts, c.callbacks.compileExpr(expr))
	if out != bytecode.NoopOperand && out != dst {
		c.ctx.Emitter.EmitMove(dst, out)
	}
	c.ctx.Symbols.ExitScope()

	return true
}

func (c *exprMatchCompiler) compileMatchObjectPattern(valueReg bytecode.Operand, ctx fql.IMatchObjectPatternContext, onFail core.Label) {
	if ctx == nil {
		return
	}

	props := ctx.AllMatchObjectPatternProperty()
	if len(props) == 0 {
		keys := c.emitObjectsKeys(valueReg)
		c.ctx.Emitter.EmitJumpIfNone(keys, onFail)

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

		val := c.ctx.Registers.Allocate()
		if keyOp.IsConstant() {
			c.ctx.Emitter.EmitMatchLoadPropertyConst(val, valueReg, keyOp, onFail)
		} else {
			c.ctx.Emitter.EmitJumpCompare(bytecode.OpJumpIfMissingProperty, valueReg, keyOp, onFail)
			c.ctx.Emitter.EmitABC(bytecode.OpLoadProperty, val, valueReg, keyOp)
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
			return c.ctx.Symbols.AddConstant(val)
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

	return c.ctx.Symbols.AddConstant(runtime.NewString(name))
}

func (c *exprMatchCompiler) emitObjectsKeys(scrReg bytecode.Operand) bytecode.Operand {
	scrReg = ensureOperandRegister(c.ctx, c.facts, scrReg)
	seq := c.ctx.Registers.AllocateSequence(1)
	c.ctx.Emitter.EmitMove(seq[0], scrReg)

	return c.callbacks.compileFunctionCallByNameWith(nil, runtime.NewString("KEYS"), true, seq)
}

func (c *exprMatchCompiler) declareMatchBinding(ctx antlr.ParserRuleContext, name string, valueReg bytecode.Operand) bytecode.Operand {
	valueReg = ensureOperandRegister(c.ctx, c.facts, valueReg)
	reg, ok := c.ctx.Symbols.DeclareLocal(name, core.TypeAny)
	if ok {
		c.ctx.Emitter.EmitMove(reg, valueReg)
		c.ctx.Types.Set(reg, c.facts.OperandType(valueReg))

		return reg
	}

	if ctx != nil {
		c.ctx.Errors.DuplicateMatchBinding(ctx, name)
	}

	if existing, _, found := c.ctx.Symbols.Resolve(name); found {
		return existing
	}

	return valueReg
}
