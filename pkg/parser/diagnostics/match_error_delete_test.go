package diagnostics

import (
	"testing"

	pkgdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestMatchDeleteStatementErrorsUsesPreviousDeleteToken(t *testing.T) {
	const input = "DELETE obj"

	src := source.New("delete.fql", input)
	tokens := lexDefaultTokens(input)
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}

	history := NewTokenHistory(4)
	for _, token := range tokens {
		history.Add(token)
	}

	err := &pkgdiagnostics.Diagnostic{
		Kind:    SyntaxError,
		Message: "no viable alternative at input 'DELETEobj'",
		Source:  src,
	}

	if !matchDeleteStatementErrors(src, err, history.Last()) {
		t.Fatal("expected DELETE matcher to handle invalid target")
	}

	if err.Message != deleteTargetMessage {
		t.Fatalf("message = %q, want %q", err.Message, deleteTargetMessage)
	}

	const wantHint = `Use DELETE obj.foo or DELETE obj["foo"] to remove a property.`
	if err.Hint != wantHint {
		t.Fatalf("hint = %q, want %q", err.Hint, wantHint)
	}

	if len(err.Spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(err.Spans))
	}

	span := err.Spans[0]
	if !span.Main {
		t.Fatal("expected main diagnostic span")
	}
	if span.Label != "invalid delete target" {
		t.Fatalf("span label = %q, want %q", span.Label, "invalid delete target")
	}

	wantSpan := spanFromTokenSafe(tokens[1], src)
	if span.Span != wantSpan {
		t.Fatalf("span = %#v, want %#v", span.Span, wantSpan)
	}
}
