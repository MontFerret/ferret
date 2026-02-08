package diagnostics

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/MontFerret/ferret/pkg/file"
)

func FormatError(out io.Writer, e *CompilationError, indent int) {
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

	if e.Cause != nil {
		if nested, ok := e.Cause.(*CompilationError); ok {
			fmt.Fprintf(out, "%sCaused by:\n", prefix)
			FormatError(out, nested, indent+1)
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
