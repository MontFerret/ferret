package core

type (
	Location struct {
		line   int
		column int
	}

	Source struct {
		contents []rune
		offsets  []int32
	}
)

func (l Location) Column() int {
	return l.column
}

func (l Location) Line() int {
	return l.line
}

func (l Location) Empty() bool {
	return l.column == 0 && l.line == 0
}

func (s *Source) Content() string {
	return string(s.contents)
}

func (s *Source) Snippet(line int) (string, bool) {
	charStart, found := s.getOffset(line)

	if !found || len(s.contents) == 0 {
		return "", false
	}

	charEnd, found := s.getOffset(line + 1)

	if found {
		return string(s.contents[charStart : charEnd-1]), true
	}

	return string(s.contents[charStart:]), true
}

func (s *Source) getOffset(line int) (int32, bool) {
	if line == 1 {
		return 0, true
	}

	if line > 1 && line <= len(s.offsets) {
		offset := s.offsets[line-2]

		return offset, true
	}

	return -1, false
}
