package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type bindingPatternFormatter struct {
	*engine
}

func (f *bindingPatternFormatter) formatWith(p *printer, ctx *fql.BindingPatternContext) {
	if ctx == nil {
		return
	}

	if id := ctx.BindingIdentifier(); id != nil {
		p.write(id.GetText())

		return
	}

	if id := ctx.IgnoreIdentifier(); id != nil {
		p.write(id.GetText())

		return
	}

	if structured := ctx.StructuredBindingPattern(); structured != nil {
		f.formatStructuredWith(p, structured.(*fql.StructuredBindingPatternContext))
	}
}

func (f *bindingPatternFormatter) formatStructured(ctx *fql.StructuredBindingPatternContext) {
	if ctx == nil {
		return
	}

	f.formatStructuredWith(f.p, ctx)
}

func (f *bindingPatternFormatter) formatStructuredWith(p *printer, ctx *fql.StructuredBindingPatternContext) {
	if ctx == nil {
		return
	}

	if object := ctx.ObjectBindingPattern(); object != nil {
		f.formatObjectWith(p, object.(*fql.ObjectBindingPatternContext))

		return
	}

	if array := ctx.ArrayBindingPattern(); array != nil {
		f.formatArrayWith(p, array.(*fql.ArrayBindingPatternContext))
	}
}

func (f *bindingPatternFormatter) formatObjectWith(p *printer, ctx *fql.ObjectBindingPatternContext) {
	entries := ctx.AllObjectBindingEntry()
	if len(entries) == 0 {
		f.formatEmptyWithComments(p, ctx.OpenBrace(), ctx.CloseBrace(), "{", "}")

		return
	}

	if p.forceSingleLine {
		f.formatObjectMode(p, ctx, true)

		return
	}

	if !f.hasComments(ctx) {
		inline, ok := f.renderInline(func(out *printer) {
			f.formatObjectMode(out, ctx, true)
		})

		if ok && f.inlineFitsWith(p, inline) {
			p.write(inline)

			return
		}
	}

	f.formatObjectMode(p, ctx, false)
}

func (f *bindingPatternFormatter) formatObjectMode(p *printer, ctx *fql.ObjectBindingPatternContext, inline bool) {
	entries := ctx.AllObjectBindingEntry()
	p.write("{")

	if inline {
		if f.config.BracketSpacing() {
			p.space()
		}

		for index, entry := range entries {
			f.formatObjectEntryWith(p, entry.(*fql.ObjectBindingEntryContext))
			if index < len(entries)-1 {
				p.write(",")
				p.space()
			}
		}

		if f.config.BracketSpacing() {
			p.space()
		}

		p.write("}")

		return
	}

	f.formatMultilineEntries(
		p,
		ctx.OpenBrace(),
		ctx.CloseBrace(),
		f.objectBindingEntryRules(entries),
		func(index int) {
			f.formatObjectEntryWith(p, entries[index].(*fql.ObjectBindingEntryContext))
		},
	)
	p.write("}")
}

func (f *bindingPatternFormatter) formatObjectEntryWith(p *printer, ctx *fql.ObjectBindingEntryContext) {
	if ctx == nil || ctx.BindingIdentifier() == nil {
		return
	}

	p.write(ctx.BindingIdentifier().GetText())
	if nested := ctx.BindingPattern(); nested != nil {
		p.write(":")
		p.space()
		f.formatWith(p, nested.(*fql.BindingPatternContext))
	}
}

func (f *bindingPatternFormatter) formatArrayWith(p *printer, ctx *fql.ArrayBindingPatternContext) {
	entries := ctx.AllBindingPattern()
	if len(entries) == 0 {
		f.formatEmptyWithComments(p, ctx.OpenBracket(), ctx.CloseBracket(), "[", "]")

		return
	}

	if p.forceSingleLine {
		f.formatArrayMode(p, ctx, true)

		return
	}

	if !f.hasComments(ctx) {
		inline, ok := f.renderInline(func(out *printer) {
			f.formatArrayMode(out, ctx, true)
		})

		if ok && f.inlineFitsWith(p, inline) {
			p.write(inline)

			return
		}
	}

	f.formatArrayMode(p, ctx, false)
}

func (f *bindingPatternFormatter) formatArrayMode(p *printer, ctx *fql.ArrayBindingPatternContext, inline bool) {
	entries := ctx.AllBindingPattern()
	p.write("[")

	if inline {
		for index, entry := range entries {
			f.formatWith(p, entry.(*fql.BindingPatternContext))
			if index < len(entries)-1 {
				p.write(",")
				p.space()
			}
		}

		p.write("]")

		return
	}

	f.formatMultilineEntries(
		p,
		ctx.OpenBracket(),
		ctx.CloseBracket(),
		f.bindingPatternRules(entries),
		func(index int) {
			f.formatWith(p, entries[index].(*fql.BindingPatternContext))
		},
	)
	p.write("]")
}

func (f *bindingPatternFormatter) formatMultilineEntries(
	p *printer,
	open antlr.TerminalNode,
	close antlr.TerminalNode,
	entries []antlr.ParserRuleContext,
	format func(index int),
) {
	p.withIndent(func() {
		leading := f.trivia.sliceBetween(f.trivia.tokenStop(open)+1, f.trivia.startIndex(entries[0]))
		if f.trivia.containsComment(leading) {
			f.trivia.emitListTriviaWith(p, leading)
		} else {
			p.newline()
		}

		closeStart := f.trivia.tokenStart(close)
		for index, entry := range entries {
			format(index)
			if index < len(entries)-1 {
				p.write(",")
			}

			nextStart := closeStart
			if index < len(entries)-1 {
				nextStart = f.trivia.startIndex(entries[index+1])
			}

			f.trivia.emitListTriviaWith(p, f.trivia.sliceBetween(f.trivia.stopIndex(entry)+1, nextStart))
		}
	})
}

func (f *bindingPatternFormatter) formatEmptyWithComments(
	p *printer,
	open antlr.TerminalNode,
	close antlr.TerminalNode,
	opening string,
	closing string,
) {
	p.write(opening)
	trivia := f.trivia.sliceBetween(f.trivia.tokenStop(open)+1, f.trivia.tokenStart(close))
	if f.trivia.containsComment(trivia) {
		p.withIndent(func() {
			f.trivia.emitListTriviaWith(p, trivia)
		})
	}

	p.write(closing)
}

func (f *bindingPatternFormatter) hasComments(ctx antlr.ParserRuleContext) bool {
	return f.trivia.containsComment(f.trivia.sliceBetween(f.trivia.startIndex(ctx), f.trivia.stopIndex(ctx)+1))
}

func (f *bindingPatternFormatter) objectBindingEntryRules(entries []fql.IObjectBindingEntryContext) []antlr.ParserRuleContext {
	out := make([]antlr.ParserRuleContext, len(entries))
	for index, entry := range entries {
		out[index] = entry.(antlr.ParserRuleContext)
	}

	return out
}

func (f *bindingPatternFormatter) bindingPatternRules(entries []fql.IBindingPatternContext) []antlr.ParserRuleContext {
	out := make([]antlr.ParserRuleContext, len(entries))
	for index, entry := range entries {
		out[index] = entry.(antlr.ParserRuleContext)
	}

	return out
}
