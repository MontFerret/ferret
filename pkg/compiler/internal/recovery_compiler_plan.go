package internal

import (
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func (c *RecoveryCompiler) CollectPlan(owner core.RecoveryTailOwner, options core.RecoveryPlanOptions) core.RecoveryPlan {
	if c == nil || c.ctx == nil || owner == nil {
		return core.RecoveryPlan{}
	}

	tails := owner.RecoveryTails()
	if tails == nil {
		return core.RecoveryPlan{}
	}

	var plan core.RecoveryPlan

	for _, tail := range tails.AllRecoveryTail() {
		condition, ok := c.resolveCondition(tail)
		if !ok {
			continue
		}

		if condition == core.RecoveryConditionTimeout {
			if !options.AllowTimeout {
				c.reportInvalidTail(tail, "ON TIMEOUT is only valid on WAITFOR operations", "Use ON TIMEOUT only on WAITFOR expressions that define timeout handling.")
				continue
			}

			if !options.HasTimeout {
				c.reportInvalidTail(tail, "ON TIMEOUT requires a TIMEOUT clause on WAITFOR", "Add a TIMEOUT clause before ON TIMEOUT, e.g. WAITFOR VALUE x TIMEOUT 1s ON TIMEOUT RETURN NONE.")
				continue
			}
		}

		handler, ok := c.resolveHandler(tail, condition)
		if !ok {
			continue
		}

		switch condition {
		case core.RecoveryConditionError:
			if plan.OnError != nil {
				c.reportInvalidTail(tail, "Duplicate ON ERROR handler", "Each operation may define ON ERROR at most once.")
				continue
			}

			plan.OnError = handler
		case core.RecoveryConditionTimeout:
			if plan.OnTimeout != nil {
				c.reportInvalidTail(tail, "Duplicate ON TIMEOUT handler", "Each operation may define ON TIMEOUT at most once.")
				continue
			}

			plan.OnTimeout = handler
		}
	}

	return plan
}

func (c *RecoveryCompiler) NormalizePlan(plan core.RecoveryPlan) core.RecoveryPlan {
	normalized := plan

	if recoveryHandlerRetries(normalized.OnError) && normalized.OnError.Retry != nil &&
		normalized.OnError.Retry.FinalActionKind != core.RecoveryActionReturn &&
		normalized.OnError.Retry.Count <= 0 {
		normalized.OnError = nil
	}

	return normalized
}

func (c *RecoveryCompiler) MergePlans(primary, extra core.RecoveryPlan) core.RecoveryPlan {
	merged := primary

	if extra.OnTimeout != nil {
		if merged.OnTimeout != nil {
			c.reportInvalidTail(extra.OnTimeout.TailNode, "Duplicate ON TIMEOUT handler", "Each operation may define ON TIMEOUT at most once.")
		} else {
			merged.OnTimeout = extra.OnTimeout
		}
	}

	if extra.OnError != nil {
		if merged.OnError != nil {
			c.reportInvalidTail(extra.OnError.TailNode, "Duplicate ON ERROR handler", "Each operation may define ON ERROR at most once.")
		} else {
			merged.OnError = extra.OnError
		}
	}

	return merged
}

func (c *RecoveryCompiler) WidenResultType(out bytecode.Operand, plan core.RecoveryPlan) bytecode.Operand {
	if c == nil || c.ctx == nil || !out.IsRegister() || !recoveryPlanHasReturnHandler(plan) {
		return out
	}

	c.ctx.Types.Set(out, core.TypeAny)

	return out
}

func (c *RecoveryCompiler) EnsureRegister(op bytecode.Operand) bytecode.Operand {
	if c == nil || c.ctx == nil || op == bytecode.NoopOperand || op.IsRegister() {
		return op
	}

	dst := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadConst(dst, op)
	c.ctx.Types.Set(dst, c.front.TypeFacts.OperandType(op))

	return dst
}

func (c *RecoveryCompiler) resolveCondition(tail fql.IRecoveryTailContext) (core.RecoveryCondition, bool) {
	if tail == nil {
		return core.RecoveryConditionUnknown, false
	}

	cond := tail.RecoveryCondition()
	if cond == nil {
		c.reportInvalidTail(tail, "Expected ERROR or TIMEOUT after 'ON' in recovery tail", "Complete the tail as ON ERROR FAIL, ON ERROR RETURN <expr>, ON ERROR RETRY <count>, ON TIMEOUT FAIL, or ON TIMEOUT RETURN <expr>.")
		return core.RecoveryConditionUnknown, false
	}

	switch {
	case cond.TimeoutKeyword() != nil:
		return core.RecoveryConditionTimeout, true
	case cond.ErrorKeyword() != nil:
		return core.RecoveryConditionError, true
	}

	if strings.EqualFold(cond.GetText(), core.SuppressRecoveryAction) {
		c.reportInvalidTail(cond.(antlr.ParserRuleContext), "SUPPRESS is not supported in recovery tails", "Use ON ERROR RETURN NONE instead.")
		return core.RecoveryConditionUnknown, false
	}

	c.reportInvalidTail(cond.(antlr.ParserRuleContext), "Expected ERROR or TIMEOUT after 'ON' in recovery tail", "Complete the tail as ON ERROR FAIL, ON ERROR RETURN <expr>, ON ERROR RETRY <count>, ON TIMEOUT FAIL, or ON TIMEOUT RETURN <expr>.")

	return core.RecoveryConditionUnknown, false
}

func (c *RecoveryCompiler) resolveHandler(tail fql.IRecoveryTailContext, condition core.RecoveryCondition) (*core.RecoveryHandler, bool) {
	if tail == nil {
		return nil, false
	}

	action := tail.RecoveryAction()
	if action == nil {
		c.reportInvalidTail(tail, missingRecoveryActionMessage(condition), missingRecoveryActionHint(condition))
		return nil, false
	}

	if action.RecoveryActionOrClause() != nil {
		c.reportInvalidTail(action.RecoveryActionOrClause(), "OR is only valid inside ON ERROR RETRY", "Use OR only after ON ERROR RETRY <count>, e.g. ON ERROR RETRY 3 OR RETURN NONE.")
		return nil, false
	}

	switch {
	case action.FailKeyword() != nil:
		return &core.RecoveryHandler{
			ActionKind: core.RecoveryActionFail,
			ActionNode: action.(antlr.ParserRuleContext),
			TailNode:   tail.(antlr.ParserRuleContext),
		}, true
	case action.ReturnKeyword() != nil:
		returnExpr := action.RecoveryReturnExpr()
		if returnExpr == nil || returnExpr.Expression() == nil {
			c.reportInvalidTail(action.(antlr.ParserRuleContext), "Expected expression after 'RETURN' in recovery tail", "Provide a fallback expression, e.g. ON ERROR RETURN NONE.")
			return nil, false
		}

		return &core.RecoveryHandler{
			ActionKind: core.RecoveryActionReturn,
			ActionNode: action.(antlr.ParserRuleContext),
			Expr:       returnExpr.Expression(),
			TailNode:   tail.(antlr.ParserRuleContext),
		}, true
	case action.RecoveryRetryAction() != nil:
		if condition != core.RecoveryConditionError {
			c.reportInvalidTail(action.RecoveryRetryAction(), "RETRY is only valid under ON ERROR", "Use ON ERROR RETRY <count> ... to retry failures, or ON TIMEOUT FAIL/RETURN for timeout handling.")
			return nil, false
		}

		retry, ok := c.resolveRetryPlan(action.RecoveryRetryAction())
		if !ok {
			return nil, false
		}

		return &core.RecoveryHandler{
			ActionKind: core.RecoveryActionRetry,
			ActionNode: action.(antlr.ParserRuleContext),
			Retry:      retry,
			TailNode:   tail.(antlr.ParserRuleContext),
		}, true
	}

	if strings.EqualFold(action.GetText(), core.SuppressRecoveryAction) {
		c.reportInvalidTail(action.(antlr.ParserRuleContext), "SUPPRESS is not supported in recovery tails", "Use ON ERROR RETURN NONE instead.")
		return nil, false
	}

	c.reportInvalidTail(action.(antlr.ParserRuleContext), missingRecoveryActionMessage(condition), missingRecoveryActionHint(condition))

	return nil, false
}

func (c *RecoveryCompiler) resolveRetryPlan(action fql.IRecoveryRetryActionContext) (*core.RecoveryRetryPlan, bool) {
	if action == nil {
		return nil, false
	}

	plan := &core.RecoveryRetryPlan{
		RetryPlan: core.RetryPlan{
			Backoff: core.RetryBackoffNone,
		},
		ActionNode:      action.(antlr.ParserRuleContext),
		FinalActionKind: core.RecoveryActionFail,
	}

	valid := true

	countCtx := action.RecoveryRetryCount()
	if countCtx == nil || countCtx.IntegerLiteral() == nil {
		c.reportInvalidTail(action, "Expected retry count after 'RETRY'", "Provide an integer retry count, e.g. ON ERROR RETRY 3.")
		valid = false
	} else {
		plan.CountNode = countCtx.(antlr.ParserRuleContext)
		count, err := strconv.Atoi(countCtx.IntegerLiteral().GetText())
		if err != nil {
			c.reportInvalidTail(countCtx, "Expected retry count after 'RETRY'", "Provide an integer retry count, e.g. ON ERROR RETRY 3.")
			valid = false
		} else {
			plan.Count = count
		}
	}

	delayClause := action.RecoveryRetryDelayClause()
	if delayClause != nil {
		switch {
		case delayClause.DelayKeyword() != nil:
			if delayValue := delayClause.RecoveryRetryDelayValue(); delayValue != nil {
				plan.HasDelay = true
				plan.Delay = delayValue
				plan.DelayNode = delayValue.(antlr.ParserRuleContext)
			} else {
				c.reportInvalidTail(delayClause.DelayKeyword(), "Expected value after 'DELAY' in retry policy", "Provide a duration or duration-like value, e.g. DELAY 100ms.")
				valid = false
			}

			backoffClause := delayClause.RecoveryRetryBackoffClause()
			if backoffClause != nil {
				backoff, ok := c.resolveRetryBackoff(backoffClause)
				if !ok {
					valid = false
				} else {
					plan.Backoff = backoff
					plan.BackoffNode = backoffClause.(antlr.ParserRuleContext)
				}
			}
		case delayClause.RecoveryRetryBackoffClause() != nil:
			c.reportInvalidTail(delayClause.RecoveryRetryBackoffClause(), "BACKOFF requires DELAY in retry policy", "Add DELAY before BACKOFF, e.g. ON ERROR RETRY 3 DELAY 100ms BACKOFF EXPONENTIAL.")
			valid = false

			if _, ok := c.resolveRetryBackoff(delayClause.RecoveryRetryBackoffClause()); !ok {
				valid = false
			}
		}
	}

	orClauses := action.AllRecoveryRetryOrClause()
	if len(orClauses) > 1 {
		for i := 1; i < len(orClauses); i++ {
			c.reportInvalidTail(orClauses[i], "Duplicate OR fallback in retry policy", "Each retry policy may define OR at most once.")
		}
		valid = false
	}

	if len(orClauses) > 0 {
		plan.HasOr = true

		finalKind, finalExpr, finalNode, ok := c.resolveRetryFinalAction(orClauses[0])
		if !ok {
			valid = false
		} else {
			plan.FinalActionKind = finalKind
			plan.FinalExpr = finalExpr
			plan.FinalActionNode = finalNode
		}
	}

	if !valid {
		return nil, false
	}

	return plan, true
}

func (c *RecoveryCompiler) resolveRetryFinalAction(clause fql.IRecoveryRetryOrClauseContext) (core.RecoveryActionKind, fql.IExpressionContext, antlr.ParserRuleContext, bool) {
	if clause == nil {
		return core.RecoveryActionUnknown, nil, nil, false
	}

	action := clause.RecoveryRetryFinalAction()
	if action == nil {
		c.reportInvalidTail(clause, "Expected FAIL or RETURN after 'OR' in retry fallback", "Complete the retry fallback as OR FAIL or OR RETURN <expr>.")
		return core.RecoveryActionUnknown, nil, nil, false
	}

	switch {
	case action.FailKeyword() != nil:
		return core.RecoveryActionFail, nil, action.(antlr.ParserRuleContext), true
	case action.ReturnKeyword() != nil:
		returnExpr := action.RecoveryReturnExpr()
		if returnExpr == nil || returnExpr.Expression() == nil {
			c.reportInvalidTail(action, "Expected expression after 'RETURN' in recovery tail", "Provide a fallback expression, e.g. OR RETURN NONE.")
			return core.RecoveryActionUnknown, nil, nil, false
		}

		return core.RecoveryActionReturn, returnExpr.Expression(), action.(antlr.ParserRuleContext), true
	default:
		c.reportInvalidTail(action, "Expected FAIL or RETURN after 'OR' in retry fallback", "Complete the retry fallback as OR FAIL or OR RETURN <expr>.")
		return core.RecoveryActionUnknown, nil, nil, false
	}
}

func (c *RecoveryCompiler) reportInvalidTail(node antlr.ParserRuleContext, message, hint string) {
	if c == nil || c.ctx == nil || c.ctx.Errors == nil || node == nil {
		return
	}

	err := c.ctx.Errors.Create(parserd.SyntaxError, node, message)
	err.Hint = hint
	c.ctx.Errors.Add(err)
}
