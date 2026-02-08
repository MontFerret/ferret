package vm

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/pkg/file"
)

func TestRuntimeErrorFormat(t *testing.T) {
	src := file.NewSource("script.fql", "LET num = numerator / 0")
	span := file.Span{Start: 10, End: 23} // "numerator / 0"

	err := &RuntimeError{
		Message: "division by zero",
		Hint:    "ensure the denominator is non-zero before division",
		Note:    "add a conditional check before dividing",
		Label:   "attempt to divide by zero",
		Source:  src,
		Span:    span,
	}

	formatted := err.Format()

	if !strings.Contains(formatted, "error: division by zero") {
		t.Fatalf("expected error header, got:\n%s", formatted)
	}

	if !strings.Contains(formatted, "--> script.fql:1:11") {
		t.Fatalf("expected location header, got:\n%s", formatted)
	}

	if !strings.Contains(formatted, "attempt to divide by zero") {
		t.Fatalf("expected caret label, got:\n%s", formatted)
	}

	if strings.Contains(formatted, "~") {
		t.Fatalf("expected caret to use '^', got:\n%s", formatted)
	}

	if !strings.Contains(formatted, "= help: ensure the denominator is non-zero before division") {
		t.Fatalf("expected help line, got:\n%s", formatted)
	}

	if !strings.Contains(formatted, "= note: add a conditional check before dividing") {
		t.Fatalf("expected note line, got:\n%s", formatted)
	}
}
