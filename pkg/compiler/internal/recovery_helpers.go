package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func missingRecoveryActionMessage(condition core.RecoveryCondition) string {
	switch condition {
	case core.RecoveryConditionTimeout:
		return "Expected FAIL or RETURN after 'ON TIMEOUT'"
	default:
		return "Expected FAIL, RETURN, or RETRY after 'ON ERROR'"
	}
}

func missingRecoveryActionHint(condition core.RecoveryCondition) string {
	switch condition {
	case core.RecoveryConditionTimeout:
		return "Use ON TIMEOUT FAIL to propagate timeout expiration or ON TIMEOUT RETURN <expr> to supply a fallback value."
	default:
		return "Use ON ERROR FAIL to propagate failures, ON ERROR RETURN <expr> to supply a fallback value, or ON ERROR RETRY <count> to retry."
	}
}

func hasErrorReturnNoneHandler(plan core.RecoveryPlan) bool {
	if plan.OnError == nil || plan.OnError.ActionKind != core.RecoveryActionReturn {
		return false
	}

	return isNoneLiteralExpression(plan.OnError.Expr)
}

func recoveryPlanHasReturnHandler(plan core.RecoveryPlan) bool {
	return recoveryHandlerReturns(plan.OnError) || recoveryHandlerReturns(plan.OnTimeout)
}

func allowsTailCallRecovery(plan core.RecoveryPlan) bool {
	return plan.OnError == nil || plan.OnError.ActionKind == core.RecoveryActionFail
}

func recoveryHandlerReturns(handler *core.RecoveryHandler) bool {
	switch {
	case handler == nil:
		return false
	case handler.ActionKind == core.RecoveryActionReturn:
		return true
	case handler.ActionKind == core.RecoveryActionRetry && handler.Retry != nil:
		return handler.Retry.FinalActionKind == core.RecoveryActionReturn
	default:
		return false
	}
}

func recoveryHandlerRetries(handler *core.RecoveryHandler) bool {
	return handler != nil && handler.ActionKind == core.RecoveryActionRetry && handler.Retry != nil
}

func isNoneLiteralExpression(expr fql.IExpressionContext) bool {
	if expr == nil || expr.UnaryOperator() != nil || expr.LogicalAndOperator() != nil || expr.LogicalOrOperator() != nil || expr.GetTernaryOperator() != nil {
		return false
	}

	predicate := expr.Predicate()
	if predicate == nil || predicate.EqualityOperator() != nil || predicate.ArrayOperator() != nil || predicate.InOperator() != nil || predicate.LikeOperator() != nil {
		return false
	}

	atom := predicate.ExpressionAtom()
	if atom == nil || atom.MultiplicativeOperator() != nil || atom.AdditiveOperator() != nil || atom.RegexpOperator() != nil {
		return false
	}

	literal := atom.Literal()
	return literal != nil && literal.NoneLiteral() != nil
}

func catchJumpModeForForExpression(ctx fql.IForExpressionContext) core.CatchJumpMode {
	if ctx != nil && ctx.In() != nil {
		return core.CatchJumpModeEnd
	}

	return core.CatchJumpModeNone
}

func catchJumpModeForWaitForExpression(ctx fql.IWaitForExpressionContext) core.CatchJumpMode {
	if ctx != nil && ctx.WaitForEventExpression() != nil {
		return core.CatchJumpModeEnd
	}

	return core.CatchJumpModeNone
}
