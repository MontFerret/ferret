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
