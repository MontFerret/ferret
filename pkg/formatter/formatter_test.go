package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestFormatter_TemplateLiteralDoesNotIndentInterpolation(t *testing.T) {
	input := "RETURN { foo: `line1\n${1}`, veryLongPropertyNameThatForcesMultilineFormatting: 1 }"
	src := source.NewAnonymous(input)
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
	src := source.NewAnonymous(input)
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

func TestFormatter_NestedObjectRespectsPrintWidthAtLineStart(t *testing.T) {
	input := "RETURN [{ a: 1, bb: 2 }]"
	src := source.NewAnonymous(input)
	var buf bytes.Buffer
	fmt := New(WithPrintWidth(18))

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := strings.TrimSpace(buf.String())
	expected := strings.TrimSpace(`RETURN [
    {
        a: 1,
        bb: 2
    }
]`)

	if out != expected {
		t.Fatalf("unexpected nested object formatting:\n%s", out)
	}

	for _, line := range strings.Split(out, "\n") {
		if len(line) > 18 {
			t.Fatalf("line exceeds print width 18 (%d): %q", len(line), line)
		}
	}
}

func TestFormatter_BlockCommentPreservesLeadingSpace(t *testing.T) {
	input := "LET x = 1\n/*\n * a\n * b\n */\nRETURN 2"
	src := source.NewAnonymous(input)
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

func TestFormatter_WaitForEventFilterUsesWhenAndRemainsParseable(t *testing.T) {
	input := "LET obs = []\nWAITFOR EVENT \"test\" IN obs WHEN .type == \"match\"\nRETURN 1"
	src := source.NewAnonymous(input)
	var buf bytes.Buffer
	fmt := New()

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "WHEN .type == \"match\"") {
		t.Fatalf("expected WAITFOR event filter to use WHEN; got:\n%s", out)
	}
	if strings.Contains(out, "FILTER .type == \"match\"") {
		t.Fatalf("unexpected legacy FILTER in WAITFOR event filter; got:\n%s", out)
	}

	var roundTrip bytes.Buffer
	if err := fmt.Format(&roundTrip, source.NewAnonymous(out)); err != nil {
		t.Fatalf("formatted output must remain parseable: %v\nformatted:\n%s", err, out)
	}
}
