package internal

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

const suppressRecoveryAction = "SUPPRESS"

type recoveryCondition int

const (
	recoveryConditionUnknown recoveryCondition = iota
	recoveryConditionError
	recoveryConditionTimeout
)

type recoveryActionKind int

const (
	recoveryActionUnknown recoveryActionKind = iota
	recoveryActionFail
	recoveryActionReturn
)

type recoveryPlan struct {
	onError   *recoveryHandler
	onTimeout *recoveryHandler
}

type recoveryHandler struct {
	actionNode antlr.ParserRuleContext
	expr       fql.IExpressionContext
	tailNode   antlr.ParserRuleContext
	actionKind recoveryActionKind
}

type recoveryPlanOptions struct {
	allowTimeout bool
	hasTimeout   bool
}

type recoveryTailOwner interface {
	RecoveryTails() fql.IRecoveryTailsContext
}

func collectRecoveryPlan(ctx *CompilerContext, owner recoveryTailOwner, options recoveryPlanOptions) recoveryPlan {
	if ctx == nil || owner == nil {
		return recoveryPlan{}
	}

	tails := owner.RecoveryTails()
	if tails == nil {
		return recoveryPlan{}
	}

	var plan recoveryPlan

	for _, tail := range tails.AllRecoveryTail() {
		condition, ok := resolveRecoveryCondition(ctx, tail)
		if !ok {
			continue
		}

		if condition == recoveryConditionTimeout {
			if !options.allowTimeout {
				reportInvalidRecoveryTail(ctx, tail, "ON TIMEOUT is only valid on WAITFOR operations", "Use ON TIMEOUT only on WAITFOR expressions that define timeout handling.")
				continue
			}

			if !options.hasTimeout {
				reportInvalidRecoveryTail(ctx, tail, "ON TIMEOUT requires a TIMEOUT clause on WAITFOR", "Add a TIMEOUT clause before ON TIMEOUT, e.g. WAITFOR VALUE x TIMEOUT 1s ON TIMEOUT RETURN NONE.")
				continue
			}
		}

		handler, ok := resolveRecoveryHandler(ctx, tail, condition)
		if !ok {
			continue
		}

		switch condition {
		case recoveryConditionError:
			if plan.onError != nil {
				reportInvalidRecoveryTail(ctx, tail, "Duplicate ON ERROR handler", "Each operation may define ON ERROR at most once.")
				continue
			}

			plan.onError = handler
		case recoveryConditionTimeout:
			if plan.onTimeout != nil {
				reportInvalidRecoveryTail(ctx, tail, "Duplicate ON TIMEOUT handler", "Each operation may define ON TIMEOUT at most once.")
				continue
			}

			plan.onTimeout = handler
		}
	}

	return plan
}

func resolveRecoveryCondition(ctx *CompilerContext, tail fql.IRecoveryTailContext) (recoveryCondition, bool) {
	if tail == nil {
		return recoveryConditionUnknown, false
	}

	cond := tail.RecoveryCondition()
	if cond == nil {
		reportInvalidRecoveryTail(ctx, tail, "Expected ERROR or TIMEOUT after 'ON' in recovery tail", "Complete the tail as ON ERROR FAIL, ON ERROR RETURN <expr>, ON TIMEOUT FAIL, or ON TIMEOUT RETURN <expr>.")
		return recoveryConditionUnknown, false
	}

	switch {
	case cond.TimeoutKeyword() != nil:
		return recoveryConditionTimeout, true
	case cond.ErrorKeyword() != nil:
		return recoveryConditionError, true
	}

	if strings.EqualFold(cond.GetText(), suppressRecoveryAction) {
		reportInvalidRecoveryTail(ctx, cond.(antlr.ParserRuleContext), "SUPPRESS is not supported in recovery tails", "Use ON ERROR RETURN NONE instead.")
		return recoveryConditionUnknown, false
	}

	reportInvalidRecoveryTail(ctx, cond.(antlr.ParserRuleContext), "Expected ERROR or TIMEOUT after 'ON' in recovery tail", "Complete the tail as ON ERROR FAIL, ON ERROR RETURN <expr>, ON TIMEOUT FAIL, or ON TIMEOUT RETURN <expr>.")

	return recoveryConditionUnknown, false
}

func resolveRecoveryHandler(ctx *CompilerContext, tail fql.IRecoveryTailContext, condition recoveryCondition) (*recoveryHandler, bool) {
	if tail == nil {
		return nil, false
	}

	action := tail.RecoveryAction()
	if action == nil {
		reportInvalidRecoveryTail(ctx, tail, missingRecoveryActionMessage(condition), missingRecoveryActionHint(condition))
		return nil, false
	}

	switch {
	case action.FailKeyword() != nil:
		return &recoveryHandler{
			actionKind: recoveryActionFail,
			actionNode: action.(antlr.ParserRuleContext),
			tailNode:   tail.(antlr.ParserRuleContext),
		}, true
	case action.ReturnKeyword() != nil:
		returnExpr := action.RecoveryReturnExpr()
		if returnExpr == nil || returnExpr.Expression() == nil {
			reportInvalidRecoveryTail(ctx, action.(antlr.ParserRuleContext), "Expected expression after 'RETURN' in recovery tail", "Provide a fallback expression, e.g. ON ERROR RETURN NONE.")
			return nil, false
		}

		return &recoveryHandler{
			actionKind: recoveryActionReturn,
			actionNode: action.(antlr.ParserRuleContext),
			expr:       returnExpr.Expression(),
			tailNode:   tail.(antlr.ParserRuleContext),
		}, true
	}

	if strings.EqualFold(action.GetText(), suppressRecoveryAction) {
		reportInvalidRecoveryTail(ctx, action.(antlr.ParserRuleContext), "SUPPRESS is not supported in recovery tails", "Use ON ERROR RETURN NONE instead.")
		return nil, false
	}

	reportInvalidRecoveryTail(ctx, action.(antlr.ParserRuleContext), missingRecoveryActionMessage(condition), missingRecoveryActionHint(condition))

	return nil, false
}

func missingRecoveryActionMessage(condition recoveryCondition) string {
	switch condition {
	case recoveryConditionTimeout:
		return "Expected FAIL or RETURN after 'ON TIMEOUT'"
	default:
		return "Expected FAIL or RETURN after 'ON ERROR'"
	}
}

func missingRecoveryActionHint(condition recoveryCondition) string {
	switch condition {
	case recoveryConditionTimeout:
		return "Use ON TIMEOUT FAIL to propagate timeout expiration or ON TIMEOUT RETURN <expr> to supply a fallback value."
	default:
		return "Use ON ERROR FAIL to propagate failures or ON ERROR RETURN <expr> to supply a fallback value."
	}
}

func reportInvalidRecoveryTail(ctx *CompilerContext, node antlr.ParserRuleContext, message, hint string) {
	if ctx == nil || ctx.Errors == nil || node == nil {
		return
	}

	err := ctx.Errors.Create(parserd.SyntaxError, node, message)
	err.Hint = hint
	ctx.Errors.Add(err)
}

func hasErrorReturnNoneHandler(plan recoveryPlan) bool {
	if plan.onError == nil || plan.onError.actionKind != recoveryActionReturn {
		return false
	}

	return isNoneLiteralExpression(plan.onError.expr)
}

func allowsTailCallRecovery(plan recoveryPlan) bool {
	return plan.onError == nil || plan.onError.actionKind == recoveryActionFail
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
