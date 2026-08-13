package internal

import (
	"unicode/utf8"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func functionParameterByID(decl fql.IFunctionDeclarationContext, id core.BindingID) antlr.ParserRuleContext {
	if decl == nil || decl.FunctionParameterList() == nil {
		return nil
	}

	for _, param := range decl.FunctionParameterList().AllFunctionParameter() {
		ctx, ok := param.(antlr.ParserRuleContext)
		if ok && bindingIDFromRule(ctx) == id {
			return ctx
		}
	}

	return nil
}

func ruleSpan(ctx antlr.ParserRuleContext) source.Span {
	if ctx == nil {
		return source.Span{}
	}

	return parserd.SpanFromRuleContext(ctx)
}

func validSpan(span source.Span) bool {
	return span.Start >= 0 && span.End > span.Start
}

func semanticByteOffsets(src *source.Source) []int {
	if src == nil {
		return []int{0}
	}

	text := src.Content()
	offsets := make([]int, 0, utf8.RuneCountInString(text)+1)

	for offset := range text {
		offsets = append(offsets, offset)
	}

	offsets = append(offsets, len(text))

	return offsets
}
