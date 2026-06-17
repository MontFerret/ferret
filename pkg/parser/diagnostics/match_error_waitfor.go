package diagnostics

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func matchWaitForErrors(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if err == nil || offending == nil {
		return false
	}

	if has(err.Message, "waitforpredicate failed predicate") {
		if keyword, spanNode := waitForPredicateKeyword(offending); keyword != "" {
			span := spanFromTokenSafe(spanNode.Token(), src)
			err.Message = fmt.Sprintf("Expected expression after '%s' in WAITFOR predicate", keyword)
			err.Hint = fmt.Sprintf("Provide an expression after %s, e.g. WAITFOR %s x.", keyword, keyword)
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing expression"),
			}
			return true
		}
	}

	if keyword, spanNode := waitForPredicateKeyword(offending); keyword != "" {
		if is(offending, "RETURN") || isMissing(err.Message) || isNoAlternative(err.Message) {
			span := spanFromTokenSafe(spanNode.Token(), src)
			err.Message = fmt.Sprintf("Expected expression after '%s' in WAITFOR predicate", keyword)
			err.Hint = fmt.Sprintf("Provide an expression after %s, e.g. WAITFOR %s x.", keyword, keyword)
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing expression"),
			}
			return true
		}
	}

	if spanNode := waitForTriggerInlineWaitfor(offending); spanNode != nil {
		span := spanFromTokenSafe(spanNode.Token(), src)
		err.Message = "Nested WAITFOR in TRIGGER shorthand must use a parenthesized block"
		err.Hint = "Use TRIGGER (...), e.g. TRIGGER (WAITFOR EVENT \"ready\" IN target)."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "parenthesize nested wait"),
		}
		return true
	}

	if spanNode := waitForTriggerInvalidBody(offending); spanNode != nil {
		span := spanFromTokenSafe(spanNode.Token(), src)
		err.Message = "Expected trigger statement after 'TRIGGER' in WAITFOR EVENT"
		err.Hint = "Use a side-effect statement or TRIGGER (...), e.g. TRIGGER target <- \"click\"."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing trigger statement"),
		}
		return true
	}

	if span, ok := waitForEventEveryClause(src, offending); ok {
		err.Message = "EVERY is not valid for WAITFOR EVENT"
		err.Hint = "Remove EVERY; event waits subscribe to the event stream and use TIMEOUT as the wait deadline. Use WAITFOR VALUE ... EVERY ... for polling expressions."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "unsupported clause"),
		}
		return true
	}

	if clause, spanNode := waitForMissingClauseValue(offending); clause != "" {
		span := spanFromTokenSafe(spanNode.Token(), src)
		err.Message = fmt.Sprintf("Expected value after '%s' in WAITFOR clause", clause)
		switch clause {
		case "BACKOFF":
			err.Hint = "Provide a backoff strategy, e.g. BACKOFF LINEAR."
		case "JITTER":
			err.Hint = "Provide a jitter value between 0 and 1, e.g. JITTER 0.2."
		default:
			err.Hint = fmt.Sprintf("Provide a duration, e.g. %s 100ms.", clause)
		}
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing value"),
		}
		return true
	}

	return false
}

func waitForPredicateKeyword(offending *TokenNode) (string, *TokenNode) {
	if offending == nil {
		return "", nil
	}

	if is(offending, "EXISTS") {
		if is(offending.Prev(), "NOT") {
			return "NOT EXISTS", offending
		}
		return "EXISTS", offending
	}

	if is(offending, "VALUE") {
		return "VALUE", offending
	}

	if is(offending.Prev(), "EXISTS") {
		if is(offending.Prev().Prev(), "NOT") {
			return "NOT EXISTS", offending.Prev()
		}
		return "EXISTS", offending.Prev()
	}

	if is(offending.Prev(), "VALUE") {
		return "VALUE", offending.Prev()
	}

	return "", nil
}

func waitForTriggerInlineWaitfor(offending *TokenNode) *TokenNode {
	if offending == nil {
		return nil
	}

	if is(offending, "TRIGGER") && is(offending.Next(), "WAITFOR") && hasWaitforBefore(offending) {
		return offending
	}

	for curr := offending; curr != nil; curr = curr.Prev() {
		if is(curr, "TRIGGER") {
			return nil
		}

		if is(curr, "WAITFOR") && is(curr.Prev(), "TRIGGER") && hasWaitforBefore(curr.Prev()) {
			return curr.Prev()
		}
	}

	return nil
}

func waitForTriggerInvalidBody(offending *TokenNode) *TokenNode {
	if offending == nil {
		return nil
	}

	if is(offending, "TRIGGER") && hasWaitforBefore(offending) {
		return offending
	}

	prev := offending.Prev()
	if is(prev, "TRIGGER") && hasWaitforBefore(prev) {
		return prev
	}

	return nil
}

func waitForEventEveryClause(src *source.Source, offending *TokenNode) (source.Span, bool) {
	if src == nil || offending == nil || !hasWaitForEventBefore(offending) {
		return source.Span{}, false
	}

	if node := waitForEventEveryNode(offending); node != nil {
		return spanFromTokenSafe(node.Token(), src), true
	}

	return waitForEventEverySpanAfter(src, offending)
}

func waitForEventEveryNode(offending *TokenNode) *TokenNode {
	for curr := offending; curr != nil; curr = curr.Prev() {
		if is(curr, "WAITFOR") {
			return nil
		}

		if is(curr, "EVERY") {
			return curr
		}
	}

	return nil
}

func waitForEventEverySpanAfter(src *source.Source, offending *TokenNode) (source.Span, bool) {
	tok := offending.Token()
	if tok == nil {
		return source.Span{}, false
	}

	content := src.Content()
	start := tok.GetStop() + 1
	if start < 0 {
		start = tok.GetStart()
	}
	if start < 0 {
		start = 0
	}
	if start >= len(content) {
		return source.Span{}, false
	}

	lexer := fql.NewFqlLexer(antlr.NewInputStream(asciiUpper(content[start:])))
	for i := 0; i < 16; {
		next := lexer.NextToken()
		if next == nil || next.GetTokenType() == antlr.TokenEOF {
			return source.Span{}, false
		}
		if next.GetChannel() != antlr.TokenDefaultChannel {
			continue
		}

		i++
		if next.GetTokenType() == fql.FqlLexerEvery {
			return source.Span{Start: start + next.GetStart(), End: start + next.GetStop() + 1}, true
		}
		if stopsWaitForEventEveryScan(next) {
			return source.Span{}, false
		}
	}

	return source.Span{}, false
}

func stopsWaitForEventEveryScan(tok antlr.Token) bool {
	switch tok.GetTokenType() {
	case fql.FqlLexerSemiColon,
		fql.FqlLexerReturn,
		fql.FqlLexerLet,
		fql.FqlLexerVar,
		fql.FqlLexerFor,
		fql.FqlLexerWaitfor,
		fql.FqlLexerDispatch,
		fql.FqlLexerDelete,
		fql.FqlLexerQuery:
		return true
	case fql.FqlLexerIdentifier:
		return strings.EqualFold(tok.GetText(), "ON")
	default:
		return false
	}
}

func asciiUpper(input string) string {
	out := []byte(input)
	for i, ch := range out {
		if ch >= 'a' && ch <= 'z' {
			out[i] = ch - ('a' - 'A')
		}
	}

	return string(out)
}

func waitForMissingClauseValue(offending *TokenNode) (string, *TokenNode) {
	if offending == nil {
		return "", nil
	}

	if is(offending, "TIMEOUT") || is(offending, "EVERY") || is(offending, "BACKOFF") || is(offending, "JITTER") {
		if hasWaitforBefore(offending) {
			return strings.ToUpper(offending.GetText()), offending
		}
	}

	prev := offending.Prev()
	if prev != nil {
		if is(prev, "TIMEOUT") || is(prev, "EVERY") || is(prev, "BACKOFF") || is(prev, "JITTER") {
			if hasWaitforBefore(prev) {
				return strings.ToUpper(prev.GetText()), prev
			}
		}
	}

	return "", nil
}

func hasWaitForEventBefore(node *TokenNode) bool {
	for curr := node; curr != nil; curr = curr.Prev() {
		if is(curr, "WAITFOR") {
			return is(curr.Next(), "EVENT")
		}
	}

	return false
}

func hasWaitforBefore(node *TokenNode) bool {
	for curr := node.Prev(); curr != nil; curr = curr.Prev() {
		if is(curr, "WAITFOR") {
			return true
		}
	}

	return false
}
