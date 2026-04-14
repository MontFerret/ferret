package diagnostics

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
)

type runtimeDiagnosticSpec struct {
	Cause   error
	Kind    diagnostics.Kind
	Hint    string
	Label   string
	Message string
	Note    string
	Span    source.Span
}

func newRuntimeErrorWithSpec(
	program *bytecode.Program,
	callStack []frame.TraceEntry,
	spec runtimeDiagnosticSpec,
) *RuntimeError {
	if spec.Hint == "" {
		spec.Hint = synthesizeRuntimeHint(spec)
	}

	return &RuntimeError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    spec.Kind,
			Message: spec.Message,
			Hint:    spec.Hint,
			Note:    appendStackNote(spec.Note, callStack),
			Source:  runtimeErrorSource(program),
			Spans:   buildSpans(program, callStack, spec.Span, spec.Label),
			Cause:   spec.Cause,
		},
	}
}

func runtimeErrorSource(program *bytecode.Program) *source.Source {
	if program == nil {
		return nil
	}

	return program.Source
}

func unwrapRuntimeDetail(err error) (string, error) {
	if err == nil {
		return "", nil
	}

	wrapped, detail := diagnostics.Unwrap(err)
	if wrapped == nil || detail == nil {
		return "", err
	}

	text := strings.TrimSpace(detail.Error())
	if text == "" {
		return "", wrapped
	}

	return text, wrapped
}

func detailNote(detail string) string {
	return strings.TrimSpace(detail)
}

func argumentDetailNote(index int, detail string) string {
	detail = detailNote(detail)
	if detail == "" {
		return ""
	}

	if strings.HasPrefix(detail, "expected ") {
		return fmt.Sprintf("argument %d expects %s", index, strings.TrimPrefix(detail, "expected "))
	}

	if strings.HasPrefix(detail, "expects ") || strings.HasPrefix(detail, "must ") {
		return fmt.Sprintf("argument %d %s", index, detail)
	}

	return fmt.Sprintf("argument %d: %s", index, detail)
}

func fallbackRuntimeMessage(err error) string {
	_, cause := unwrapRuntimeDetail(err)
	if cause == nil {
		return "runtime error"
	}

	return strings.TrimSpace(cause.Error())
}

func panicValueNote(r any) string {
	return fmt.Sprintf("panic value: %v", r)
}
