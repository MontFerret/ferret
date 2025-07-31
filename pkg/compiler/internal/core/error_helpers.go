package core

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/file"
)

func SpanFromRuleContext(ctx antlr.ParserRuleContext) file.Span {
	start := ctx.GetStart()
	stop := ctx.GetStop()

	if start == nil || stop == nil {
		return file.Span{Start: 0, End: 0}
	}

	return file.Span{Start: start.GetStart(), End: stop.GetStop() + 1}
}

func SpanFromToken(tok antlr.Token) file.Span {
	if tok == nil {
		return file.Span{Start: 0, End: 0}
	}

	return file.Span{Start: tok.GetStart(), End: tok.GetStop() + 1}
}

func SpanFromTokenSafe(tok antlr.Token, src *file.Source) file.Span {
	if tok == nil {
		return file.Span{Start: 0, End: 1}
	}

	start := tok.GetStart()
	end := tok.GetStop() + 1 // exclusive end

	if start < 0 {
		start = 0
	}

	if end <= start {
		end = start + 1
	}

	// clamp to source length
	maxLen := len(src.Content())

	if end > maxLen {
		end = maxLen
	}
	if start > maxLen {
		start = maxLen - 1
	}

	return file.Span{Start: start, End: end}
}
