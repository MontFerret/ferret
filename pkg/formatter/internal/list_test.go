package internal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestListFormatter_TemplateLiteralNewlineForcesMultiline(t *testing.T) {
	input := "RETURN [`line1\n${1}`]"
	program := parseProgram(t, input)
	array := mustFirst[*fql.ArrayLiteralContext](t, program)

	var buf bytes.Buffer
	opts := DefaultOptions()
	opts.printWidth = 200
	e := newEngine(source.NewAnonymous(input), &buf, opts)

	e.list.formatArrayLiteral(array)
	out := buf.String()
	if !strings.Contains(out, "line1\n${1}") {
		t.Fatalf("expected newline in template literal; got:\n%s", out)
	}
	if strings.Contains(out, "line1 ${1}") {
		t.Fatalf("unexpected newline collapse in template literal; got:\n%s", out)
	}
}
