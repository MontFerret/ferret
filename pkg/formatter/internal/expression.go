package internal

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

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
	case ctx.MatchExpression() != nil:
		f.formatMatchExpression(ctx.MatchExpression().(*fql.MatchExpressionContext))
	case ctx.QueryExpression() != nil:
		f.formatQueryExpression(ctx.QueryExpression().(*fql.QueryExpressionContext))
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

func (f *expressionFormatter) formatMatchExpression(ctx *fql.MatchExpressionContext) {
	if ctx == nil {
		return
	}

	hasComments := f.matchArmsHaveComments(ctx)

	if f.p.forceSingleLine {
		if hasComments {
			f.formatMatchExpressionWith(f.p, ctx, false)
			return
		}

		f.formatMatchExpressionWith(f.p, ctx, true)
		return
	}

	if !hasComments {
		inline, ok := f.renderInline(func(p *printer) {
			f.formatMatchExpressionWith(p, ctx, true)
		})

		if ok && f.inlineFits(inline) {
			f.p.write(inline)
			return
		}
	}

	f.formatMatchExpressionWith(f.p, ctx, false)
}

func (f *expressionFormatter) formatMatchExpressionWith(p *printer, ctx *fql.MatchExpressionContext, inline bool) {
	if ctx == nil {
		return
	}

	f.writeKeywordWith(p, keywordMatch)

	if expr := ctx.Expression(); expr != nil {
		p.space()
		f.formatExpressionWith(p, expr.(*fql.ExpressionContext))
	}

	p.space()
	p.write("(")

	if inline {
		p.space()
		f.formatMatchArmsInline(p, ctx)
		p.space()
		p.write(")")
		return
	}

	p.newline()
	p.withIndent(func() {
		f.formatMatchArmsMultiline(p, ctx)
	})
	p.write(")")
}

func (f *expressionFormatter) formatMatchArmsInline(p *printer, ctx *fql.MatchExpressionContext) {
	if ctx == nil {
		return
	}

	if arms := ctx.MatchPatternArms(); arms != nil {
		list := arms.MatchPatternArmList()
		if list != nil {
			armList := list.AllMatchPatternArm()
			for i, arm := range armList {
				f.formatMatchPatternArmWith(p, arm.(*fql.MatchPatternArmContext))
				if i < len(armList)-1 || arms.MatchDefaultArm() != nil {
					p.write(",")
					p.space()
				}
			}
		}

		if def := arms.MatchDefaultArm(); def != nil {
			f.formatMatchDefaultArmWith(p, def.(*fql.MatchDefaultArmContext))
		}

		return
	}

	if arms := ctx.MatchGuardArms(); arms != nil {
		list := arms.MatchGuardArmList()
		if list != nil {
			armList := list.AllMatchGuardArm()
			for i, arm := range armList {
				f.formatMatchGuardArmWith(p, arm.(*fql.MatchGuardArmContext))
				if i < len(armList)-1 || arms.MatchDefaultArm() != nil {
					p.write(",")
					p.space()
				}
			}
		}

		if def := arms.MatchDefaultArm(); def != nil {
			f.formatMatchDefaultArmWith(p, def.(*fql.MatchDefaultArmContext))
		}
	}
}

func (f *expressionFormatter) formatMatchArmsMultiline(p *printer, ctx *fql.MatchExpressionContext) {
	if ctx == nil {
		return
	}

	arms, closeParen := f.matchArmContexts(ctx)
	if len(arms) == 0 {
		return
	}

	closeStart := f.trivia.tokenStart(closeParen)

	for i, arm := range arms {
		switch v := arm.(type) {
		case *fql.MatchPatternArmContext:
			f.formatMatchPatternArmWith(p, v)
		case *fql.MatchGuardArmContext:
			f.formatMatchGuardArmWith(p, v)
		case *fql.MatchDefaultArmContext:
			f.formatMatchDefaultArmWith(p, v)
		}

		p.write(",")

		nextStart := closeStart
		if i < len(arms)-1 {
			nextStart = f.trivia.startIndex(arms[i+1])
		}

		f.trivia.emitListTriviaWith(p, f.trivia.sliceBetween(f.trivia.stopIndex(arm)+1, nextStart))
	}
}

func (f *expressionFormatter) matchArmsHaveComments(ctx *fql.MatchExpressionContext) bool {
	arms, closeParen := f.matchArmContexts(ctx)
	if len(arms) == 0 {
		return false
	}

	closeStart := f.trivia.tokenStart(closeParen)
	for i, arm := range arms {
		start := f.trivia.stopIndex(arm) + 1
		end := closeStart
		if i < len(arms)-1 {
			end = f.trivia.startIndex(arms[i+1])
		}

		if f.trivia.containsComment(f.trivia.sliceBetween(start, end)) {
			return true
		}
	}

	return false
}

func (f *expressionFormatter) matchArmContexts(ctx *fql.MatchExpressionContext) ([]antlr.ParserRuleContext, antlr.TerminalNode) {
	if ctx == nil {
		return nil, nil
	}

	if arms := ctx.MatchPatternArms(); arms != nil {
		out := make([]antlr.ParserRuleContext, 0, 4)
		if list := arms.MatchPatternArmList(); list != nil {
			for _, arm := range list.AllMatchPatternArm() {
				out = append(out, arm.(antlr.ParserRuleContext))
			}
		}
		if def := arms.MatchDefaultArm(); def != nil {
			out = append(out, def.(antlr.ParserRuleContext))
		}
		return out, arms.CloseParen()
	}

	if arms := ctx.MatchGuardArms(); arms != nil {
		out := make([]antlr.ParserRuleContext, 0, 4)
		if list := arms.MatchGuardArmList(); list != nil {
			for _, arm := range list.AllMatchGuardArm() {
				out = append(out, arm.(antlr.ParserRuleContext))
			}
		}
		if def := arms.MatchDefaultArm(); def != nil {
			out = append(out, def.(antlr.ParserRuleContext))
		}
		return out, arms.CloseParen()
	}

	return nil, nil
}

func (f *expressionFormatter) formatMatchPatternArmWith(p *printer, ctx *fql.MatchPatternArmContext) {
	if ctx == nil {
		return
	}

	if pattern := ctx.MatchPattern(); pattern != nil {
		f.formatMatchPatternWith(p, pattern.(*fql.MatchPatternContext))
	}

	if guard := ctx.MatchPatternGuard(); guard != nil {
		p.space()
		f.writeKeywordWith(p, keywordWhen)
		p.space()
		if expr := guard.Expression(); expr != nil {
			f.formatExpressionWith(p, expr.(*fql.ExpressionContext))
		}
	}

	p.space()
	p.write("=>")
	p.space()

	if expr := ctx.Expression(); expr != nil {
		f.formatExpressionWith(p, expr.(*fql.ExpressionContext))
	}
}

func (f *expressionFormatter) formatMatchGuardArmWith(p *printer, ctx *fql.MatchGuardArmContext) {
	if ctx == nil {
		return
	}

	f.writeKeywordWith(p, keywordWhen)
	p.space()

	exprs := ctx.AllExpression()
	if len(exprs) > 0 {
		f.formatExpressionWith(p, exprs[0].(*fql.ExpressionContext))
	}

	p.space()
	p.write("=>")
	p.space()

	if len(exprs) > 1 {
		f.formatExpressionWith(p, exprs[1].(*fql.ExpressionContext))
	}
}

func (f *expressionFormatter) formatMatchDefaultArmWith(p *printer, ctx *fql.MatchDefaultArmContext) {
	if ctx == nil {
		return
	}

	p.write("_")
	p.space()
	p.write("=>")
	p.space()
	if expr := ctx.Expression(); expr != nil {
		f.formatExpressionWith(p, expr.(*fql.ExpressionContext))
	}
}

func (f *expressionFormatter) formatMatchPatternWith(p *printer, ctx *fql.MatchPatternContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.MatchLiteralPattern() != nil:
		f.formatMatchLiteralPatternWith(p, ctx.MatchLiteralPattern().(*fql.MatchLiteralPatternContext))
	case ctx.MatchBindingPattern() != nil:
		f.formatMatchBindingPatternWith(p, ctx.MatchBindingPattern().(*fql.MatchBindingPatternContext))
	case ctx.MatchObjectPattern() != nil:
		f.formatMatchObjectPatternWith(p, ctx.MatchObjectPattern().(*fql.MatchObjectPatternContext))
	}
}

func (f *expressionFormatter) formatMatchLiteralPatternWith(p *printer, ctx *fql.MatchLiteralPatternContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.NoneLiteral() != nil:
		if nl := ctx.NoneLiteral().(*fql.NoneLiteralContext); nl != nil {
			if nl.Null() != nil {
				p.write(applyCase(f.opts.caseMode, nl.Null().GetText()))
			} else if nl.None() != nil {
				p.write(applyCase(f.opts.caseMode, nl.None().GetText()))
			}
		}
	case ctx.BooleanLiteral() != nil:
		if bl := ctx.BooleanLiteral().(*fql.BooleanLiteralContext); bl != nil && bl.BooleanLiteral() != nil {
			p.write(applyCase(f.opts.caseMode, bl.BooleanLiteral().GetText()))
		}
	case ctx.StringLiteral() != nil:
		f.literal.formatStringLiteralNodeWith(p, ctx.StringLiteral())
	case ctx.FloatLiteral() != nil:
		p.write(ctx.FloatLiteral().GetText())
	case ctx.IntegerLiteral() != nil:
		p.write(ctx.IntegerLiteral().GetText())
	}
}

func (f *expressionFormatter) formatMatchBindingPatternWith(p *printer, ctx *fql.MatchBindingPatternContext) {
	if ctx == nil {
		return
	}

	if id := ctx.Identifier(); id != nil {
		p.write(id.GetText())
		return
	}

	if srw := ctx.SafeReservedWord(); srw != nil {
		p.write(srw.GetText())
	}
}

func (f *expressionFormatter) formatMatchObjectPatternWith(p *printer, ctx *fql.MatchObjectPatternContext) {
	if ctx == nil {
		return
	}

	props := ctx.AllMatchObjectPatternProperty()
	if len(props) == 0 {
		p.write("{}")
		return
	}

	if p.forceSingleLine {
		f.formatMatchObjectPatternWithMode(p, ctx, true)
		return
	}

	inline, ok := f.renderInline(func(out *printer) {
		f.formatMatchObjectPatternWithMode(out, ctx, true)
	})
	if ok && f.inlineFitsWith(p, inline) {
		p.write(inline)
		return
	}

	f.formatMatchObjectPatternWithMode(p, ctx, false)
}

func (f *expressionFormatter) formatMatchObjectPatternWithMode(p *printer, ctx *fql.MatchObjectPatternContext, inline bool) {
	props := ctx.AllMatchObjectPatternProperty()
	p.write("{")

	if inline {
		if f.opts.bracketSpacing {
			p.space()
		}

		for i, prop := range props {
			f.formatMatchObjectPatternPropertyWith(p, prop.(*fql.MatchObjectPatternPropertyContext))
			if i < len(props)-1 {
				p.write(",")
				p.space()
			}
		}

		if f.opts.bracketSpacing {
			p.space()
		}

		p.write("}")
		return
	}

	p.newline()
	p.withIndent(func() {
		closeStart := f.trivia.tokenStart(ctx.CloseBrace())

		for i, prop := range props {
			propCtx := prop.(*fql.MatchObjectPatternPropertyContext)
			f.formatMatchObjectPatternPropertyWith(p, propCtx)

			if i < len(props)-1 {
				p.write(",")
			}

			nextStart := closeStart
			if i < len(props)-1 {
				nextStart = f.trivia.startIndex(props[i+1].(antlr.ParserRuleContext))
			}

			f.trivia.emitListTriviaWith(p, f.trivia.sliceBetween(f.trivia.stopIndex(propCtx)+1, nextStart))
		}
	})

	p.write("}")
}

func (f *expressionFormatter) formatMatchObjectPatternPropertyWith(p *printer, ctx *fql.MatchObjectPatternPropertyContext) {
	if ctx == nil {
		return
	}

	if key := ctx.MatchObjectPatternKey(); key != nil {
		f.formatMatchObjectPatternKeyWith(p, key.(*fql.MatchObjectPatternKeyContext))
	}

	p.write(":")
	p.space()

	if pattern := ctx.MatchPattern(); pattern != nil {
		f.formatMatchPatternWith(p, pattern.(*fql.MatchPatternContext))
	}
}

func (f *expressionFormatter) formatMatchObjectPatternKeyWith(p *printer, ctx *fql.MatchObjectPatternKeyContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.Identifier() != nil:
		p.write(ctx.Identifier().GetText())
	case ctx.StringLiteral() != nil:
		f.literal.formatStringLiteralNodeWith(p, ctx.StringLiteral())
	case ctx.SafeReservedWord() != nil:
		p.write(ctx.SafeReservedWord().GetText())
	case ctx.UnsafeReservedWord() != nil:
		p.write(ctx.UnsafeReservedWord().GetText())
	}
}

func (f *expressionFormatter) formatQueryExpression(ctx *fql.QueryExpressionContext) {
	if ctx == nil {
		return
	}

	if ctx.QueryWithOpt() == nil || f.p.forceSingleLine {
		f.formatQueryExpressionWith(f.p, ctx, true)
		return
	}

	inline, ok := f.renderInline(func(p *printer) {
		f.formatQueryExpressionWith(p, ctx, true)
	})

	if ok && f.inlineFits(inline) {
		f.p.write(inline)
		return
	}

	f.formatQueryExpressionWith(f.p, ctx, false)
}

func (f *expressionFormatter) formatQueryExpressionWith(p *printer, ctx *fql.QueryExpressionContext, inline bool) {
	if ctx == nil {
		return
	}

	f.writeKeywordWith(p, keywordQuery)
	p.space()
	f.writeQueryModifierWith(p, ctx.QueryModifier())

	if payload := ctx.QueryPayload(); payload != nil {
		if lit := payload.StringLiteral(); lit != nil {
			f.literal.formatStringLiteralNodeWith(p, lit)
		} else if param := payload.Param(); param != nil {
			f.formatParamWith(p, param.(*fql.ParamContext))
		} else if variable := payload.Variable(); variable != nil {
			f.formatVariableWith(p, variable.(*fql.VariableContext))
		}
	}

	p.space()
	f.writeKeywordWith(p, keywordIn)
	p.space()

	if expr := ctx.Expression(); expr != nil {
		f.formatExpressionWith(p, expr.(*fql.ExpressionContext))
	}

	p.space()
	f.writeKeywordWith(p, keywordUsing)
	p.space()

	if id := ctx.GetDialect(); id != nil {
		p.write(id.GetText())
	}

	if with := ctx.QueryWithOpt(); with != nil {
		if expr := with.Expression(); expr != nil {
			if inline {
				p.space()
				f.writeKeywordWith(p, keywordWith)
				p.space()
				f.formatExpressionWith(p, expr.(*fql.ExpressionContext))
			} else {
				p.newline()
				p.withIndent(func() {
					f.writeKeywordWith(p, keywordWith)
					p.space()
					f.formatExpressionWith(p, expr.(*fql.ExpressionContext))
				})
			}
		}
	}
}

func (f *expressionFormatter) writeQueryModifierWith(p *printer, modifier fql.IQueryModifierContext) {
	if modifier == nil {
		return
	}

	switch text := modifier.GetText(); {
	case strings.EqualFold(text, keywordExists):
		f.writeKeywordWith(p, keywordExists)
	case strings.EqualFold(text, keywordCount):
		f.writeKeywordWith(p, keywordCount)
	case strings.EqualFold(text, keywordAny):
		f.writeKeywordWith(p, keywordAny)
	case strings.EqualFold(text, keywordValue):
		f.writeKeywordWith(p, keywordValue)
	case strings.EqualFold(text, keywordOne):
		f.writeKeywordWith(p, keywordOne)
	default:
		return
	}

	p.space()
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

func (f *expressionFormatter) writeKeywordWith(p *printer, val string) {
	if p == nil {
		return
	}

	p.write(applyCase(f.opts.caseMode, val))
}
