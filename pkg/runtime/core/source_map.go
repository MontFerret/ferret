package core

import "fmt"

type SourceMap struct {
	text   string
	line   int
	column int
}

func NewSourceMap(text string, line, col int) SourceMap {
	return SourceMap{text, line, col}
}

func (s SourceMap) Line() int {
	return s.line
}

func (s SourceMap) Column() int {
	return s.column
}

func (s SourceMap) String() string {
	return fmt.Sprintf("%s at %d:%d", s.text, s.line, s.column)
}
