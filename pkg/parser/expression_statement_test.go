package parser_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func TestGeneralExpressionStatementsParseInEveryStatementScope(t *testing.T) {
	expressions := []string{
		`42`,
		`1 + 2 * 3`,
		`value.member`,
		`QUERY ("." + suffix) IN page`,
		`MATCH value { 1 => true, _ => false }`,
		`(((value)))`,
		`(FOR item IN [1] { RETURN item })`,
		`(FOR item IN [1] {})`,
	}

	for _, expr := range expressions {
		for _, scope := range []string{"script", "udf", "for body"} {
			t.Run(scope+"/"+expr, func(t *testing.T) {
				input := expressionStatementInput(scope, expr)
				program, errors := parseQueryPayloadProgram(input)
				if errors.HasErrors() {
					t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
				}

				statements := collectUDFMemberNodes[*fql.ExpressionStatementContext](program)
				if got, want := len(statements), 1; got != want {
					t.Fatalf("expression statement count = %d, want %d", got, want)
				}

				if got, want := statements[0].Expression().GetText(), stripSpaces(expr); got != want {
					t.Fatalf("expression text = %q, want %q", got, want)
				}
			})
		}
	}
}

func TestStatementClassificationRemainsContextual(t *testing.T) {
	const input = `VAR value = 0
value = 1
value == 1
FOR item IN [1] {
  FILTER item > 0
  SORT item
  LIMIT 1
  COLLECT group = item
  FILTER()
  SORT.member
  LIMIT + item
  COLLECT == item
  RETURN group
}
RETURN value`

	program, errors := parseQueryPayloadProgram(input)
	if errors.HasErrors() {
		t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
	}

	body := program.Body().(*fql.BodyContext)
	if got, want := len(body.AllBodyStatement()), 4; got != want {
		t.Fatalf("top-level statement count = %d, want %d", got, want)
	}

	if body.BodyStatement(1).AssignmentStatement() == nil {
		t.Fatal("value = 1 was not classified as an assignment")
	}

	if body.BodyStatement(2).ExpressionStatement() == nil {
		t.Fatal("value == 1 was not classified as an expression statement")
	}

	if body.BodyExpression() == nil || body.BodyExpression().ReturnExpression() == nil {
		t.Fatal("terminal RETURN was not classified as the body expression")
	}

	loop := body.BodyStatement(3).ForExpression().(*fql.ForExpressionContext)
	entries := loop.AllForExpressionBody()
	if got, want := len(entries), 8; got != want {
		t.Fatalf("loop body count = %d, want %d", got, want)
	}

	for index := range 4 {
		if entries[index].ForExpressionClause() == nil {
			t.Fatalf("loop entry %d was not classified as a clause", index)
		}
	}

	for index := 4; index < 8; index++ {
		stmt := entries[index].ForExpressionStatement()
		if stmt == nil || stmt.ExpressionStatement() == nil {
			t.Fatalf("loop entry %d was not classified as an expression statement", index)
		}
	}
}

func TestGroupedForReturnClassification(t *testing.T) {
	for _, loop := range []string{
		`FOR value IN [1] { RETURN value }`,
		`FOR value IN [1] {}`,
	} {
		t.Run(loop, func(t *testing.T) {
			program, errors := parseQueryPayloadProgram("RETURN (" + loop + ")")
			if errors.HasErrors() {
				t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
			}

			ret := program.Body().(*fql.BodyContext).BodyExpression().ReturnExpression()
			if ret.ReturnValue().ForExpression() != nil {
				t.Fatal("grouped FOR unexpectedly used the direct bare-FOR return alternative")
			}

			if ret.ReturnValue().Expression() == nil {
				t.Fatal("grouped FOR did not parse as a return expression")
			}
		})
	}
}

func expressionStatementInput(scope, expr string) string {
	switch scope {
	case "script":
		return expr + "\nRETURN NONE"
	case "udf":
		return fmt.Sprintf("FUNC effect() {\n%s\nRETURN NONE\n}\nRETURN effect()", expr)
	case "for body":
		return fmt.Sprintf("FOR outer IN [1] {\n%s\nRETURN outer\n}\nRETURN NONE", expr)
	default:
		return ""
	}
}
