package formatter

import (
	"io"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type visitor struct {
	*fql.BaseFqlParserVisitor
	opts *options
	p    *printer
	src  *file.Source
}

func newVisitor(src *file.Source, out io.Writer, opts *options) *visitor {
	v := &visitor{
		BaseFqlParserVisitor: new(fql.BaseFqlParserVisitor),
		opts:                 opts,
		p:                    newPrinter(out, opts),
		src:                  src,
	}

	return v
}

func (v *visitor) Err() error {
	return v.p.Err()
}

func (v *visitor) startIndex(ctx antlr.ParserRuleContext) int {
	if ctx == nil {
		return 0
	}
	if tok := ctx.GetStart(); tok != nil {
		return tok.GetStart()
	}
	return 0
}

func (v *visitor) stopIndex(ctx antlr.ParserRuleContext) int {
	if ctx == nil {
		return 0
	}
	if tok := ctx.GetStop(); tok != nil {
		return tok.GetStop()
	}
	return 0
}

func (v *visitor) sliceBetween(start, end int) string {
	if v.src == nil {
		return ""
	}
	text := v.src.Content()
	if start < 0 {
		start = 0
	}
	if end > len(text) {
		end = len(text)
	}
	if end <= start {
		return ""
	}
	return text[start:end]
}

func (v *visitor) emitTrivia(text string) {
	if text == "" {
		return
	}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			if !v.p.atLineStart {
				v.p.space()
			}
			v.p.write(trimmed)
		}
		if i < len(lines)-1 {
			v.p.newline()
		}
	}
}

func (v *visitor) emitBetween(prev, next antlr.ParserRuleContext) {
	if prev == nil || next == nil {
		return
	}
	start := v.stopIndex(prev) + 1
	end := v.startIndex(next)
	v.emitBetweenIndices(start, end)
}

func (v *visitor) emitBetweenIndices(start, end int) {
	text := v.sliceBetween(start, end)
	if text == "" {
		v.p.newline()
		return
	}
	if strings.TrimSpace(text) == "" && !strings.Contains(text, "\n") {
		v.p.newline()
		return
	}

	v.emitTrivia(text)
	if !strings.Contains(text, "\n") {
		v.p.newline()
	}
}

func (v *visitor) emitLeading(next antlr.ParserRuleContext) {
	if next == nil {
		return
	}
	v.emitTrivia(v.sliceBetween(0, v.startIndex(next)))
}

func (v *visitor) emitTrailing(prev antlr.ParserRuleContext) {
	if prev == nil {
		return
	}
	start := v.stopIndex(prev) + 1
	v.emitTrivia(v.sliceBetween(start, len(v.src.Content())))
}

func (v *visitor) bodyFirstElement(ctx *fql.BodyContext) antlr.ParserRuleContext {
	if ctx == nil {
		return nil
	}
	stmts := ctx.AllBodyStatement()
	if len(stmts) > 0 {
		return stmts[0].(antlr.ParserRuleContext)
	}
	if expr := ctx.BodyExpression(); expr != nil {
		return expr.(antlr.ParserRuleContext)
	}
	return nil
}

func (v *visitor) bodyLastElement(ctx *fql.BodyContext) antlr.ParserRuleContext {
	if ctx == nil {
		return nil
	}
	if expr := ctx.BodyExpression(); expr != nil {
		return expr.(antlr.ParserRuleContext)
	}
	stmts := ctx.AllBodyStatement()
	if len(stmts) > 0 {
		return stmts[len(stmts)-1].(antlr.ParserRuleContext)
	}
	return nil
}

func (v *visitor) forHeaderStopIndex(ctx *fql.ForExpressionContext) int {
	if ctx == nil {
		return 0
	}
	switch {
	case ctx.In() != nil:
		if src := ctx.ForExpressionSource(); src != nil {
			return v.stopIndex(src.(antlr.ParserRuleContext))
		}
	case ctx.Step() != nil:
		if expr := ctx.GetStepUpdateExp(); expr != nil {
			return v.stopIndex(expr.(antlr.ParserRuleContext))
		}
		if tok := ctx.GetStepUpdate(); tok != nil {
			return tok.GetStop()
		}
		if tok := ctx.GetStepVariable(); tok != nil {
			return tok.GetStop()
		}
	case ctx.While() != nil:
		if expr := ctx.Expression(0); expr != nil {
			return v.stopIndex(expr.(antlr.ParserRuleContext))
		}
	}
	return v.stopIndex(ctx)
}

func (v *visitor) VisitProgram(ctx *fql.ProgramContext) interface{} {
	v.formatProgram(ctx)
	return nil
}

func (v *visitor) formatProgram(ctx *fql.ProgramContext) {
	if ctx == nil {
		return
	}

	heads := ctx.AllHead()
	var first antlr.ParserRuleContext
	if len(heads) > 0 {
		first = heads[0].(antlr.ParserRuleContext)
	} else if body := ctx.Body(); body != nil {
		first = v.bodyFirstElement(body.(*fql.BodyContext))
	}
	v.emitLeading(first)

	for i, head := range heads {
		if i > 0 {
			v.emitBetween(heads[i-1].(antlr.ParserRuleContext), head.(antlr.ParserRuleContext))
		}
		v.formatHead(head.(*fql.HeadContext))
	}

	if ctx.Body() != nil {
		body := ctx.Body().(*fql.BodyContext)
		if len(heads) > 0 {
			v.emitBetween(heads[len(heads)-1].(antlr.ParserRuleContext), v.bodyFirstElement(body))
		}
		v.formatBody(body)
		last := v.bodyLastElement(body)
		if last == nil && len(heads) > 0 {
			last = heads[len(heads)-1].(antlr.ParserRuleContext)
		}
		v.emitTrailing(last)
	} else if len(heads) > 0 {
		v.emitTrailing(heads[len(heads)-1].(antlr.ParserRuleContext))
	}
}

func (v *visitor) formatHead(ctx *fql.HeadContext) {
	if ctx == nil {
		return
	}

	if useExpr := ctx.UseExpression(); useExpr != nil {
		v.formatUseExpression(useExpr.(*fql.UseExpressionContext))
	}
}

func (v *visitor) formatUseExpression(ctx *fql.UseExpressionContext) {
	if ctx == nil {
		return
	}

	if use := ctx.Use(); use != nil {
		v.formatUse(use.(*fql.UseContext))
	}
}

func (v *visitor) formatUse(ctx *fql.UseContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("USE")
	v.p.space()
	if ns := ctx.NamespaceIdentifier(); ns != nil {
		v.p.write(ns.GetText())
	}
	v.p.space()
	v.writeKeyword("AS")
	v.p.space()
	if alias := ctx.GetAlias(); alias != nil {
		v.p.write(alias.GetText())
	}
}

func (v *visitor) formatBody(ctx *fql.BodyContext) {
	if ctx == nil {
		return
	}

	stmts := ctx.AllBodyStatement()
	for i, stmt := range stmts {
		v.formatBodyStatement(stmt.(*fql.BodyStatementContext))
		if i < len(stmts)-1 {
			v.emitBetween(stmt.(antlr.ParserRuleContext), stmts[i+1].(antlr.ParserRuleContext))
		}
	}

	if expr := ctx.BodyExpression(); expr != nil {
		if len(stmts) > 0 {
			v.emitBetween(stmts[len(stmts)-1].(antlr.ParserRuleContext), expr.(antlr.ParserRuleContext))
		}
		v.formatBodyExpression(expr.(*fql.BodyExpressionContext))
	}
}

func (v *visitor) formatBodyStatement(ctx *fql.BodyStatementContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.VariableDeclaration() != nil:
		v.formatVariableDeclaration(ctx.VariableDeclaration().(*fql.VariableDeclarationContext))
	case ctx.FunctionCallExpression() != nil:
		v.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.WaitForExpression() != nil:
		v.formatWaitForExpression(ctx.WaitForExpression().(*fql.WaitForExpressionContext))
	case ctx.DispatchExpression() != nil:
		v.formatDispatchExpression(ctx.DispatchExpression().(*fql.DispatchExpressionContext))
	}
}

func (v *visitor) formatBodyExpression(ctx *fql.BodyExpressionContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.ReturnExpression() != nil:
		v.formatReturnExpression(ctx.ReturnExpression().(*fql.ReturnExpressionContext))
	case ctx.ForExpression() != nil:
		v.formatForExpression(ctx.ForExpression().(*fql.ForExpressionContext))
	}
}

func (v *visitor) formatVariableDeclaration(ctx *fql.VariableDeclarationContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("LET")
	v.p.space()
	if id := ctx.Identifier(); id != nil {
		v.p.write(id.GetText())
	} else if id := ctx.IgnoreIdentifier(); id != nil {
		v.p.write(id.GetText())
	} else if id := ctx.SafeReservedWord(); id != nil {
		v.p.write(id.GetText())
	}
	v.p.space()
	v.p.write("=")
	v.p.space()
	if expr := ctx.Expression(); expr != nil {
		v.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (v *visitor) formatReturnExpression(ctx *fql.ReturnExpressionContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("RETURN")
	if ctx.Distinct() != nil {
		v.p.space()
		v.writeKeyword("DISTINCT")
	}
	v.p.space()
	if expr := ctx.Expression(); expr != nil {
		v.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (v *visitor) formatForExpression(ctx *fql.ForExpressionContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("FOR")
	v.p.space()
	if tok := ctx.GetValueVariable(); tok != nil {
		v.p.write(tok.GetText())
	}
	if tok := ctx.GetCounterVariable(); tok != nil {
		v.p.write(",")
		v.p.space()
		v.p.write(tok.GetText())
	}

	switch {
	case ctx.In() != nil:
		v.p.space()
		v.writeKeyword("IN")
		v.p.space()
		v.formatForExpressionSource(ctx.ForExpressionSource().(*fql.ForExpressionSourceContext))
	case ctx.Step() != nil:
		v.p.space()
		v.p.write("=")
		v.p.space()
		v.formatExpression(ctx.GetStepInit().(*fql.ExpressionContext))
		v.p.space()
		v.writeKeyword("WHILE")
		v.p.space()
		v.formatExpression(ctx.GetStepCondition().(*fql.ExpressionContext))
		v.p.space()
		v.writeKeyword("STEP")
		v.p.space()
		if tok := ctx.GetStepVariable(); tok != nil {
			v.p.write(tok.GetText())
		}
		if tok := ctx.GetStepUpdate(); tok != nil {
			v.p.write(tok.GetText())
		} else if expr := ctx.GetStepUpdateExp(); expr != nil {
			v.p.space()
			v.p.write("=")
			v.p.space()
			v.formatExpression(expr.(*fql.ExpressionContext))
		}
	case ctx.While() != nil:
		if ctx.Do() != nil {
			v.p.space()
			v.writeKeyword("DO")
		}
		v.p.space()
		v.writeKeyword("WHILE")
		v.p.space()
		v.formatExpression(ctx.Expression(0).(*fql.ExpressionContext))
	}

	bodies := ctx.AllForExpressionBody()
	ret := ctx.ForExpressionReturn()
	if len(bodies) == 0 && ret == nil {
		return
	}

	headerStop := v.forHeaderStopIndex(ctx)
	var first antlr.ParserRuleContext
	if len(bodies) > 0 {
		first = bodies[0].(antlr.ParserRuleContext)
	} else if ret != nil {
		first = ret.(antlr.ParserRuleContext)
	}

	v.p.withIndent(func() {
		if first != nil {
			v.emitBetweenIndices(headerStop+1, v.startIndex(first))
		} else {
			v.p.newline()
		}

		for i, body := range bodies {
			v.formatForExpressionBody(body.(*fql.ForExpressionBodyContext))
			if i < len(bodies)-1 {
				v.emitBetween(body.(antlr.ParserRuleContext), bodies[i+1].(antlr.ParserRuleContext))
			}
		}
		if ret != nil {
			if len(bodies) > 0 {
				v.emitBetween(bodies[len(bodies)-1].(antlr.ParserRuleContext), ret.(antlr.ParserRuleContext))
			}
			v.formatForExpressionReturn(ret.(*fql.ForExpressionReturnContext))
		}
	})
}

func (v *visitor) formatForExpressionSource(ctx *fql.ForExpressionSourceContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.FunctionCallExpression() != nil:
		v.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.ArrayLiteral() != nil:
		v.formatArrayLiteral(ctx.ArrayLiteral().(*fql.ArrayLiteralContext))
	case ctx.ObjectLiteral() != nil:
		v.formatObjectLiteral(ctx.ObjectLiteral().(*fql.ObjectLiteralContext))
	case ctx.Variable() != nil:
		v.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.MemberExpression() != nil:
		v.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	case ctx.RangeOperator() != nil:
		v.formatRangeOperator(ctx.RangeOperator().(*fql.RangeOperatorContext))
	case ctx.Param() != nil:
		v.formatParam(ctx.Param().(*fql.ParamContext))
	}
}

func (v *visitor) formatForExpressionBody(ctx *fql.ForExpressionBodyContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.ForExpressionStatement() != nil:
		stmt := ctx.ForExpressionStatement().(*fql.ForExpressionStatementContext)
		switch {
		case stmt.VariableDeclaration() != nil:
			v.formatVariableDeclaration(stmt.VariableDeclaration().(*fql.VariableDeclarationContext))
		case stmt.FunctionCallExpression() != nil:
			v.formatFunctionCallExpression(stmt.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
		}
	case ctx.ForExpressionClause() != nil:
		clause := ctx.ForExpressionClause().(*fql.ForExpressionClauseContext)
		switch {
		case clause.FilterClause() != nil:
			v.formatFilterClause(clause.FilterClause().(*fql.FilterClauseContext))
		case clause.LimitClause() != nil:
			v.formatLimitClause(clause.LimitClause().(*fql.LimitClauseContext))
		case clause.SortClause() != nil:
			v.formatSortClause(clause.SortClause().(*fql.SortClauseContext))
		case clause.CollectClause() != nil:
			v.formatCollectClause(clause.CollectClause().(*fql.CollectClauseContext))
		}
	}
}

func (v *visitor) formatForExpressionReturn(ctx *fql.ForExpressionReturnContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.ReturnExpression() != nil:
		v.formatReturnExpression(ctx.ReturnExpression().(*fql.ReturnExpressionContext))
	case ctx.ForExpression() != nil:
		v.formatForExpression(ctx.ForExpression().(*fql.ForExpressionContext))
	}
}

func (v *visitor) formatFilterClause(ctx *fql.FilterClauseContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("FILTER")
	v.p.space()
	if expr := ctx.Expression(); expr != nil {
		v.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (v *visitor) formatLimitClause(ctx *fql.LimitClauseContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("LIMIT")
	v.p.space()
	values := ctx.AllLimitClauseValue()
	for i, val := range values {
		v.formatLimitClauseValue(val.(*fql.LimitClauseValueContext))
		if i < len(values)-1 {
			v.p.write(",")
			v.p.space()
		}
	}
}

func (v *visitor) formatLimitClauseValue(ctx *fql.LimitClauseValueContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.IntegerLiteral() != nil:
		v.p.write(ctx.IntegerLiteral().GetText())
	case ctx.Param() != nil:
		v.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.Variable() != nil:
		v.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.FunctionCallExpression() != nil:
		v.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.MemberExpression() != nil:
		v.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	}
}

func (v *visitor) formatSortClause(ctx *fql.SortClauseContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("SORT")
	v.p.space()
	exprs := ctx.AllSortClauseExpression()
	for i, expr := range exprs {
		v.formatSortClauseExpression(expr.(*fql.SortClauseExpressionContext))
		if i < len(exprs)-1 {
			v.p.write(",")
			v.p.space()
		}
	}
}

func (v *visitor) formatSortClauseExpression(ctx *fql.SortClauseExpressionContext) {
	if ctx == nil {
		return
	}

	if expr := ctx.Expression(); expr != nil {
		v.formatExpression(expr.(*fql.ExpressionContext))
	}
	if ctx.SortDirection() != nil {
		v.p.space()
		v.p.write(applyCase(v.opts.caseMode, ctx.SortDirection().GetText()))
	}
}

func (v *visitor) formatCollectClause(ctx *fql.CollectClauseContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("COLLECT")
	if grouping := ctx.CollectGrouping(); grouping != nil {
		v.p.space()
		v.formatCollectGrouping(grouping.(*fql.CollectGroupingContext))
	}
	if aggregator := ctx.CollectAggregator(); aggregator != nil {
		v.p.space()
		v.formatCollectAggregator(aggregator.(*fql.CollectAggregatorContext))
	}
	if counter := ctx.CollectCounter(); counter != nil {
		v.p.space()
		v.formatCollectCounter(counter.(*fql.CollectCounterContext))
	}
	if projection := ctx.CollectGroupProjection(); projection != nil {
		v.p.space()
		v.formatCollectGroupProjection(projection.(*fql.CollectGroupProjectionContext))
	}
}

func (v *visitor) formatCollectGrouping(ctx *fql.CollectGroupingContext) {
	if ctx == nil {
		return
	}

	selectors := ctx.AllCollectSelector()
	for i, sel := range selectors {
		v.formatCollectSelector(sel.(*fql.CollectSelectorContext))
		if i < len(selectors)-1 {
			v.p.write(",")
			v.p.space()
		}
	}
}

func (v *visitor) formatCollectSelector(ctx *fql.CollectSelectorContext) {
	if ctx == nil {
		return
	}

	if id := ctx.Identifier(); id != nil {
		v.p.write(id.GetText())
	}
	v.p.space()
	v.p.write("=")
	v.p.space()
	if expr := ctx.Expression(); expr != nil {
		v.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (v *visitor) formatCollectAggregator(ctx *fql.CollectAggregatorContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("AGGREGATE")
	v.p.space()
	selectors := ctx.AllCollectAggregateSelector()
	for i, sel := range selectors {
		v.formatCollectAggregateSelector(sel.(*fql.CollectAggregateSelectorContext))
		if i < len(selectors)-1 {
			v.p.write(",")
			v.p.space()
		}
	}
}

func (v *visitor) formatCollectAggregateSelector(ctx *fql.CollectAggregateSelectorContext) {
	if ctx == nil {
		return
	}

	if id := ctx.Identifier(); id != nil {
		v.p.write(id.GetText())
	}
	v.p.space()
	v.p.write("=")
	v.p.space()
	if call := ctx.FunctionCallExpression(); call != nil {
		v.formatFunctionCallExpression(call.(*fql.FunctionCallExpressionContext))
	}
}

func (v *visitor) formatCollectGroupProjection(ctx *fql.CollectGroupProjectionContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("INTO")
	v.p.space()
	if sel := ctx.CollectSelector(); sel != nil {
		v.formatCollectSelector(sel.(*fql.CollectSelectorContext))
		return
	}
	if id := ctx.Identifier(); id != nil {
		v.p.write(id.GetText())
		if filter := ctx.CollectGroupProjectionFilter(); filter != nil {
			v.p.space()
			v.formatCollectGroupProjectionFilter(filter.(*fql.CollectGroupProjectionFilterContext))
		}
	}
}

func (v *visitor) formatCollectGroupProjectionFilter(ctx *fql.CollectGroupProjectionFilterContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("KEEP")
	v.p.space()
	ids := ctx.AllIdentifier()
	for i, id := range ids {
		v.p.write(id.GetText())
		if i < len(ids)-1 {
			v.p.write(",")
			v.p.space()
		}
	}
}

func (v *visitor) formatCollectCounter(ctx *fql.CollectCounterContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("WITH")
	v.p.space()
	if id := ctx.Identifier(0); id != nil {
		v.p.write(id.GetText())
	}
	v.p.space()
	v.writeKeyword("INTO")
	v.p.space()
	if id := ctx.Identifier(1); id != nil {
		v.p.write(id.GetText())
	}
}

func (v *visitor) formatWaitForExpression(ctx *fql.WaitForExpressionContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("WAITFOR")
	v.p.space()
	if event := ctx.WaitForEventExpression(); event != nil {
		v.formatWaitForEventExpression(event.(*fql.WaitForEventExpressionContext))
	} else if pred := ctx.WaitForPredicateExpression(); pred != nil {
		v.formatWaitForPredicateExpression(pred.(*fql.WaitForPredicateExpressionContext))
	}

	if ctx.WaitForOrThrowClause() != nil {
		v.p.space()
		v.writeKeyword("OR")
		v.p.space()
		v.writeKeyword("THROW")
	}
}

func (v *visitor) formatWaitForEventExpression(ctx *fql.WaitForEventExpressionContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("EVENT")
	v.p.space()
	if name := ctx.WaitForEventName(); name != nil {
		v.formatWaitForEventName(name.(*fql.WaitForEventNameContext))
	}
	v.p.space()
	v.writeKeyword("IN")
	v.p.space()
	if src := ctx.WaitForEventSource(); src != nil {
		v.formatWaitForEventSource(src.(*fql.WaitForEventSourceContext))
	}
	if opt := ctx.OptionsClause(); opt != nil {
		v.p.space()
		v.formatOptionsClause(opt.(*fql.OptionsClauseContext))
	}
	if filter := ctx.FilterClause(); filter != nil {
		v.p.space()
		v.formatFilterClause(filter.(*fql.FilterClauseContext))
	}
	if timeout := ctx.TimeoutClause(); timeout != nil {
		v.p.space()
		v.formatTimeoutClause(timeout.(*fql.TimeoutClauseContext))
	}
}

func (v *visitor) formatWaitForPredicateExpression(ctx *fql.WaitForPredicateExpressionContext) {
	if ctx == nil {
		return
	}

	if pred := ctx.WaitForPredicate(); pred != nil {
		v.formatWaitForPredicate(pred.(*fql.WaitForPredicateContext))
	}
	if timeout := ctx.TimeoutClause(); timeout != nil {
		v.p.space()
		v.formatTimeoutClause(timeout.(*fql.TimeoutClauseContext))
	}
	if every := ctx.EveryClause(); every != nil {
		v.p.space()
		v.formatEveryClause(every.(*fql.EveryClauseContext))
	}
	if backoff := ctx.BackoffClause(); backoff != nil {
		v.p.space()
		v.formatBackoffClause(backoff.(*fql.BackoffClauseContext))
	}
	if jitter := ctx.JitterClause(); jitter != nil {
		v.p.space()
		v.formatJitterClause(jitter.(*fql.JitterClauseContext))
	}
}

func (v *visitor) formatWaitForPredicate(ctx *fql.WaitForPredicateContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.Not() != nil && ctx.Exists() != nil:
		v.writeKeyword("NOT")
		v.p.space()
		v.writeKeyword("EXISTS")
		v.p.space()
		v.formatExpression(ctx.Expression().(*fql.ExpressionContext))
	case ctx.Exists() != nil:
		v.writeKeyword("EXISTS")
		v.p.space()
		v.formatExpression(ctx.Expression().(*fql.ExpressionContext))
	case ctx.Value() != nil:
		v.writeKeyword("VALUE")
		v.p.space()
		v.formatExpression(ctx.Expression().(*fql.ExpressionContext))
	case ctx.Expression() != nil:
		v.formatExpression(ctx.Expression().(*fql.ExpressionContext))
	}
}

func (v *visitor) formatOptionsClause(ctx *fql.OptionsClauseContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("OPTIONS")
	v.p.space()
	if obj := ctx.ObjectLiteral(); obj != nil {
		v.formatObjectLiteral(obj.(*fql.ObjectLiteralContext))
	}
}

func (v *visitor) formatTimeoutClause(ctx *fql.TimeoutClauseContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("TIMEOUT")
	v.p.space()
	switch {
	case ctx.DurationLiteral() != nil:
		v.p.write(ctx.DurationLiteral().GetText())
	case ctx.IntegerLiteral() != nil:
		v.p.write(ctx.IntegerLiteral().GetText())
	case ctx.FloatLiteral() != nil:
		v.p.write(ctx.FloatLiteral().GetText())
	case ctx.Variable() != nil:
		v.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.Param() != nil:
		v.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.MemberExpression() != nil:
		v.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	case ctx.FunctionCall() != nil:
		v.formatFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext))
	}
}

func (v *visitor) formatEveryClause(ctx *fql.EveryClauseContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("EVERY")
	v.p.space()
	values := ctx.AllEveryClauseValue()
	for i, val := range values {
		v.formatEveryClauseValue(val.(*fql.EveryClauseValueContext))
		if i < len(values)-1 {
			v.p.write(",")
			v.p.space()
		}
	}
}

func (v *visitor) formatEveryClauseValue(ctx *fql.EveryClauseValueContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.DurationLiteral() != nil:
		v.p.write(ctx.DurationLiteral().GetText())
	case ctx.IntegerLiteral() != nil:
		v.p.write(ctx.IntegerLiteral().GetText())
	case ctx.FloatLiteral() != nil:
		v.p.write(ctx.FloatLiteral().GetText())
	case ctx.Variable() != nil:
		v.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.Param() != nil:
		v.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.MemberExpression() != nil:
		v.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	case ctx.FunctionCall() != nil:
		v.formatFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext))
	}
}

func (v *visitor) formatBackoffClause(ctx *fql.BackoffClauseContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("BACKOFF")
	v.p.space()
	if strat := ctx.BackoffStrategy(); strat != nil {
		v.formatBackoffStrategy(strat.(*fql.BackoffStrategyContext))
	}
}

func (v *visitor) formatBackoffStrategy(ctx *fql.BackoffStrategyContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.Identifier() != nil:
		v.p.write(ctx.Identifier().GetText())
	case ctx.StringLiteral() != nil:
		v.p.write(formatStringLiteral(ctx.StringLiteral().GetText(), v.opts.singleQuote))
	case ctx.None() != nil:
		v.p.write(applyCase(v.opts.caseMode, ctx.None().GetText()))
	}
}

func (v *visitor) formatJitterClause(ctx *fql.JitterClauseContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("JITTER")
	v.p.space()
	if val := ctx.JitterClauseValue(); val != nil {
		v.formatJitterClauseValue(val.(*fql.JitterClauseValueContext))
	}
}

func (v *visitor) formatJitterClauseValue(ctx *fql.JitterClauseValueContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.FloatLiteral() != nil:
		v.p.write(ctx.FloatLiteral().GetText())
	case ctx.IntegerLiteral() != nil:
		v.p.write(ctx.IntegerLiteral().GetText())
	case ctx.Variable() != nil:
		v.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.Param() != nil:
		v.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.MemberExpression() != nil:
		v.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	case ctx.FunctionCall() != nil:
		v.formatFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext))
	}
}

func (v *visitor) formatDispatchExpression(ctx *fql.DispatchExpressionContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("DISPATCH")
	v.p.space()
	if name := ctx.DispatchEventName(); name != nil {
		v.formatDispatchEventName(name.(*fql.DispatchEventNameContext))
	}
	v.p.space()
	v.writeKeyword("IN")
	v.p.space()
	if tgt := ctx.DispatchTarget(); tgt != nil {
		v.formatDispatchTarget(tgt.(*fql.DispatchTargetContext))
	}
	if with := ctx.DispatchWithClause(); with != nil {
		v.p.space()
		v.formatDispatchWithClause(with.(*fql.DispatchWithClauseContext))
	}
	if opt := ctx.DispatchOptionsClause(); opt != nil {
		v.p.space()
		v.formatDispatchOptionsClause(opt.(*fql.DispatchOptionsClauseContext))
	}
}

func (v *visitor) formatDispatchEventName(ctx *fql.DispatchEventNameContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.StringLiteral() != nil:
		v.p.write(formatStringLiteral(ctx.StringLiteral().GetText(), v.opts.singleQuote))
	case ctx.Variable() != nil:
		v.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.Param() != nil:
		v.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.MemberExpression() != nil:
		v.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	case ctx.FunctionCall() != nil:
		v.formatFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext))
	}
}

func (v *visitor) formatDispatchTarget(ctx *fql.DispatchTargetContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.FunctionCallExpression() != nil:
		v.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.Variable() != nil:
		v.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.Param() != nil:
		v.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.MemberExpression() != nil:
		v.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	}
}

func (v *visitor) formatDispatchWithClause(ctx *fql.DispatchWithClauseContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("WITH")
	v.p.space()
	if expr := ctx.Expression(); expr != nil {
		v.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (v *visitor) formatDispatchOptionsClause(ctx *fql.DispatchOptionsClauseContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("OPTIONS")
	v.p.space()
	if expr := ctx.Expression(); expr != nil {
		v.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (v *visitor) formatWaitForEventName(ctx *fql.WaitForEventNameContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.StringLiteral() != nil:
		v.p.write(formatStringLiteral(ctx.StringLiteral().GetText(), v.opts.singleQuote))
	case ctx.Variable() != nil:
		v.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.Param() != nil:
		v.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.MemberExpression() != nil:
		v.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	case ctx.FunctionCall() != nil:
		v.formatFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext))
	}
}

func (v *visitor) formatWaitForEventSource(ctx *fql.WaitForEventSourceContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.FunctionCallExpression() != nil:
		v.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.Variable() != nil:
		v.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.MemberExpression() != nil:
		v.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	}
}

func (v *visitor) formatExpression(ctx *fql.ExpressionContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.UnaryOperator() != nil:
		v.formatUnaryOperator(ctx.UnaryOperator().(*fql.UnaryOperatorContext))
		v.formatExpression(ctx.GetRight().(*fql.ExpressionContext))
	case ctx.LogicalAndOperator() != nil:
		v.formatExpression(ctx.GetLeft().(*fql.ExpressionContext))
		v.p.space()
		v.formatLogicalAndOperator(ctx.LogicalAndOperator().(*fql.LogicalAndOperatorContext))
		v.p.space()
		v.formatExpression(ctx.GetRight().(*fql.ExpressionContext))
	case ctx.LogicalOrOperator() != nil:
		v.formatExpression(ctx.GetLeft().(*fql.ExpressionContext))
		v.p.space()
		v.formatLogicalOrOperator(ctx.LogicalOrOperator().(*fql.LogicalOrOperatorContext))
		v.p.space()
		v.formatExpression(ctx.GetRight().(*fql.ExpressionContext))
	case ctx.GetTernaryOperator() != nil:
		v.formatExpression(ctx.GetCondition().(*fql.ExpressionContext))
		v.p.space()
		v.p.write("?")
		v.p.space()
		if ctx.GetOnTrue() != nil {
			v.formatExpression(ctx.GetOnTrue().(*fql.ExpressionContext))
		}
		v.p.space()
		v.p.write(":")
		v.p.space()
		v.formatExpression(ctx.GetOnFalse().(*fql.ExpressionContext))
	case ctx.Predicate() != nil:
		v.formatPredicate(ctx.Predicate().(*fql.PredicateContext))
	}
}

func (v *visitor) formatPredicate(ctx *fql.PredicateContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.EqualityOperator() != nil:
		v.formatPredicate(ctx.GetLeft().(*fql.PredicateContext))
		v.p.space()
		v.formatEqualityOperator(ctx.EqualityOperator().(*fql.EqualityOperatorContext))
		v.p.space()
		v.formatPredicate(ctx.GetRight().(*fql.PredicateContext))
	case ctx.ArrayOperator() != nil:
		v.formatPredicate(ctx.GetLeft().(*fql.PredicateContext))
		v.p.space()
		v.formatArrayOperator(ctx.ArrayOperator().(*fql.ArrayOperatorContext))
		v.p.space()
		v.formatPredicate(ctx.GetRight().(*fql.PredicateContext))
	case ctx.InOperator() != nil:
		v.formatPredicate(ctx.GetLeft().(*fql.PredicateContext))
		v.p.space()
		v.formatInOperator(ctx.InOperator().(*fql.InOperatorContext))
		v.p.space()
		v.formatPredicate(ctx.GetRight().(*fql.PredicateContext))
	case ctx.LikeOperator() != nil:
		v.formatPredicate(ctx.GetLeft().(*fql.PredicateContext))
		v.p.space()
		v.formatLikeOperator(ctx.LikeOperator().(*fql.LikeOperatorContext))
		v.p.space()
		v.formatPredicate(ctx.GetRight().(*fql.PredicateContext))
	case ctx.ExpressionAtom() != nil:
		v.formatExpressionAtom(ctx.ExpressionAtom().(*fql.ExpressionAtomContext))
	}
}

func (v *visitor) formatExpressionAtom(ctx *fql.ExpressionAtomContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.MultiplicativeOperator() != nil:
		v.formatExpressionAtom(ctx.GetLeft().(*fql.ExpressionAtomContext))
		v.p.space()
		v.formatMultiplicativeOperator(ctx.MultiplicativeOperator().(*fql.MultiplicativeOperatorContext))
		v.p.space()
		v.formatExpressionAtom(ctx.GetRight().(*fql.ExpressionAtomContext))
	case ctx.AdditiveOperator() != nil:
		v.formatExpressionAtom(ctx.GetLeft().(*fql.ExpressionAtomContext))
		v.p.space()
		v.formatAdditiveOperator(ctx.AdditiveOperator().(*fql.AdditiveOperatorContext))
		v.p.space()
		v.formatExpressionAtom(ctx.GetRight().(*fql.ExpressionAtomContext))
	case ctx.RegexpOperator() != nil:
		v.formatExpressionAtom(ctx.GetLeft().(*fql.ExpressionAtomContext))
		v.p.space()
		v.formatRegexpOperator(ctx.RegexpOperator().(*fql.RegexpOperatorContext))
		v.p.space()
		v.formatExpressionAtom(ctx.GetRight().(*fql.ExpressionAtomContext))
	case ctx.FunctionCallExpression() != nil:
		v.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.RangeOperator() != nil:
		v.formatRangeOperator(ctx.RangeOperator().(*fql.RangeOperatorContext))
	case ctx.Literal() != nil:
		v.formatLiteral(ctx.Literal().(*fql.LiteralContext))
	case ctx.Variable() != nil:
		v.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.MemberExpression() != nil:
		v.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	case ctx.Param() != nil:
		v.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.DispatchExpression() != nil:
		v.formatDispatchExpression(ctx.DispatchExpression().(*fql.DispatchExpressionContext))
	case ctx.WaitForExpression() != nil:
		v.formatWaitForExpression(ctx.WaitForExpression().(*fql.WaitForExpressionContext))
	case ctx.OpenParen() != nil:
		v.formatParenthesizedExpression(ctx)
	}
}

func (v *visitor) formatParenthesizedExpression(ctx *fql.ExpressionAtomContext) {
	if ctx == nil {
		return
	}

	v.p.write("(")

	if fe := ctx.ForExpression(); fe != nil {
		v.p.newline()
		v.p.withIndent(func() {
			v.formatForExpression(fe.(*fql.ForExpressionContext))
		})
		v.p.newline()
		v.p.write(")")
		if ctx.ErrorOperator() != nil {
			v.p.write("?")
		}
		return
	}

	if we := ctx.WaitForExpression(); we != nil {
		v.formatWaitForExpression(we.(*fql.WaitForExpressionContext))
		v.p.write(")")
		if ctx.ErrorOperator() != nil {
			v.p.write("?")
		}
		return
	}

	if expr := ctx.Expression(); expr != nil {
		v.formatExpression(expr.(*fql.ExpressionContext))
	}
	v.p.write(")")
	if ctx.ErrorOperator() != nil {
		v.p.write("?")
	}
}

func (v *visitor) formatUnaryOperator(ctx *fql.UnaryOperatorContext) {
	if ctx == nil {
		return
	}

	op := ctx.GetText()
	if op == "NOT" || op == "!" {
		v.p.write(applyCase(v.opts.caseMode, op))
		v.p.space()
		return
	}
	v.p.write(op)
}

func (v *visitor) formatLogicalAndOperator(ctx *fql.LogicalAndOperatorContext) {
	if ctx == nil {
		return
	}

	v.p.write(applyCase(v.opts.caseMode, ctx.GetText()))
}

func (v *visitor) formatLogicalOrOperator(ctx *fql.LogicalOrOperatorContext) {
	if ctx == nil {
		return
	}

	v.p.write(applyCase(v.opts.caseMode, ctx.GetText()))
}

func (v *visitor) formatEqualityOperator(ctx *fql.EqualityOperatorContext) {
	if ctx == nil {
		return
	}

	v.p.write(ctx.GetText())
}

func (v *visitor) formatArrayOperator(ctx *fql.ArrayOperatorContext) {
	if ctx == nil {
		return
	}

	if op := ctx.GetOperator(); op != nil {
		v.p.write(applyCase(v.opts.caseMode, op.GetText()))
	}
	v.p.space()
	if in := ctx.InOperator(); in != nil {
		v.formatInOperator(in.(*fql.InOperatorContext))
	} else if eq := ctx.EqualityOperator(); eq != nil {
		v.formatEqualityOperator(eq.(*fql.EqualityOperatorContext))
	}
}

func (v *visitor) formatInOperator(ctx *fql.InOperatorContext) {
	if ctx == nil {
		return
	}

	if ctx.Not() != nil {
		v.p.write(applyCase(v.opts.caseMode, ctx.Not().GetText()))
		v.p.space()
	}
	v.p.write(applyCase(v.opts.caseMode, ctx.In().GetText()))
}

func (v *visitor) formatLikeOperator(ctx *fql.LikeOperatorContext) {
	if ctx == nil {
		return
	}

	if ctx.Not() != nil {
		v.p.write(applyCase(v.opts.caseMode, ctx.Not().GetText()))
		v.p.space()
	}
	v.p.write(applyCase(v.opts.caseMode, ctx.Like().GetText()))
}

func (v *visitor) formatMultiplicativeOperator(ctx *fql.MultiplicativeOperatorContext) {
	if ctx == nil {
		return
	}

	v.p.write(ctx.GetText())
}

func (v *visitor) formatAdditiveOperator(ctx *fql.AdditiveOperatorContext) {
	if ctx == nil {
		return
	}

	v.p.write(ctx.GetText())
}

func (v *visitor) formatRegexpOperator(ctx *fql.RegexpOperatorContext) {
	if ctx == nil {
		return
	}

	v.p.write(ctx.GetText())
}

func (v *visitor) formatRangeOperator(ctx *fql.RangeOperatorContext) {
	if ctx == nil {
		return
	}

	if left := ctx.GetLeft(); left != nil {
		v.formatRangeOperand(left.(*fql.RangeOperandContext))
	}
	v.p.write("..")
	if right := ctx.GetRight(); right != nil {
		v.formatRangeOperand(right.(*fql.RangeOperandContext))
	}
}

func (v *visitor) formatRangeOperand(ctx *fql.RangeOperandContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.IntegerLiteral() != nil:
		v.p.write(ctx.IntegerLiteral().GetText())
	case ctx.Variable() != nil:
		v.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.Param() != nil:
		v.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.FunctionCallExpression() != nil:
		v.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.MemberExpression() != nil:
		v.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	}
}

func (v *visitor) formatLiteral(ctx *fql.LiteralContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.ArrayLiteral() != nil:
		v.formatArrayLiteral(ctx.ArrayLiteral().(*fql.ArrayLiteralContext))
	case ctx.ObjectLiteral() != nil:
		v.formatObjectLiteral(ctx.ObjectLiteral().(*fql.ObjectLiteralContext))
	case ctx.BooleanLiteral() != nil:
		v.formatBooleanLiteral(ctx.BooleanLiteral().(*fql.BooleanLiteralContext))
	case ctx.StringLiteral() != nil:
		v.p.write(formatStringLiteral(ctx.StringLiteral().GetText(), v.opts.singleQuote))
	case ctx.FloatLiteral() != nil:
		v.p.write(ctx.FloatLiteral().GetText())
	case ctx.IntegerLiteral() != nil:
		v.p.write(ctx.IntegerLiteral().GetText())
	case ctx.NoneLiteral() != nil:
		v.formatNoneLiteral(ctx.NoneLiteral().(*fql.NoneLiteralContext))
	}
}

func (v *visitor) formatBooleanLiteral(ctx *fql.BooleanLiteralContext) {
	if ctx == nil || ctx.BooleanLiteral() == nil {
		return
	}

	v.p.write(applyCase(v.opts.caseMode, ctx.BooleanLiteral().GetText()))
}

func (v *visitor) formatNoneLiteral(ctx *fql.NoneLiteralContext) {
	if ctx == nil {
		return
	}

	if ctx.Null() != nil {
		v.p.write(applyCase(v.opts.caseMode, ctx.Null().GetText()))
		return
	}
	if ctx.None() != nil {
		v.p.write(applyCase(v.opts.caseMode, ctx.None().GetText()))
	}
}

func (v *visitor) formatArrayLiteral(ctx *fql.ArrayLiteralContext) {
	if ctx == nil {
		return
	}

	if ctx.ArgumentList() == nil {
		v.p.write("[]")
		return
	}

	if v.p.forceSingleLine {
		v.formatArrayLiteralInline(ctx)
		return
	}

	inline := v.renderInline(func(p *printer) {
		v.formatArrayLiteralWith(p, ctx, true)
	})
	if len(inline) <= int(v.opts.printWidth) {
		v.p.write(inline)
		return
	}

	v.formatArrayLiteralWith(v.p, ctx, false)
}

func (v *visitor) formatArrayLiteralInline(ctx *fql.ArrayLiteralContext) {
	v.formatArrayLiteralWith(v.p, ctx, true)
}

func (v *visitor) formatArrayLiteralWith(p *printer, ctx *fql.ArrayLiteralContext, inline bool) {
	args := ctx.ArgumentList().AllExpression()
	p.write("[")
	if !inline {
		p.newline()
		p.withIndent(func() {
			for i, expr := range args {
				v.formatExpressionWith(p, expr.(*fql.ExpressionContext))
				if i < len(args)-1 {
					p.write(",")
				}
				p.newline()
			}
		})
		p.write("]")
		return
	}
	for i, expr := range args {
		v.formatExpressionWith(p, expr.(*fql.ExpressionContext))
		if i < len(args)-1 {
			p.write(",")
			p.space()
		}
	}
	p.write("]")
}

func (v *visitor) formatObjectLiteral(ctx *fql.ObjectLiteralContext) {
	if ctx == nil {
		return
	}

	props := ctx.AllPropertyAssignment()
	if len(props) == 0 {
		v.p.write("{}")
		return
	}

	if v.p.forceSingleLine {
		v.formatObjectLiteralInline(ctx)
		return
	}

	inline := v.renderInline(func(p *printer) {
		v.formatObjectLiteralWith(p, ctx, true)
	})
	if len(inline) <= int(v.opts.printWidth) {
		v.p.write(inline)
		return
	}

	v.formatObjectLiteralWith(v.p, ctx, false)
}

func (v *visitor) formatObjectLiteralInline(ctx *fql.ObjectLiteralContext) {
	v.formatObjectLiteralWith(v.p, ctx, true)
}

func (v *visitor) formatObjectLiteralWith(p *printer, ctx *fql.ObjectLiteralContext, inline bool) {
	props := ctx.AllPropertyAssignment()
	p.write("{")
	if inline {
		if v.opts.bracketSpacing {
			p.space()
		}
		for i, prop := range props {
			v.formatPropertyAssignmentWith(p, prop.(*fql.PropertyAssignmentContext))
			if i < len(props)-1 {
				p.write(",")
				p.space()
			}
		}
		if v.opts.bracketSpacing {
			p.space()
		}
		p.write("}")
		return
	}

	p.newline()
	p.withIndent(func() {
		for i, prop := range props {
			v.formatPropertyAssignmentWith(p, prop.(*fql.PropertyAssignmentContext))
			if i < len(props)-1 {
				p.write(",")
			}
			p.newline()
		}
	})
	p.write("}")
}

func (v *visitor) formatPropertyAssignmentWith(p *printer, ctx *fql.PropertyAssignmentContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.PropertyName() != nil:
		v.formatPropertyNameWith(p, ctx.PropertyName().(*fql.PropertyNameContext))
		p.write(":")
		p.space()
		v.formatExpressionWith(p, ctx.Expression().(*fql.ExpressionContext))
	case ctx.ComputedPropertyName() != nil:
		v.formatComputedPropertyNameWith(p, ctx.ComputedPropertyName().(*fql.ComputedPropertyNameContext))
		p.write(":")
		p.space()
		v.formatExpressionWith(p, ctx.Expression().(*fql.ExpressionContext))
	case ctx.Variable() != nil:
		v.formatVariableWith(p, ctx.Variable().(*fql.VariableContext))
	}
}

func (v *visitor) formatPropertyNameWith(p *printer, ctx *fql.PropertyNameContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.Identifier() != nil:
		p.write(ctx.Identifier().GetText())
	case ctx.StringLiteral() != nil:
		p.write(formatStringLiteral(ctx.StringLiteral().GetText(), v.opts.singleQuote))
	case ctx.Param() != nil:
		v.formatParamWith(p, ctx.Param().(*fql.ParamContext))
	case ctx.SafeReservedWord() != nil:
		p.write(ctx.SafeReservedWord().GetText())
	case ctx.UnsafeReservedWord() != nil:
		p.write(ctx.UnsafeReservedWord().GetText())
	}
}

func (v *visitor) formatComputedPropertyNameWith(p *printer, ctx *fql.ComputedPropertyNameContext) {
	if ctx == nil {
		return
	}
	p.write("[")
	if expr := ctx.Expression(); expr != nil {
		v.formatExpressionWith(p, expr.(*fql.ExpressionContext))
	}
	p.write("]")
}

func (v *visitor) formatFunctionCallExpression(ctx *fql.FunctionCallExpressionContext) {
	if ctx == nil {
		return
	}

	v.formatFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext))
	if ctx.ErrorOperator() != nil {
		v.p.write("?")
	}
}

func (v *visitor) formatFunctionCall(ctx *fql.FunctionCallContext) {
	if ctx == nil {
		return
	}

	if ns := ctx.Namespace(); ns != nil {
		v.p.write(ns.GetText())
	}
	if fn := ctx.FunctionName(); fn != nil {
		v.p.write(fn.GetText())
	}
	v.p.write("(")
	if args := ctx.ArgumentList(); args != nil {
		v.formatArgumentList(args.(*fql.ArgumentListContext))
	}
	v.p.write(")")
}

func (v *visitor) formatArgumentList(ctx *fql.ArgumentListContext) {
	if ctx == nil {
		return
	}

	args := ctx.AllExpression()
	if len(args) == 0 {
		return
	}

	if v.p.forceSingleLine {
		v.formatArgumentListInline(ctx)
		return
	}

	inline := v.renderInline(func(p *printer) {
		v.formatArgumentListWith(p, ctx, true)
	})
	if len(inline) <= int(v.opts.printWidth) {
		v.p.write(inline)
		return
	}

	v.formatArgumentListWith(v.p, ctx, false)
}

func (v *visitor) formatArgumentListInline(ctx *fql.ArgumentListContext) {
	v.formatArgumentListWith(v.p, ctx, true)
}

func (v *visitor) formatArgumentListWith(p *printer, ctx *fql.ArgumentListContext, inline bool) {
	args := ctx.AllExpression()
	if inline {
		for i, arg := range args {
			v.formatExpressionWith(p, arg.(*fql.ExpressionContext))
			if i < len(args)-1 {
				p.write(",")
				p.space()
			}
		}
		return
	}

	p.newline()
	p.withIndent(func() {
		for i, arg := range args {
			v.formatExpressionWith(p, arg.(*fql.ExpressionContext))
			if i < len(args)-1 {
				p.write(",")
			}
			p.newline()
		}
	})
}

func (v *visitor) formatMemberExpression(ctx *fql.MemberExpressionContext) {
	if ctx == nil {
		return
	}

	v.formatMemberExpressionSource(ctx.MemberExpressionSource().(*fql.MemberExpressionSourceContext))
	for _, path := range ctx.AllMemberExpressionPath() {
		v.formatMemberExpressionPath(path.(*fql.MemberExpressionPathContext))
	}
}

func (v *visitor) formatMemberExpressionSource(ctx *fql.MemberExpressionSourceContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.Variable() != nil:
		v.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.Param() != nil:
		v.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.ArrayLiteral() != nil:
		v.formatArrayLiteral(ctx.ArrayLiteral().(*fql.ArrayLiteralContext))
	case ctx.ObjectLiteral() != nil:
		v.formatObjectLiteral(ctx.ObjectLiteral().(*fql.ObjectLiteralContext))
	case ctx.FunctionCall() != nil:
		v.formatFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext))
	case ctx.OpenParen() != nil:
		v.p.write("(")
		if fe := ctx.ForExpression(); fe != nil {
			v.p.newline()
			v.p.withIndent(func() {
				v.formatForExpression(fe.(*fql.ForExpressionContext))
			})
			v.p.newline()
			v.p.write(")")
			return
		}
		if we := ctx.WaitForExpression(); we != nil {
			v.formatWaitForExpression(we.(*fql.WaitForExpressionContext))
			v.p.write(")")
			return
		}
		if expr := ctx.Expression(); expr != nil {
			v.formatExpression(expr.(*fql.ExpressionContext))
		}
		v.p.write(")")
	}
}

func (v *visitor) formatMemberExpressionPath(ctx *fql.MemberExpressionPathContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.PropertyName() != nil:
		if ctx.ErrorOperator() != nil {
			v.p.write("?")
		}
		v.p.write(".")
		v.formatPropertyNameWith(v.p, ctx.PropertyName().(*fql.PropertyNameContext))
	case ctx.ComputedPropertyName() != nil:
		if ctx.ErrorOperator() != nil {
			v.p.write("?")
			v.p.write(".")
		}
		v.formatComputedPropertyNameWith(v.p, ctx.ComputedPropertyName().(*fql.ComputedPropertyNameContext))
	case ctx.ArrayContraction() != nil:
		v.formatArrayContraction(ctx.ArrayContraction().(*fql.ArrayContractionContext))
	case ctx.ArrayExpansion() != nil:
		v.formatArrayExpansion(ctx.ArrayExpansion().(*fql.ArrayExpansionContext))
	case ctx.ArrayQuestionMark() != nil:
		v.formatArrayQuestionMark(ctx.ArrayQuestionMark().(*fql.ArrayQuestionMarkContext))
	case ctx.ArrayApply() != nil:
		v.formatArrayApply(ctx.ArrayApply().(*fql.ArrayApplyContext))
	}
}

func (v *visitor) formatArrayExpansion(ctx *fql.ArrayExpansionContext) {
	if ctx == nil {
		return
	}

	v.p.write("[")
	v.p.write("*")
	if inline := ctx.InlineExpression(); inline != nil {
		v.p.space()
		v.formatInlineExpression(inline.(*fql.InlineExpressionContext))
	}
	v.p.write("]")
}

func (v *visitor) formatArrayContraction(ctx *fql.ArrayContractionContext) {
	if ctx == nil {
		return
	}

	v.p.write("[")
	stars := ctx.AllMulti()
	for range stars {
		v.p.write("*")
	}
	if inline := ctx.InlineExpression(); inline != nil {
		v.p.space()
		v.formatInlineExpression(inline.(*fql.InlineExpressionContext))
	}
	v.p.write("]")
}

func (v *visitor) formatArrayQuestionMark(ctx *fql.ArrayQuestionMarkContext) {
	if ctx == nil {
		return
	}

	v.p.write("[")
	v.p.write("?")
	if quant := ctx.ArrayQuestionQuantifier(); quant != nil {
		v.p.space()
		v.formatArrayQuestionQuantifier(quant.(*fql.ArrayQuestionQuantifierContext))
	}
	v.p.space()
	v.writeKeyword("FILTER")
	v.p.space()
	if expr := ctx.Expression(); expr != nil {
		v.formatExpression(expr.(*fql.ExpressionContext))
	}
	v.p.write("]")
}

func (v *visitor) formatArrayQuestionQuantifier(ctx *fql.ArrayQuestionQuantifierContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.Any() != nil:
		v.p.write(applyCase(v.opts.caseMode, ctx.Any().GetText()))
	case ctx.All() != nil:
		v.p.write(applyCase(v.opts.caseMode, ctx.All().GetText()))
	case ctx.None() != nil:
		v.p.write(applyCase(v.opts.caseMode, ctx.None().GetText()))
	case ctx.At() != nil && ctx.Least() != nil:
		v.p.write(applyCase(v.opts.caseMode, ctx.At().GetText()))
		v.p.space()
		v.p.write(applyCase(v.opts.caseMode, ctx.Least().GetText()))
		v.p.space()
		v.p.write("(")
		if val := ctx.ArrayQuestionQuantifierValue(0); val != nil {
			v.formatArrayQuestionQuantifierValue(val.(*fql.ArrayQuestionQuantifierValueContext))
		}
		v.p.write(")")
	case ctx.Range() != nil:
		v.formatArrayQuestionQuantifierValue(ctx.ArrayQuestionQuantifierValue(0).(*fql.ArrayQuestionQuantifierValueContext))
		v.p.write("..")
		v.formatArrayQuestionQuantifierValue(ctx.ArrayQuestionQuantifierValue(1).(*fql.ArrayQuestionQuantifierValueContext))
	case ctx.ArrayQuestionQuantifierValue(0) != nil:
		v.formatArrayQuestionQuantifierValue(ctx.ArrayQuestionQuantifierValue(0).(*fql.ArrayQuestionQuantifierValueContext))
	}
}

func (v *visitor) formatArrayQuestionQuantifierValue(ctx *fql.ArrayQuestionQuantifierValueContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.IntegerLiteral() != nil:
		v.p.write(ctx.IntegerLiteral().GetText())
	case ctx.Param() != nil:
		v.formatParam(ctx.Param().(*fql.ParamContext))
	}
}

func (v *visitor) formatArrayApply(ctx *fql.ArrayApplyContext) {
	if ctx == nil {
		return
	}

	v.p.write("[")
	v.p.write("~")
	v.p.space()
	if q := ctx.QueryLiteral(); q != nil {
		v.formatQueryLiteral(q.(*fql.QueryLiteralContext))
	}
	v.p.write("]")
}

func (v *visitor) formatInlineExpression(ctx *fql.InlineExpressionContext) {
	if ctx == nil {
		return
	}

	if filter := ctx.InlineFilter(); filter != nil {
		v.formatInlineFilter(filter.(*fql.InlineFilterContext))
	}
	if limit := ctx.InlineLimit(); limit != nil {
		if ctx.InlineFilter() != nil {
			v.p.space()
		}
		v.formatInlineLimit(limit.(*fql.InlineLimitContext))
	}
	if ret := ctx.InlineReturn(); ret != nil {
		if ctx.InlineFilter() != nil || ctx.InlineLimit() != nil {
			v.p.space()
		}
		v.formatInlineReturn(ret.(*fql.InlineReturnContext))
	}
}

func (v *visitor) formatInlineFilter(ctx *fql.InlineFilterContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("FILTER")
	v.p.space()
	if expr := ctx.Expression(); expr != nil {
		v.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (v *visitor) formatInlineLimit(ctx *fql.InlineLimitContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("LIMIT")
	v.p.space()
	values := ctx.AllLimitClauseValue()
	for i, val := range values {
		v.formatLimitClauseValue(val.(*fql.LimitClauseValueContext))
		if i < len(values)-1 {
			v.p.write(",")
			v.p.space()
		}
	}
}

func (v *visitor) formatInlineReturn(ctx *fql.InlineReturnContext) {
	if ctx == nil {
		return
	}

	v.writeKeyword("RETURN")
	v.p.space()
	if expr := ctx.Expression(); expr != nil {
		v.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (v *visitor) formatQueryLiteral(ctx *fql.QueryLiteralContext) {
	if ctx == nil {
		return
	}

	if id := ctx.Identifier(); id != nil {
		v.p.write(id.GetText())
	}
	if sl := ctx.StringLiteral(); sl != nil {
		v.p.space()
		v.p.write(formatStringLiteral(sl.GetText(), v.opts.singleQuote))
		if ctx.OpenParen() != nil && ctx.Expression() != nil {
			v.p.write("(")
			v.formatExpression(ctx.Expression().(*fql.ExpressionContext))
			v.p.write(")")
		}
	}
}

func (v *visitor) formatParam(ctx *fql.ParamContext) {
	v.formatParamWith(v.p, ctx)
}

func (v *visitor) formatParamWith(p *printer, ctx *fql.ParamContext) {
	if ctx == nil {
		return
	}

	p.write("@")
	if id := ctx.Identifier(); id != nil {
		p.write(id.GetText())
	} else if id := ctx.SafeReservedWord(); id != nil {
		p.write(id.GetText())
	}
}

func (v *visitor) formatVariable(ctx *fql.VariableContext) {
	v.formatVariableWith(v.p, ctx)
}

func (v *visitor) formatVariableWith(p *printer, ctx *fql.VariableContext) {
	if ctx == nil {
		return
	}

	if id := ctx.Identifier(); id != nil {
		p.write(id.GetText())
	} else if id := ctx.SafeReservedWord(); id != nil {
		p.write(id.GetText())
	}
}

func (v *visitor) formatExpressionWith(p *printer, ctx *fql.ExpressionContext) {
	if p == v.p {
		v.formatExpression(ctx)
		return
	}

	orig := v.p
	v.p = p
	v.formatExpression(ctx)
	v.p = orig
}

func (v *visitor) formatPropertyAssignment(ctx *fql.PropertyAssignmentContext) {
	v.formatPropertyAssignmentWith(v.p, ctx)
}

func (v *visitor) formatPropertyName(ctx *fql.PropertyNameContext) {
	v.formatPropertyNameWith(v.p, ctx)
}

func (v *visitor) formatComputedPropertyName(ctx *fql.ComputedPropertyNameContext) {
	v.formatComputedPropertyNameWith(v.p, ctx)
}

func (v *visitor) writeKeyword(val string) {
	v.p.write(applyCase(v.opts.caseMode, val))
}

func (v *visitor) renderInline(fn func(p *printer)) string {
	var b strings.Builder
	p := newPrinter(&b, v.opts)
	p.forceSingleLine = true
	fn(p)
	return b.String()
}
