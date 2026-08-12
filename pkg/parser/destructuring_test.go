package parser_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func TestDestructuringBindingsParse(t *testing.T) {
	tests := []string{
		`LET { name, age: years } = user RETURN [name, years]`,
		`VAR [first, _, { value: nested }] = values RETURN [first, nested]`,
		`LET {} = value RETURN 1`,
		`LET [] = value RETURN 1`,
		`LET { name, nested: [first, _], } = value RETURN [name, first]`,
		`FOR { name, score: points } IN users RETURN [name, points]`,
		`FOR [first, _,] IN rows RETURN first`,
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			program, errors := parseQueryPayloadProgram(input)
			if errors.HasErrors() {
				t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
			}

			if mustFindFirst[*fql.StructuredBindingPatternContext](t, program) == nil {
				t.Fatal("expected structured binding pattern")
			}
		})
	}
}

func TestInvalidDestructuringBindingsFailParse(t *testing.T) {
	tests := []string{
		`LET [first,,third] = values RETURN first`,
		`LET [first = 1] = values RETURN first`,
		`LET [...rest] = values RETURN rest`,
		`LET { "name": value } = user RETURN value`,
		`LET { [key]: value } = user RETURN value`,
		`LET { value WHEN value > 0 } = user RETURN value`,
		`VAR target = {} { value } = target RETURN value`,
		`FUNC read({ value }) => value RETURN read({ value: 1 })`,
		`FOR { value } WHILE false RETURN value`,
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, errors := parseQueryPayloadProgram(input)
			if !errors.HasErrors() {
				t.Fatalf("expected parse errors for %q", input)
			}
		})
	}
}
