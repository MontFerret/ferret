package diagnostics

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func SpanAt(program *bytecode.Program, pc int) source.Span {
	if program == nil {
		return source.Span{Start: -1, End: -1}
	}

	if pc < 0 || pc >= len(program.Metadata.DebugSpans) {
		return source.Span{Start: -1, End: -1}
	}

	return program.Metadata.DebugSpans[pc]
}
