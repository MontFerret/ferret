package internal

import (
	"math"
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
	if ctx == nil {
		return false
	}

	if ev := ctx.WaitForEventExpression(); ev != nil && ev.TimeoutClause() != nil {
		return true
	}

	if pred := ctx.WaitForPredicateExpression(); pred != nil && pred.TimeoutClause() != nil {
		return true
	}

	return false
}

func parseDurationLiteral(text string) (runtime.Value, error) {
	raw := normalizeDurationLiteral(text)
	if raw == "" {
		return runtime.None, strconv.ErrSyntax
	}

	number, unit, ok := splitDurationLiteral(raw)
	if !ok || number == "" {
		return runtime.None, strconv.ErrSyntax
	}

	value, err := parseDurationLiteralNumber(number)
	if err != nil {
		return runtime.None, err
	}

	multiplier, ok := durationUnitMultiplier(unit)
	if !ok {
		return runtime.None, strconv.ErrSyntax
	}

	ms := value * multiplier
	return parseDurationMillisecondsValue(ms)
}

func normalizeDurationLiteral(text string) string {
	return strings.ToUpper(strings.TrimSpace(text))
}

func splitDurationLiteral(raw string) (string, string, bool) {
	switch {
	case strings.HasSuffix(raw, "MS"):
		return strings.TrimSuffix(raw, "MS"), "MS", true
	case strings.HasSuffix(raw, "S"):
		return strings.TrimSuffix(raw, "S"), "S", true
	case strings.HasSuffix(raw, "M"):
		return strings.TrimSuffix(raw, "M"), "M", true
	case strings.HasSuffix(raw, "H"):
		return strings.TrimSuffix(raw, "H"), "H", true
	case strings.HasSuffix(raw, "D"):
		return strings.TrimSuffix(raw, "D"), "D", true
	default:
		return "", "", false
	}
}

func parseDurationLiteralNumber(raw string) (float64, error) {
	return strconv.ParseFloat(raw, 64)
}

func durationUnitMultiplier(unit string) (float64, bool) {
	switch unit {
	case "MS":
		return 1, true
	case "S":
		return 1000, true
	case "M":
		return 60000, true
	case "H":
		return 3600000, true
	case "D":
		return 86400000, true
	default:
		return 0, false
	}
}

func parseDurationMillisecondsValue(ms float64) (runtime.Value, error) {
	if err := validateDurationMilliseconds(ms); err != nil {
		return runtime.None, err
	}

	return durationValueFromMilliseconds(ms), nil
}

func validateDurationMilliseconds(ms float64) error {
	if math.IsNaN(ms) || math.IsInf(ms, 0) {
		return strconv.ErrRange
	}

	const (
		maxInt64Float = float64(1<<63 - 1)
		minInt64Float = -float64(1 << 63)
	)

	if ms < minInt64Float || ms > maxInt64Float {
		return strconv.ErrRange
	}

	return nil
}

func durationValueFromMilliseconds(ms float64) runtime.Value {
	if frac := math.Mod(ms, 1); frac == 0 {
		return runtime.NewInt64(int64(ms))
	}

	return runtime.NewFloat(ms)
}
