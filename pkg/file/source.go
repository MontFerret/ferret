package file

import "strings"

type Source struct {
	name  string
	text  string
	lines []string
}

func NewSource(name, text string) *Source {
	lines := strings.Split(text, "\n")

	return &Source{
		name:  name,
		text:  text,
		lines: lines,
	}
}

func (s *Source) Name() string {
	if s == nil {
		return ""
	}

	return s.name
}

func (s *Source) Empty() bool {
	if s == nil || s.text == "" || len(s.lines) == 0 {
		return true
	}

	return false
}

func (s *Source) Content() string {
	return s.text
}

func (s *Source) LocationAt(span Span) (line, column int) {
	if s.Empty() || span.Start < 0 || span.End > len(s.text) {
		return 0, 0
	}

	offset := span.Start
	total := 0

	for i, l := range s.lines {
		lineLen := len(l) + 1 // +1 for '\n'
		lineStart := total
		lineEndWithNL := total + lineLen

		// If offset is exactly at the start of this line (not the very first line),
		// treat it as the end of the previous line.
		if offset == lineStart && i > 0 {
			prev := s.lines[i-1]
			return i, len(prev) + 1
		}

		if lineEndWithNL > offset {
			// Normal case: offset lives on this line
			return i + 1, offset - total + 1
		}

		total = lineEndWithNL
	}

	// If we somehow fell through, clamp to last line end
	if len(s.lines) > 0 {
		last := s.lines[len(s.lines)-1]
		return len(s.lines), len(last) + 1
	}

	return 0, 0
}

func (s *Source) Snippet(span Span) []Snippet {
	if s.Empty() {
		return nil
	}

	lineNum, _ := s.LocationAt(span)
	lines := s.lines
	var result []Snippet

	// Show previous Line if it exists
	if lineNum > 1 {
		result = append(result, NewSnippet(lines, lineNum-1))
	}

	result = append(result, NewSnippetWithCaret(lines, span, lineNum))

	// Show next Line if it exists
	if lineNum < len(lines) {
		result = append(result, NewSnippet(lines, lineNum+1))
	}

	return result
}

func computeVisualOffset(line string, column int) int {
	runes := []rune(line)
	offset := 0
	tabWidth := 4

	for i := 0; i < column-1 && i < len(runes); i++ {
		if runes[i] == '\t' {
			offset += tabWidth - (offset % tabWidth)
		} else {
			offset += 1
		}
	}

	return offset
}
