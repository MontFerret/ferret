package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/file"
)

func TestFormatter_TemplateLiteralDoesNotIndentInterpolation(t *testing.T) {
	input := "RETURN { foo: `line1\n${1}`, veryLongPropertyNameThatForcesMultilineFormatting: 1 }"
	src := file.NewAnonymousSource(input)
	var buf bytes.Buffer
	fmt := New(WithPrintWidth(10))

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "line1\n${1}") {
		t.Fatalf("expected interpolation to start immediately after newline; got:\n%s", out)
	}
	if strings.Contains(out, "line1\n    ${1}") {
		t.Fatalf("unexpected indentation injected before interpolation; got:\n%s", out)
	}
}

func TestFormatter_ArrayTemplateLiteralNewlineForcesMultiline(t *testing.T) {
	input := "RETURN [`line1\n${1}`]"
	src := file.NewAnonymousSource(input)
	var buf bytes.Buffer
	fmt := New(WithPrintWidth(200))

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "line1\n${1}") {
		t.Fatalf("expected newline in template literal; got:\n%s", out)
	}
	if strings.Contains(out, "line1 ${1}") {
		t.Fatalf("unexpected newline collapse in template literal; got:\n%s", out)
	}
}

func TestFormatter_BlockCommentPreservesLeadingSpace(t *testing.T) {
	input := "RETURN 1\n/*\n * a\n * b\n */\nRETURN 2"
	src := file.NewAnonymousSource(input)
	var buf bytes.Buffer
	fmt := New()

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "\n * a\n") {
		t.Fatalf("expected leading space in block comment line; got:\n%s", out)
	}
	if !strings.Contains(out, "\n * b\n") {
		t.Fatalf("expected leading space in block comment line; got:\n%s", out)
	}
}
