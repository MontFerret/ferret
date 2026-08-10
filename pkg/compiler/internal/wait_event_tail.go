package internal

import "github.com/MontFerret/ferret/v2/pkg/parser/fql"

func waitForEventTriggerClause(ctx fql.IWaitForEventExpressionContext) fql.IWaitForTriggerClauseContext {
	if ctx == nil {
		return nil
	}

	tail := ctx.WaitForEventTail()
	if tail == nil {
		return nil
	}

	return tail.WaitForTriggerClause()
}

func waitForEventTimeoutClause(ctx fql.IWaitForEventExpressionContext) fql.ITimeoutClauseContext {
	if ctx == nil {
		return nil
	}

	tail := ctx.WaitForEventTail()
	if tail == nil {
		return nil
	}

	return tail.TimeoutClause()
}

func waitForEventGroupTriggerClause(ctx fql.IWaitForEventGroupExpressionContext) fql.IWaitForTriggerClauseContext {
	if ctx == nil || ctx.WaitForEventTail() == nil {
		return nil
	}

	return ctx.WaitForEventTail().WaitForTriggerClause()
}

func waitForEventGroupTimeoutClause(ctx fql.IWaitForEventGroupExpressionContext) fql.ITimeoutClauseContext {
	if ctx == nil || ctx.WaitForEventTail() == nil {
		return nil
	}

	return ctx.WaitForEventTail().TimeoutClause()
}
