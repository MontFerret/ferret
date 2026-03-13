package diagnostics

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func SpanAt(program *bytecode.Program, pc int) file.Span {
	if program == nil {
		return file.Span{Start: -1, End: -1}
	}

	if pc < 0 || pc >= len(program.Metadata.DebugSpans) {
		return file.Span{Start: -1, End: -1}
	}

	return program.Metadata.DebugSpans[pc]
}
