package parser_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func TestBracedReturnlessForParses(t *testing.T) {
	tests := []string{
		`FOR value IN [1] {}`,
		`FOR value IN [1] { LET copy = value }`,
		`FOR WHILE false {}`,
		`FOR DO WHILE false {}`,
		`FOR outer IN [1] { FOR inner IN [outer] {} LET after = outer }`,
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			program, errors := parseQueryPayloadProgram(input)
			if errors.HasErrors() {
				t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
			}

			loop := mustFindFirst[*fql.ForExpressionContext](t, program)
			if loop.OpenBrace() == nil || loop.CloseBrace() == nil {
				t.Fatal("returnless FOR is not braced")
			}

			if loop.ReturnExpression() != nil || loop.ForExpressionReturn() != nil {
				t.Fatal("returnless FOR unexpectedly has a terminal result")
			}
		})
	}
}

func TestBracedForAllowsNestedLoopBeforeLaterStatement(t *testing.T) {
	const input = `FOR outer IN [1] {
  FOR inner IN [outer] {}
  LET after = outer
}`

	program, errors := parseQueryPayloadProgram(input)
	if errors.HasErrors() {
		t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
	}

	loop := mustFindFirst[*fql.ForExpressionContext](t, program)
	body := loop.AllForExpressionBody()
	if got, want := len(body), 2; got != want {
		t.Fatalf("body count = %d, want %d", got, want)
	}

	if body[0].ForExpressionStatement().ForExpression() == nil {
		t.Fatal("first body item is not a nested FOR statement")
	}

	if body[1].ForExpressionStatement().VariableDeclaration() == nil {
		t.Fatal("second body item is not a variable declaration")
	}
}

func TestUnbracedForStillRequiresTerminalResult(t *testing.T) {
	const collecting = `FOR value IN [1] RETURN value`

	program, errors := parseQueryPayloadProgram(collecting)
	if errors.HasErrors() {
		t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
	}

	loop := mustFindFirst[*fql.ForExpressionContext](t, program)
	if loop.ForExpressionReturn() == nil {
		t.Fatal("legacy unbraced FOR lost its terminal result")
	}

	const passThrough = `RETURN FOR outer IN [1]
  FOR inner IN [outer]
    RETURN inner`

	program, errors = parseQueryPayloadProgram(passThrough)
	if errors.HasErrors() {
		t.Fatalf("unexpected pass-through parse errors:\n%s", errors.Errors().Format())
	}

	loop = mustFindFirst[*fql.ForExpressionContext](t, program)
	if terminal := loop.ForExpressionReturn(); terminal == nil || terminal.ForExpression() == nil {
		t.Fatal("legacy unbraced FOR lost terminal pass-through nesting")
	}

	_, errors = parseQueryPayloadProgram(`FOR value IN [1] LET copy = value`)
	if !errors.HasErrors() {
		t.Fatal("unbraced returnless FOR unexpectedly parsed")
	}
}
