package diagnostics

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/file"
)

func TestFormatError_Basic(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")

	err := &Diagnostic{
		Kind:    UnexpectedError,
		Message: "test error message",
		Hint:    "test hint",
		Source:  src,
	}

	var output strings.Builder
	FormatDiagnostic(&output, err, 0)

	result := output.String()

	if !strings.Contains(result, "UnexpectedError: test error message") {
		t.Error("FormatDiagnostic should include error kind and message")
	}

	if !strings.Contains(result, "Hint: test hint") {
		t.Error("FormatDiagnostic should include hint")
	}
}

func TestFormatError_WithSpans(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")

	err := &Diagnostic{
		Kind:    UnexpectedError,
		Message: "test error",
		Source:  src,
		Spans: []ErrorSpan{
			NewMainErrorSpan(file.Span{Start: 0, End: 3}, "main error"),
			NewSecondaryErrorSpan(file.Span{Start: 4, End: 5}, "secondary error"),
		},
	}

	var output strings.Builder
	FormatDiagnostic(&output, err, 0)

	result := output.String()

	if !strings.Contains(result, "UnexpectedError: test error") {
		t.Error("FormatDiagnostic should include error kind and message")
	}
}

func TestFormatError_WithCause(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")

	cause := &Diagnostic{
		Kind:    UnexpectedError,
		Message: "internal error",
		Source:  src,
	}

	err := &Diagnostic{
		Kind:    UnexpectedError,
		Message: "syntax error",
		Source:  src,
		Cause:   cause,
	}

	var output strings.Builder
	FormatDiagnostic(&output, err, 0)

	result := output.String()

	if !strings.Contains(result, "UnexpectedError: syntax error") {
		t.Error("FormatDiagnostic should include main error")
	}

	if !strings.Contains(result, "Caused by:") {
		t.Error("FormatDiagnostic should include caused by section")
	}

	if !strings.Contains(result, "UnexpectedError: internal error") {
		t.Error("FormatDiagnostic should include cause error")
	}
}

func TestFormatError_WithIndent(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")

	err := &Diagnostic{
		Kind:    UnexpectedError,
		Message: "test error",
		Source:  src,
	}

	var output strings.Builder
	FormatDiagnostic(&output, err, 1)

	result := output.String()

	if !strings.Contains(result, "  UnexpectedError: test error") {
		t.Error("FormatDiagnostic should include proper indentation")
	}
}

func TestFormatError_NilSpansHandling(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")

	err := &Diagnostic{
		Kind:    UnexpectedError,
		Message: "test error",
		Source:  src,
		Spans:   nil, // Test with nil spans
	}

	var output strings.Builder
	FormatDiagnostic(&output, err, 0)

	result := output.String()

	if !strings.Contains(result, "UnexpectedError: test error") {
		t.Error("FormatDiagnostic should handle nil spans gracefully")
	}
}
