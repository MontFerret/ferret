// Package sourceoffset converts parser rune boundaries to native source bytes.
package sourceoffset

import (
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

// Offsets maps rune boundaries, including EOF, to UTF-8 byte offsets.
// A nil table is the ASCII identity map. Tables are immutable after construction.
type Offsets []int

// New builds a linear-size table in linear time, without allocating for ASCII.
func New(text string) Offsets {
	count := utf8.RuneCountInString(text)
	if count == len(text) {
		return nil
	}

	offsets := make(Offsets, 0, count+1)
	for offset := range text {
		offsets = append(offsets, offset)
	}

	return append(offsets, len(text))
}

// Len returns the number of runes for the original source's byte length.
func (o Offsets) Len(byteLength int) int {
	if o == nil {
		return byteLength
	}

	return len(o) - 1
}

// Byte converts a rune boundary. Unconvertible positions remain unchanged.
func (o Offsets) Byte(offset int) int {
	if o == nil || offset < 0 || offset >= len(o) {
		return offset
	}

	return o[offset]
}

// Span converts both boundaries of a half-open rune span. If either boundary
// cannot be converted, the entire span is preserved, matching compiler behavior.
// Callers must distinguish absent annotations from valid zero-width spans.
func (o Offsets) Span(span source.Span) source.Span {
	if o == nil || span.Start < 0 || span.End < 0 || span.Start >= len(o) || span.End >= len(o) {
		return span
	}

	return source.Span{Start: o[span.Start], End: o[span.End]}
}
