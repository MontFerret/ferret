package internal

import (
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func waitForSpan(src antlr.RuleContext, fallback antlr.RuleContext) source.Span {
	span := source.Span{Start: -1, End: -1}

	if src != nil {
		if prc, ok := src.(antlr.ParserRuleContext); ok {
			span = parserd.SpanFromRuleContext(prc)
			return span
		}
	}

	if fallback != nil {
		if prc, ok := fallback.(antlr.ParserRuleContext); ok {
			span = parserd.SpanFromRuleContext(prc)
		}
	}

	return span
}

func legacyWaitForOrThrowNode(expr fql.IExpressionContext) antlr.ParserRuleContext {
	if expr == nil || expr.LogicalOrOperator() == nil {
		return nil
	}

	return bareThrowExpressionNode(expr.GetRight())
}

func bareThrowExpressionNode(expr fql.IExpressionContext) antlr.ParserRuleContext {
	if expr == nil {
		return nil
	}

	if expr.UnaryOperator() != nil || expr.LogicalAndOperator() != nil || expr.LogicalOrOperator() != nil || expr.GetTernaryOperator() != nil {
		return nil
	}

	return bareThrowPredicateNode(expr.Predicate())
}

func bareThrowPredicateNode(pred fql.IPredicateContext) antlr.ParserRuleContext {
	if pred == nil {
		return nil
	}

	if pred.EqualityOperator() != nil || pred.ArrayOperator() != nil || pred.InOperator() != nil || pred.LikeOperator() != nil {
		return nil
	}

	return bareThrowAtomNode(pred.ExpressionAtom())
}

func bareThrowAtomNode(atom fql.IExpressionAtomContext) antlr.ParserRuleContext {
	if atom == nil {
		return nil
	}

	if atom.MultiplicativeOperator() != nil || atom.AdditiveOperator() != nil || atom.RegexpOperator() != nil {
		return nil
	}

	variable := atom.Variable()
	if variable == nil || !strings.EqualFold(matchVariableName(variable), "THROW") {
		return nil
	}

	node, ok := variable.(antlr.ParserRuleContext)
	if !ok {
		return nil
	}

	return node
}

func resolveWaitPredicateMode(hasValue, hasExists, hasNot bool) waitForPredicateMode {
	if hasValue {
		return waitForPredicateModeValue
	}

	if hasExists {
		if hasNot {
			return waitForPredicateModeNotExists
		}

		return waitForPredicateModeExists
	}

	return waitForPredicateModeBool
}

func waitPredicateWhenExpressions(ctxs []fql.IWaitForPredicateWhenClauseContext) []fql.IExpressionContext {
	if len(ctxs) == 0 {
		return nil
	}

	exprs := make([]fql.IExpressionContext, 0, len(ctxs))
	for _, ctx := range ctxs {
		if ctx == nil || ctx.Expression() == nil {
			continue
		}

		exprs = append(exprs, ctx.Expression())
	}

	return exprs
}

func literalFromExpression(ctx fql.IExpressionContext) fql.ILiteralContext {
	if ctx == nil {
		return nil
	}

	predicate := ctx.Predicate()
	if predicate == nil {
		return nil
	}

	atom := predicate.ExpressionAtom()
	if atom == nil {
		return nil
	}

	return atom.Literal()
}

func literalPresentFromExpression(ctx fql.IExpressionContext) (bool, bool) {
	lit := literalFromExpression(ctx)
	if lit == nil {
		return false, false
	}

	return lit.NoneLiteral() == nil, true
}

func literalExistsFromExpression(ctx fql.IExpressionContext) (bool, bool) {
	lit := literalFromExpression(ctx)
	if lit == nil {
		return false, false
	}

	switch {
	case lit.NoneLiteral() != nil:
		return false, true
	case lit.StringLiteral() != nil:
		if str, ok := parseStringLiteralConst(lit.StringLiteral()); ok {
			return str.String() != "", true
		}
		return false, false
	case lit.ArrayLiteral() != nil:
		arr := lit.ArrayLiteral()
		return arr.ArgumentList() != nil, true
	case lit.ObjectLiteral() != nil:
		obj := lit.ObjectLiteral()
		return len(obj.AllPropertyAssignment()) > 0, true
	default:
		return true, true
	}
}

func literalTruthinessFromExpression(ctx fql.IExpressionContext) (bool, bool) {
	lit := literalFromExpression(ctx)
	if lit == nil {
		return false, false
	}

	switch {
	case lit.NoneLiteral() != nil:
		return false, true
	case lit.BooleanLiteral() != nil:
		return strings.ToLower(lit.BooleanLiteral().GetText()) == "true", true
	case lit.IntegerLiteral() != nil:
		val, err := strconv.Atoi(lit.IntegerLiteral().GetText())
		if err != nil {
			return false, false
		}
		return val != 0, true
	case lit.FloatLiteral() != nil:
		val, err := strconv.ParseFloat(lit.FloatLiteral().GetText(), 64)
		if err != nil {
			return false, false
		}
		return val != 0, true
	case lit.StringLiteral() != nil:
		if str, ok := parseStringLiteralConst(lit.StringLiteral()); ok {
			return str.String() != "", true
		}
		return false, false
	default:
		return true, true
	}
}

func waitForHasExplicitTimeoutClause(ctx fql.IWaitForExpressionContext) bool {
	return waitForTimeoutClause(ctx) != nil
}

func waitForTimeoutClause(ctx fql.IWaitForExpressionContext) fql.ITimeoutClauseContext {
	if ctx == nil {
		return nil
	}

	if ev := ctx.WaitForEventExpression(); ev != nil {
		return waitForEventTimeoutClause(ev)
	}

	if pred := ctx.WaitForPredicateExpression(); pred != nil {
		return pred.TimeoutClause()
	}

	return nil
}

func parseDurationLiteral(text string) (runtime.Value, error) {
	return runtime.ParseDuration(text)
}
