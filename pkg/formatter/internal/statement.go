package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type statementFormatter struct {
	*engine
}

func (f *statementFormatter) formatBodyStatement(ctx *fql.BodyStatementContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.VariableDeclaration() != nil:
		f.formatVariableDeclaration(ctx.VariableDeclaration().(*fql.VariableDeclarationContext))
	case ctx.AssignmentStatement() != nil:
		f.formatAssignmentStatement(ctx.AssignmentStatement().(*fql.AssignmentStatementContext))
	case ctx.DeleteStatement() != nil:
		f.formatDeleteStatement(ctx.DeleteStatement().(*fql.DeleteStatementContext))
	case ctx.FunctionDeclaration() != nil:
		f.formatFunctionDeclaration(ctx.FunctionDeclaration().(*fql.FunctionDeclarationContext))
	case ctx.FunctionCallExpression() != nil:
		f.expression.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.WaitForExpression() != nil:
		f.formatWaitForExpression(ctx.WaitForExpression().(*fql.WaitForExpressionContext))
	case ctx.DispatchExpression() != nil:
		f.formatDispatchExpression(ctx.DispatchExpression().(*fql.DispatchExpressionContext))
	case ctx.ForExpression() != nil:
		f.formatForExpression(ctx.ForExpression().(*fql.ForExpressionContext))
	}
}

func (f *statementFormatter) formatBodyExpression(ctx *fql.BodyExpressionContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.ReturnExpression() != nil:
		f.formatReturnExpression(ctx.ReturnExpression().(*fql.ReturnExpressionContext))
	}
}

func (f *statementFormatter) formatVariableDeclaration(ctx *fql.VariableDeclarationContext) {
	if ctx == nil {
		return
	}

	if ctx.Var() != nil {
		f.writeKeyword(keywordVar)
	} else {
		f.writeKeyword(keywordLet)
	}

	f.p.space()

	if pattern := ctx.StructuredBindingPattern(); pattern != nil {
		f.bindings.formatStructured(pattern.(*fql.StructuredBindingPatternContext))
	} else if id := ctx.Identifier(); id != nil {
		f.p.write(id.GetText())
	} else if id := ctx.IgnoreIdentifier(); id != nil {
		f.p.write(id.GetText())
	} else if id := ctx.SafeReservedWord(); id != nil {
		f.p.write(id.GetText())
	} else if id := ctx.BindingIdentifier(); id != nil {
		f.p.write(id.GetText())
	}

	f.p.space()
	f.p.write("=")
	f.p.space()

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (f *statementFormatter) formatAssignmentStatement(ctx *fql.AssignmentStatementContext) {
	if ctx == nil {
		return
	}

	if target := ctx.AssignmentTarget(); target != nil {
		if target.BindingIdentifier() != nil {
			f.p.write(target.BindingIdentifier().GetText())
		}

		for _, path := range target.AllAssignmentTargetPath() {
			f.member.formatAssignmentTargetPath(path.(*fql.AssignmentTargetPathContext))
		}
	}

	f.p.space()

	if op := ctx.AssignmentOperator(); op != nil {
		f.p.write(op.GetText())
	} else {
		f.p.write("=")
	}

	f.p.space()

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (f *statementFormatter) formatDeleteStatement(ctx *fql.DeleteStatementContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordDelete)
	f.p.space()

	if target := ctx.AssignmentTarget(); target != nil {
		if target.BindingIdentifier() != nil {
			f.p.write(target.BindingIdentifier().GetText())
		}

		for _, path := range target.AllAssignmentTargetPath() {
			f.member.formatAssignmentTargetPath(path.(*fql.AssignmentTargetPathContext))
		}
	}
}

func (f *statementFormatter) formatReturnExpression(ctx *fql.ReturnExpressionContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordReturn)

	if ctx.Distinct() != nil {
		f.p.space()
		f.writeKeyword(keywordDistinct)
	}

	f.p.space()

	if value := ctx.ReturnValue(); value != nil {
		f.formatReturnValue(value.(*fql.ReturnValueContext))
	}
}

func (f *statementFormatter) formatFunctionDeclaration(ctx *fql.FunctionDeclarationContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordFunc)
	f.p.space()

	if name := ctx.FunctionName(); name != nil {
		f.p.write(name.GetText())
	}

	f.p.write("(")

	if params := ctx.FunctionParameterList(); params != nil {
		f.formatFunctionParameterList(params.(*fql.FunctionParameterListContext))
	}

	f.p.write(")")

	body := ctx.FunctionBody()
	if body == nil {
		return
	}

	funcBody := body.(*fql.FunctionBodyContext)
	if arrow := funcBody.FunctionArrow(); arrow != nil {
		f.p.space()
		f.p.write("=>")
		f.p.space()

		if expr := arrow.Expression(); expr != nil {
			f.expression.formatExpression(expr.(*fql.ExpressionContext))
		}

		return
	}

	block := funcBody.FunctionBlock()
	if block == nil {
		return
	}

	stmts := block.AllFunctionStatement()
	ret := block.FunctionReturn()

	f.p.space()
	f.p.write("{")

	headerStop := f.functionHeaderStopIndex(ctx)

	var first antlr.ParserRuleContext
	if len(stmts) > 0 {
		first = stmts[0].(antlr.ParserRuleContext)
	} else if ret != nil {
		first = ret.(antlr.ParserRuleContext)
	}

	f.p.withIndent(func() {
		if first != nil {
			f.trivia.emitListTriviaWith(
				f.p,
				f.trivia.blockLeadingTrivia(headerStop, block.OpenBrace(), f.trivia.startIndex(first)),
			)
		} else if block.CloseBrace() != nil {
			f.trivia.emitListTriviaWith(
				f.p,
				f.trivia.blockLeadingTrivia(headerStop, block.OpenBrace(), f.trivia.tokenStart(block.CloseBrace())),
			)
		} else {
			f.p.newline()
		}

		for i, stmt := range stmts {
			f.formatFunctionStatement(stmt.(*fql.FunctionStatementContext))

			if i < len(stmts)-1 {
				f.trivia.emitBetween(stmt.(antlr.ParserRuleContext), stmts[i+1].(antlr.ParserRuleContext))
			}
		}

		if ret != nil {
			if len(stmts) > 0 {
				f.trivia.emitBetween(stmts[len(stmts)-1].(antlr.ParserRuleContext), ret.(antlr.ParserRuleContext))
			}

			f.formatFunctionReturn(ret.(*fql.FunctionReturnContext))
		}

		var last antlr.ParserRuleContext
		switch {
		case ret != nil:
			last = ret.(antlr.ParserRuleContext)
		case len(stmts) > 0:
			last = stmts[len(stmts)-1].(antlr.ParserRuleContext)
		}

		if last != nil && block.CloseBrace() != nil {
			f.trivia.emitBetweenIndices(f.trivia.stopIndex(last)+1, f.trivia.tokenStart(block.CloseBrace()))
		}
	})

	if !f.p.atLineStart {
		f.p.newline()
	}

	f.p.write("}")
}

func (f *statementFormatter) formatFunctionParameterList(ctx *fql.FunctionParameterListContext) {
	if ctx == nil {
		return
	}

	params := ctx.AllFunctionParameter()
	for i, param := range params {
		pctx, ok := param.(*fql.FunctionParameterContext)
		if !ok || pctx == nil {
			continue
		}

		if id := pctx.BindingIdentifier(); id != nil {
			f.p.write(id.GetText())
		}

		if i < len(params)-1 {
			f.p.write(",")
			f.p.space()
		}
	}
}

func (f *statementFormatter) formatFunctionStatement(ctx *fql.FunctionStatementContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.VariableDeclaration() != nil:
		f.formatVariableDeclaration(ctx.VariableDeclaration().(*fql.VariableDeclarationContext))
	case ctx.AssignmentStatement() != nil:
		f.formatAssignmentStatement(ctx.AssignmentStatement().(*fql.AssignmentStatementContext))
	case ctx.DeleteStatement() != nil:
		f.formatDeleteStatement(ctx.DeleteStatement().(*fql.DeleteStatementContext))
	case ctx.FunctionDeclaration() != nil:
		f.formatFunctionDeclaration(ctx.FunctionDeclaration().(*fql.FunctionDeclarationContext))
	case ctx.FunctionCallExpression() != nil:
		f.expression.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.WaitForExpression() != nil:
		f.formatWaitForExpression(ctx.WaitForExpression().(*fql.WaitForExpressionContext))
	case ctx.DispatchExpression() != nil:
		f.formatDispatchExpression(ctx.DispatchExpression().(*fql.DispatchExpressionContext))
	case ctx.ForExpression() != nil:
		f.formatForExpression(ctx.ForExpression().(*fql.ForExpressionContext))
	case ctx.ExpressionStatement() != nil:
		f.formatExpressionStatement(ctx.ExpressionStatement().(*fql.ExpressionStatementContext))
	}
}

func (f *statementFormatter) formatFunctionReturn(ctx *fql.FunctionReturnContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordReturn)

	if ctx.Distinct() != nil {
		f.p.space()
		f.writeKeyword(keywordDistinct)
	}

	f.p.space()

	if value := ctx.ReturnValue(); value != nil {
		f.formatReturnValue(value.(*fql.ReturnValueContext))
	}
}

func (f *statementFormatter) formatReturnValue(ctx *fql.ReturnValueContext) {
	if ctx == nil {
		return
	}

	if loop := ctx.ForExpression(); loop != nil {
		f.formatForExpression(loop.(*fql.ForExpressionContext))

		return
	}

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (f *statementFormatter) formatExpressionStatement(ctx *fql.ExpressionStatementContext) {
	if ctx == nil {
		return
	}

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (f *statementFormatter) functionHeaderStopIndex(ctx *fql.FunctionDeclarationContext) int {
	if ctx == nil {
		return 0
	}

	if tok := ctx.GetToken(fql.FqlParserCloseParen, 0); tok != nil {
		if sym := tok.GetSymbol(); sym != nil {
			return sym.GetStop()
		}
	}

	if params := ctx.FunctionParameterList(); params != nil {
		return f.trivia.stopIndex(params.(antlr.ParserRuleContext))
	}

	if name := ctx.FunctionName(); name != nil {
		return f.trivia.stopIndex(name.(antlr.ParserRuleContext))
	}

	return f.trivia.stopIndex(ctx)
}

func (f *statementFormatter) formatForExpression(ctx *fql.ForExpressionContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordFor)

	if pattern := ctx.GetValuePattern(); pattern != nil {
		f.p.space()
		f.bindings.formatStructured(pattern.(*fql.StructuredBindingPatternContext))
	} else {
		writeValueVariable := true
		if tok := ctx.GetValueVariable(); ctx.In() == nil && tok != nil && tok.GetText() == "_" {
			writeValueVariable = false
		}

		if tok := ctx.GetValueVariable(); tok != nil && writeValueVariable {
			f.p.space()
			f.p.write(tok.GetText())
		}
	}

	if tok := ctx.GetCounterVariable(); tok != nil {
		f.p.write(",")
		f.p.space()
		f.p.write(tok.GetText())
	}

	switch {
	case ctx.In() != nil:
		f.p.space()
		f.writeKeyword(keywordIn)
		f.p.space()
		f.formatForExpressionSource(ctx.ForExpressionSource().(*fql.ForExpressionSourceContext))
	case ctx.While() != nil:
		if ctx.Do() != nil {
			f.p.space()
			f.writeKeyword(keywordDo)
		}

		f.p.space()
		f.writeKeyword(keywordWhile)
		f.p.space()
		f.expression.formatExpression(ctx.Expression().(*fql.ExpressionContext))
	}

	bodies := ctx.AllForExpressionBody()
	ret := ctx.ForExpressionReturn()
	braced := ctx.OpenBrace() != nil

	if len(bodies) == 0 && ret == nil {
		return
	}

	headerStop := f.forHeaderStopIndex(ctx)
	if braced {
		f.p.space()
		f.p.write("{")
	}

	var first antlr.ParserRuleContext

	if len(bodies) > 0 {
		first = bodies[0].(antlr.ParserRuleContext)
	} else if ret != nil {
		first = ret.(antlr.ParserRuleContext)
	}

	f.p.withIndent(func() {
		if first != nil {
			f.trivia.emitListTriviaWith(
				f.p,
				f.trivia.blockLeadingTrivia(headerStop, ctx.OpenBrace(), f.trivia.startIndex(first)),
			)
		} else {
			f.p.newline()
		}

		for i, body := range bodies {
			f.formatForExpressionBody(body.(*fql.ForExpressionBodyContext))

			if i < len(bodies)-1 {
				f.trivia.emitBetween(body.(antlr.ParserRuleContext), bodies[i+1].(antlr.ParserRuleContext))
			}
		}

		if ret != nil {
			if len(bodies) > 0 {
				f.trivia.emitBetween(bodies[len(bodies)-1].(antlr.ParserRuleContext), ret.(antlr.ParserRuleContext))
			}

			f.formatForExpressionReturn(ret.(*fql.ForExpressionReturnContext))
		}

		if braced && ret != nil && ctx.CloseBrace() != nil {
			f.trivia.emitBetweenIndices(
				f.trivia.stopIndex(ret.(antlr.ParserRuleContext))+1,
				f.trivia.tokenStart(ctx.CloseBrace()),
			)
		}
	})

	if braced {
		if !f.p.atLineStart {
			f.p.newline()
		}

		f.p.write("}")
	}
}

func (f *statementFormatter) forHeaderStopIndex(ctx *fql.ForExpressionContext) int {
	if ctx == nil {
		return 0
	}

	switch {
	case ctx.In() != nil:
		if src := ctx.ForExpressionSource(); src != nil {
			return f.trivia.stopIndex(src.(antlr.ParserRuleContext))
		}
	case ctx.While() != nil:
		if expr := ctx.Expression(); expr != nil {
			return f.trivia.stopIndex(expr.(antlr.ParserRuleContext))
		}
	}

	return f.trivia.stopIndex(ctx)
}

func (f *statementFormatter) formatForExpressionSource(ctx *fql.ForExpressionSourceContext) {
	if ctx == nil {
		return
	}

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (f *statementFormatter) formatForExpressionBody(ctx *fql.ForExpressionBodyContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.ForExpressionStatement() != nil:
		stmt := ctx.ForExpressionStatement().(*fql.ForExpressionStatementContext)

		switch {
		case stmt.VariableDeclaration() != nil:
			f.formatVariableDeclaration(stmt.VariableDeclaration().(*fql.VariableDeclarationContext))
		case stmt.AssignmentStatement() != nil:
			f.formatAssignmentStatement(stmt.AssignmentStatement().(*fql.AssignmentStatementContext))
		case stmt.DeleteStatement() != nil:
			f.formatDeleteStatement(stmt.DeleteStatement().(*fql.DeleteStatementContext))
		case stmt.FunctionCallExpression() != nil:
			f.expression.formatFunctionCallExpression(stmt.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
		case stmt.WaitForExpression() != nil:
			f.formatWaitForExpression(stmt.WaitForExpression().(*fql.WaitForExpressionContext))
		case stmt.DispatchExpression() != nil:
			f.formatDispatchExpression(stmt.DispatchExpression().(*fql.DispatchExpressionContext))
		}
	case ctx.ForExpressionClause() != nil:
		clause := ctx.ForExpressionClause().(*fql.ForExpressionClauseContext)

		switch {
		case clause.FilterClause() != nil:
			f.clause.formatFilterClause(clause.FilterClause().(*fql.FilterClauseContext))
		case clause.LimitClause() != nil:
			f.clause.formatLimitClause(clause.LimitClause().(*fql.LimitClauseContext))
		case clause.SortClause() != nil:
			f.clause.formatSortClause(clause.SortClause().(*fql.SortClauseContext))
		case clause.CollectClause() != nil:
			f.clause.formatCollectClause(clause.CollectClause().(*fql.CollectClauseContext))
		}
	}
}

func (f *statementFormatter) formatForExpressionReturn(ctx *fql.ForExpressionReturnContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.ReturnExpression() != nil:
		f.formatReturnExpression(ctx.ReturnExpression().(*fql.ReturnExpressionContext))
	case ctx.ForExpression() != nil:
		f.formatForExpression(ctx.ForExpression().(*fql.ForExpressionContext))
	}
}

func (f *statementFormatter) formatWaitForExpression(ctx *fql.WaitForExpressionContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordWaitFor)
	f.p.space()

	groupedStop := -1
	if event := ctx.WaitForEventGroupExpression(); event != nil {
		eventCtx := event.(*fql.WaitForEventGroupExpressionContext)
		f.formatWaitForEventGroupExpression(eventCtx)
		groupedStop = f.trivia.stopIndex(eventCtx)
	} else if pred := ctx.WaitForPredicateGroupExpression(); pred != nil {
		predCtx := pred.(*fql.WaitForPredicateGroupExpressionContext)
		f.formatWaitForPredicateGroupExpression(predCtx)
		groupedStop = f.trivia.stopIndex(predCtx)
	} else if event := ctx.WaitForEventExpression(); event != nil {
		f.formatWaitForEventExpression(event.(*fql.WaitForEventExpressionContext))
	} else if pred := ctx.WaitForPredicateExpression(); pred != nil {
		f.formatWaitForPredicateExpression(pred.(*fql.WaitForPredicateExpressionContext))
	}

	if groupedStop >= 0 {
		f.expression.formatRecoveryTailsMultiline(ctx.RecoveryTails(), groupedStop)
	} else {
		f.expression.formatRecoveryTails(ctx.RecoveryTails())
	}
}

func (f *statementFormatter) formatWaitForEventGroupExpression(ctx *fql.WaitForEventGroupExpressionContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordEvent)
	f.p.space()
	f.formatWaitForSynchronization(ctx.WaitForSynchronization().(*fql.WaitForSynchronizationContext))
	f.p.space()
	f.p.write("{")

	entries := ctx.AllWaitForEventGroupEntry()
	f.p.withIndent(func() {
		if len(entries) == 0 {
			f.p.newline()
			return
		}

		headerStop := ctx.WaitForSynchronization().GetStop().GetStop()
		first := entries[0].(antlr.ParserRuleContext)
		f.trivia.emitListTriviaWith(
			f.p,
			f.trivia.blockLeadingTrivia(headerStop, ctx.OpenBrace(), f.trivia.startIndex(first)),
		)

		for idx, entry := range entries {
			f.formatWaitForEventGroupEntry(entry.(*fql.WaitForEventGroupEntryContext))
			if idx < len(entries)-1 {
				f.trivia.emitBetween(entry.(antlr.ParserRuleContext), entries[idx+1].(antlr.ParserRuleContext))
			}
		}

		f.trivia.emitBetweenIndices(
			f.trivia.stopIndex(entries[len(entries)-1].(antlr.ParserRuleContext))+1,
			f.trivia.tokenStart(ctx.CloseBrace()),
		)
	})

	if !f.p.atLineStart {
		f.p.newline()
	}
	f.p.write("}")
	previousStop := f.trivia.tokenStop(ctx.CloseBrace())

	if tail := ctx.WaitForEventTail(); tail != nil {
		if trigger := tail.WaitForTriggerClause(); trigger != nil {
			triggerCtx := trigger.(*fql.WaitForTriggerClauseContext)
			f.trivia.emitClauseBoundary(previousStop+1, f.trivia.startIndex(triggerCtx))
			f.formatWaitForTriggerClause(triggerCtx)
			previousStop = f.trivia.stopIndex(triggerCtx)
		}

		if timeout := tail.TimeoutClause(); timeout != nil {
			timeoutCtx := timeout.(*fql.TimeoutClauseContext)
			f.trivia.emitClauseBoundary(previousStop+1, f.trivia.startIndex(timeoutCtx))
			f.clause.formatTimeoutClause(timeoutCtx)
		}
	}
}

func (f *statementFormatter) formatWaitForEventGroupEntry(ctx *fql.WaitForEventGroupEntryContext) {
	if ctx == nil {
		return
	}

	f.formatWaitForEventName(ctx.WaitForEventName().(*fql.WaitForEventNameContext))
	f.p.space()
	f.writeKeyword(keywordIn)
	f.p.space()
	f.formatWaitForEventSource(ctx.WaitForEventSource().(*fql.WaitForEventSourceContext))
	if options := ctx.OptionsClause(); options != nil {
		f.p.space()
		f.clause.formatOptionsClause(options.(*fql.OptionsClauseContext))
	}

	previous := antlr.ParserRuleContext(ctx.WaitForEventSource().(antlr.ParserRuleContext))
	if options := ctx.OptionsClause(); options != nil {
		previous = options.(antlr.ParserRuleContext)
	}
	for _, filter := range ctx.AllEventFilterClause() {
		filterCtx := filter.(*fql.EventFilterClauseContext)
		f.p.withIndent(func() {
			f.trivia.emitBetween(previous, filterCtx)
			f.clause.formatEventFilterClause(filterCtx)
		})
		previous = filterCtx
	}
}

func (f *statementFormatter) formatWaitForPredicateGroupExpression(ctx *fql.WaitForPredicateGroupExpressionContext) {
	if ctx == nil {
		return
	}

	if mode := ctx.WaitForPredicateGroupMode(); mode != nil {
		f.formatWaitForPredicateGroupMode(mode.(*fql.WaitForPredicateGroupModeContext))
		f.p.space()
	}
	f.formatWaitForSynchronization(ctx.WaitForSynchronization().(*fql.WaitForSynchronizationContext))
	f.p.space()
	f.p.write("{")

	entries := ctx.AllWaitForPredicateGroupEntry()
	f.p.withIndent(func() {
		if len(entries) == 0 {
			f.p.newline()
			return
		}

		headerStop := ctx.WaitForSynchronization().GetStop().GetStop()
		first := entries[0].(antlr.ParserRuleContext)
		f.trivia.emitListTriviaWith(
			f.p,
			f.trivia.blockLeadingTrivia(headerStop, ctx.OpenBrace(), f.trivia.startIndex(first)),
		)

		for idx, entry := range entries {
			f.formatWaitForPredicateGroupEntry(entry.(*fql.WaitForPredicateGroupEntryContext))
			if idx < len(entries)-1 {
				f.trivia.emitBetween(entry.(antlr.ParserRuleContext), entries[idx+1].(antlr.ParserRuleContext))
			}
		}

		f.trivia.emitBetweenIndices(
			f.trivia.stopIndex(entries[len(entries)-1].(antlr.ParserRuleContext))+1,
			f.trivia.tokenStart(ctx.CloseBrace()),
		)
	})

	if !f.p.atLineStart {
		f.p.newline()
	}
	f.p.write("}")
	previousStop := f.trivia.tokenStop(ctx.CloseBrace())

	if timeout := ctx.TimeoutClause(); timeout != nil {
		timeoutCtx := timeout.(*fql.TimeoutClauseContext)
		f.trivia.emitClauseBoundary(previousStop+1, f.trivia.startIndex(timeoutCtx))
		f.clause.formatTimeoutClause(timeoutCtx)
		previousStop = f.trivia.stopIndex(timeoutCtx)
	}

	if every := ctx.EveryClause(); every != nil {
		everyCtx := every.(*fql.EveryClauseContext)
		f.trivia.emitClauseBoundary(previousStop+1, f.trivia.startIndex(everyCtx))
		f.clause.formatEveryClause(everyCtx)
		previousStop = f.trivia.stopIndex(everyCtx)
	}

	if backoff := ctx.BackoffClause(); backoff != nil {
		backoffCtx := backoff.(*fql.BackoffClauseContext)
		f.trivia.emitClauseBoundary(previousStop+1, f.trivia.startIndex(backoffCtx))
		f.clause.formatBackoffClause(backoffCtx)
		previousStop = f.trivia.stopIndex(backoffCtx)
	}

	if jitter := ctx.JitterClause(); jitter != nil {
		jitterCtx := jitter.(*fql.JitterClauseContext)
		f.trivia.emitClauseBoundary(previousStop+1, f.trivia.startIndex(jitterCtx))
		f.clause.formatJitterClause(jitterCtx)
	}
}

func (f *statementFormatter) formatWaitForPredicateGroupEntry(ctx *fql.WaitForPredicateGroupEntryContext) {
	if ctx == nil || ctx.Expression() == nil {
		return
	}

	expression := ctx.Expression().(*fql.ExpressionContext)
	f.expression.formatExpression(expression)
	previous := antlr.ParserRuleContext(expression)
	for _, when := range ctx.AllWaitForPredicateWhenClause() {
		whenCtx := when.(*fql.WaitForPredicateWhenClauseContext)
		f.p.withIndent(func() {
			f.trivia.emitBetween(previous, whenCtx)
			f.clause.formatWaitForPredicateWhenClause(whenCtx)
		})
		previous = whenCtx
	}
}

func (f *statementFormatter) formatWaitForPredicateGroupMode(ctx *fql.WaitForPredicateGroupModeContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.Not() != nil:
		f.writeKeyword(keywordNot)
		f.p.space()
		f.writeKeyword(keywordExists)
	case ctx.Exists() != nil:
		f.writeKeyword(keywordExists)
	case ctx.Value() != nil:
		f.writeKeyword(keywordValue)
	}
}

func (f *statementFormatter) formatWaitForSynchronization(ctx *fql.WaitForSynchronizationContext) {
	if ctx == nil {
		return
	}

	if ctx.All() != nil {
		f.writeKeyword(keywordAll)
		return
	}

	f.writeKeyword(keywordAny)
}

func (f *statementFormatter) formatWaitForEventExpression(ctx *fql.WaitForEventExpressionContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordEvent)
	f.p.space()

	if name := ctx.WaitForEventName(); name != nil {
		f.formatWaitForEventName(name.(*fql.WaitForEventNameContext))
	}

	f.p.space()
	f.writeKeyword(keywordIn)
	f.p.space()

	if src := ctx.WaitForEventSource(); src != nil {
		f.formatWaitForEventSource(src.(*fql.WaitForEventSourceContext))
	}

	if opt := ctx.OptionsClause(); opt != nil {
		f.p.space()
		f.clause.formatOptionsClause(opt.(*fql.OptionsClauseContext))
	}

	for _, filter := range ctx.AllEventFilterClause() {
		f.p.space()
		f.clause.formatEventFilterClause(filter.(*fql.EventFilterClauseContext))
	}

	if tail := ctx.WaitForEventTail(); tail != nil {
		if trigger := tail.WaitForTriggerClause(); trigger != nil {
			f.p.space()
			f.formatWaitForTriggerClause(trigger.(*fql.WaitForTriggerClauseContext))
		}

		if timeout := tail.TimeoutClause(); timeout != nil {
			f.p.space()
			f.clause.formatTimeoutClause(timeout.(*fql.TimeoutClauseContext))
		}
	}
}

func (f *statementFormatter) formatWaitForTriggerClause(ctx *fql.WaitForTriggerClauseContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordTrigger)
	f.p.space()

	if inline := ctx.WaitForTriggerInlineStatement(); inline != nil {
		f.formatWaitForTriggerInlineStatement(inline.(*fql.WaitForTriggerInlineStatementContext))
		return
	}

	f.p.write("(")

	stmts := ctx.AllWaitForTriggerStatement()
	if len(stmts) == 0 {
		f.p.write(")")

		return
	}

	start := f.trivia.stopIndex(ctx) + 1
	if openParen := ctx.OpenParen(); openParen != nil {
		if sym := openParen.GetSymbol(); sym != nil {
			start = sym.GetStop() + 1
		}
	}

	first := stmts[0].(antlr.ParserRuleContext)

	f.p.withIndent(func() {
		f.trivia.emitBetweenIndices(start, f.trivia.startIndex(first))

		for i, stmt := range stmts {
			f.formatWaitForTriggerStatement(stmt.(*fql.WaitForTriggerStatementContext))

			if i < len(stmts)-1 {
				f.trivia.emitBetween(stmt.(antlr.ParserRuleContext), stmts[i+1].(antlr.ParserRuleContext))
			}
		}
	})

	if !f.p.atLineStart {
		f.p.newline()
	}

	f.p.write(")")
}

func (f *statementFormatter) formatWaitForTriggerStatement(ctx *fql.WaitForTriggerStatementContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.VariableDeclaration() != nil:
		f.formatVariableDeclaration(ctx.VariableDeclaration().(*fql.VariableDeclarationContext))
	case ctx.AssignmentStatement() != nil:
		f.formatAssignmentStatement(ctx.AssignmentStatement().(*fql.AssignmentStatementContext))
	case ctx.DeleteStatement() != nil:
		f.formatDeleteStatement(ctx.DeleteStatement().(*fql.DeleteStatementContext))
	case ctx.FunctionCallExpression() != nil:
		f.expression.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.WaitForExpression() != nil:
		f.formatWaitForExpression(ctx.WaitForExpression().(*fql.WaitForExpressionContext))
	case ctx.DispatchExpression() != nil:
		f.formatDispatchExpression(ctx.DispatchExpression().(*fql.DispatchExpressionContext))
	}
}

func (f *statementFormatter) formatWaitForTriggerInlineStatement(ctx *fql.WaitForTriggerInlineStatementContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.VariableDeclaration() != nil:
		f.formatVariableDeclaration(ctx.VariableDeclaration().(*fql.VariableDeclarationContext))
	case ctx.AssignmentStatement() != nil:
		f.formatAssignmentStatement(ctx.AssignmentStatement().(*fql.AssignmentStatementContext))
	case ctx.DeleteStatement() != nil:
		f.formatDeleteStatement(ctx.DeleteStatement().(*fql.DeleteStatementContext))
	case ctx.FunctionCallNoRecoveryExpression() != nil:
		f.expression.formatFunctionCallNoRecoveryExpression(
			ctx.FunctionCallNoRecoveryExpression().(*fql.FunctionCallNoRecoveryExpressionContext),
		)
	case ctx.WaitForTriggerInlineDispatchStatement() != nil:
		f.formatWaitForTriggerInlineDispatchStatement(
			ctx.WaitForTriggerInlineDispatchStatement().(*fql.WaitForTriggerInlineDispatchStatementContext),
		)
	}
}

func (f *statementFormatter) formatWaitForTriggerInlineDispatchStatement(ctx *fql.WaitForTriggerInlineDispatchStatementContext) {
	if ctx == nil {
		return
	}

	if ctx.Dispatch() != nil {
		f.writeKeyword(keywordDispatch)
		f.p.space()

		if name := ctx.DispatchEventName(); name != nil {
			f.formatDispatchEventName(name.(*fql.DispatchEventNameContext))
		}

		f.p.space()
		f.writeKeyword(keywordIn)
		f.p.space()

		if tgt := ctx.DispatchTarget(); tgt != nil {
			f.formatDispatchTarget(tgt.(*fql.DispatchTargetContext))
		}

		if with := ctx.DispatchWithClause(); with != nil {
			f.p.space()
			f.formatDispatchWithClause(with.(*fql.DispatchWithClauseContext))
		}

		if opt := ctx.DispatchOptionsClause(); opt != nil {
			f.p.space()
			f.formatDispatchOptionsClause(opt.(*fql.DispatchOptionsClauseContext))
		}

		return
	}

	if tgt := ctx.DispatchTarget(); tgt != nil {
		f.formatDispatchTarget(tgt.(*fql.DispatchTargetContext))
	}

	f.p.space()
	f.p.write("<-")
	f.p.space()

	if name := ctx.DispatchEventName(); name != nil {
		f.formatDispatchEventName(name.(*fql.DispatchEventNameContext))
	}
}

func (f *statementFormatter) formatWaitForPredicateExpression(ctx *fql.WaitForPredicateExpressionContext) {
	if ctx == nil {
		return
	}

	if pred := ctx.WaitForPredicate(); pred != nil {
		f.formatWaitForPredicate(pred.(*fql.WaitForPredicateContext))
	}

	for _, when := range ctx.AllWaitForPredicateWhenClause() {
		f.p.space()
		f.clause.formatWaitForPredicateWhenClause(when.(*fql.WaitForPredicateWhenClauseContext))
	}

	if timeout := ctx.TimeoutClause(); timeout != nil {
		f.p.space()
		f.clause.formatTimeoutClause(timeout.(*fql.TimeoutClauseContext))
	}

	if every := ctx.EveryClause(); every != nil {
		f.p.space()
		f.clause.formatEveryClause(every.(*fql.EveryClauseContext))
	}

	if backoff := ctx.BackoffClause(); backoff != nil {
		f.p.space()
		f.clause.formatBackoffClause(backoff.(*fql.BackoffClauseContext))
	}

	if jitter := ctx.JitterClause(); jitter != nil {
		f.p.space()
		f.clause.formatJitterClause(jitter.(*fql.JitterClauseContext))
	}
}

func (f *statementFormatter) formatWaitForPredicate(ctx *fql.WaitForPredicateContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.Not() != nil && ctx.Exists() != nil:
		f.writeKeyword(keywordNot)
		f.p.space()
		f.writeKeyword(keywordExists)
		f.p.space()
		f.expression.formatExpression(ctx.Expression().(*fql.ExpressionContext))
	case ctx.Exists() != nil:
		f.writeKeyword(keywordExists)
		f.p.space()
		f.expression.formatExpression(ctx.Expression().(*fql.ExpressionContext))
	case ctx.Value() != nil:
		f.writeKeyword(keywordValue)
		f.p.space()
		f.expression.formatExpression(ctx.Expression().(*fql.ExpressionContext))
	case ctx.Expression() != nil:
		f.expression.formatExpression(ctx.Expression().(*fql.ExpressionContext))
	}
}

func (f *statementFormatter) formatDispatchExpression(ctx *fql.DispatchExpressionContext) {
	if ctx == nil {
		return
	}

	if ctx.Dispatch() != nil {
		f.writeKeyword(keywordDispatch)
		f.p.space()

		if name := ctx.DispatchEventName(); name != nil {
			f.formatDispatchEventName(name.(*fql.DispatchEventNameContext))
		}

		f.p.space()
		f.writeKeyword(keywordIn)
		f.p.space()

		if tgt := ctx.DispatchTarget(); tgt != nil {
			f.formatDispatchTarget(tgt.(*fql.DispatchTargetContext))
		}

		if with := ctx.DispatchWithClause(); with != nil {
			f.p.space()
			f.formatDispatchWithClause(with.(*fql.DispatchWithClauseContext))
		}

		if opt := ctx.DispatchOptionsClause(); opt != nil {
			f.p.space()
			f.formatDispatchOptionsClause(opt.(*fql.DispatchOptionsClauseContext))
		}

		f.expression.formatRecoveryTails(ctx.RecoveryTails())

		return
	}

	if tgt := ctx.DispatchTarget(); tgt != nil {
		f.formatDispatchTarget(tgt.(*fql.DispatchTargetContext))
	}

	f.p.space()
	f.p.write("<-")
	f.p.space()

	if name := ctx.DispatchEventName(); name != nil {
		f.formatDispatchEventName(name.(*fql.DispatchEventNameContext))
	}

	f.expression.formatRecoveryTails(ctx.RecoveryTails())
}

func (f *statementFormatter) formatDispatchEventName(ctx *fql.DispatchEventNameContext) {
	if ctx == nil {
		return
	}

	f.values.formatStringOrRef(
		ctx.StringLiteral(),
		ctx.Variable(),
		ctx.Param(),
		ctx.MemberExpression(),
		ctx.FunctionCall(),
	)
}

func (f *statementFormatter) formatDispatchTarget(ctx *fql.DispatchTargetContext) {
	if ctx == nil {
		return
	}

	if expr := ctx.Expression(); expr != nil {
		f.p.write("(")
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
		f.p.write(")")

		return
	}

	f.values.formatRefValueWithCallExpr(
		ctx.FunctionCallExpression(),
		ctx.Variable(),
		ctx.Param(),
		ctx.MemberExpression(),
	)
}

func (f *statementFormatter) formatDispatchWithClause(ctx *fql.DispatchWithClauseContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordWith)
	f.p.space()

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (f *statementFormatter) formatDispatchOptionsClause(ctx *fql.DispatchOptionsClauseContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordOptions)
	f.p.space()

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (f *statementFormatter) formatWaitForEventName(ctx *fql.WaitForEventNameContext) {
	if ctx == nil {
		return
	}

	f.values.formatStringOrRef(
		ctx.StringLiteral(),
		ctx.Variable(),
		ctx.Param(),
		ctx.MemberExpression(),
		ctx.FunctionCall(),
	)
}

func (f *statementFormatter) formatWaitForEventSource(ctx *fql.WaitForEventSourceContext) {
	if ctx == nil {
		return
	}

	f.values.formatRefValueWithCallExpr(
		ctx.FunctionCallExpression(),
		ctx.Variable(),
		nil,
		ctx.MemberExpression(),
	)
}
