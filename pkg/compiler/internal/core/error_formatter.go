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
	lineNum := e.Location.Line()
	lineNoWidth := len(fmt.Sprintf("%d", lineNum))

	// Pipe line
	fmt.Fprintf(out, "%s%s\n", prefix, strings.Repeat(" ", lineNoWidth)+" |")

	// Code line
	lineText, caret := e.Source.Snippet(*e.Location)
	fmt.Fprintf(out, "%s%*d | %s\n", prefix, lineNoWidth, lineNum, lineText)

	// Caret line
	fmt.Fprintf(out, "%s%s | %s\n", prefix, strings.Repeat(" ", lineNoWidth), caret)

	// Hint
	if e.Hint != "" {
		fmt.Fprintf(out, "%sHint: %s\n", prefix, e.Hint)
	}

	// Cause
	if e.Cause != nil {
		if nested, ok := e.Cause.(*CompilationError); ok {
			fmt.Fprintf(out, "%sCaused by:\n", prefix)
			FormatError(out, nested, indent+1)
		} else {
			fmt.Fprintf(out, "%sCaused by: %s\n", prefix, e.Cause.Error())
		}
	}
}
