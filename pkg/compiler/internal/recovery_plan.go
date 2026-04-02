package internal

import (
	"strconv"
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
	recoveryActionRetry
)

type recoveryPlan struct {
	onError   *recoveryHandler
	onTimeout *recoveryHandler
}

type recoveryHandler struct {
	actionNode antlr.ParserRuleContext
	expr       fql.IExpressionContext
	retry      *recoveryRetryPlan
	tailNode   antlr.ParserRuleContext
	actionKind recoveryActionKind
}

type recoveryRetryPlan struct {
	delay           durationClause
	actionNode      antlr.ParserRuleContext
	countNode       antlr.ParserRuleContext
	delayNode       antlr.ParserRuleContext
	backoffNode     antlr.ParserRuleContext
	finalActionNode antlr.ParserRuleContext
	finalExpr       fql.IExpressionContext
	count           int
	backoff         waitForBackoff
	finalActionKind recoveryActionKind
	hasDelay        bool
	hasOr           bool
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
		reportInvalidRecoveryTail(ctx, tail, "Expected ERROR or TIMEOUT after 'ON' in recovery tail", "Complete the tail as ON ERROR FAIL, ON ERROR RETURN <expr>, ON ERROR RETRY <count>, ON TIMEOUT FAIL, or ON TIMEOUT RETURN <expr>.")
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

	reportInvalidRecoveryTail(ctx, cond.(antlr.ParserRuleContext), "Expected ERROR or TIMEOUT after 'ON' in recovery tail", "Complete the tail as ON ERROR FAIL, ON ERROR RETURN <expr>, ON ERROR RETRY <count>, ON TIMEOUT FAIL, or ON TIMEOUT RETURN <expr>.")

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

	if action.RecoveryActionOrClause() != nil {
		reportInvalidRecoveryTail(ctx, action.RecoveryActionOrClause(), "OR is only valid inside ON ERROR RETRY", "Use OR only after ON ERROR RETRY <count>, e.g. ON ERROR RETRY 3 OR RETURN NONE.")
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
	case action.RecoveryRetryAction() != nil:
		if condition != recoveryConditionError {
			reportInvalidRecoveryTail(ctx, action.RecoveryRetryAction(), "RETRY is only valid under ON ERROR", "Use ON ERROR RETRY <count> ... to retry failures, or ON TIMEOUT FAIL/RETURN for timeout handling.")
			return nil, false
		}

		retry, ok := resolveRecoveryRetryPlan(ctx, action.RecoveryRetryAction())
		if !ok {
			return nil, false
		}

		return &recoveryHandler{
			actionKind: recoveryActionRetry,
			actionNode: action.(antlr.ParserRuleContext),
			retry:      retry,
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

func resolveRecoveryRetryPlan(ctx *CompilerContext, action fql.IRecoveryRetryActionContext) (*recoveryRetryPlan, bool) {
	if action == nil {
		return nil, false
	}

	plan := &recoveryRetryPlan{
		actionNode:      action.(antlr.ParserRuleContext),
		backoff:         waitForBackoffNone,
		finalActionKind: recoveryActionFail,
	}

	valid := true

	countCtx := action.RecoveryRetryCount()
	if countCtx == nil || countCtx.IntegerLiteral() == nil {
		reportInvalidRecoveryTail(ctx, action, "Expected retry count after 'RETRY'", "Provide an integer retry count, e.g. ON ERROR RETRY 3.")
		valid = false
	} else {
		plan.countNode = countCtx.(antlr.ParserRuleContext)
		count, err := strconv.Atoi(countCtx.IntegerLiteral().GetText())
		if err != nil {
			reportInvalidRecoveryTail(ctx, countCtx, "Expected retry count after 'RETRY'", "Provide an integer retry count, e.g. ON ERROR RETRY 3.")
			valid = false
		} else {
			plan.count = count
		}
	}

	delayClause := action.RecoveryRetryDelayClause()
	if delayClause != nil {
		switch {
		case delayClause.DelayKeyword() != nil:
			if delayValue := delayClause.RecoveryRetryDelayValue(); delayValue != nil {
				plan.hasDelay = true
				plan.delay = delayValue
				plan.delayNode = delayValue.(antlr.ParserRuleContext)
			} else {
				reportInvalidRecoveryTail(ctx, delayClause.DelayKeyword(), "Expected value after 'DELAY' in retry policy", "Provide a duration or duration-like value, e.g. DELAY 100ms.")
				valid = false
			}

			backoffClause := delayClause.RecoveryRetryBackoffClause()
			if backoffClause != nil {
				backoff, ok := resolveRecoveryRetryBackoff(ctx, backoffClause)
				if !ok {
					valid = false
				} else {
					plan.backoff = backoff
					plan.backoffNode = backoffClause.(antlr.ParserRuleContext)
				}
			}
		case delayClause.RecoveryRetryBackoffClause() != nil:
			reportInvalidRecoveryTail(ctx, delayClause.RecoveryRetryBackoffClause(), "BACKOFF requires DELAY in retry policy", "Add DELAY before BACKOFF, e.g. ON ERROR RETRY 3 DELAY 100ms BACKOFF EXPONENTIAL.")
			valid = false

			if _, ok := resolveRecoveryRetryBackoff(ctx, delayClause.RecoveryRetryBackoffClause()); !ok {
				valid = false
			}
		}
	}

	orClauses := action.AllRecoveryRetryOrClause()
	if len(orClauses) > 1 {
		for i := 1; i < len(orClauses); i++ {
			reportInvalidRecoveryTail(ctx, orClauses[i], "Duplicate OR fallback in retry policy", "Each retry policy may define OR at most once.")
		}
		valid = false
	}

	if len(orClauses) > 0 {
		plan.hasOr = true

		finalKind, finalExpr, finalNode, ok := resolveRecoveryRetryFinalAction(ctx, orClauses[0])
		if !ok {
			valid = false
		} else {
			plan.finalActionKind = finalKind
			plan.finalExpr = finalExpr
			plan.finalActionNode = finalNode
		}
	}

	if !valid {
		return nil, false
	}

	return plan, true
}

func resolveRecoveryRetryBackoff(ctx *CompilerContext, clause fql.IRecoveryRetryBackoffClauseContext) (waitForBackoff, bool) {
	if clause == nil {
		return waitForBackoffNone, true
	}

	kind := clause.RecoveryRetryBackoffKind()
	if kind == nil {
		reportInvalidRecoveryTail(ctx, clause, "Expected backoff kind after 'BACKOFF' in retry policy", "Use BACKOFF CONSTANT, BACKOFF LINEAR, or BACKOFF EXPONENTIAL.")
		return waitForBackoffNone, false
	}

	raw := ""

	switch {
	case kind.Identifier() != nil:
		raw = kind.Identifier().GetText()
	case kind.StringLiteral() != nil:
		if parsed, ok := parseStringLiteralConst(kind.StringLiteral()); ok {
			raw = parsed.String()
		}
	case kind.None() != nil:
		raw = kind.None().GetText()
	}

	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "CONSTANT":
		return waitForBackoffNone, true
	case "LINEAR":
		return waitForBackoffLinear, true
	case "EXPONENTIAL":
		return waitForBackoffExponential, true
	default:
		reportInvalidRecoveryTail(ctx, kind, "Unknown BACKOFF strategy", "Use one of: CONSTANT, LINEAR, EXPONENTIAL.")
		return waitForBackoffNone, false
	}
}

func resolveRecoveryRetryFinalAction(ctx *CompilerContext, clause fql.IRecoveryRetryOrClauseContext) (recoveryActionKind, fql.IExpressionContext, antlr.ParserRuleContext, bool) {
	if clause == nil {
		return recoveryActionUnknown, nil, nil, false
	}

	action := clause.RecoveryRetryFinalAction()
	if action == nil {
		reportInvalidRecoveryTail(ctx, clause, "Expected FAIL or RETURN after 'OR' in retry fallback", "Complete the retry fallback as OR FAIL or OR RETURN <expr>.")
		return recoveryActionUnknown, nil, nil, false
	}

	switch {
	case action.FailKeyword() != nil:
		return recoveryActionFail, nil, action.(antlr.ParserRuleContext), true
	case action.ReturnKeyword() != nil:
		returnExpr := action.RecoveryReturnExpr()
		if returnExpr == nil || returnExpr.Expression() == nil {
			reportInvalidRecoveryTail(ctx, action, "Expected expression after 'RETURN' in recovery tail", "Provide a fallback expression, e.g. OR RETURN NONE.")
			return recoveryActionUnknown, nil, nil, false
		}

		return recoveryActionReturn, returnExpr.Expression(), action.(antlr.ParserRuleContext), true
	default:
		reportInvalidRecoveryTail(ctx, action, "Expected FAIL or RETURN after 'OR' in retry fallback", "Complete the retry fallback as OR FAIL or OR RETURN <expr>.")
		return recoveryActionUnknown, nil, nil, false
	}
}

func missingRecoveryActionMessage(condition recoveryCondition) string {
	switch condition {
	case recoveryConditionTimeout:
		return "Expected FAIL or RETURN after 'ON TIMEOUT'"
	default:
		return "Expected FAIL, RETURN, or RETRY after 'ON ERROR'"
	}
}

func missingRecoveryActionHint(condition recoveryCondition) string {
	switch condition {
	case recoveryConditionTimeout:
		return "Use ON TIMEOUT FAIL to propagate timeout expiration or ON TIMEOUT RETURN <expr> to supply a fallback value."
	default:
		return "Use ON ERROR FAIL to propagate failures, ON ERROR RETURN <expr> to supply a fallback value, or ON ERROR RETRY <count> to retry."
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

func recoveryPlanHasReturnHandler(plan recoveryPlan) bool {
	return recoveryHandlerReturns(plan.onError) || recoveryHandlerReturns(plan.onTimeout)
}

func mergeRecoveryPlans(ctx *CompilerContext, primary, extra recoveryPlan) recoveryPlan {
	merged := primary

	if extra.onTimeout != nil {
		if merged.onTimeout != nil {
			reportInvalidRecoveryTail(ctx, extra.onTimeout.tailNode, "Duplicate ON TIMEOUT handler", "Each operation may define ON TIMEOUT at most once.")
		} else {
			merged.onTimeout = extra.onTimeout
		}
	}

	if extra.onError != nil {
		if merged.onError != nil {
			reportInvalidRecoveryTail(ctx, extra.onError.tailNode, "Duplicate ON ERROR handler", "Each operation may define ON ERROR at most once.")
		} else {
			merged.onError = extra.onError
		}
	}

	return merged
}

func allowsTailCallRecovery(plan recoveryPlan) bool {
	return plan.onError == nil || plan.onError.actionKind == recoveryActionFail
}

func recoveryHandlerReturns(handler *recoveryHandler) bool {
	switch {
	case handler == nil:
		return false
	case handler.actionKind == recoveryActionReturn:
		return true
	case handler.actionKind == recoveryActionRetry && handler.retry != nil:
		return handler.retry.finalActionKind == recoveryActionReturn
	default:
		return false
	}
}

func recoveryHandlerRetries(handler *recoveryHandler) bool {
	return handler != nil && handler.actionKind == recoveryActionRetry && handler.retry != nil
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
