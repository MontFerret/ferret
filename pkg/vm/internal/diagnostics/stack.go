package diagnostics

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
)

func attachCallStack(err *RuntimeError, program *bytecode.Program, callStack []frame.TraceEntry) *RuntimeError {
	if err == nil || len(callStack) == 0 || err.Diagnostic == nil {
		return err
	}

	if !hasStackTraceLabel(err.Spans) {
		err.Spans = append(callStackSpans(program, callStack), err.Spans...)
	}

	err.Note = appendStackNote(err.Note, callStack)

	return err
}

func hasStackTraceLabel(spans []diagnostics.ErrorSpan) bool {
	for _, span := range spans {
		if strings.HasPrefix(span.Label, "called from ") {
			return true
		}
	}

	return false
}

func callStackSpans(program *bytecode.Program, callStack []frame.TraceEntry) []diagnostics.ErrorSpan {
	if len(callStack) == 0 {
		return nil
	}

	spans := make([]diagnostics.ErrorSpan, 0, len(callStack))

	for i, entry := range callStack {
		span := SpanAt(program, entry.CallSitePC)
		if span.Start < 0 || span.End < 0 {
			continue
		}

		label := fmt.Sprintf("called from (#%d)", i+1)
		if name := strings.TrimSpace(entry.FnName); name != "" {
			label = fmt.Sprintf("called from %s (#%d)", name, i+1)
		}

		spans = append(spans, diagnostics.NewSecondaryErrorSpan(span, label))
	}

	return spans
}

func buildSpans(program *bytecode.Program, callStack []frame.TraceEntry, mainSpan source.Span, label string) []diagnostics.ErrorSpan {
	spans := callStackSpans(program, callStack)
	spans = append(spans, diagnostics.NewMainErrorSpan(mainSpan, label))

	return spans
}

func stackNote(callStack []frame.TraceEntry) string {
	if len(callStack) == 0 {
		return ""
	}

	names := make([]string, 0, len(callStack))

	// callStack is nearest -> farthest; render as outer -> ... -> inner
	for i := len(callStack) - 1; i >= 0; i-- {
		entry := callStack[i]
		name := strings.TrimSpace(entry.FnName)
		if name == "" {
			if entry.FnID < 0 {
				continue
			}

			name = fmt.Sprintf("#%d", entry.FnID)
		}

		names = append(names, name)
	}

	if len(names) == 0 {
		return ""
	}

	return vmStackNotePrefix + strings.Join(names, " -> ")
}

func appendStackNote(note string, callStack []frame.TraceEntry) string {
	stack := stackNote(callStack)
	if stack == "" {
		return note
	}

	if strings.Contains(note, vmStackNotePrefix) || strings.Contains(note, stack) {
		return note
	}

	if note == "" {
		return stack
	}

	return note + "\n" + stack
}
