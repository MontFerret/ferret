package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type listFormatter struct {
	*engine
}

func (l *listFormatter) arrayHasComments(ctx *fql.ArrayLiteralContext) bool {
	if ctx == nil || ctx.ArgumentList() == nil {
		return false
	}

	args := ctx.ArgumentList().AllExpression()
	if len(args) == 0 {
		return false
	}

	closeStart := l.trivia.tokenStart(ctx.CloseBracket())

	for i, arg := range args {
		start := l.trivia.stopIndex(arg.(antlr.ParserRuleContext)) + 1
		end := closeStart

		if i < len(args)-1 {
			end = l.trivia.startIndex(args[i+1].(antlr.ParserRuleContext))
		}

		if l.trivia.containsComment(l.trivia.sliceBetween(start, end)) {
			return true
		}
	}

	return false
}

func (l *listFormatter) objectHasComments(ctx *fql.ObjectLiteralContext) bool {
	if ctx == nil {
		return false
	}

	props := ctx.AllPropertyAssignment()
	if len(props) == 0 {
		return false
	}

	closeStart := l.trivia.tokenStart(ctx.CloseBrace())

	for i, prop := range props {
		start := l.trivia.stopIndex(prop.(antlr.ParserRuleContext)) + 1
		end := closeStart

		if i < len(props)-1 {
			end = l.trivia.startIndex(props[i+1].(antlr.ParserRuleContext))
		}

		if l.trivia.containsComment(l.trivia.sliceBetween(start, end)) {
			return true
		}
	}

	return false
}

func (l *listFormatter) argumentListClose(ctx *fql.ArgumentListContext) antlr.TerminalNode {
	if ctx == nil {
		return nil
	}

	if parent, ok := ctx.GetParent().(*fql.FunctionCallContext); ok {
		return parent.CloseParen()
	}

	return nil
}

func (l *listFormatter) argumentListHasComments(ctx *fql.ArgumentListContext) bool {
	if ctx == nil {
		return false
	}

	args := ctx.AllExpression()

	if len(args) == 0 {
		return false
	}

	closeStart := l.trivia.tokenStart(l.argumentListClose(ctx))

	for i, arg := range args {
		start := l.trivia.stopIndex(arg.(antlr.ParserRuleContext)) + 1
		end := closeStart

		if i < len(args)-1 {
			end = l.trivia.startIndex(args[i+1].(antlr.ParserRuleContext))
		}

		if l.trivia.containsComment(l.trivia.sliceBetween(start, end)) {
			return true
		}
	}

	return false
}

func (l *listFormatter) formatArrayLiteral(ctx *fql.ArrayLiteralContext) {
	if ctx == nil {
		return
	}

	if ctx.ArgumentList() == nil {
		l.p.write("[]")

		return
	}

	hasComments := l.arrayHasComments(ctx)

	if l.p.forceSingleLine {
		if hasComments {
			l.formatArrayLiteralWith(l.p, ctx, false)

			return
		}

		l.formatArrayLiteralInline(ctx)

		return
	}

	if !hasComments {
		inline, ok := l.renderInline(func(p *printer) {
			l.formatArrayLiteralWith(p, ctx, true)
		})

		if ok && len(inline) <= int(l.opts.printWidth) {
			l.p.write(inline)

			return
		}
	}

	l.formatArrayLiteralWith(l.p, ctx, false)
}

func (l *listFormatter) formatArrayLiteralInline(ctx *fql.ArrayLiteralContext) {
	l.formatArrayLiteralWith(l.p, ctx, true)
}

func (l *listFormatter) formatArrayLiteralWith(p *printer, ctx *fql.ArrayLiteralContext, inline bool) {
	args := ctx.ArgumentList().AllExpression()
	p.write("[")

	if !inline {
		p.newline()
		p.withIndent(func() {
			closeStart := l.trivia.tokenStart(ctx.CloseBracket())

			for i, expr := range args {
				exprCtx := expr.(*fql.ExpressionContext)
				l.expression.formatExpressionWith(p, exprCtx)

				if i < len(args)-1 {
					p.write(",")
				}

				nextStart := closeStart

				if i < len(args)-1 {
					nextStart = l.trivia.startIndex(args[i+1].(antlr.ParserRuleContext))
				}

				l.trivia.emitListTriviaWith(p, l.trivia.sliceBetween(l.trivia.stopIndex(exprCtx)+1, nextStart))
			}
		})

		p.write("]")

		return
	}

	for i, expr := range args {
		l.expression.formatExpressionWith(p, expr.(*fql.ExpressionContext))

		if i < len(args)-1 {
			p.write(",")
			p.space()
		}
	}

	p.write("]")
}

func (l *listFormatter) formatObjectLiteral(ctx *fql.ObjectLiteralContext) {
	if ctx == nil {
		return
	}

	props := ctx.AllPropertyAssignment()
	if len(props) == 0 {
		l.p.write("{}")

		return
	}

	hasComments := l.objectHasComments(ctx)

	if l.p.forceSingleLine {
		if hasComments {
			l.formatObjectLiteralWith(l.p, ctx, false)
			return
		}

		l.formatObjectLiteralInline(ctx)

		return
	}

	if !hasComments {
		inline, ok := l.renderInline(func(p *printer) {
			l.formatObjectLiteralWith(p, ctx, true)
		})

		if ok && len(inline) <= int(l.opts.printWidth) {
			l.p.write(inline)

			return
		}
	}

	l.formatObjectLiteralWith(l.p, ctx, false)
}

func (l *listFormatter) formatObjectLiteralInline(ctx *fql.ObjectLiteralContext) {
	l.formatObjectLiteralWith(l.p, ctx, true)
}

func (l *listFormatter) formatObjectLiteralWith(p *printer, ctx *fql.ObjectLiteralContext, inline bool) {
	props := ctx.AllPropertyAssignment()
	p.write("{")

	if inline {
		if l.opts.bracketSpacing {
			p.space()
		}

		for i, prop := range props {
			l.literal.formatPropertyAssignmentWith(p, prop.(*fql.PropertyAssignmentContext))

			if i < len(props)-1 {
				p.write(",")
				p.space()
			}
		}

		if l.opts.bracketSpacing {
			p.space()
		}

		p.write("}")

		return
	}

	p.newline()
	p.withIndent(func() {
		closeStart := l.trivia.tokenStart(ctx.CloseBrace())

		for i, prop := range props {
			propCtx := prop.(*fql.PropertyAssignmentContext)
			l.literal.formatPropertyAssignmentWith(p, propCtx)

			if i < len(props)-1 {
				p.write(",")
			}

			nextStart := closeStart

			if i < len(props)-1 {
				nextStart = l.trivia.startIndex(props[i+1].(antlr.ParserRuleContext))
			}

			l.trivia.emitListTriviaWith(p, l.trivia.sliceBetween(l.trivia.stopIndex(propCtx)+1, nextStart))
		}
	})

	p.write("}")
}

func (l *listFormatter) formatArgumentList(ctx *fql.ArgumentListContext) {
	if ctx == nil {
		return
	}

	args := ctx.AllExpression()
	if len(args) == 0 {
		return
	}

	hasComments := l.argumentListHasComments(ctx)

	if l.p.forceSingleLine {
		if hasComments {
			l.formatArgumentListWith(l.p, ctx, false)

			return
		}

		l.formatArgumentListInline(ctx)

		return
	}

	if !hasComments {
		inline, ok := l.renderInline(func(p *printer) {
			l.formatArgumentListWith(p, ctx, true)
		})

		if ok && len(inline) <= int(l.opts.printWidth) {
			l.p.write(inline)

			return
		}
	}

	l.formatArgumentListWith(l.p, ctx, false)
}

func (l *listFormatter) formatArgumentListInline(ctx *fql.ArgumentListContext) {
	l.formatArgumentListWith(l.p, ctx, true)
}

func (l *listFormatter) formatArgumentListWith(p *printer, ctx *fql.ArgumentListContext, inline bool) {
	args := ctx.AllExpression()

	if inline {
		for i, arg := range args {
			l.expression.formatExpressionWith(p, arg.(*fql.ExpressionContext))

			if i < len(args)-1 {
				p.write(",")
				p.space()
			}
		}

		return
	}

	p.newline()
	p.withIndent(func() {
		closeStart := l.trivia.tokenStart(l.argumentListClose(ctx))

		for i, arg := range args {
			argCtx := arg.(*fql.ExpressionContext)
			l.expression.formatExpressionWith(p, argCtx)

			if i < len(args)-1 {
				p.write(",")
			}

			nextStart := closeStart

			if i < len(args)-1 {
				nextStart = l.trivia.startIndex(args[i+1].(antlr.ParserRuleContext))
			}

			l.trivia.emitListTriviaWith(p, l.trivia.sliceBetween(l.trivia.stopIndex(argCtx)+1, nextStart))
		}
	})
}
