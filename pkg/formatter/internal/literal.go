package internal

import "github.com/MontFerret/ferret/v2/pkg/parser/fql"

type literalFormatter struct {
	*engine
}

func (f *literalFormatter) formatLiteral(ctx *fql.LiteralContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.ArrayLiteral() != nil:
		f.list.formatArrayLiteral(ctx.ArrayLiteral().(*fql.ArrayLiteralContext))
	case ctx.ObjectLiteral() != nil:
		f.list.formatObjectLiteral(ctx.ObjectLiteral().(*fql.ObjectLiteralContext))
	case ctx.BooleanLiteral() != nil:
		f.formatBooleanLiteral(ctx.BooleanLiteral().(*fql.BooleanLiteralContext))
	case ctx.StringLiteral() != nil:
		f.formatStringLiteralNode(ctx.StringLiteral())
	case ctx.FloatLiteral() != nil:
		f.p.write(ctx.FloatLiteral().GetText())
	case ctx.IntegerLiteral() != nil:
		f.p.write(ctx.IntegerLiteral().GetText())
	case ctx.NoneLiteral() != nil:
		f.formatNoneLiteral(ctx.NoneLiteral().(*fql.NoneLiteralContext))
	}
}

func (f *literalFormatter) formatBooleanLiteral(ctx *fql.BooleanLiteralContext) {
	if ctx == nil || ctx.BooleanLiteral() == nil {
		return
	}

	f.p.write(applyCase(f.opts.caseMode, ctx.BooleanLiteral().GetText()))
}

func (f *literalFormatter) formatNoneLiteral(ctx *fql.NoneLiteralContext) {
	if ctx == nil {
		return
	}

	if ctx.Null() != nil {
		f.p.write(applyCase(f.opts.caseMode, ctx.Null().GetText()))

		return
	}

	if ctx.None() != nil {
		f.p.write(applyCase(f.opts.caseMode, ctx.None().GetText()))
	}
}

func (f *literalFormatter) formatStringLiteralNode(ctx fql.IStringLiteralContext) {
	f.formatStringLiteralNodeWith(f.p, ctx)
}

func (f *literalFormatter) formatStringLiteralNodeWith(p *printer, ctx fql.IStringLiteralContext) {
	if ctx == nil {
		return
	}

	if tmpl := ctx.TemplateLiteral(); tmpl != nil {
		f.formatTemplateLiteralWith(p, tmpl)

		return
	}

	if tok := ctx.StringLiteral(); tok != nil {
		p.write(formatStringLiteral(tok.GetText(), f.opts.singleQuote))
	}
}

func (f *literalFormatter) formatTemplateLiteral(ctx fql.ITemplateLiteralContext) {
	f.formatTemplateLiteralWith(f.p, ctx)
}

func (f *literalFormatter) formatTemplateLiteralWith(p *printer, ctx fql.ITemplateLiteralContext) {
	if ctx == nil {
		return
	}

	p.write("`")
	for _, el := range ctx.AllTemplateElement() {

		if el == nil {
			continue
		}

		if chunk := el.TemplateChars(); chunk != nil {
			p.writeRaw(chunk.GetText())

			continue
		}

		if expr := el.Expression(); expr != nil {
			p.writeRaw("${")
			f.expression.formatExpressionWith(p, expr.(*fql.ExpressionContext))
			p.writeRaw("}")
		}
	}

	p.write("`")
}

func (f *literalFormatter) formatPropertyAssignmentWith(p *printer, ctx *fql.PropertyAssignmentContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.PropertyName() != nil:
		f.formatPropertyNameWith(p, ctx.PropertyName().(*fql.PropertyNameContext))
		p.write(":")
		p.space()
		f.expression.formatExpressionWith(p, ctx.Expression().(*fql.ExpressionContext))
	case ctx.ComputedPropertyName() != nil:
		f.formatComputedPropertyNameWith(p, ctx.ComputedPropertyName().(*fql.ComputedPropertyNameContext))
		p.write(":")
		p.space()
		f.expression.formatExpressionWith(p, ctx.Expression().(*fql.ExpressionContext))
	case ctx.Variable() != nil:
		f.expression.formatVariableWith(p, ctx.Variable().(*fql.VariableContext))
	}
}

func (f *literalFormatter) formatPropertyNameWith(p *printer, ctx *fql.PropertyNameContext) {
	if ctx == nil {
		return
	}

	switch {
	case ctx.Identifier() != nil:
		p.write(ctx.Identifier().GetText())
	case ctx.StringLiteral() != nil:
		f.formatStringLiteralNodeWith(p, ctx.StringLiteral())
	case ctx.Param() != nil:
		f.expression.formatParamWith(p, ctx.Param().(*fql.ParamContext))
	case ctx.SafeReservedWord() != nil:
		p.write(ctx.SafeReservedWord().GetText())
	case ctx.UnsafeReservedWord() != nil:
		p.write(ctx.UnsafeReservedWord().GetText())
	}
}

func (f *literalFormatter) formatComputedPropertyNameWith(p *printer, ctx *fql.ComputedPropertyNameContext) {
	if ctx == nil {
		return
	}

	p.write("[")

	if expr := ctx.Expression(); expr != nil {
		f.expression.formatExpressionWith(p, expr.(*fql.ExpressionContext))
	}

	p.write("]")
}

func (f *literalFormatter) formatPropertyAssignment(ctx *fql.PropertyAssignmentContext) {
	f.formatPropertyAssignmentWith(f.p, ctx)
}

func (f *literalFormatter) formatPropertyName(ctx *fql.PropertyNameContext) {
	f.formatPropertyNameWith(f.p, ctx)
}

func (f *literalFormatter) formatComputedPropertyName(ctx *fql.ComputedPropertyNameContext) {
	f.formatComputedPropertyNameWith(f.p, ctx)
}
