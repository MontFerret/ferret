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

func (s *Source) Snippet(loc Location) (line string, caret string) {
	if s.Empty() || loc.Line() <= 0 || loc.Line() > len(s.lines) {
		return "", ""
	}

	srcLine := s.lines[loc.Line()-1]
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
	visualOffset := s.computeVisualOffset(srcLine, column)

	caretLine := strings.Repeat(" ", visualOffset) + "^"

	return srcLine, caretLine
}

func (s *Source) computeVisualOffset(line string, column int) int {
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
