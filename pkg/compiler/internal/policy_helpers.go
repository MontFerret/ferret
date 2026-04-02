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

func collectRecoveryPlan(ctx *CompilerContext, owner core.RecoveryTailOwner, options core.RecoveryPlanOptions) core.RecoveryPlan {
	if ctx == nil || owner == nil {
		return core.RecoveryPlan{}
	}

	tails := owner.RecoveryTails()
	if tails == nil {
		return core.RecoveryPlan{}
	}

	var plan core.RecoveryPlan

	for _, tail := range tails.AllRecoveryTail() {
		condition, ok := resolveRecoveryCondition(ctx, tail)
		if !ok {
			continue
		}

		if condition == core.RecoveryConditionTimeout {
			if !options.AllowTimeout {
				reportInvalidRecoveryTail(ctx, tail, "ON TIMEOUT is only valid on WAITFOR operations", "Use ON TIMEOUT only on WAITFOR expressions that define timeout handling.")
				continue
			}

			if !options.HasTimeout {
				reportInvalidRecoveryTail(ctx, tail, "ON TIMEOUT requires a TIMEOUT clause on WAITFOR", "Add a TIMEOUT clause before ON TIMEOUT, e.g. WAITFOR VALUE x TIMEOUT 1s ON TIMEOUT RETURN NONE.")
				continue
			}
		}

		handler, ok := resolveRecoveryHandler(ctx, tail, condition)
		if !ok {
			continue
		}

		switch condition {
		case core.RecoveryConditionError:
			if plan.OnError != nil {
				reportInvalidRecoveryTail(ctx, tail, "Duplicate ON ERROR handler", "Each operation may define ON ERROR at most once.")
				continue
			}

			plan.OnError = handler
		case core.RecoveryConditionTimeout:
			if plan.OnTimeout != nil {
				reportInvalidRecoveryTail(ctx, tail, "Duplicate ON TIMEOUT handler", "Each operation may define ON TIMEOUT at most once.")
				continue
			}

			plan.OnTimeout = handler
		}
	}

	return plan
}

func resolveRecoveryCondition(ctx *CompilerContext, tail fql.IRecoveryTailContext) (core.RecoveryCondition, bool) {
	if tail == nil {
		return core.RecoveryConditionUnknown, false
	}

	cond := tail.RecoveryCondition()
	if cond == nil {
		reportInvalidRecoveryTail(ctx, tail, "Expected ERROR or TIMEOUT after 'ON' in recovery tail", "Complete the tail as ON ERROR FAIL, ON ERROR RETURN <expr>, ON ERROR RETRY <count>, ON TIMEOUT FAIL, or ON TIMEOUT RETURN <expr>.")
		return core.RecoveryConditionUnknown, false
	}

	switch {
	case cond.TimeoutKeyword() != nil:
		return core.RecoveryConditionTimeout, true
	case cond.ErrorKeyword() != nil:
		return core.RecoveryConditionError, true
	}

	if strings.EqualFold(cond.GetText(), core.SuppressRecoveryAction) {
		reportInvalidRecoveryTail(ctx, cond.(antlr.ParserRuleContext), "SUPPRESS is not supported in recovery tails", "Use ON ERROR RETURN NONE instead.")
		return core.RecoveryConditionUnknown, false
	}

	reportInvalidRecoveryTail(ctx, cond.(antlr.ParserRuleContext), "Expected ERROR or TIMEOUT after 'ON' in recovery tail", "Complete the tail as ON ERROR FAIL, ON ERROR RETURN <expr>, ON ERROR RETRY <count>, ON TIMEOUT FAIL, or ON TIMEOUT RETURN <expr>.")

	return core.RecoveryConditionUnknown, false
}

func resolveRecoveryHandler(ctx *CompilerContext, tail fql.IRecoveryTailContext, condition core.RecoveryCondition) (*core.RecoveryHandler, bool) {
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
		return &core.RecoveryHandler{
			ActionKind: core.RecoveryActionFail,
			ActionNode: action.(antlr.ParserRuleContext),
			TailNode:   tail.(antlr.ParserRuleContext),
		}, true
	case action.ReturnKeyword() != nil:
		returnExpr := action.RecoveryReturnExpr()
		if returnExpr == nil || returnExpr.Expression() == nil {
			reportInvalidRecoveryTail(ctx, action.(antlr.ParserRuleContext), "Expected expression after 'RETURN' in recovery tail", "Provide a fallback expression, e.g. ON ERROR RETURN NONE.")
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
			reportInvalidRecoveryTail(ctx, action.RecoveryRetryAction(), "RETRY is only valid under ON ERROR", "Use ON ERROR RETRY <count> ... to retry failures, or ON TIMEOUT FAIL/RETURN for timeout handling.")
			return nil, false
		}

		retry, ok := resolveRecoveryRetryPlan(ctx, action.RecoveryRetryAction())
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
		reportInvalidRecoveryTail(ctx, action.(antlr.ParserRuleContext), "SUPPRESS is not supported in recovery tails", "Use ON ERROR RETURN NONE instead.")
		return nil, false
	}

	reportInvalidRecoveryTail(ctx, action.(antlr.ParserRuleContext), missingRecoveryActionMessage(condition), missingRecoveryActionHint(condition))

	return nil, false
}

func resolveRecoveryRetryPlan(ctx *CompilerContext, action fql.IRecoveryRetryActionContext) (*core.RecoveryRetryPlan, bool) {
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
		reportInvalidRecoveryTail(ctx, action, "Expected retry count after 'RETRY'", "Provide an integer retry count, e.g. ON ERROR RETRY 3.")
		valid = false
	} else {
		plan.CountNode = countCtx.(antlr.ParserRuleContext)
		count, err := strconv.Atoi(countCtx.IntegerLiteral().GetText())
		if err != nil {
			reportInvalidRecoveryTail(ctx, countCtx, "Expected retry count after 'RETRY'", "Provide an integer retry count, e.g. ON ERROR RETRY 3.")
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
				reportInvalidRecoveryTail(ctx, delayClause.DelayKeyword(), "Expected value after 'DELAY' in retry policy", "Provide a duration or duration-like value, e.g. DELAY 100ms.")
				valid = false
			}

			backoffClause := delayClause.RecoveryRetryBackoffClause()
			if backoffClause != nil {
				backoff, ok := resolveRetryBackoff(ctx, backoffClause)
				if !ok {
					valid = false
				} else {
					plan.Backoff = backoff
					plan.BackoffNode = backoffClause.(antlr.ParserRuleContext)
				}
			}
		case delayClause.RecoveryRetryBackoffClause() != nil:
			reportInvalidRecoveryTail(ctx, delayClause.RecoveryRetryBackoffClause(), "BACKOFF requires DELAY in retry policy", "Add DELAY before BACKOFF, e.g. ON ERROR RETRY 3 DELAY 100ms BACKOFF EXPONENTIAL.")
			valid = false

			if _, ok := resolveRetryBackoff(ctx, delayClause.RecoveryRetryBackoffClause()); !ok {
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
		plan.HasOr = true

		finalKind, finalExpr, finalNode, ok := resolveRecoveryRetryFinalAction(ctx, orClauses[0])
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

func resolveRecoveryRetryFinalAction(ctx *CompilerContext, clause fql.IRecoveryRetryOrClauseContext) (core.RecoveryActionKind, fql.IExpressionContext, antlr.ParserRuleContext, bool) {
	if clause == nil {
		return core.RecoveryActionUnknown, nil, nil, false
	}

	action := clause.RecoveryRetryFinalAction()
	if action == nil {
		reportInvalidRecoveryTail(ctx, clause, "Expected FAIL or RETURN after 'OR' in retry fallback", "Complete the retry fallback as OR FAIL or OR RETURN <expr>.")
		return core.RecoveryActionUnknown, nil, nil, false
	}

	switch {
	case action.FailKeyword() != nil:
		return core.RecoveryActionFail, nil, action.(antlr.ParserRuleContext), true
	case action.ReturnKeyword() != nil:
		returnExpr := action.RecoveryReturnExpr()
		if returnExpr == nil || returnExpr.Expression() == nil {
			reportInvalidRecoveryTail(ctx, action, "Expected expression after 'RETURN' in recovery tail", "Provide a fallback expression, e.g. OR RETURN NONE.")
			return core.RecoveryActionUnknown, nil, nil, false
		}

		return core.RecoveryActionReturn, returnExpr.Expression(), action.(antlr.ParserRuleContext), true
	default:
		reportInvalidRecoveryTail(ctx, action, "Expected FAIL or RETURN after 'OR' in retry fallback", "Complete the retry fallback as OR FAIL or OR RETURN <expr>.")
		return core.RecoveryActionUnknown, nil, nil, false
	}
}

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

func reportInvalidRecoveryTail(ctx *CompilerContext, node antlr.ParserRuleContext, message, hint string) {
	if ctx == nil || ctx.Errors == nil || node == nil {
		return
	}

	err := ctx.Errors.Create(parserd.SyntaxError, node, message)
	err.Hint = hint
	ctx.Errors.Add(err)
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

func mergeRecoveryPlans(ctx *CompilerContext, primary, extra core.RecoveryPlan) core.RecoveryPlan {
	merged := primary

	if extra.OnTimeout != nil {
		if merged.OnTimeout != nil {
			reportInvalidRecoveryTail(ctx, extra.OnTimeout.TailNode, "Duplicate ON TIMEOUT handler", "Each operation may define ON TIMEOUT at most once.")
		} else {
			merged.OnTimeout = extra.OnTimeout
		}
	}

	if extra.OnError != nil {
		if merged.OnError != nil {
			reportInvalidRecoveryTail(ctx, extra.OnError.TailNode, "Duplicate ON ERROR handler", "Each operation may define ON ERROR at most once.")
		} else {
			merged.OnError = extra.OnError
		}
	}

	return merged
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

func emitRecoveryRetryDelay(ctx *CompilerContext, retry *core.RecoveryRetryPlan, state core.RetryDelayState) {
	if ctx == nil || retry == nil || !retry.HasDelay {
		return
	}

	delayReady := ctx.Emitter.NewLabel("recovery", "retry", "delay", "ready")
	ctx.Emitter.EmitJumpIfTrue(state.ReadyReg, delayReady)

	delayValue := ensureRecoveryRegister(ctx, ctx.WaitCompiler.compileDurationClause(retry.Delay))
	emitMoveAuto(ctx, state.BaseReg, delayValue)
	emitMoveAuto(ctx, state.CurrentReg, state.BaseReg)
	ctx.Emitter.EmitBoolean(state.ReadyReg, true)
	ctx.Emitter.MarkLabel(delayReady)

	ctx.Emitter.EmitA(bytecode.OpSleep, state.CurrentReg)

	if retry.Backoff != core.RetryBackoffNone {
		ctx.WaitCompiler.emitBackoffUpdate(retry.Backoff, state.CurrentReg, state.BaseReg)
	}
}

func widenRecoveryResultType(ctx *CompilerContext, out bytecode.Operand, plan core.RecoveryPlan) bytecode.Operand {
	if ctx == nil || !out.IsRegister() || !recoveryPlanHasReturnHandler(plan) {
		return out
	}

	ctx.Types.Set(out, core.TypeAny)

	return out
}

func ensureRecoveryRegister(ctx *CompilerContext, op bytecode.Operand) bytecode.Operand {
	if ctx == nil || op == bytecode.NoopOperand || op.IsRegister() {
		return op
	}

	dst := ctx.Registers.Allocate()
	ctx.Emitter.EmitLoadConst(dst, op)
	ctx.Types.Set(dst, operandType(ctx, op))

	return dst
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
