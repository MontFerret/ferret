package internal

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func parseProgram(t *testing.T, input string) *fql.ProgramContext {
	t.Helper()

	stream := antlr.NewInputStream(input)
	lexer := fql.NewFqlLexer(stream)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := fql.NewFqlParser(tokens)
	parser.BuildParseTrees = true

	return parser.Program().(*fql.ProgramContext)
}

func mustFirst[T any](t *testing.T, tree antlr.Tree) T {
	t.Helper()

	value, ok := findFirst[T](tree)
	if !ok {
		t.Fatalf("failed to find node %T", value)
	}
	return value
}

func findFirst[T any](tree antlr.Tree) (T, bool) {
	var zero T
	if tree == nil {
		return zero, false
	}
	if value, ok := tree.(T); ok {
		return value, true
	}
	for i := 0; i < tree.GetChildCount(); i++ {
		child := tree.GetChild(i)
		if value, ok := findFirst[T](child); ok {
			return value, true
		}
	}
	return zero, false
}
