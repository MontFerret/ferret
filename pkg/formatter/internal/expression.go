package internal

import "github.com/MontFerret/ferret/v2/pkg/parser/fql"

type expressionFormatter struct {
	*engine
}

func (f *expressionFormatter) formatExpression(ctx *fql.ExpressionContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.UnaryOperator() != nil:
		f.formatUnaryOperator(ctx.UnaryOperator().(*fql.UnaryOperatorContext))
		f.formatExpression(ctx.GetRight().(*fql.ExpressionContext))
	case ctx.LogicalAndOperator() != nil:
		f.formatExpression(ctx.GetLeft().(*fql.ExpressionContext))
		f.p.space()
		f.formatLogicalAndOperator(ctx.LogicalAndOperator().(*fql.LogicalAndOperatorContext))
		f.p.space()
		f.formatExpression(ctx.GetRight().(*fql.ExpressionContext))
	case ctx.LogicalOrOperator() != nil:
		f.formatExpression(ctx.GetLeft().(*fql.ExpressionContext))
		f.p.space()
		f.formatLogicalOrOperator(ctx.LogicalOrOperator().(*fql.LogicalOrOperatorContext))
		f.p.space()
		f.formatExpression(ctx.GetRight().(*fql.ExpressionContext))
	case ctx.GetTernaryOperator() != nil:
		f.formatExpression(ctx.GetCondition().(*fql.ExpressionContext))
		f.p.space()
		f.p.write("?")
		f.p.space()

		if ctx.GetOnTrue() != nil {
			f.formatExpression(ctx.GetOnTrue().(*fql.ExpressionContext))
		}

		f.p.space()
		f.p.write(":")
		f.p.space()
		f.formatExpression(ctx.GetOnFalse().(*fql.ExpressionContext))
	case ctx.Predicate() != nil:
		f.formatPredicate(ctx.Predicate().(*fql.PredicateContext))
	}
}

func (f *expressionFormatter) formatPredicate(ctx *fql.PredicateContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.EqualityOperator() != nil:
		f.formatPredicate(ctx.GetLeft().(*fql.PredicateContext))
		f.p.space()
		f.formatEqualityOperator(ctx.EqualityOperator().(*fql.EqualityOperatorContext))
		f.p.space()
		f.formatPredicate(ctx.GetRight().(*fql.PredicateContext))
	case ctx.ArrayOperator() != nil:
		f.formatPredicate(ctx.GetLeft().(*fql.PredicateContext))
		f.p.space()
		f.formatArrayOperator(ctx.ArrayOperator().(*fql.ArrayOperatorContext))
		f.p.space()
		f.formatPredicate(ctx.GetRight().(*fql.PredicateContext))
	case ctx.InOperator() != nil:
		f.formatPredicate(ctx.GetLeft().(*fql.PredicateContext))
		f.p.space()
		f.formatInOperator(ctx.InOperator().(*fql.InOperatorContext))
		f.p.space()
		f.formatPredicate(ctx.GetRight().(*fql.PredicateContext))
	case ctx.LikeOperator() != nil:
		f.formatPredicate(ctx.GetLeft().(*fql.PredicateContext))
		f.p.space()
		f.formatLikeOperator(ctx.LikeOperator().(*fql.LikeOperatorContext))
		f.p.space()
		f.formatPredicate(ctx.GetRight().(*fql.PredicateContext))
	case ctx.ExpressionAtom() != nil:
		f.formatExpressionAtom(ctx.ExpressionAtom().(*fql.ExpressionAtomContext))
	}
}

func (f *expressionFormatter) formatExpressionAtom(ctx *fql.ExpressionAtomContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.MultiplicativeOperator() != nil:
		f.formatExpressionAtom(ctx.GetLeft().(*fql.ExpressionAtomContext))
		f.p.space()
		f.formatMultiplicativeOperator(ctx.MultiplicativeOperator().(*fql.MultiplicativeOperatorContext))
		f.p.space()
		f.formatExpressionAtom(ctx.GetRight().(*fql.ExpressionAtomContext))
	case ctx.AdditiveOperator() != nil:
		f.formatExpressionAtom(ctx.GetLeft().(*fql.ExpressionAtomContext))
		f.p.space()
		f.formatAdditiveOperator(ctx.AdditiveOperator().(*fql.AdditiveOperatorContext))
		f.p.space()
		f.formatExpressionAtom(ctx.GetRight().(*fql.ExpressionAtomContext))
	case ctx.RegexpOperator() != nil:
		f.formatExpressionAtom(ctx.GetLeft().(*fql.ExpressionAtomContext))
		f.p.space()
		f.formatRegexpOperator(ctx.RegexpOperator().(*fql.RegexpOperatorContext))
		f.p.space()
		f.formatExpressionAtom(ctx.GetRight().(*fql.ExpressionAtomContext))
	case ctx.FunctionCallExpression() != nil:
		f.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.RangeOperator() != nil:
		f.formatRangeOperator(ctx.RangeOperator().(*fql.RangeOperatorContext))
	case ctx.Literal() != nil:
		f.literal.formatLiteral(ctx.Literal().(*fql.LiteralContext))
	case ctx.Variable() != nil:
		f.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.ImplicitCurrentExpression() != nil:
		f.p.write(".")
	case ctx.ImplicitMemberExpression() != nil:
		f.member.formatImplicitMemberExpression(ctx.ImplicitMemberExpression().(*fql.ImplicitMemberExpressionContext))
	case ctx.MemberExpression() != nil:
		f.member.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	case ctx.Param() != nil:
		f.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.DispatchExpression() != nil:
		f.statement.formatDispatchExpression(ctx.DispatchExpression().(*fql.DispatchExpressionContext))
	case ctx.WaitForExpression() != nil:
		f.statement.formatWaitForExpression(ctx.WaitForExpression().(*fql.WaitForExpressionContext))
	case ctx.OpenParen() != nil:
		f.formatParenthesizedExpression(ctx)
	}
}

func (f *expressionFormatter) formatParenthesizedExpression(ctx *fql.ExpressionAtomContext) {
	if ctx == nil {
		return
	}

	f.p.write("(")

	if fe := ctx.ForExpression(); fe != nil {
		f.p.newline()
		f.p.withIndent(func() {
			f.statement.formatForExpression(fe.(*fql.ForExpressionContext))
		})
		f.p.newline()
		f.p.write(")")

		if ctx.ErrorOperator() != nil {
			f.p.write("?")
		}

		return
	}

	if we := ctx.WaitForExpression(); we != nil {
		f.statement.formatWaitForExpression(we.(*fql.WaitForExpressionContext))
		f.p.write(")")

		if ctx.ErrorOperator() != nil {
			f.p.write("?")
		}

		return
	}

	if expr := ctx.Expression(); expr != nil {
		f.formatExpression(expr.(*fql.ExpressionContext))
	}

	f.p.write(")")

	if ctx.ErrorOperator() != nil {
		f.p.write("?")
	}
}

func (f *expressionFormatter) formatUnaryOperator(ctx *fql.UnaryOperatorContext) {
	if ctx == nil {
		return
	}

	op := ctx.GetText()
	if op == keywordNot || op == "!" {
		if op == "!" {
			f.p.write(op)
		} else {
			f.writeKeyword(keywordNot)
		}
		f.p.space()

		return
	}

	f.p.write(op)
}

func (f *expressionFormatter) formatLogicalAndOperator(ctx *fql.LogicalAndOperatorContext) {
	if ctx == nil {
		return
	}

	f.p.write(applyCase(f.opts.caseMode, ctx.GetText()))
}

func (f *expressionFormatter) formatLogicalOrOperator(ctx *fql.LogicalOrOperatorContext) {
	if ctx == nil {
		return
	}

	f.p.write(applyCase(f.opts.caseMode, ctx.GetText()))
}

func (f *expressionFormatter) formatEqualityOperator(ctx *fql.EqualityOperatorContext) {
	if ctx == nil {
		return
	}

	f.p.write(ctx.GetText())
}

func (f *expressionFormatter) formatArrayOperator(ctx *fql.ArrayOperatorContext) {
	if ctx == nil {
		return
	}

	if op := ctx.GetOperator(); op != nil {
		f.p.write(applyCase(f.opts.caseMode, op.GetText()))
	}

	f.p.space()

	if in := ctx.InOperator(); in != nil {
		f.formatInOperator(in.(*fql.InOperatorContext))
	} else if eq := ctx.EqualityOperator(); eq != nil {
		f.formatEqualityOperator(eq.(*fql.EqualityOperatorContext))
	}
}

func (f *expressionFormatter) formatInOperator(ctx *fql.InOperatorContext) {
	if ctx == nil {
		return
	}

	if ctx.Not() != nil {
		f.p.write(applyCase(f.opts.caseMode, ctx.Not().GetText()))
		f.p.space()
	}

	f.p.write(applyCase(f.opts.caseMode, ctx.In().GetText()))
}

func (f *expressionFormatter) formatLikeOperator(ctx *fql.LikeOperatorContext) {
	if ctx == nil {
		return
	}

	if ctx.Not() != nil {
		f.p.write(applyCase(f.opts.caseMode, ctx.Not().GetText()))
		f.p.space()
	}

	f.p.write(applyCase(f.opts.caseMode, ctx.Like().GetText()))
}

func (f *expressionFormatter) formatMultiplicativeOperator(ctx *fql.MultiplicativeOperatorContext) {
	if ctx == nil {
		return
	}

	f.p.write(ctx.GetText())
}

func (f *expressionFormatter) formatAdditiveOperator(ctx *fql.AdditiveOperatorContext) {
	if ctx == nil {
		return
	}

	f.p.write(ctx.GetText())
}

func (f *expressionFormatter) formatRegexpOperator(ctx *fql.RegexpOperatorContext) {
	if ctx == nil {
		return
	}

	f.p.write(ctx.GetText())
}

func (f *expressionFormatter) formatRangeOperator(ctx *fql.RangeOperatorContext) {
	if ctx == nil {
		return
	}

	if left := ctx.GetLeft(); left != nil {
		f.formatRangeOperand(left.(*fql.RangeOperandContext))
	}

	f.p.write("..")

	if right := ctx.GetRight(); right != nil {
		f.formatRangeOperand(right.(*fql.RangeOperandContext))
	}
}

func (f *expressionFormatter) formatRangeOperand(ctx *fql.RangeOperandContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.IntegerLiteral() != nil:
		f.p.write(ctx.IntegerLiteral().GetText())
	case ctx.Variable() != nil:
		f.formatVariable(ctx.Variable().(*fql.VariableContext))
	case ctx.Param() != nil:
		f.formatParam(ctx.Param().(*fql.ParamContext))
	case ctx.FunctionCallExpression() != nil:
		f.formatFunctionCallExpression(ctx.FunctionCallExpression().(*fql.FunctionCallExpressionContext))
	case ctx.ImplicitMemberExpression() != nil:
		f.member.formatImplicitMemberExpression(ctx.ImplicitMemberExpression().(*fql.ImplicitMemberExpressionContext))
	case ctx.MemberExpression() != nil:
		f.member.formatMemberExpression(ctx.MemberExpression().(*fql.MemberExpressionContext))
	}
}

func (f *expressionFormatter) formatFunctionCallExpression(ctx *fql.FunctionCallExpressionContext) {
	if ctx == nil {
		return
	}

	f.formatFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext))

	if ctx.ErrorOperator() != nil {
		f.p.write("?")
	}
}

func (f *expressionFormatter) formatFunctionCall(ctx *fql.FunctionCallContext) {
	if ctx == nil {
		return
	}

	if ns := ctx.Namespace(); ns != nil {
		f.p.write(ns.GetText())
	}

	if fn := ctx.FunctionName(); fn != nil {
		f.p.write(fn.GetText())
	}

	f.p.write("(")

	if args := ctx.ArgumentList(); args != nil {
		f.list.formatArgumentList(args.(*fql.ArgumentListContext))
	}

	f.p.write(")")
}

func (f *expressionFormatter) formatParam(ctx *fql.ParamContext) {
	f.formatParamWith(f.p, ctx)
}

func (f *expressionFormatter) formatParamWith(p *printer, ctx *fql.ParamContext) {
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

func (f *expressionFormatter) formatVariable(ctx *fql.VariableContext) {
	f.formatVariableWith(f.p, ctx)
}

func (f *expressionFormatter) formatVariableWith(p *printer, ctx *fql.VariableContext) {
	if ctx == nil {
		return
	}

	if id := ctx.Identifier(); id != nil {
		p.write(id.GetText())
	} else if id := ctx.SafeReservedWord(); id != nil {
		p.write(id.GetText())
	}
}

func (f *expressionFormatter) formatExpressionWith(p *printer, ctx *fql.ExpressionContext) {
	if p == f.p {
		f.formatExpression(ctx)
		return
	}

	orig := f.p
	f.p = p
	f.formatExpression(ctx)
	f.p = orig
}
