package internal

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

type triviaEmitter struct {
	*engine
}

func (t *triviaEmitter) startIndex(ctx antlr.ParserRuleContext) int {
	if ctx == nil {
		return 0
	}

	if tok := ctx.GetStart(); tok != nil {
		return tok.GetStart()
	}

	return 0
}

func (t *triviaEmitter) stopIndex(ctx antlr.ParserRuleContext) int {
	if ctx == nil {
		return 0
	}

	if tok := ctx.GetStop(); tok != nil {
		return tok.GetStop()
	}

	return 0
}

func (t *triviaEmitter) sliceBetween(start, end int) string {
	if t.src == nil {
		return ""
	}

	text := t.src.Content()

	if start < 0 {
		start = 0
	}

	if end > len(text) {
		end = len(text)
	}

	if end <= start {
		return ""
	}

	return text[start:end]
}

func (t *triviaEmitter) tokenStart(node antlr.TerminalNode) int {
	if node == nil {
		return 0
	}

	if sym := node.GetSymbol(); sym != nil {
		return sym.GetStart()
	}

	return 0
}

func (t *triviaEmitter) isInlineComment(text string) bool {
	if text == "" || strings.Contains(text, "\n") {
		return false
	}

	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return false
	}

	return strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*")
}

func (t *triviaEmitter) extractInlineComment(line string) string {
	if line == "" {
		return ""
	}

	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return ""
	}

	trimmed = strings.TrimLeft(trimmed, " \t,")
	if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
		return trimmed
	}

	return ""
}

func (t *triviaEmitter) containsComment(text string) bool {
	return strings.Contains(text, "//") || strings.Contains(text, "/*")
}

func (t *triviaEmitter) emitBetween(prev, next antlr.ParserRuleContext) {
	if prev == nil || next == nil {
		return
	}

	start := t.stopIndex(prev) + 1
	end := t.startIndex(next)
	t.emitBetweenIndices(start, end)
}

func (t *triviaEmitter) emitBetweenIndices(start, end int) {
	text := t.sliceBetween(start, end)
	if text == "" {
		t.p.newline()
		return
	}

	if t.isInlineComment(text) {
		if !t.p.atLineStart {
			t.p.space()
		}

		t.p.write(strings.TrimSpace(text))
		t.p.newline()

		return
	}

	if strings.TrimSpace(text) == "" {
		count := strings.Count(text, "\n")
		if count >= 2 {
			t.p.newline()
		}

		t.p.newline()

		return
	}

	t.emitTrivia(text, false, true)
}

func (t *triviaEmitter) emitTrivia(text string, trimLeading bool, hasPrevLine bool) {
	if text == "" {
		return
	}

	lines := strings.Split(text, "\n")
	emptyCount := 0
	seenNonEmpty := false

	for _, line := range lines {
		formatted, ok := t.formatTriviaLine(t.p, line)

		if ok {
			if seenNonEmpty {
				if emptyCount >= 1 {
					t.p.newline()
				}

				t.p.newline()
			} else if hasPrevLine {
				if emptyCount >= 2 {
					t.p.newline()
				}

				t.p.newline()
			} else if !trimLeading {
				if emptyCount >= 2 {
					t.p.newline()
				}

				if emptyCount >= 1 {
					t.p.newline()
				}
			}

			emptyCount = 0
			seenNonEmpty = true

			if !t.p.atLineStart {
				t.p.space()
			}

			t.p.write(formatted)

			continue
		}

		if trimLeading && !seenNonEmpty && !hasPrevLine {
			continue
		}

		emptyCount++
	}

	if seenNonEmpty || hasPrevLine {
		if emptyCount >= 2 {
			t.p.newline()
		}

		if emptyCount >= 1 || (seenNonEmpty && emptyCount == 0) {
			t.p.newline()
		}
	}
}

func (t *triviaEmitter) emitListTriviaWith(p *printer, text string) {
	if text == "" {
		p.newline()

		return
	}

	if strings.TrimSpace(text) == "" {
		if strings.Count(text, "\n") >= 2 {
			p.newline()
		}

		p.newline()

		return
	}

	line := text
	rest := ""

	if idx := strings.IndexByte(text, '\n'); idx >= 0 {
		line = text[:idx]
		rest = text[idx+1:]
	}

	if inline := t.extractInlineComment(line); inline != "" {
		if !p.atLineStart {
			p.space()
		}

		p.write(inline)
		p.newline()
	} else {
		p.newline()
	}

	t.emitListCommentLinesWith(p, rest)
}

func (t *triviaEmitter) emitListCommentLinesWith(p *printer, text string) {
	if text == "" {
		return
	}

	lines := strings.Split(text, "\n")
	emptyCount := 0
	wroteLine := false

	for _, line := range lines {
		formatted, ok := t.formatTriviaLine(p, line)

		if !ok {
			emptyCount++
			continue
		}

		if emptyCount > 0 {
			p.newline()
			emptyCount = 0
		}

		if !p.atLineStart {
			p.space()
		}

		p.write(formatted)
		p.newline()
		wroteLine = true
	}

	if emptyCount >= 2 {
		p.newline()
		return
	}

	if !wroteLine && emptyCount > 0 {
		p.newline()
	}
}

func (t *triviaEmitter) emitLeading(next antlr.ParserRuleContext) {
	if next == nil {
		return
	}

	t.emitTrivia(t.sliceBetween(0, t.startIndex(next)), true, false)
}

func (t *triviaEmitter) emitTrailing(prev antlr.ParserRuleContext) {
	if prev == nil {
		return
	}

	start := t.stopIndex(prev) + 1
	text := t.sliceBetween(start, len(t.src.Content()))

	if t.isInlineComment(text) {
		if !t.p.atLineStart {
			t.p.space()
		}

		t.p.write(strings.TrimSpace(text))

		return
	}

	t.emitTrivia(text, false, true)
}

func (t *triviaEmitter) formatTriviaLine(p *printer, line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return "", false
	}

	rightTrimmed := strings.TrimRight(line, " \t")
	leftTrimmed := strings.TrimLeft(rightTrimmed, " \t")

	if t.isCommentLine(leftTrimmed) {
		return t.trimIndentPrefix(p, rightTrimmed), true
	}

	return trimmed, true
}

func (t *triviaEmitter) isCommentLine(trimmed string) bool {
	return strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "*")
}

func (t *triviaEmitter) trimIndentPrefix(p *printer, line string) string {
	if p == nil {
		return strings.TrimLeft(line, " \t")
	}

	max := int(p.opts.tabWidth) * p.indent
	if max <= 0 {
		return line
	}

	idx := 0
	for idx < len(line) && max > 0 {
		switch line[idx] {
		case ' ', '\t':
			idx++
			max--
		default:
			max = 0
		}
	}

	return line[idx:]
}
