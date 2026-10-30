package compiler

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/internal/sourceoffset"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

// ANTLR uses rune offsets internally. Published spans use byte offsets, as do
// source indexing and debugger locations. A nil table is the ASCII identity map.
func sourceByteOffsets(src source.Source) sourceoffset.Offsets {
	return sourceoffset.New(src.Content())
}

func sourceByteSpan(offsets sourceoffset.Offsets, span source.Span) source.Span {
	return offsets.Span(span)
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
