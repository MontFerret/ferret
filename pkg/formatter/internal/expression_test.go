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

func TestExpressionFormatter_QueryExpressionInline(t *testing.T) {
	input := "RETURN QUERY `.items` IN doc USING css WITH { limit: 10 }"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(file.NewAnonymousSource(input), &buf, DefaultOptions())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "QUERY `.items` IN doc USING css WITH { limit: 10 }" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionParamPayload(t *testing.T) {
	input := "RETURN QUERY @q IN doc USING css"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(file.NewAnonymousSource(input), &buf, DefaultOptions())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "QUERY @q IN doc USING css" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionCountModifier(t *testing.T) {
	input := "RETURN QUERY COUNT `.items` IN doc USING css"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(file.NewAnonymousSource(input), &buf, DefaultOptions())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "QUERY COUNT `.items` IN doc USING css" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionOneModifierWithMultiline(t *testing.T) {
	input := "RETURN QUERY ONE `.items` IN doc USING css WITH { limit: 10, timeout: 5 }"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	opts := DefaultOptions()
	opts.printWidth = 20
	e := newEngine(file.NewAnonymousSource(input), &buf, opts)

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "QUERY ONE `.items` IN doc USING css\n    WITH {\n        limit: 10,\n        timeout: 5\n    }" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_MatchExpressionInline(t *testing.T) {
	input := "RETURN MATCH x(1=>10,_=>0)"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(file.NewAnonymousSource(input), &buf, DefaultOptions())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "MATCH x ( 1 => 10, _ => 0 )" {
		t.Fatalf("unexpected MATCH inline formatting: %q", got)
	}
}

func TestExpressionFormatter_MatchExpressionGuardMultiline(t *testing.T) {
	input := "RETURN MATCH(WHEN a>0=>a,WHEN a<0=>-a,_=>0)"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	opts := DefaultOptions()
	opts.printWidth = 10
	e := newEngine(file.NewAnonymousSource(input), &buf, opts)

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "MATCH (\n    WHEN a > 0 => a,\n    WHEN a < 0 => -a,\n    _ => 0,\n)" {
		t.Fatalf("unexpected MATCH guard multiline formatting: %q", got)
	}
}

func TestExpressionFormatter_MatchExpressionObjectPattern(t *testing.T) {
	input := `RETURN MATCH obj({ "a": 1, b: v }=>v,_=>0)`
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(file.NewAnonymousSource(input), &buf, DefaultOptions())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != `MATCH obj ( { "a": 1, b: v } => v, _ => 0 )` {
		t.Fatalf("unexpected MATCH object pattern formatting: %q", got)
	}
}
