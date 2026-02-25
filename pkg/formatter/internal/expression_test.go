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

func TestExpressionFormatter_ImplicitMemberExpression(t *testing.T) {
	input := "RETURN [1][* RETURN .name]"
	program := parseProgram(t, input)
	inlineRet := mustFirst[*fql.InlineReturnContext](t, program)
	expr := inlineRet.Expression().(*fql.ExpressionContext)

	var buf bytes.Buffer
	e := newEngine(file.NewAnonymousSource(input), &buf, DefaultOptions())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != ".name" {
		t.Fatalf("unexpected implicit member formatting: %q", got)
	}
}

func TestExpressionFormatter_ImplicitMemberExpressionOptional(t *testing.T) {
	input := "RETURN [1][* RETURN ?.name]"
	program := parseProgram(t, input)
	inlineRet := mustFirst[*fql.InlineReturnContext](t, program)
	expr := inlineRet.Expression().(*fql.ExpressionContext)

	var buf bytes.Buffer
	e := newEngine(file.NewAnonymousSource(input), &buf, DefaultOptions())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "?.name" {
		t.Fatalf("unexpected implicit optional member formatting: %q", got)
	}
}

func TestExpressionFormatter_ImplicitCurrentExpression(t *testing.T) {
	input := "RETURN [1][* RETURN .]"
	program := parseProgram(t, input)
	inlineRet := mustFirst[*fql.InlineReturnContext](t, program)
	expr := inlineRet.Expression().(*fql.ExpressionContext)

	var buf bytes.Buffer
	e := newEngine(file.NewAnonymousSource(input), &buf, DefaultOptions())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "." {
		t.Fatalf("unexpected implicit current formatting: %q", got)
	}
}
