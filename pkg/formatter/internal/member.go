package internal

import "github.com/MontFerret/ferret/v2/pkg/parser/fql"

type memberFormatter struct {
	*engine
}

func (f *memberFormatter) formatMemberExpression(ctx *fql.MemberExpressionContext) {
	if ctx == nil {
		return
	}

	f.formatMemberExpressionSource(ctx.MemberExpressionSource().(*fql.MemberExpressionSourceContext))

	for _, path := range ctx.AllMemberExpressionPath() {
		f.formatMemberExpressionPath(path.(*fql.MemberExpressionPathContext))
	}
}

func (f *memberFormatter) formatImplicitMemberExpression(ctx *fql.ImplicitMemberExpressionContext) {
	if ctx == nil {
		return
	}

	f.formatImplicitMemberExpressionStart(ctx.ImplicitMemberExpressionStart().(*fql.ImplicitMemberExpressionStartContext))

	for _, path := range ctx.AllMemberExpressionPath() {
		f.formatMemberExpressionPath(path.(*fql.MemberExpressionPathContext))
	}
}

func (f *memberFormatter) formatImplicitMemberExpressionStart(ctx *fql.ImplicitMemberExpressionStartContext) {
	if ctx == nil {
		return
	}

	if ctx.ErrorOperator() != nil {
		f.p.write("?")
	}

	f.p.write(".")

	if ctx.PropertyName() != nil {
		f.literal.formatPropertyNameWith(f.p, ctx.PropertyName().(*fql.PropertyNameContext))
		return
	}

	if ctx.ComputedPropertyName() != nil {
		f.literal.formatComputedPropertyNameWith(f.p, ctx.ComputedPropertyName().(*fql.ComputedPropertyNameContext))
	}
}

func (f *memberFormatter) formatMemberExpressionSource(ctx *fql.MemberExpressionSourceContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.Variable() != nil:
		f.expression.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.Param() != nil:
		f.expression.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.ArrayLiteral() != nil:
		f.list.formatArrayLiteral(ctx.ArrayLiteral().(*fql.ArrayLiteralContext))
	case ctx.ObjectLiteral() != nil:
		f.list.formatObjectLiteral(ctx.ObjectLiteral().(*fql.ObjectLiteralContext))
	case ctx.FunctionCall() != nil:
		f.expression.formatFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext))
	case ctx.OpenParen() != nil:
		f.p.write("(")

		if fe := ctx.ForExpression(); fe != nil {
			f.p.newline()
			f.p.withIndent(func() {
				f.statement.formatForExpression(fe.(*fql.ForExpressionContext))
			})
			f.p.newline()
			f.p.write(")")

			return
		}
		if we := ctx.WaitForExpression(); we != nil {
			f.statement.formatWaitForExpression(we.(*fql.WaitForExpressionContext))
			f.p.write(")")

			return
		}
		if expr := ctx.Expression(); expr != nil {
			f.expression.formatExpression(expr.(*fql.ExpressionContext))
		}

		f.p.write(")")
	}
}

func (f *memberFormatter) formatMemberExpressionPath(ctx *fql.MemberExpressionPathContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.PropertyName() != nil:
		if ctx.ErrorOperator() != nil {
			f.p.write("?")
		}

		f.p.write(".")
		f.literal.formatPropertyNameWith(f.p, ctx.PropertyName().(*fql.PropertyNameContext))
	case ctx.ComputedPropertyName() != nil:
		if ctx.ErrorOperator() != nil {
			f.p.write("?")
			f.p.write(".")
		}

		f.literal.formatComputedPropertyNameWith(f.p, ctx.ComputedPropertyName().(*fql.ComputedPropertyNameContext))
	case ctx.ArrayContraction() != nil:
		f.formatArrayContraction(ctx.ArrayContraction().(*fql.ArrayContractionContext))
	case ctx.ArrayExpansion() != nil:
		f.formatArrayExpansion(ctx.ArrayExpansion().(*fql.ArrayExpansionContext))
	case ctx.ArrayQuestionMark() != nil:
		f.formatArrayQuestionMark(ctx.ArrayQuestionMark().(*fql.ArrayQuestionMarkContext))
	case ctx.ArrayApply() != nil:
		f.formatArrayApply(ctx.ArrayApply().(*fql.ArrayApplyContext))
	}
}

func (f *memberFormatter) formatArrayExpansion(ctx *fql.ArrayExpansionContext) {
	if ctx == nil {
		return
	}

	f.p.write("[")
	f.p.write("*")

	if inline := ctx.InlineExpression(); inline != nil {
		f.p.space()
		f.formatInlineExpression(inline.(*fql.InlineExpressionContext))
	}

	f.p.write("]")
}

func (f *memberFormatter) formatArrayContraction(ctx *fql.ArrayContractionContext) {
	if ctx == nil {
		return
	}

	f.p.write("[")

	stars := ctx.AllMulti()
	for range stars {
		f.p.write("*")
	}

	if inline := ctx.InlineExpression(); inline != nil {
		f.p.space()
		f.formatInlineExpression(inline.(*fql.InlineExpressionContext))
	}

	f.p.write("]")
}

func (f *memberFormatter) formatArrayQuestionMark(ctx *fql.ArrayQuestionMarkContext) {
	if ctx == nil {
		return
	}

	f.p.write("[")
	f.p.write("?")

	if quant := ctx.ArrayQuestionQuantifier(); quant != nil {
		f.p.space()
		f.formatArrayQuestionQuantifier(quant.(*fql.ArrayQuestionQuantifierContext))
	}

	f.p.space()
	f.writeKeyword(keywordFilter)
	f.p.space()

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
	}

	f.p.write("]")
}

func (f *memberFormatter) formatArrayQuestionQuantifier(ctx *fql.ArrayQuestionQuantifierContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.Any() != nil:
		f.writeKeyword(keywordAny)
	case ctx.All() != nil:
		f.writeKeyword(keywordAll)
	case ctx.None() != nil:
		f.writeKeyword(keywordNone)
	case ctx.At() != nil && ctx.Least() != nil:
		f.writeKeyword(keywordAt)
		f.p.space()
		f.writeKeyword(keywordLeast)
		f.p.space()
		f.p.write("(")

		if val := ctx.ArrayQuestionQuantifierValue(0); val != nil {
			f.formatArrayQuestionQuantifierValue(val.(*fql.ArrayQuestionQuantifierValueContext))
		}

		f.p.write(")")
	case ctx.Range() != nil:
		f.formatArrayQuestionQuantifierValue(ctx.ArrayQuestionQuantifierValue(0).(*fql.ArrayQuestionQuantifierValueContext))
		f.p.write("..")
		f.formatArrayQuestionQuantifierValue(ctx.ArrayQuestionQuantifierValue(1).(*fql.ArrayQuestionQuantifierValueContext))
	case ctx.ArrayQuestionQuantifierValue(0) != nil:
		f.formatArrayQuestionQuantifierValue(ctx.ArrayQuestionQuantifierValue(0).(*fql.ArrayQuestionQuantifierValueContext))
	}
}

func (f *memberFormatter) formatArrayQuestionQuantifierValue(ctx *fql.ArrayQuestionQuantifierValueContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.IntegerLiteral() != nil:
		f.p.write(ctx.IntegerLiteral().GetText())
	case ctx.Param() != nil:
		f.expression.formatParam(ctx.Param().(*fql.ParamContext))
	}
}

func (f *memberFormatter) formatArrayApply(ctx *fql.ArrayApplyContext) {
	if ctx == nil {
		return
	}

	f.p.write("[")
	f.p.write("~")
	f.p.space()

	if q := ctx.QueryLiteral(); q != nil {
		f.formatQueryLiteral(q.(*fql.QueryLiteralContext))
	}

	f.p.write("]")
}

func (f *memberFormatter) formatInlineExpression(ctx *fql.InlineExpressionContext) {
	if ctx == nil {
		return
	}

	if filter := ctx.InlineFilter(); filter != nil {
		f.formatInlineFilter(filter.(*fql.InlineFilterContext))
	}

	if limit := ctx.InlineLimit(); limit != nil {
		if ctx.InlineFilter() != nil {
			f.p.space()
		}

		f.formatInlineLimit(limit.(*fql.InlineLimitContext))
	}

	if ret := ctx.InlineReturn(); ret != nil {
		if ctx.InlineFilter() != nil || ctx.InlineLimit() != nil {
			f.p.space()
		}

		f.formatInlineReturn(ret.(*fql.InlineReturnContext))
	}
}

func (f *memberFormatter) formatInlineFilter(ctx *fql.InlineFilterContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordFilter)
	f.p.space()

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (f *memberFormatter) formatInlineLimit(ctx *fql.InlineLimitContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordLimit)
	f.p.space()
	values := ctx.AllLimitClauseValue()

	for i, val := range values {
		f.clause.formatLimitClauseValue(val.(*fql.LimitClauseValueContext))

		if i < len(values)-1 {
			f.p.write(",")
			f.p.space()
		}
	}
}

func (f *memberFormatter) formatInlineReturn(ctx *fql.InlineReturnContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword(keywordReturn)
	f.p.space()

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpression(expr.(*fql.ExpressionContext))
	}
}

func (f *memberFormatter) formatQueryLiteral(ctx *fql.QueryLiteralContext) {
	if ctx == nil {
		return
	}

	if id := ctx.Identifier(); id != nil {
		f.p.write(id.GetText())
	}

	if sl := ctx.StringLiteral(); sl != nil {
		f.p.space()
		f.literal.formatStringLiteralNode(sl)

		if ctx.OpenParen() != nil && ctx.Expression() != nil {
			f.p.write("(")
			f.expression.formatExpression(ctx.Expression().(*fql.ExpressionContext))
			f.p.write(")")
		}
	}
}
