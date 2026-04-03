package internal

import (
	"strconv"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (c *ExprCompiler) compileMatchExpression(ctx fql.IMatchExpressionContext) bytecode.Operand {
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

		scrReg := c.ensureRegister(c.Compile(scrutinee))
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

func (c *ExprCompiler) compileMatchPatternArms(scrReg bytecode.Operand, ctx fql.IMatchPatternArmsContext, dst bytecode.Operand, end core.Label) {
	if ctx == nil {
		return
	}

	arms := collectMatchPatternArms(ctx)
	mergeLabels, mergeGroups := collectMatchResultMerges(c, arms)
	defaultLabel, hasDefaultLabel := c.matchMergeDefaultLabel(mergeGroups)

	for idx, arm := range arms {
		c.compileMatchPatternArm(scrReg, arm, idx, mergeLabels, dst, end)
	}

	if hasDefaultLabel {
		c.compileMatchMergedResults(mergeGroups, defaultLabel, dst, end)
	}

	c.compileMatchPatternDefaultArm(ctx.MatchDefaultArm(), dst)
}

func collectMatchPatternArms(ctx fql.IMatchPatternArmsContext) []fql.IMatchPatternArmContext {
	if ctx == nil {
		return nil
	}

	list := ctx.MatchPatternArmList()
	if list == nil {
		return nil
	}

	return list.AllMatchPatternArm()
}

func (c *ExprCompiler) matchMergeDefaultLabel(groups []matchResultGroup) (core.Label, bool) {
	if len(groups) == 0 {
		return core.Label{}, false
	}

	return c.ctx.Emitter.NewLabel("match.default"), true
}

func (c *ExprCompiler) compileMatchPatternArm(scrReg bytecode.Operand, arm fql.IMatchPatternArmContext, idx int, mergeLabels map[int]core.Label, dst bytecode.Operand, end core.Label) {
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

func (c *ExprCompiler) compileMatchPatternArmConditions(scrReg bytecode.Operand, arm fql.IMatchPatternArmContext, next core.Label) {
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

	c.emitConditionJump(expr, next, false)
}

func (c *ExprCompiler) compileMatchPatternArmResult(arm fql.IMatchPatternArmContext, idx int, mergeLabels map[int]core.Label, dst bytecode.Operand, end core.Label) {
	result := arm.Expression()
	if result == nil {
		c.ctx.Emitter.EmitJump(end)
		return
	}

	if label, ok := mergeLabels[idx]; ok {
		c.ctx.Emitter.EmitJump(label)
		return
	}

	out := c.ensureRegister(c.Compile(result))
	if out != bytecode.NoopOperand && out != dst {
		c.ctx.Emitter.EmitMove(dst, out)
	}

	c.ctx.Emitter.EmitJump(end)
}

func (c *ExprCompiler) compileMatchMergedResults(groups []matchResultGroup, defaultLabel core.Label, dst bytecode.Operand, end core.Label) {
	c.ctx.Emitter.EmitJump(defaultLabel)

	for _, group := range groups {
		c.ctx.Emitter.MarkLabel(group.label)
		c.ctx.Symbols.EnterScope()
		out := c.ensureRegister(c.Compile(group.result))
		if out != bytecode.NoopOperand && out != dst {
			c.ctx.Emitter.EmitMove(dst, out)
		}
		c.ctx.Symbols.ExitScope()
		c.ctx.Emitter.EmitJump(end)
	}

	c.ctx.Emitter.MarkLabel(defaultLabel)
}

func (c *ExprCompiler) compileMatchPatternDefaultArm(def fql.IMatchDefaultArmContext, dst bytecode.Operand) {
	if def == nil {
		return
	}

	c.ctx.Symbols.EnterScope()
	result := def.Expression()
	if result != nil {
		out := c.ensureRegister(c.Compile(result))
		if out != bytecode.NoopOperand && out != dst {
			c.ctx.Emitter.EmitMove(dst, out)
		}
	}
	c.ctx.Symbols.ExitScope()
}

func (c *ExprCompiler) compileMatchGuardArms(ctx fql.IMatchGuardArmsContext, dst bytecode.Operand, end core.Label) {
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
			c.emitConditionJump(exprs[0], next, false)
		}

		if len(exprs) > 1 {
			out := c.ensureRegister(c.Compile(exprs[1]))
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
			out := c.ensureRegister(c.Compile(result))
			if out != bytecode.NoopOperand && out != dst {
				c.ctx.Emitter.EmitMove(dst, out)
			}
		}
		c.ctx.Symbols.ExitScope()
	}
}

func collectMatchResultMerges(c *ExprCompiler, arms []fql.IMatchPatternArmContext) (map[int]core.Label, []matchResultGroup) {
	if len(arms) == 0 {
		return nil, nil
	}

	groups := make(map[string]*matchResultGroup)
	order := make([]*matchResultGroup, 0)

	for idx, arm := range arms {
		if arm == nil {
			continue
		}
		if arm.MatchPatternGuard() != nil {
			continue
		}
		pattern := arm.MatchPattern()
		if pattern == nil || pattern.MatchLiteralPattern() == nil {
			continue
		}
		result := arm.Expression()
		if result == nil {
			continue
		}
		key, ok := matchPureResultKey(result)
		if !ok {
			continue
		}

		group, exists := groups[key]
		if !exists {
			group = &matchResultGroup{result: result}
			groups[key] = group
			order = append(order, group)
		}
		group.arms = append(group.arms, idx)
	}

	if len(order) == 0 {
		return nil, nil
	}

	labels := make(map[int]core.Label)
	merged := make([]matchResultGroup, 0)

	for _, group := range order {
		if len(group.arms) < 2 {
			continue
		}
		label := c.ctx.Emitter.NewLabel("match.result")
		group.label = label
		for _, idx := range group.arms {
			labels[idx] = label
		}
		merged = append(merged, *group)
	}

	if len(merged) == 0 {
		return nil, nil
	}

	return labels, merged
}

func (c *ExprCompiler) compileMatchPatternValue(valueReg bytecode.Operand, ctx fql.IMatchPatternContext, onFail core.Label) {
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

func (c *ExprCompiler) compileMatchLiteralOperand(ctx fql.IMatchLiteralPatternContext) bytecode.Operand {
	return compileScalarLiteralOperand(c.ctx, c.literals, ctx)
}

func (c *ExprCompiler) tryCompileMatchConstantFold(scrutinee fql.IExpressionContext, armsCtx fql.IMatchPatternArmsContext, dst bytecode.Operand) bool {
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

func matchConstantFoldPreconditions(scrutinee fql.IExpressionContext, armsCtx fql.IMatchPatternArmsContext) (runtime.Value, []fql.IMatchPatternArmContext, bool) {
	if scrutinee == nil || armsCtx == nil {
		return nil, nil, false
	}

	scrutineeVal, ok := matchLiteralValueFromExpression(scrutinee)
	if !ok {
		return nil, nil, false
	}

	list := armsCtx.MatchPatternArmList()
	if list == nil {
		return nil, nil, false
	}

	arms := list.AllMatchPatternArm()
	if len(arms) == 0 {
		return nil, nil, false
	}

	return scrutineeVal, arms, true
}

func selectMatchConstantFoldExpression(scrutinee runtime.Value, arms []fql.IMatchPatternArmContext, defaultArm fql.IMatchDefaultArmContext) (fql.IExpressionContext, bool) {
	var selected fql.IExpressionContext

	for _, arm := range arms {
		if arm == nil {
			continue
		}

		patternValue, expression, ok := matchConstantFoldArmExpression(arm)
		if !ok {
			return nil, false
		}

		if runtime.CompareValues(scrutinee, patternValue) == 0 {
			selected = expression
			break
		}
	}

	if selected == nil && defaultArm != nil {
		selected = defaultArm.Expression()
	}

	return selected, selected != nil
}

func matchConstantFoldArmExpression(arm fql.IMatchPatternArmContext) (runtime.Value, fql.IExpressionContext, bool) {
	if arm.MatchPatternGuard() != nil {
		return nil, nil, false
	}

	pattern := arm.MatchPattern()
	if pattern == nil {
		return nil, nil, false
	}

	literalPattern := pattern.MatchLiteralPattern()
	if literalPattern == nil {
		return nil, nil, false
	}

	patternValue, ok := literalValueFromMatchLiteral(literalPattern)
	if !ok {
		return nil, nil, false
	}

	return patternValue, arm.Expression(), true
}

func (c *ExprCompiler) emitMatchConstantFoldExpression(expr fql.IExpressionContext, dst bytecode.Operand) bool {
	if expr == nil {
		return false
	}

	c.ctx.Symbols.EnterScope()
	out := c.ensureRegister(c.Compile(expr))
	if out != bytecode.NoopOperand && out != dst {
		c.ctx.Emitter.EmitMove(dst, out)
	}
	c.ctx.Symbols.ExitScope()

	return true
}

func matchLiteralValueFromExpression(expr fql.IExpressionContext) (runtime.Value, bool) {
	atom, ok := matchPureResultAtom(expr)
	if !ok {
		return nil, false
	}

	literal := atom.Literal()
	if literal == nil {
		return nil, false
	}

	return literalValueFromLiteral(literal)
}

func literalValueFromLiteral(lit fql.ILiteralContext) (runtime.Value, bool) {
	return scalarLiteralValue(lit)
}

func literalValueFromMatchLiteral(lit fql.IMatchLiteralPatternContext) (runtime.Value, bool) {
	return scalarLiteralValue(lit)
}

func matchPureResultKey(expr fql.IExpressionContext) (string, bool) {
	atom, ok := matchPureResultAtom(expr)
	if !ok {
		return "", false
	}

	return matchPureResultKeyFromAtom(atom)
}

func matchPureResultAtom(expr fql.IExpressionContext) (fql.IExpressionAtomContext, bool) {
	if expr == nil {
		return nil, false
	}

	if expr.UnaryOperator() != nil || expr.LogicalAndOperator() != nil || expr.LogicalOrOperator() != nil || expr.GetTernaryOperator() != nil {
		return nil, false
	}

	pred := expr.Predicate()
	if pred == nil {
		return nil, false
	}

	if pred.EqualityOperator() != nil || pred.ArrayOperator() != nil || pred.InOperator() != nil || pred.LikeOperator() != nil {
		return nil, false
	}

	atom := pred.ExpressionAtom()
	if atom == nil {
		return nil, false
	}

	if atom.ExpressionAtom(0) != nil || atom.ExpressionAtom(1) != nil {
		return nil, false
	}

	if atom.AdditiveOperator() != nil || atom.MultiplicativeOperator() != nil || atom.RegexpOperator() != nil || atom.ErrorOperator() != nil || atom.RecoveryTails() != nil || atom.RangeOperator() != nil {
		return nil, false
	}

	return atom, true
}

func matchPureResultKeyFromAtom(atom fql.IExpressionAtomContext) (string, bool) {
	if atom == nil {
		return "", false
	}

	if inner := atom.Expression(); inner != nil {
		return matchPureResultKey(inner)
	}

	if lit := atom.Literal(); lit != nil {
		val, ok := literalValueFromLiteral(lit)
		if !ok {
			return "", false
		}
		return matchPureLiteralKey(val)
	}

	if v := atom.Variable(); v != nil {
		name := matchVariableName(v)
		if name == "" {
			return "", false
		}
		return "var:" + name, true
	}

	if p := atom.Param(); p != nil {
		name := matchParamName(p)
		if name == "" {
			return "", false
		}
		return "param:" + name, true
	}

	if m := atom.MemberExpression(); m != nil {
		return matchMemberExpressionKey(m)
	}

	if im := atom.ImplicitMemberExpression(); im != nil {
		return matchImplicitMemberExpressionKey(im)
	}

	return "", false
}

func matchPureLiteralKey(val runtime.Value) (string, bool) {
	if val == nil {
		return "", false
	}

	if val == runtime.None {
		return "lit:none", true
	}

	switch v := val.(type) {
	case runtime.Boolean:
		if v {
			return "lit:true", true
		}
		return "lit:false", true
	case runtime.Int:
		return "lit:int:" + strconv.FormatInt(int64(v), 10), true
	case runtime.Float:
		return "lit:float:" + strconv.FormatFloat(float64(v), 'g', -1, 64), true
	case runtime.String:
		return "lit:str:" + strconv.Quote(string(v)), true
	default:
		return "", false
	}
}

func matchVariableName(ctx fql.IVariableContext) string {
	if ctx == nil {
		return ""
	}

	if id := ctx.Identifier(); id != nil {
		return id.GetText()
	}

	if srw := ctx.SafeReservedWord(); srw != nil {
		return srw.GetText()
	}

	return ""
}

func matchParamName(ctx fql.IParamContext) string {
	if ctx == nil {
		return ""
	}

	var name string
	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if srw := ctx.SafeReservedWord(); srw != nil {
		name = srw.GetText()
	}

	if name == "" {
		return ""
	}

	return "@" + name
}

func matchMemberExpressionKey(ctx fql.IMemberExpressionContext) (string, bool) {
	if ctx == nil {
		return "", false
	}

	source := ctx.MemberExpressionSource()
	if source == nil {
		return "", false
	}

	var base string
	if v := source.Variable(); v != nil {
		name := matchVariableName(v)
		if name == "" {
			return "", false
		}
		base = "var:" + name
	} else if p := source.Param(); p != nil {
		name := matchParamName(p)
		if name == "" {
			return "", false
		}
		base = "param:" + name
	} else {
		return "", false
	}

	paths := ctx.AllMemberExpressionPath()
	if len(paths) != 1 {
		return "", false
	}

	prop, ok := matchMemberExpressionPathKey(paths[0])
	if !ok {
		return "", false
	}

	return "member:" + base + "." + prop, true
}

func matchImplicitMemberExpressionKey(ctx fql.IImplicitMemberExpressionContext) (string, bool) {
	if ctx == nil {
		return "", false
	}

	if len(ctx.AllMemberExpressionPath()) > 0 {
		return "", false
	}

	start := ctx.ImplicitMemberExpressionStart()
	if start == nil {
		return "", false
	}

	if start.ErrorOperator() != nil || start.ComputedPropertyName() != nil || start.ArrayExpansion() != nil || start.ArrayContraction() != nil || start.ArrayQuestionMark() != nil || start.ArrayApply() != nil {
		return "", false
	}

	prop, ok := matchPropertyNameKey(start.PropertyName())
	if !ok {
		return "", false
	}

	return "member:implicit:." + prop, true
}

func matchMemberExpressionPathKey(ctx fql.IMemberExpressionPathContext) (string, bool) {
	if ctx == nil {
		return "", false
	}

	if ctx.ErrorOperator() != nil || ctx.ComputedPropertyName() != nil || ctx.ArrayContraction() != nil || ctx.ArrayExpansion() != nil || ctx.ArrayQuestionMark() != nil || ctx.ArrayApply() != nil {
		return "", false
	}

	return matchPropertyNameKey(ctx.PropertyName())
}

func matchPropertyNameKey(ctx fql.IPropertyNameContext) (string, bool) {
	if ctx == nil {
		return "", false
	}

	var name string

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if srw := ctx.SafeReservedWord(); srw != nil {
		name = srw.GetText()
	} else if urw := ctx.UnsafeReservedWord(); urw != nil {
		name = urw.GetText()
	} else if sl := ctx.StringLiteral(); sl != nil {
		if val, ok := parseStringLiteralConst(sl); ok {
			name = string(val)
		}
	}

	if name == "" {
		return "", false
	}

	return strconv.Quote(name), true
}

func (c *ExprCompiler) compileMatchObjectPattern(valueReg bytecode.Operand, ctx fql.IMatchObjectPatternContext, onFail core.Label) {
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

func (c *ExprCompiler) compileMatchObjectPatternKey(ctx fql.IMatchObjectPatternKeyContext) bytecode.Operand {
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

func (c *ExprCompiler) emitObjectsKeys(scrReg bytecode.Operand) bytecode.Operand {
	scrReg = c.ensureRegister(scrReg)
	seq := c.ctx.Registers.AllocateSequence(1)
	c.ctx.Emitter.EmitMove(seq[0], scrReg)

	return c.CompileFunctionCallByNameWith(nil, runtime.NewString("KEYS"), true, seq)
}

func (c *ExprCompiler) declareMatchBinding(ctx antlr.ParserRuleContext, name string, valueReg bytecode.Operand) bytecode.Operand {
	valueReg = c.ensureRegister(valueReg)
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
