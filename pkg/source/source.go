package source

import (
	"errors"
	"strings"

	"github.com/goccy/go-json"
)

type (
	Source struct {
		name  string
		text  string
		lines []string
		id    ID
	}

	sourceJSON struct {
		Name string `json:"name"`
		Text string `json:"text"`
	}
)

func New(name, text string) Source {
	lines := strings.Split(text, "\n")

	return Source{
		id:    generateID(name, text),
		name:  name,
		text:  text,
		lines: lines,
	}
}

func NewAnonymous(text string) Source {
	return New("anonymous", text)
}

func (f Source) ID() ID {
	return f.id
}

func (s Source) Name() string {
	return s.name
}

func (s Source) Empty() bool {
	if s.text == "" || len(s.lines) == 0 {
		return true
	}

	return false
}

func (s Source) Content() string {
	return s.text
}

func (s Source) Length() int {
	return len(s.text)
}

func (s Source) LocationAt(span Span) (line, column int) {
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
		// treat it as the end of the previous line only for zero-length spans.
		if offset == lineStart && i > 0 && span.End <= offset {
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

func (s Source) Snippet(span Span) []Snippet {
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

func (s Source) MarshalJSON() ([]byte, error) {
	return json.Marshal(sourceJSON{
		Name: s.Name(),
		Text: s.Content(),
	})
}

func (s *Source) UnmarshalJSON(bytes []byte) error {
	if s == nil {
		return errors.New("source: UnmarshalJSON on nil source")
	}

	var data sourceJSON

	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	s.id = generateID(data.Name, data.Text)
	s.name = data.Name
	s.text = data.Text
	s.lines = strings.Split(data.Text, "\n")

	return nil
}
