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
