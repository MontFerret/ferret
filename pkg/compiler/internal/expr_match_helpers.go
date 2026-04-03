package internal

import (
	"strconv"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

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

func collectMatchResultMerges(ctx *CompilationSession, arms []fql.IMatchPatternArmContext) (map[int]core.Label, []matchResultGroup) {
	if len(arms) == 0 || ctx == nil {
		return nil, nil
	}

	groups := make(map[string]*matchResultGroup)
	order := make([]*matchResultGroup, 0)

	for idx, arm := range arms {
		if arm == nil || arm.MatchPatternGuard() != nil {
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

		label := ctx.Emitter.NewLabel("match.result")
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
