package diagnostics

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/MontFerret/ferret/pkg/file"
)

func FormatDiagnostic(out io.Writer, e *Diagnostic, indent int) {
	prefix := strings.Repeat("  ", indent)

	fmt.Fprintf(out, "%s%s: %s\n", prefix, e.Kind, e.Message)

	// Sort spans by Start for rendering order
	spans := append([]ErrorSpan{}, e.Spans...)
	sort.SliceStable(spans, func(i, j int) bool {
		return spans[i].Span.Start < spans[j].Span.Start
	})

	// Group by file (assumes single file now)
	mainSpan := ErrorSpan{}
	for _, s := range spans {
		if s.Main {
			mainSpan = s
			continue
		}
		renderErrorSpan(out, prefix, e.Source, s)
	}

	// Render primary span last
	if mainSpan.Span.End > 0 {
		renderErrorSpan(out, prefix, e.Source, mainSpan)
	}

	if e.Hint != "" {
		fmt.Fprintf(out, "%sHint: %s\n", prefix, e.Hint)
	}

	if e.Note != "" {
		fmt.Fprintf(out, "%sNote: %s\n", prefix, e.Note)
	}

	if e.Cause != nil {
		var nested *Diagnostic

		if errors.As(e.Cause, &nested) {
			fmt.Fprintf(out, "%sCaused by:\n", prefix)
			FormatDiagnostic(out, nested, indent+1)
		} else {
			fmt.Fprintf(out, "%sCaused by: %s\n", prefix, e.Cause.Error())
		}
	}
}

func renderErrorSpan(out io.Writer, prefix string, src *file.Source, s ErrorSpan) {
	renderer := SpanRenderer{
		Prefix:             prefix,
		CaretChar:          '^',
		ShowTrailingGutter: false,
	}

	renderer.Render(out, src, s.Span, s.Label)
}
