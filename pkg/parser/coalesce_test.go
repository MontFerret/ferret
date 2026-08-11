package parser_test

import (
	"slices"
	"testing"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func TestCoalesceLexerKeepsQuestionMarkFormsDistinct(t *testing.T) {
	lexer := fql.NewFqlLexer(antlr.NewInputStream("a ?? b ? : c?.d?"))
	tokens := lexer.GetAllTokens()

	gotTypes := make([]int, 0, len(tokens))
	gotText := make([]string, 0, len(tokens))

	for _, token := range tokens {
		if token.GetChannel() != antlr.TokenDefaultChannel {
			continue
		}

		gotTypes = append(gotTypes, token.GetTokenType())
		gotText = append(gotText, token.GetText())
	}

	wantTypes := []int{
		fql.FqlLexerIdentifier,
		fql.FqlLexerCoalesce,
		fql.FqlLexerIdentifier,
		fql.FqlLexerQuestionMark,
		fql.FqlLexerColon,
		fql.FqlLexerIdentifier,
		fql.FqlLexerQuestionMark,
		fql.FqlLexerDot,
		fql.FqlLexerIdentifier,
		fql.FqlLexerQuestionMark,
	}
	wantText := []string{"a", "??", "b", "?", ":", "c", "?", ".", "d", "?"}

	if !slices.Equal(gotTypes, wantTypes) {
		t.Fatalf("unexpected token types: got %v, want %v", gotTypes, wantTypes)
	}

	if !slices.Equal(gotText, wantText) {
		t.Fatalf("unexpected token text: got %v, want %v", gotText, wantText)
	}
}

func TestCoalesceExpressionIsRightAssociative(t *testing.T) {
	expr := parseCoalesceExpression(t, "RETURN a ?? b ?? c")
	if expr.GetCoalesceOperator() == nil {
		t.Fatal("expected root coalescing expression")
	}

	if got := expr.GetLeft().GetText(); got != "a" {
		t.Fatalf("unexpected root left operand: got %q, want %q", got, "a")
	}

	right := coalesceExpressionContext(t, expr.GetRight())
	if right.GetCoalesceOperator() == nil {
		t.Fatal("expected coalescing expression on the right")
	}

	if got := right.GetLeft().GetText(); got != "b" {
		t.Fatalf("unexpected nested left operand: got %q, want %q", got, "b")
	}

	if got := right.GetRight().GetText(); got != "c" {
		t.Fatalf("unexpected nested right operand: got %q, want %q", got, "c")
	}
}

func TestCoalesceExpressionPrecedence(t *testing.T) {
	t.Run("logical OR binds tighter", func(t *testing.T) {
		expr := parseCoalesceExpression(t, "RETURN a OR b ?? c")
		if expr.GetCoalesceOperator() == nil {
			t.Fatal("expected coalescing expression at the root")
		}

		left := coalesceExpressionContext(t, expr.GetLeft())
		if left.LogicalOrOperator() == nil {
			t.Fatal("expected logical OR in the left operand")
		}
	})

	t.Run("logical OR binds inside fallback", func(t *testing.T) {
		expr := parseCoalesceExpression(t, "RETURN a ?? b OR c")
		if expr.GetCoalesceOperator() == nil {
			t.Fatal("expected coalescing expression at the root")
		}

		right := coalesceExpressionContext(t, expr.GetRight())
		if right.LogicalOrOperator() == nil {
			t.Fatal("expected logical OR in the right operand")
		}
	})

	t.Run("coalesce binds tighter than ternary", func(t *testing.T) {
		expr := parseCoalesceExpression(t, "RETURN condition ? value ?? fallback : other")
		if expr.GetTernaryOperator() == nil {
			t.Fatal("expected ternary expression at the root")
		}

		onTrue := coalesceExpressionContext(t, expr.GetOnTrue())
		if onTrue.GetCoalesceOperator() == nil {
			t.Fatal("expected coalescing expression in the true branch")
		}
	})
}

func TestCoalesceExpressionComposesWithQuestionMarkSyntax(t *testing.T) {
	tests := []string{
		`RETURN user?.profile?.name ?? "Unknown"`,
		`RETURN FAIL()? ?? "fallback"`,
		"RETURN value ? : fallback",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			parseCoalesceExpression(t, input)
		})
	}
}

func parseCoalesceExpression(t *testing.T, input string) *fql.ExpressionContext {
	t.Helper()

	program, errors := parseQueryPayloadProgram(input)
	if errors.HasErrors() {
		t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
	}

	return mustFindFirst[*fql.ExpressionContext](t, program)
}

func coalesceExpressionContext(t *testing.T, ctx fql.IExpressionContext) *fql.ExpressionContext {
	t.Helper()

	expr, ok := ctx.(*fql.ExpressionContext)
	if !ok || expr == nil {
		t.Fatalf("expected expression context, got %T", ctx)
	}

	return expr
}
