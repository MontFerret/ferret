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
	case ctx.GetCoalesceOperator() != nil:
		f.formatCoalesce(ctx)
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

func (f *expressionFormatter) formatCoalesce(ctx *fql.ExpressionContext) {
	if ctx == nil {
		return
	}

	if f.p.forceSingleLine {
		f.formatCoalesceInlineWith(f.p, ctx)

		return
	}

	inline, ok := f.renderInline(func(p *printer) {
		f.formatCoalesceInlineWith(p, ctx)
	})
	if ok && f.inlineFits(inline) {
		f.p.write(inline)

		return
	}

	operands := f.coalesceOperands(ctx)
	if len(operands) == 0 {
		return
	}

	f.formatExpression(operands[0])
	f.p.withIndent(func() {
		for _, operand := range operands[1:] {
			f.p.newline()
			f.p.write("??")
			f.p.space()
			f.formatExpression(operand)
		}
	})
}

func (f *expressionFormatter) formatCoalesceInlineWith(p *printer, ctx *fql.ExpressionContext) {
	if p == nil || ctx == nil {
		return
	}

	f.formatExpressionWith(p, ctx.GetLeft().(*fql.ExpressionContext))
	p.space()
	p.write("??")
	p.space()
	f.formatExpressionWith(p, ctx.GetRight().(*fql.ExpressionContext))
}

func (f *expressionFormatter) coalesceOperands(ctx *fql.ExpressionContext) []*fql.ExpressionContext {
	operands := make([]*fql.ExpressionContext, 0, 3)
	current := ctx

	for current != nil && current.GetCoalesceOperator() != nil {
		left, ok := current.GetLeft().(*fql.ExpressionContext)
		if !ok {
			return nil
		}

		operands = append(operands, left)

		right, ok := current.GetRight().(*fql.ExpressionContext)
		if !ok {
			return nil
		}

		if right.GetCoalesceOperator() == nil {
			operands = append(operands, right)

			return operands
		}

		current = right
	}

	return operands
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

	hasComments := f.matchHasComments(ctx)

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
	p.write("{")

	if inline {
		p.space()
		f.formatMatchArmsInline(p, ctx)
		p.space()
		p.write("}")

		return
	}

	arms, openBrace, _ := f.matchArmContexts(ctx)
	p.withIndent(func() {
		if len(arms) > 0 {
			headerStop := ctx.GetStart().GetStop()

			if expr := ctx.Expression(); expr != nil {
				headerStop = expr.GetStop().GetStop()
			}

			f.trivia.emitListTriviaWith(
				p,
				f.trivia.blockLeadingTrivia(headerStop, openBrace, f.trivia.startIndex(arms[0])),
			)
		} else {
			p.newline()
		}

		f.formatMatchArmsMultiline(p, ctx)
	})

	p.write("}")
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

	arms, _, closeBrace := f.matchArmContexts(ctx)
	if len(arms) == 0 {
		return
	}

	closeStart := f.trivia.tokenStart(closeBrace)

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

func (f *expressionFormatter) matchHasComments(ctx *fql.MatchExpressionContext) bool {
	arms, openBrace, closeBrace := f.matchArmContexts(ctx)
	if len(arms) == 0 {
		return false
	}

	headerStop := ctx.GetStart().GetStop()
	if expr := ctx.Expression(); expr != nil {
		headerStop = expr.GetStop().GetStop()
	}

	if f.trivia.containsComment(f.trivia.blockLeadingTrivia(headerStop, openBrace, f.trivia.startIndex(arms[0]))) {
		return true
	}

	closeStart := f.trivia.tokenStart(closeBrace)

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

func (f *expressionFormatter) matchArmContexts(ctx *fql.MatchExpressionContext) ([]antlr.ParserRuleContext, antlr.TerminalNode, antlr.TerminalNode) {
	if ctx == nil {
		return nil, nil, nil
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

		return out, arms.OpenBrace(), arms.CloseBrace()
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

		return out, arms.OpenBrace(), arms.CloseBrace()
	}

	return nil, nil, nil
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
	case ctx.DurationLiteral() != nil:
		p.write(ctx.DurationLiteral().GetText())
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

	if (ctx.QueryWithOpt() == nil && ctx.QueryOptionsOpt() == nil) || f.p.forceSingleLine {
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
		if expr := payload.Expression(); expr != nil {
			p.write("(")
			f.formatExpressionWith(p, expr.(*fql.ExpressionContext))
			p.write(")")
		} else if member := payload.MemberExpression(); member != nil {
			f.formatMemberExpressionWith(p, member.(*fql.MemberExpressionContext))
		} else if lit := payload.Literal(); lit != nil {
			f.literal.formatLiteralWith(p, lit.(*fql.LiteralContext))
		} else if call := payload.FunctionCallExpression(); call != nil {
			f.formatFunctionCallExpressionWith(p, call.(*fql.FunctionCallExpressionContext))
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

	if id := ctx.GetDialect(); id != nil {
		p.space()
		f.writeKeywordWith(p, keywordUsing)
		p.space()
		p.write(id.GetText())
	}

	if with := ctx.QueryWithOpt(); with != nil {
		if expr := with.Expression(); expr != nil {
			f.formatQueryClauseWith(p, keywordWith, expr, inline)
		}
	}

	if options := ctx.QueryOptionsOpt(); options != nil {
		if expr := options.Expression(); expr != nil {
			f.formatQueryClauseWith(p, keywordOptions, expr, inline)
		}
	}

	f.formatRecoveryTailsWith(p, ctx.RecoveryTails())
}

func (f *expressionFormatter) formatQueryClauseWith(p *printer, keyword string, expr fql.IExpressionContext, inline bool) {
	if inline {
		p.space()
		f.writeKeywordWith(p, keyword)
		p.space()
		f.formatExpressionWith(p, expr.(*fql.ExpressionContext))

		return
	}

	p.newline()
	p.withIndent(func() {
		f.writeKeywordWith(p, keyword)
		p.space()
		f.formatExpressionWith(p, expr.(*fql.ExpressionContext))
	})
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

		f.formatExpressionAtomErrorTail(ctx)

		return
	}

	if we := ctx.WaitForExpression(); we != nil {
		f.statement.formatWaitForExpression(we.(*fql.WaitForExpressionContext))
		f.p.write(")")

		f.formatExpressionAtomErrorTail(ctx)

		return
	}

	if expr := ctx.Expression(); expr != nil {
		f.formatExpression(expr.(*fql.ExpressionContext))
	}

	f.p.write(")")
	f.formatExpressionAtomErrorTail(ctx)
}

func (f *expressionFormatter) formatUnaryOperator(ctx *fql.UnaryOperatorContext) {
	f.formatUnaryOperatorWith(f.p, ctx)
}

func (f *expressionFormatter) formatUnaryOperatorWith(p *printer, ctx *fql.UnaryOperatorContext) {
	if ctx == nil {
		return
	}

	op := ctx.GetText()
	if op == keywordNot || op == "!" {
		if op == "!" {
			p.write(op)
		} else {
			f.writeKeywordWith(p, keywordNot)
		}

		p.space()

		return
	}

	p.write(op)
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
	case ctx.ImplicitCurrentExpression() != nil:
		f.p.write(".")
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

		return
	}

	f.formatRecoveryTails(ctx.RecoveryTails())
}

func (f *expressionFormatter) formatFunctionCallNoRecoveryExpression(ctx *fql.FunctionCallNoRecoveryExpressionContext) {
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

func (f *expressionFormatter) formatExpressionAtomErrorTail(ctx *fql.ExpressionAtomContext) {
	if ctx == nil {
		return
	}

	if ctx.ErrorOperator() != nil {
		f.p.write("?")

		return
	}

	f.formatRecoveryTails(ctx.RecoveryTails())
}

func (f *expressionFormatter) formatRecoveryTails(ctx fql.IRecoveryTailsContext) {
	f.formatRecoveryTailsWith(f.p, ctx)
}

func (f *expressionFormatter) formatRecoveryTailsWith(p *printer, ctx fql.IRecoveryTailsContext) {
	for _, tail := range f.orderedRecoveryTails(ctx) {
		f.formatRecoveryTailWith(p, tail)
	}
}

func (f *expressionFormatter) formatRecoveryTailsMultiline(ctx fql.IRecoveryTailsContext, previousStop int) {
	if ctx == nil {
		return
	}

	tails := ctx.AllRecoveryTail()
	if f.recoveryTailsContainComments(previousStop, tails) {
		for _, tail := range tails {
			f.trivia.emitClauseBoundary(previousStop+1, f.trivia.startIndex(tail))
			f.formatRecoveryTailWith(f.p, tail)
			previousStop = f.trivia.stopIndex(tail)
		}

		return
	}

	for _, tail := range f.orderedRecoveryTails(ctx) {
		f.p.newline()
		f.formatRecoveryTailWith(f.p, tail)
	}
}

func (f *expressionFormatter) recoveryTailsContainComments(previousStop int, tails []fql.IRecoveryTailContext) bool {
	for _, tail := range tails {
		if f.trivia.containsComment(f.trivia.sliceBetween(previousStop+1, f.trivia.startIndex(tail))) {
			return true
		}

		previousStop = f.trivia.stopIndex(tail)
	}

	return false
}

func (f *expressionFormatter) orderedRecoveryTails(ctx fql.IRecoveryTailsContext) []fql.IRecoveryTailContext {
	if ctx == nil {
		return nil
	}

	var timeoutTail fql.IRecoveryTailContext
	var errorTail fql.IRecoveryTailContext
	var otherTails []fql.IRecoveryTailContext

	for _, tail := range ctx.AllRecoveryTail() {
		cond := tail.RecoveryCondition()
		if cond == nil {
			otherTails = append(otherTails, tail)

			continue
		}

		switch {
		case cond.TimeoutKeyword() != nil || strings.EqualFold(cond.GetText(), keywordTimeout):
			if timeoutTail == nil {
				timeoutTail = tail
			}
		case cond.ErrorKeyword() != nil || strings.EqualFold(cond.GetText(), keywordError):
			if errorTail == nil {
				errorTail = tail
			}
		default:
			otherTails = append(otherTails, tail)
		}
	}

	ordered := make([]fql.IRecoveryTailContext, 0, len(otherTails)+2)
	for _, tail := range []fql.IRecoveryTailContext{timeoutTail, errorTail} {
		if tail != nil {
			ordered = append(ordered, tail)
		}
	}
	ordered = append(ordered, otherTails...)

	return ordered
}

func (f *expressionFormatter) formatRecoveryTailWith(p *printer, ctx fql.IRecoveryTailContext) {
	if ctx == nil {
		return
	}

	p.space()
	f.writeKeywordWith(p, keywordOn)

	if cond := ctx.RecoveryCondition(); cond != nil {
		p.space()

		switch {
		case cond.TimeoutKeyword() != nil || strings.EqualFold(cond.GetText(), keywordTimeout):
			f.writeKeywordWith(p, keywordTimeout)
		case cond.ErrorKeyword() != nil || strings.EqualFold(cond.GetText(), keywordError):
			f.writeKeywordWith(p, keywordError)
		default:
			p.write(applyCase(f.opts.caseMode, cond.GetText()))
		}
	}

	action := ctx.RecoveryAction()
	if action == nil {
		return
	}

	p.space()
	switch {
	case action.RecoveryRetryAction() != nil:
		f.formatRecoveryRetryActionWith(p, action.RecoveryRetryAction())
	case action.ReturnKeyword() != nil || strings.EqualFold(action.GetText(), keywordReturn):
		f.writeKeywordWith(p, keywordReturn)

		if expr := action.RecoveryReturnExpr(); expr != nil && expr.Expression() != nil {
			p.space()
			f.formatExpressionWith(p, expr.Expression().(*fql.ExpressionContext))
		}
	case action.FailKeyword() != nil || strings.EqualFold(action.GetText(), keywordFail):
		f.writeKeywordWith(p, keywordFail)
	default:
		p.write(applyCase(f.opts.caseMode, action.GetText()))
	}
}

func (f *expressionFormatter) formatRecoveryRetryActionWith(p *printer, ctx fql.IRecoveryRetryActionContext) {
	if ctx == nil {
		return
	}

	f.writeKeywordWith(p, keywordRetry)

	if count := ctx.RecoveryRetryCount(); count != nil {
		p.space()
		p.write(count.GetText())
	}

	if delayClause := ctx.RecoveryRetryDelayClause(); delayClause != nil {
		f.formatRecoveryRetryDelayClauseWith(p, delayClause)
	}

	for _, orClause := range ctx.AllRecoveryRetryOrClause() {
		f.formatRecoveryRetryOrClauseWith(p, orClause)
	}
}

func (f *expressionFormatter) formatRecoveryRetryDelayClauseWith(p *printer, ctx fql.IRecoveryRetryDelayClauseContext) {
	if ctx == nil {
		return
	}

	if ctx.DelayKeyword() != nil {
		p.space()
		f.writeKeywordWith(p, keywordDelay)

		if value := ctx.RecoveryRetryDelayValue(); value != nil {
			p.space()

			if unary := value.UnaryOperator(); unary != nil {
				f.formatUnaryOperatorWith(p, unary.(*fql.UnaryOperatorContext))
			}

			if predicate := value.Predicate(); predicate != nil {
				f.formatPredicateWith(p, predicate.(*fql.PredicateContext))
			}
		}
	}

	if backoff := ctx.RecoveryRetryBackoffClause(); backoff != nil {
		p.space()
		f.writeKeywordWith(p, keywordBackoff)

		if kind := backoff.RecoveryRetryBackoffKind(); kind != nil {
			p.space()

			switch {
			case kind.None() != nil || kind.Identifier() != nil:
				p.write(applyCase(f.opts.caseMode, kind.GetText()))
			default:
				p.write(kind.GetText())
			}
		}
	}
}

func (f *expressionFormatter) formatRecoveryRetryOrClauseWith(p *printer, ctx fql.IRecoveryRetryOrClauseContext) {
	if ctx == nil {
		return
	}

	p.space()
	f.writeKeywordWith(p, keywordOr)

	if action := ctx.RecoveryRetryFinalAction(); action != nil {
		p.space()

		switch {
		case action.FailKeyword() != nil || strings.EqualFold(action.GetText(), keywordFail):
			f.writeKeywordWith(p, keywordFail)
		case action.ReturnKeyword() != nil || strings.EqualFold(action.GetText(), keywordReturn):
			f.writeKeywordWith(p, keywordReturn)

			if expr := action.RecoveryReturnExpr(); expr != nil && expr.Expression() != nil {
				p.space()
				f.formatExpressionWith(p, expr.Expression().(*fql.ExpressionContext))
			}
		default:
			p.write(applyCase(f.opts.caseMode, action.GetText()))
		}
	}
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

func (f *expressionFormatter) formatPredicateWith(p *printer, ctx *fql.PredicateContext) {
	if p == f.p {
		f.formatPredicate(ctx)

		return
	}

	orig := f.p
	f.p = p
	f.formatPredicate(ctx)
	f.p = orig
}

func (f *expressionFormatter) formatFunctionCallExpressionWith(p *printer, ctx *fql.FunctionCallExpressionContext) {
	if p == f.p {
		f.formatFunctionCallExpression(ctx)

		return
	}

	orig := f.p
	f.p = p
	f.formatFunctionCallExpression(ctx)
	f.p = orig
}

func (f *expressionFormatter) formatMemberExpressionWith(p *printer, ctx *fql.MemberExpressionContext) {
	if p == f.p {
		f.member.formatMemberExpression(ctx)

		return
	}

	orig := f.p
	f.p = p
	f.member.formatMemberExpression(ctx)
	f.p = orig
}

func (f *expressionFormatter) writeKeywordWith(p *printer, val string) {
	if p == nil {
		return
	}

	p.write(applyCase(f.opts.caseMode, val))
}
