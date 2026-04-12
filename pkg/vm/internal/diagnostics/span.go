package diagnostics

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func SpanAt(program *bytecode.Program, pc int) source.Span {
	if program == nil {
		return invalidSpan()
	}

	if pc < 0 || pc >= len(program.Metadata.DebugSpans) {
		return invalidSpan()
	}

	return program.Metadata.DebugSpans[pc]
}

func CallArgumentSpanAt(program *bytecode.Program, pc int, pos int) source.Span {
	if program == nil || pos < 0 {
		return invalidSpan()
	}

	if pc < 0 || pc >= len(program.Metadata.CallArgumentSpans) {
		return invalidSpan()
	}

	spans := program.Metadata.CallArgumentSpans[pc]
	if pos >= len(spans) {
		return invalidSpan()
	}

	return spans[pos]
}

func invalidSpan() source.Span {
	return source.Span{Start: -1, End: -1}
}
