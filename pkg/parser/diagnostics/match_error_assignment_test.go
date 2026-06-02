package diagnostics

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"

	pkgdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestMatchFilterAssignmentExpressionFlexibleTargets(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		operator      string
		offending     string
		expectedHint  string
		expectedLabel string
	}{
		{
			name:          "direct path",
			input:         "FOR user IN users FILTER user.active = true RETURN user",
			operator:      "=",
			expectedHint:  filterComparisonAssignmentHint,
			expectedLabel: "use '==' for comparison",
		},
		{
			name:          "left path offending token",
			input:         "FOR user IN users FILTER user.active = true RETURN user",
			operator:      "=",
			offending:     ".",
			expectedHint:  filterComparisonAssignmentHint,
			expectedLabel: "use '==' for comparison",
		},
		{
			name:          "clause boundary offending token",
			input:         "FOR user IN users FILTER user.active = true RETURN user",
			operator:      "=",
			offending:     "RETURN",
			expectedHint:  filterComparisonAssignmentHint,
			expectedLabel: "use '==' for comparison",
		},
		{
			name:          "grouped path",
			input:         "FOR user IN users FILTER (user.active = true) RETURN user",
			operator:      "=",
			expectedHint:  filterComparisonAssignmentHint,
			expectedLabel: "use '==' for comparison",
		},
		{
			name:          "logical prefix path",
			input:         "FOR user IN users FILTER user.name != NONE AND (user.active = true) RETURN user",
			operator:      "=",
			expectedHint:  filterComparisonAssignmentHint,
			expectedLabel: "use '==' for comparison",
		},
		{
			name:          "implicit current inline filter",
			input:         "RETURN users[* FILTER (.active = true)]",
			operator:      "=",
			expectedHint:  filterComparisonAssignmentHint,
			expectedLabel: "use '==' for comparison",
		},
		{
			name:          "augmented assignment",
			input:         "FOR user IN users FILTER (user.active += true) RETURN user",
			operator:      "+=",
			expectedHint:  filterStatementAssignmentHint,
			expectedLabel: "assignment is not an expression",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := source.New("filter_assignment.fql", tt.input)
			tokens := lexDefaultTokens(tt.input)
			operatorIdx := findTestTokenIndex(t, tokens, tt.operator)
			offendingIdx := operatorIdx
			if tt.offending != "" {
				offendingIdx = findTestTokenIndex(t, tokens, tt.offending)
			}
			err := &pkgdiagnostics.Diagnostic{
				Kind:    SyntaxError,
				Message: "no viable alternative at input",
				Source:  src,
			}

			if !matchFilterAssignmentExpression(src, err, &TokenNode{token: tokens[offendingIdx]}) {
				t.Fatal("expected FILTER assignment matcher to handle diagnostic")
			}

			if err.Message != filterAssignmentMessage {
				t.Fatalf("message = %q, want %q", err.Message, filterAssignmentMessage)
			}
			if err.Hint != tt.expectedHint {
				t.Fatalf("hint = %q, want %q", err.Hint, tt.expectedHint)
			}
			if len(err.Spans) != 1 {
				t.Fatalf("expected 1 span, got %d", len(err.Spans))
			}

			span := err.Spans[0]
			if !span.Main {
				t.Fatal("expected main span")
			}
			if span.Label != tt.expectedLabel {
				t.Fatalf("label = %q, want %q", span.Label, tt.expectedLabel)
			}

			expectedSpan := spanFromTokenSafe(tokens[operatorIdx], src)
			if span.Span != expectedSpan {
				t.Fatalf("span = %#v, want %#v", span.Span, expectedSpan)
			}
		})
	}
}

func TestMatchFilterAssignmentExpressionRejectsNonCandidates(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		operator string
	}{
		{
			name:     "missing value",
			input:    "FOR user IN users FILTER user.active =",
			operator: "=",
		},
		{
			name:     "non path left operand",
			input:    "FOR user IN users FILTER (1 = 1) RETURN user",
			operator: "=",
		},
		{
			name:     "nested declaration assignment",
			input:    "FOR user IN users FILTER (FOR item IN user.items LET active = true RETURN active) RETURN user",
			operator: "=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := source.New("filter_assignment.fql", tt.input)
			tokens := lexDefaultTokens(tt.input)
			operatorIdx := findTestTokenIndex(t, tokens, tt.operator)
			err := &pkgdiagnostics.Diagnostic{
				Kind:    SyntaxError,
				Message: "no viable alternative at input",
				Source:  src,
			}

			if matchFilterAssignmentExpression(src, err, &TokenNode{token: tokens[operatorIdx]}) {
				t.Fatal("expected FILTER assignment matcher to ignore diagnostic")
			}
		})
	}
}

func findTestTokenIndex(t *testing.T, tokens []antlr.Token, text string) int {
	t.Helper()

	for idx, token := range tokens {
		if isTokenText(token, text) {
			return idx
		}
	}

	t.Fatalf("token %q not found", text)
	return -1
}
