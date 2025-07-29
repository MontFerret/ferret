package core

import (
	"fmt"
	"io"
	"strings"
)

func FormatError(out io.Writer, e *CompilationError, indent int) {
	prefix := strings.Repeat("  ", indent)

	// Header
	fmt.Fprintf(out, "%s%s: %s\n", prefix, e.Kind, e.Message)
	fmt.Fprintf(out, "%s --> %s:%d:%d\n", prefix, e.Source.Name(), e.Location.Line(), e.Location.Column())

	// Determine padding width for line number column
	lineNoWidth := len(fmt.Sprintf("%d", e.Location.Line()))
	fmt.Fprintf(out, "%s%s\n", prefix, strings.Repeat(" ", lineNoWidth)+" |")

	// Multi-line snippet with context
	snippetLines := e.Source.Snippet(*e.Location)

	for _, sl := range snippetLines {
		fmt.Fprintf(out, "%s%*d | %s\n", prefix, lineNoWidth, sl.Line, sl.Text)

		if sl.Caret != "" {
			fmt.Fprintf(out, "%s%s | %s\n", prefix, strings.Repeat(" ", lineNoWidth), sl.Caret)
		}
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
