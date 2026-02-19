package internal

import (
	"bytes"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func TestExpressionFormatter_UnaryNot(t *testing.T) {
	input := "RETURN NOT a"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(file.NewAnonymousSource(input), &buf, DefaultOptions())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "NOT a" {
		t.Fatalf("unexpected unary operator formatting: %q", got)
	}
}
