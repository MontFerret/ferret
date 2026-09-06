package compiler

import (
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

// ANTLR uses rune offsets internally. Published spans use byte offsets, as do
// source indexing and debugger locations. A nil table is the ASCII identity map.
func sourceByteOffsets(src source.Source) []int {
	text := src.Content()
	count := utf8.RuneCountInString(text)
	if count == len(text) {
		return nil
	}

	offsets := make([]int, 0, count+1)
	for offset := range text {
		offsets = append(offsets, offset)
	}

	return append(offsets, len(text))
}

func sourceByteSpan(offsets []int, span source.Span) source.Span {
	if offsets == nil || span.Start < 0 || span.End < 0 || span.Start >= len(offsets) || span.End >= len(offsets) {
		return span
	}

	return source.Span{Start: offsets[span.Start], End: offsets[span.End]}
}

func normalizeProgramSpans(program *bytecode.Program) {
	offsets := sourceByteOffsets(program.Source)
	if offsets == nil {
		return
	}

	metadata := &program.Metadata
	for i := range metadata.DebugSpans {
		metadata.DebugSpans[i] = sourceByteSpan(offsets, metadata.DebugSpans[i])
	}

	for i := range metadata.DebugPoints {
		metadata.DebugPoints[i].Span = sourceByteSpan(offsets, metadata.DebugPoints[i].Span)
	}

	for _, spans := range metadata.CallArgumentSpans {
		for i := range spans {
			spans[i] = sourceByteSpan(offsets, spans[i])
		}
	}
}
