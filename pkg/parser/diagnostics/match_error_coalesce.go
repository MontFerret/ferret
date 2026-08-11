package diagnostics

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func matchCoalesceErrors(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if src == nil || err == nil || offending == nil {
		return false
	}

	if !isNoAlternative(err.Message) && !isMissing(err.Message) && !isMismatched(err.Message) {
		return false
	}

	operator := offending.Prev()
	var span source.Span
	if operator != nil && is(operator, "??") {
		span = spanFromTokenSafe(operator.Token(), src)
	} else {
		var ok bool
		span, ok = trailingCoalesceSpan(src)
		if !ok {
			return false
		}
	}

	err.Message = "Expected right-hand expression after '??'"
	err.Hint = "Provide a fallback expression, e.g. value ?? fallback."
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "missing expression"),
	}

	return true
}

func trailingCoalesceSpan(src *source.Source) (source.Span, bool) {
	lexer := fql.NewFqlLexer(antlr.NewInputStream(src.Content()))
	lexer.RemoveErrorListeners()
	tokens := lexer.GetAllTokens()

	for i := len(tokens) - 1; i >= 0; i-- {
		token := tokens[i]
		if token.GetChannel() != antlr.TokenDefaultChannel {
			continue
		}

		if token.GetTokenType() != fql.FqlLexerCoalesce {
			return source.Span{}, false
		}

		return spanFromTokenSafe(token, src), true
	}

	return source.Span{}, false
}
