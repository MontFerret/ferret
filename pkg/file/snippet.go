package file

import "strings"

type Snippet struct {
	Line  int
	Text  string
	Caret string
}

func NewSnippet(src []string, line int) Snippet {
	text := src[line-1]

	return Snippet{
		Line: line,
		Text: text,
	}
}

func NewSnippetWithCaret(src []string, loc Location) Snippet {
	if loc.line <= 0 || loc.line > len(src) {
		return Snippet{}
	}

	srcLine := src[loc.Line()-1]
	runes := []rune(srcLine)
	column := loc.Column()

	// Clamp column to within bounds (1-based)
	if column < 1 {
		column = 1
	}

	if column > len(runes)+1 {
		column = len(runes) + 1
	}

	// Caret must align with visual column (accounting for tabs)
	visualOffset := computeVisualOffset(srcLine, column)
	caretLine := strings.Repeat("_", visualOffset) + "^"

	return Snippet{
		Line:  loc.line,
		Text:  srcLine,
		Caret: caretLine,
	}
}
