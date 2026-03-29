package diagnostics

import (
	"fmt"
	"io"
	"strings"
)

// SpanRenderer renders a source span with line numbers and a caret marker.
type SpanRenderer struct {
	Prefix             string
	CaretChar          rune
	ShowTrailingGutter bool
}

// Render prints a span diagnostic block. It returns false when no location can be rendered.
func (r SpanRenderer) Render(out io.Writer, src *source.Source, span source.Span, label string) bool {
	if out == nil || src == nil || src.Empty() {
		return false
	}

	if span.Start < 0 || span.End <= span.Start || span.End > len(src.Content()) {
		return false
	}

	line, col := src.LocationAt(span)
	if line == 0 || col == 0 {
		return false
	}

	caretChar := r.CaretChar
	if caretChar == 0 {
		caretChar = '^'
	}

	fmt.Fprintf(out, "%s --> %s:%d:%d\n", r.Prefix, src.Name(), line, col)

	lines := src.Snippet(span)
	if len(lines) == 0 {
		return false
	}

	lineNoWidth := len(fmt.Sprintf("%d", lines[len(lines)-1].Line))
	fmt.Fprintf(out, "%s%s\n", r.Prefix, strings.Repeat(" ", lineNoWidth)+" |")

	for _, sl := range lines {
		fmt.Fprintf(out, "%s%*d | %s\n", r.Prefix, lineNoWidth, sl.Line, sl.Text)

		if sl.Caret != "" {
			caretLine := normalizeCaret(sl.Caret, caretChar)

			if label != "" {
				caretLine += " " + label
			}

			fmt.Fprintf(out, "%s%s | %s\n", r.Prefix, strings.Repeat(" ", lineNoWidth), caretLine)
		}
	}

	if r.ShowTrailingGutter {
		fmt.Fprintf(out, "%s%s\n", r.Prefix, strings.Repeat(" ", lineNoWidth)+" |")
	}

	return true
}

func normalizeCaret(caret string, caretChar rune) string {
	if caret == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(caret))

	for _, r := range caret {
		if r == ' ' {
			b.WriteRune(' ')

			continue
		}

		b.WriteRune(caretChar)
	}

	return b.String()
}
