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
	case ctx.FunctionDeclaration() != nil:
		f.formatFunctionDeclaration(ctx.FunctionDeclaration().(*fql.FunctionDeclarationContext))
	case ctx.FunctionCallExpression() != nil:
		f.expression.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.WaitForExpression() != nil:
		f.formatWaitForExpression(ctx.WaitForExpression().(*fql.WaitForExpressionContext))
	case ctx.DispatchExpression() != nil:
		f.formatDispatchExpression(ctx.DispatchExpression().(*fql.DispatchExpressionContext))
	}
}

func (f *statementFormatter) formatBodyExpression(ctx *fql.BodyExpressionContext) {
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

	if id := ctx.Identifier(); id != nil {
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
		switch {
		case target.BindingIdentifier() != nil:
			f.p.write(target.BindingIdentifier().GetText())
		case target.MemberExpression() != nil:
			f.member.formatMemberExpression(target.MemberExpression().(*fql.MemberExpressionContext))
		}
	}

	f.p.space()
	f.p.write("=")
	f.p.space()

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
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

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
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

	if len(stmts) == 0 && ret == nil {
		return
	}

	f.p.space()
	f.p.write("(")

	headerStop := f.functionHeaderStopIndex(ctx)
	start := headerStop + 1
	if openParen := block.OpenParen(); openParen != nil {
		if sym := openParen.GetSymbol(); sym != nil {
			start = sym.GetStop() + 1
		}
	}

	var first antlr.ParserRuleContext
	if len(stmts) > 0 {
		first = stmts[0].(antlr.ParserRuleContext)
	} else if ret != nil {
		first = ret.(antlr.ParserRuleContext)
	}

	f.p.withIndent(func() {
		if first != nil {
			f.trivia.emitBetweenIndices(start, f.trivia.startIndex(first))
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
	})

	if !f.p.atLineStart {
		f.p.newline()
	}

	f.p.write(")")
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

		if id := pctx.Identifier(); id != nil {
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
	case ctx.FunctionDeclaration() != nil:
		f.formatFunctionDeclaration(ctx.FunctionDeclaration().(*fql.FunctionDeclarationContext))
	case ctx.FunctionCallExpression() != nil:
		f.expression.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.WaitForExpression() != nil:
		f.formatWaitForExpression(ctx.WaitForExpression().(*fql.WaitForExpressionContext))
	case ctx.DispatchExpression() != nil:
		f.formatDispatchExpression(ctx.DispatchExpression().(*fql.DispatchExpressionContext))
	case ctx.ExpressionStatement() != nil:
		f.formatExpressionStatement(ctx.ExpressionStatement().(*fql.ExpressionStatementContext))
	}
}

func (f *statementFormatter) formatFunctionReturn(ctx *fql.FunctionReturnContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordReturn)
	f.p.space()

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

	writeValueVariable := true
	if ctx.In() == nil {
		if tok := ctx.GetValueVariable(); tok != nil && tok.GetText() == "_" {
			writeValueVariable = false
		}
	}

	if tok := ctx.GetValueVariable(); tok != nil && writeValueVariable {
		f.p.space()
		f.p.write(tok.GetText())
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

	if len(bodies) == 0 && ret == nil {
		return
	}

	headerStop := f.forHeaderStopIndex(ctx)

	var first antlr.ParserRuleContext

	if len(bodies) > 0 {
		first = bodies[0].(antlr.ParserRuleContext)
	} else if ret != nil {
		first = ret.(antlr.ParserRuleContext)
	}

	f.p.withIndent(func() {
		if first != nil {
			f.trivia.emitBetweenIndices(headerStop+1, f.trivia.startIndex(first))
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
	})
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

	switch {
	case ctx.FunctionCallExpression() != nil:
		f.expression.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.ArrayLiteral() != nil:
		f.list.formatArrayLiteral(ctx.ArrayLiteral().(*fql.ArrayLiteralContext))
	case ctx.ObjectLiteral() != nil:
		f.list.formatObjectLiteral(ctx.ObjectLiteral().(*fql.ObjectLiteralContext))
	case ctx.Variable() != nil:
		f.expression.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.MemberExpression() != nil:
		f.member.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	case ctx.RangeOperator() != nil:
		f.expression.formatRangeOperator(ctx.RangeOperator().(*fql.RangeOperatorContext))
	case ctx.Param() != nil:
		f.expression.formatParam(ctx.Param().(*fql.ParamContext))
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
		case stmt.FunctionCallExpression() != nil:
			f.expression.formatFunctionCallExpression(stmt.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
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

	if event := ctx.WaitForEventExpression(); event != nil {
		f.formatWaitForEventExpression(event.(*fql.WaitForEventExpressionContext))
	} else if pred := ctx.WaitForPredicateExpression(); pred != nil {
		f.formatWaitForPredicateExpression(pred.(*fql.WaitForPredicateExpressionContext))
	}

	if ctx.WaitForOrThrowClause() != nil {
		f.p.space()
		f.writeKeyword(keywordOr)
		f.p.space()
		f.writeKeyword(keywordThrow)
	}
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

	if filter := ctx.EventFilterClause(); filter != nil {
		f.p.space()
		f.clause.formatEventFilterClause(filter.(*fql.EventFilterClauseContext))
	}

	if timeout := ctx.TimeoutClause(); timeout != nil {
		f.p.space()
		f.clause.formatTimeoutClause(timeout.(*fql.TimeoutClauseContext))
	}
}

func (f *statementFormatter) formatWaitForPredicateExpression(ctx *fql.WaitForPredicateExpressionContext) {
	if ctx == nil {
		return
	}

	if pred := ctx.WaitForPredicate(); pred != nil {
		f.formatWaitForPredicate(pred.(*fql.WaitForPredicateContext))
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
