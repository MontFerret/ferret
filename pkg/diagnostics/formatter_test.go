package diagnostics

import (
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestFormatError_Basic(t *testing.T) {
	src := source.New("test.fql", "LET x = 1")

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
	src := source.New("test.fql", "LET x = 1")

	err := &Diagnostic{
		Kind:    UnexpectedError,
		Message: "test error",
		Source:  src,
		Spans: []ErrorSpan{
			NewMainErrorSpan(source.Span{Start: 0, End: 3}, "main error"),
			NewSecondaryErrorSpan(source.Span{Start: 4, End: 5}, "secondary error"),
		},
	}

	var output strings.Builder
	FormatDiagnostic(&output, err, 0)

	want := "UnexpectedError: test error\n" +
		" --> test.fql:1:5\n  |\n1 | LET x = 1\n  |     ^ secondary error\n" +
		" --> test.fql:1:1\n  |\n1 | LET x = 1\n  | ^^^ main error\n"
	if got := output.String(); got != want {
		t.Fatalf("formatted diagnostic = %q, want %q", got, want)
	}
}

func TestFormatDiagnostic_PrimaryInsertionPoints(t *testing.T) {
	for _, test := range []struct {
		name  string
		query string
		want  string
		span  source.Span
	}{
		{
			name:  "start",
			query: "RETURN",
			span:  source.Span{Start: 0, End: 0},
			want:  " --> test.fql:1:1\n  |\n1 | RETURN\n  | ^ missing value\n",
		},
		{
			name:  "middle",
			query: "RETURN value",
			span:  source.Span{Start: 6, End: 6},
			want:  " --> test.fql:1:7\n  |\n1 | RETURN value\n  |       ^ missing value\n",
		},
		{
			name:  "EOF",
			query: "RETURN",
			span:  source.Span{Start: 6, End: 6},
			want:  " --> test.fql:1:7\n  |\n1 | RETURN\n  |       ^ missing value\n",
		},
		{
			name:  "multiline EOF",
			query: "let arr = []\n\nreturn",
			span:  source.Span{Start: 20, End: 20},
			want:  " --> test.fql:3:7\n  |\n2 | \n3 | return\n  |       ^ missing value\n",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			diag := &Diagnostic{
				Kind:    UnexpectedError,
				Message: "test error",
				Source:  source.New("test.fql", test.query),
				Spans:   []ErrorSpan{NewMainErrorSpan(test.span, "missing value")},
			}

			var output strings.Builder
			FormatDiagnostic(&output, diag, 0)

			want := "UnexpectedError: test error\n" + test.want
			if got := output.String(); got != want {
				t.Fatalf("formatted diagnostic = %q, want %q", got, want)
			}
		})
	}
}

func TestFormatError_WithCause(t *testing.T) {
	src := source.New("test.fql", "LET x = 1")

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
	src := source.New("test.fql", "LET x = 1")

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

func TestFormatDiagnostic_WithoutPrimarySpan(t *testing.T) {
	for _, test := range []struct {
		name  string
		want  string
		spans []ErrorSpan
	}{
		{name: "nil spans"},
		{name: "empty spans", spans: []ErrorSpan{}},
		{
			name:  "secondary only",
			spans: []ErrorSpan{NewSecondaryErrorSpan(source.Span{Start: 4, End: 5}, "secondary error")},
			want:  " --> test.fql:1:5\n  |\n1 | LET x = 1\n  |     ^ secondary error\n",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			diag := &Diagnostic{
				Kind:    UnexpectedError,
				Message: "test error",
				Source:  source.New("test.fql", "LET x = 1"),
				Spans:   test.spans,
				Note:    "test note",
				Hint:    "test hint",
				Cause:   errors.New("underlying error"),
			}

			var output strings.Builder
			FormatDiagnostic(&output, diag, 0)

			want := "UnexpectedError: test error\n" + test.want +
				"Note: test note\nHint: test hint\nCaused by: underlying error\n"
			if got := output.String(); got != want {
				t.Fatalf("formatted diagnostic = %q, want %q", got, want)
			}
		})
	}
}

func TestFormatDiagnostic_NestedPrimarySpan(t *testing.T) {
	diag := &Diagnostic{
		Kind:    UnexpectedError,
		Message: "outer error",
		Cause: &Diagnostic{
			Kind:    UnexpectedError,
			Message: "nested error",
			Source:  source.New("nested.fql", "RETURN"),
			Spans:   []ErrorSpan{NewMainErrorSpan(source.Span{Start: 0, End: 0}, "missing value")},
		},
	}

	var output strings.Builder
	FormatDiagnostic(&output, diag, 0)

	want := "UnexpectedError: outer error\nCaused by:\n" +
		"  UnexpectedError: nested error\n" +
		"   --> nested.fql:1:1\n    |\n  1 | RETURN\n    | ^ missing value\n"
	if got := output.String(); got != want {
		t.Fatalf("formatted diagnostic = %q, want %q", got, want)
	}
}
