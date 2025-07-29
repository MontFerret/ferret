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

func (s *Source) Snippet(loc Location) []Snippet {
	if s.Empty() {
		return []Snippet{}
	}

	lineNum := loc.Line()
	lines := s.lines
	var result []Snippet

	// Show previous line if it exists
	if lineNum > 1 {
		result = append(result, NewSnippet(lines, lineNum-1))
	}

	result = append(result, NewSnippetWithCaret(lines, loc))

	// Show next line if it exists
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
