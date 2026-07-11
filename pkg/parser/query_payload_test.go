package parser_test

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestQueryPayloadDirectAtomicFormsParse(t *testing.T) {
	tests := []struct {
		assert  func(*testing.T, *fql.QueryPayloadContext)
		name    string
		payload string
	}{
		{
			name:    "string literal",
			payload: `"div"`,
			assert:  assertQueryPayloadLiteral,
		},
		{
			name:    "variable",
			payload: "selector",
			assert:  assertQueryPayloadVariable,
		},
		{
			name:    "param",
			payload: "@selector",
			assert:  assertQueryPayloadParam,
		},
		{
			name:    "member",
			payload: "config.selector",
			assert:  assertQueryPayloadMember,
		},
		{
			name:    "index",
			payload: "selectors[index]",
			assert:  assertQueryPayloadMember,
		},
		{
			name:    "function call",
			payload: "GET_SELECTOR()",
			assert:  assertQueryPayloadFunctionCall,
		},
		{
			name:    "function result member",
			payload: "factory().selector",
			assert:  assertQueryPayloadMember,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := parseFirstQueryPayload(t, "RETURN QUERY "+tt.payload+" IN page")
			tt.assert(t, payload)
		})
	}
}

func TestQueryPayloadParenthesizedExpressionsParse(t *testing.T) {
	tests := []string{
		"prefix + selector",
		"enabled ? primary : fallback",
		"selector IN selectors",
		"BUILD_SELECTOR(options)",
	}

	for _, expr := range tests {
		t.Run(expr, func(t *testing.T) {
			payload := parseFirstQueryPayload(t, "RETURN QUERY ("+expr+") IN page")
			if payload.Expression() == nil {
				t.Fatalf("expected parenthesized payload expression for %q", expr)
			}
			if payload.Expression().GetText() != stripSpaces(expr) {
				t.Fatalf("unexpected payload expression: got %q, want %q", payload.Expression().GetText(), stripSpaces(expr))
			}
		})
	}
}

func TestQueryPayloadUnparenthesizedCompoundFormsFailParse(t *testing.T) {
	tests := []string{
		"RETURN QUERY prefix + selector IN page",
		"RETURN QUERY enabled ? primary : fallback IN page",
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

func TestQueryPayloadParenthesizedExpressionDoesNotConsumeQueryClauses(t *testing.T) {
	tests := []string{
		"RETURN QUERY (prefix + selector) IN page WITH { visible: true }",
		"RETURN QUERY (prefix + selector) IN page OPTIONS { timeout: 1000 }",
		"RETURN QUERY (prefix + selector) IN page USING css",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			query := parseFirstQuery(t, input)
			if query.QueryPayload() == nil || query.QueryPayload().Expression() == nil {
				t.Fatal("expected parenthesized query payload expression")
			}
			if query.Expression() == nil || query.Expression().GetText() != "page" {
				t.Fatalf("expected source expression to remain page, got %q", query.Expression().GetText())
			}
		})
	}
}

func assertQueryPayloadLiteral(t *testing.T, ctx *fql.QueryPayloadContext) {
	t.Helper()

	if lit := ctx.Literal(); lit == nil {
		t.Fatal("expected literal payload")
	}
}

func assertQueryPayloadMember(t *testing.T, ctx *fql.QueryPayloadContext) {
	t.Helper()

	if ctx.MemberExpression() == nil {
		t.Fatal("expected member payload")
	}
}

func assertQueryPayloadFunctionCall(t *testing.T, ctx *fql.QueryPayloadContext) {
	t.Helper()

	if ctx.FunctionCallExpression() == nil {
		t.Fatal("expected function call payload")
	}
}

func assertQueryPayloadParam(t *testing.T, ctx *fql.QueryPayloadContext) {
	t.Helper()

	if ctx.Param() == nil {
		t.Fatal("expected param payload")
	}
}

func assertQueryPayloadVariable(t *testing.T, ctx *fql.QueryPayloadContext) {
	t.Helper()

	if ctx.Variable() == nil {
		t.Fatal("expected variable payload")
	}
}

func parseFirstQueryPayload(t *testing.T, input string) *fql.QueryPayloadContext {
	t.Helper()

	query := parseFirstQuery(t, input)
	payload, ok := query.QueryPayload().(*fql.QueryPayloadContext)
	if !ok || payload == nil {
		t.Fatal("expected query payload context")
	}

	return payload
}

func parseFirstQuery(t *testing.T, input string) *fql.QueryExpressionContext {
	t.Helper()

	program, errors := parseQueryPayloadProgram(input)
	if errors.HasErrors() {
		t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
	}

	return mustFindFirst[*fql.QueryExpressionContext](t, program)
}

func parseQueryPayloadProgram(input string) (*fql.ProgramContext, *parserd.ErrorHandler) {
	src := source.NewAnonymous(input)
	stream := antlr.NewInputStream(input)
	lexer := fql.NewFqlLexer(stream)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	history := parserd.NewTokenHistory(20)
	errors := parserd.NewErrorHandler(src, 5)
	parser := fql.NewFqlParser(parserd.NewTrackingTokenStream(tokens, history))
	parser.BuildParseTrees = true
	parser.RemoveErrorListeners()
	parser.AddErrorListener(parserd.NewErrorListener(src, errors, history))

	return parser.Program().(*fql.ProgramContext), errors
}

func mustFindFirst[T any](t *testing.T, tree antlr.Tree) T {
	t.Helper()

	value, ok := findFirstQueryPayloadNode[T](tree)
	if !ok {
		t.Fatalf("failed to find node %T", value)
	}

	return value
}

func findFirstQueryPayloadNode[T any](tree antlr.Tree) (T, bool) {
	var zero T
	if tree == nil {
		return zero, false
	}
	if value, ok := tree.(T); ok {
		return value, true
	}
	for i := 0; i < tree.GetChildCount(); i++ {
		if value, ok := findFirstQueryPayloadNode[T](tree.GetChild(i)); ok {
			return value, true
		}
	}
	return zero, false
}

func stripSpaces(input string) string {
	out := make([]rune, 0, len(input))
	for _, ch := range input {
		if ch != ' ' {
			out = append(out, ch)
		}
	}
	return string(out)
}
