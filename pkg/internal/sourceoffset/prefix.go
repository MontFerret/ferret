package sourceoffset

import (
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

// Prefix maps rune boundaries on demand for lookups during parser recovery.
// Its zero value is ready to use. Use the same source text for every lookup.
type Prefix struct {
	offsets Offsets
	bytes   int
}

// Span extends the map only through the requested boundaries. Later lookups
// can revisit earlier spans without scanning the source again.
func (p *Prefix) Span(text string, span source.Span) source.Span {
	end := max(span.Start, span.End)
	if span.Start < 0 || span.End < 0 || end > len(text) {
		return span
	}

	for i := p.offsets.Len(p.bytes); i < end && p.bytes < len(text); i++ {
		_, width := utf8.DecodeRuneInString(text[p.bytes:])
		if width > 1 && p.offsets == nil {
			// The boundaries before this rune were all byte-identical.
			p.offsets = make(Offsets, i+1, end+1)
			for j := range p.offsets {
				p.offsets[j] = j
			}
		}
		p.bytes += width
		if p.offsets != nil {
			p.offsets = append(p.offsets, p.bytes)
		}
	}

	return p.offsets.Span(span)
}
