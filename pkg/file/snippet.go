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

func NewSnippetWithCaret(lines []string, span Span, line int) Snippet {
	if line <= 0 || line > len(lines) {
		return Snippet{}
	}

	srcLine := lines[line-1]
	lineStartOffset := 0

	// Compute actual start-of-line offset
	for i := 0; i < line-1; i++ {
		lineStartOffset += len(lines[i]) + 1 // +1 for \n
	}

	startCol := computeVisualOffset(srcLine, span.Start-lineStartOffset+1)
	endCol := computeVisualOffset(srcLine, span.End-lineStartOffset+1)

	caret := ""
	if endCol <= startCol+1 {
		caret = strings.Repeat(" ", startCol) + "^"
	} else {
		caret = strings.Repeat(" ", startCol) + strings.Repeat("~", endCol-startCol)
	}

	return Snippet{
		Line:  line,
		Text:  srcLine,
		Caret: caret,
	}
}
