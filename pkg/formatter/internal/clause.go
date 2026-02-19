package internal

import "github.com/MontFerret/ferret/v2/pkg/parser/fql"

type clauseFormatter struct {
	*engine
}

func (f *clauseFormatter) formatFilterClause(ctx *fql.FilterClauseContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("FILTER")
	f.p.space()

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (f *clauseFormatter) formatLimitClause(ctx *fql.LimitClauseContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("LIMIT")
	f.p.space()
	values := ctx.AllLimitClauseValue()

	for i, val := range values {
		f.formatLimitClauseValue(val.(*fql.LimitClauseValueContext))

		if i < len(values)-1 {
			f.p.write(",")
			f.p.space()
		}
	}
}

func (f *clauseFormatter) formatLimitClauseValue(ctx *fql.LimitClauseValueContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.IntegerLiteral() != nil:
		f.p.write(ctx.IntegerLiteral().GetText())
	case ctx.Param() != nil:
		f.expression.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.Variable() != nil:
		f.expression.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.FunctionCallExpression() != nil:
		f.expression.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.MemberExpression() != nil:
		f.member.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	}
}

func (f *clauseFormatter) formatSortClause(ctx *fql.SortClauseContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("SORT")
	f.p.space()
	exprs := ctx.AllSortClauseExpression()

	for i, expr := range exprs {
		f.formatSortClauseExpression(expr.(*fql.SortClauseExpressionContext))
		if i < len(exprs)-1 {
			f.p.write(",")
			f.p.space()
		}
	}
}

func (f *clauseFormatter) formatSortClauseExpression(ctx *fql.SortClauseExpressionContext) {
	if ctx == nil {
		return
	}

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
	}

	if ctx.SortDirection() != nil {
		f.p.space()
		f.p.write(applyCase(f.opts.caseMode, ctx.SortDirection().GetText()))
	}
}

func (f *clauseFormatter) formatCollectClause(ctx *fql.CollectClauseContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("COLLECT")

	if grouping := ctx.CollectGrouping(); grouping != nil {
		f.p.space()
		f.formatCollectGrouping(grouping.(*fql.CollectGroupingContext))
	}

	if aggregator := ctx.CollectAggregator(); aggregator != nil {
		f.p.space()
		f.formatCollectAggregator(aggregator.(*fql.CollectAggregatorContext))
	}

	if counter := ctx.CollectCounter(); counter != nil {
		f.p.space()
		f.formatCollectCounter(counter.(*fql.CollectCounterContext))
	}

	if projection := ctx.CollectGroupProjection(); projection != nil {
		f.p.space()
		f.formatCollectGroupProjection(projection.(*fql.CollectGroupProjectionContext))
	}
}

func (f *clauseFormatter) formatCollectGrouping(ctx *fql.CollectGroupingContext) {
	if ctx == nil {
		return
	}

	selectors := ctx.AllCollectSelector()

	for i, sel := range selectors {
		f.formatCollectSelector(sel.(*fql.CollectSelectorContext))

		if i < len(selectors)-1 {
			f.p.write(",")
			f.p.space()
		}
	}
}

func (f *clauseFormatter) formatCollectSelector(ctx *fql.CollectSelectorContext) {
	if ctx == nil {
		return
	}

	if id := ctx.Identifier(); id != nil {
		f.p.write(id.GetText())
	}

	f.p.space()
	f.p.write("=")
	f.p.space()

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (f *clauseFormatter) formatCollectAggregator(ctx *fql.CollectAggregatorContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("AGGREGATE")
	f.p.space()
	selectors := ctx.AllCollectAggregateSelector()

	for i, sel := range selectors {
		f.formatCollectAggregateSelector(sel.(*fql.CollectAggregateSelectorContext))
		if i < len(selectors)-1 {
			f.p.write(",")
			f.p.space()
		}
	}
}

func (f *clauseFormatter) formatCollectAggregateSelector(ctx *fql.CollectAggregateSelectorContext) {
	if ctx == nil {
		return
	}

	if id := ctx.Identifier(); id != nil {
		f.p.write(id.GetText())
	}

	f.p.space()
	f.p.write("=")
	f.p.space()

	if call := ctx.FunctionCallExpression(); call != nil {
		f.expression.formatFunctionCallExpression(call.(*fql.FunctionCallExpressionContext))
	}
}

func (f *clauseFormatter) formatCollectGroupProjection(ctx *fql.CollectGroupProjectionContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("INTO")
	f.p.space()

	if sel := ctx.CollectSelector(); sel != nil {
		f.formatCollectSelector(sel.(*fql.CollectSelectorContext))
		return
	}

	if id := ctx.Identifier(); id != nil {
		f.p.write(id.GetText())

		if filter := ctx.CollectGroupProjectionFilter(); filter != nil {
			f.p.space()
			f.formatCollectGroupProjectionFilter(filter.(*fql.CollectGroupProjectionFilterContext))
		}
	}
}

func (f *clauseFormatter) formatCollectGroupProjectionFilter(ctx *fql.CollectGroupProjectionFilterContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("KEEP")
	f.p.space()
	ids := ctx.AllIdentifier()

	for i, id := range ids {
		f.p.write(id.GetText())

		if i < len(ids)-1 {
			f.p.write(",")
			f.p.space()
		}
	}
}

func (f *clauseFormatter) formatCollectCounter(ctx *fql.CollectCounterContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("WITH")
	f.p.space()

	if id := ctx.Identifier(0); id != nil {
		f.p.write(id.GetText())
	}

	f.p.space()
	f.writeKeyword("INTO")
	f.p.space()

	if id := ctx.Identifier(1); id != nil {
		f.p.write(id.GetText())
	}
}

func (f *clauseFormatter) formatOptionsClause(ctx *fql.OptionsClauseContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("OPTIONS")
	f.p.space()

	if obj := ctx.ObjectLiteral(); obj != nil {
		f.list.formatObjectLiteral(obj.(*fql.ObjectLiteralContext))
	}
}

func (f *clauseFormatter) formatTimeoutClause(ctx *fql.TimeoutClauseContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("TIMEOUT")
	f.p.space()

	switch {
	case ctx.DurationLiteral() != nil:
		f.p.write(ctx.DurationLiteral().GetText())
	case ctx.IntegerLiteral() != nil:
		f.p.write(ctx.IntegerLiteral().GetText())
	case ctx.FloatLiteral() != nil:
		f.p.write(ctx.FloatLiteral().GetText())
	case ctx.Variable() != nil:
		f.expression.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.Param() != nil:
		f.expression.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.MemberExpression() != nil:
		f.member.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	case ctx.FunctionCall() != nil:
		f.expression.formatFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext))
	}
}

func (f *clauseFormatter) formatEveryClause(ctx *fql.EveryClauseContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("EVERY")
	f.p.space()
	values := ctx.AllEveryClauseValue()

	for i, val := range values {
		f.formatEveryClauseValue(val.(*fql.EveryClauseValueContext))
		if i < len(values)-1 {
			f.p.write(",")
			f.p.space()
		}
	}
}

func (f *clauseFormatter) formatEveryClauseValue(ctx *fql.EveryClauseValueContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.DurationLiteral() != nil:
		f.p.write(ctx.DurationLiteral().GetText())
	case ctx.IntegerLiteral() != nil:
		f.p.write(ctx.IntegerLiteral().GetText())
	case ctx.FloatLiteral() != nil:
		f.p.write(ctx.FloatLiteral().GetText())
	case ctx.Variable() != nil:
		f.expression.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.Param() != nil:
		f.expression.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.MemberExpression() != nil:
		f.member.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	case ctx.FunctionCall() != nil:
		f.expression.formatFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext))
	}
}

func (f *clauseFormatter) formatBackoffClause(ctx *fql.BackoffClauseContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("BACKOFF")
	f.p.space()

	if strat := ctx.BackoffStrategy(); strat != nil {
		f.formatBackoffStrategy(strat.(*fql.BackoffStrategyContext))
	}
}

func (f *clauseFormatter) formatBackoffStrategy(ctx *fql.BackoffStrategyContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.Identifier() != nil:
		f.p.write(ctx.Identifier().GetText())
	case ctx.StringLiteral() != nil:
		f.literal.formatStringLiteralNode(ctx.StringLiteral())
	case ctx.None() != nil:
		f.p.write(applyCase(f.opts.caseMode, ctx.None().GetText()))
	}
}

func (f *clauseFormatter) formatJitterClause(ctx *fql.JitterClauseContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("JITTER")
	f.p.space()

	if val := ctx.JitterClauseValue(); val != nil {
		f.formatJitterClauseValue(val.(*fql.JitterClauseValueContext))
	}
}

func (f *clauseFormatter) formatJitterClauseValue(ctx *fql.JitterClauseValueContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.FloatLiteral() != nil:
		f.p.write(ctx.FloatLiteral().GetText())
	case ctx.IntegerLiteral() != nil:
		f.p.write(ctx.IntegerLiteral().GetText())
	case ctx.Variable() != nil:
		f.expression.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.Param() != nil:
		f.expression.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.MemberExpression() != nil:
		f.member.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	case ctx.FunctionCall() != nil:
		f.expression.formatFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext))
	}
}
