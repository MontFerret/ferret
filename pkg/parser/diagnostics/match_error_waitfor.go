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

	mode, synchronization, open, close := waitForEmptyGroup(offending)
	var openSpan source.Span
	ok := open != nil
	if ok {
		first := spanFromTokenSafe(open.Token(), src)
		last := spanFromTokenSafe(close.Token(), src)
		openSpan = source.Span{Start: first.Start, End: last.End}
	} else {
		mode, synchronization, openSpan, ok = waitForEmptyGroupInSource(src, offending.Token())
	}

	if ok {
		err.Message = fmt.Sprintf("WAITFOR %s%s group requires at least one arm", mode, synchronization)
		err.Hint = "Add at least one expression or event subscription between the braces."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(openSpan, "empty synchronization group"),
		}

		return true
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

func waitForEmptyGroupInSource(src *source.Source, offending antlr.Token) (string, string, source.Span, bool) {
	if src == nil || offending == nil {
		return "", "", source.Span{}, false
	}

	lexer := fql.NewFqlLexer(antlr.NewInputStream(asciiUpper(src.Content())))
	var tokens []antlr.Token

	for {
		token := lexer.NextToken()
		if token == nil || token.GetTokenType() == antlr.TokenEOF {
			break
		}

		if token.GetChannel() == antlr.TokenDefaultChannel {
			tokens = append(tokens, token)
		}
	}

	type candidate struct {
		mode            string
		synchronization string
		span            source.Span
		distance        int
	}

	var selected *candidate
	for idx, token := range tokens {
		if token.GetTokenType() != fql.FqlLexerWaitfor {
			continue
		}

		cursor := idx + 1
		var modeParts []string
		for cursor < len(tokens) && len(modeParts) < 2 {
			tokenType := tokens[cursor].GetTokenType()
			if tokenType == fql.FqlLexerAny || tokenType == fql.FqlLexerAll {
				break
			}
			modeParts = append(modeParts, strings.ToUpper(tokens[cursor].GetText()))
			cursor++
		}

		if cursor+2 >= len(tokens) {
			continue
		}

		modeText := strings.Join(modeParts, " ")
		switch modeText {
		case "", "EVENT", "EXISTS", "VALUE", "NOT EXISTS":
		default:
			continue
		}

		synchronizationToken := tokens[cursor]
		if synchronizationToken.GetTokenType() != fql.FqlLexerAny && synchronizationToken.GetTokenType() != fql.FqlLexerAll {
			continue
		}

		open := tokens[cursor+1]
		close := tokens[cursor+2]
		if open.GetTokenType() != fql.FqlLexerOpenBrace || close.GetTokenType() != fql.FqlLexerCloseBrace {
			continue
		}

		distance, related := waitForEmptyGroupDistance(tokens, idx, cursor+2, offending)
		if !related {
			continue
		}

		mode := ""
		if modeText != "" {
			mode = modeText + " "
		}

		found := candidate{
			mode:            mode,
			synchronization: strings.ToUpper(synchronizationToken.GetText()),
			span: source.Span{
				Start: open.GetStart(),
				End:   close.GetStop() + 1,
			},
			distance: distance,
		}

		if selected == nil || found.distance < selected.distance {
			selected = &found
		}
	}

	if selected == nil {
		return "", "", source.Span{}, false
	}

	return selected.mode, selected.synchronization, selected.span, true
}

func waitForEmptyGroupDistance(tokens []antlr.Token, waitForIndex, closeIndex int, offending antlr.Token) (int, bool) {
	waitFor := tokens[waitForIndex]
	close := tokens[closeIndex]
	offendingStart := offending.GetStart()
	if offendingStart >= waitFor.GetStart() && offendingStart <= close.GetStop() {
		return 0, true
	}

	nextIndex := closeIndex + 1
	if nextIndex < len(tokens) {
		next := tokens[nextIndex]
		if next.GetTokenType() == offending.GetTokenType() && next.GetStart() == offendingStart {
			return next.GetStart() - close.GetStop(), true
		}

		return 0, false
	}

	if offending.GetTokenType() == antlr.TokenEOF {
		return offendingStart - close.GetStop(), true
	}

	return 0, false
}

func waitForEmptyGroup(offending *TokenNode) (string, string, *TokenNode, *TokenNode) {
	if offending == nil {
		return "", "", nil, nil
	}

	close := offending
	if !is(close, "}") {
		close = offending.Next()
	}

	if !is(close, "}") || !is(close.Prev(), "{") {
		return "", "", nil, nil
	}

	open := close.Prev()
	synchronizationNode := open.Prev()
	if !is(synchronizationNode, "ANY") && !is(synchronizationNode, "ALL") {
		return "", "", nil, nil
	}

	var modeParts []string
	foundWaitFor := false

	for curr := synchronizationNode.Prev(); curr != nil; curr = curr.Prev() {
		if is(curr, "WAITFOR") {
			foundWaitFor = true

			break
		}

		if len(modeParts) == 2 {
			return "", "", nil, nil
		}

		modeParts = append([]string{strings.ToUpper(curr.GetText())}, modeParts...)
	}

	if !foundWaitFor {
		return "", "", nil, nil
	}

	mode := ""

	if len(modeParts) > 0 {
		mode = strings.Join(modeParts, " ") + " "
	}

	return mode, strings.ToUpper(synchronizationNode.GetText()), open, close
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
