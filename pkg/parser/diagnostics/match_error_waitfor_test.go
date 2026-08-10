package diagnostics

import (
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestWaitForEmptyGroupInSource(t *testing.T) {
	tests := []struct {
		query           string
		mode            string
		synchronization string
	}{
		{query: "RETURN WAITFOR ANY {}", synchronization: "ANY"},
		{query: "RETURN WAITFOR VALUE ALL {}", mode: "VALUE ", synchronization: "ALL"},
		{query: "RETURN WAITFOR EVENT ANY {}", mode: "EVENT ", synchronization: "ANY"},
	}

	for _, test := range tests {
		t.Run(test.query, func(t *testing.T) {
			offending := waitForTestTokenAt(t, test.query, strings.Index(test.query, "}"))
			mode, synchronization, span, ok := waitForEmptyGroupInSource(source.NewAnonymous(test.query), offending)
			if !ok {
				t.Fatal("expected empty WAITFOR group")
			}

			if mode != test.mode || synchronization != test.synchronization {
				t.Fatalf("unexpected group: mode %q, synchronization %q", mode, synchronization)
			}

			if got := test.query[span.Start:span.End]; got != "{}" {
				t.Fatalf("unexpected empty group span %q", got)
			}
		})
	}
}

func TestWaitForEmptyGroupInSourceIgnoresUnrelatedError(t *testing.T) {
	query := "RETURN [1,,2]\nRETURN WAITFOR ANY {}"
	offending := waitForTestTokenAt(t, query, strings.Index(query, ",,"))
	if _, _, _, ok := waitForEmptyGroupInSource(source.NewAnonymous(query), offending); ok {
		t.Fatal("unrelated syntax error matched later empty WAITFOR group")
	}
}

func TestWaitForEmptyGroupInSourceSelectsOffendingGroup(t *testing.T) {
	query := "LET first = WAITFOR ANY {}\nRETURN WAITFOR VALUE ALL {}"
	secondOpen := strings.LastIndex(query, "{}")
	offending := waitForTestTokenAt(t, query, secondOpen+1)
	mode, synchronization, span, ok := waitForEmptyGroupInSource(source.NewAnonymous(query), offending)
	if !ok {
		t.Fatal("expected localized empty WAITFOR group")
	}

	if mode != "VALUE " || synchronization != "ALL" {
		t.Fatalf("unexpected group: mode %q, synchronization %q", mode, synchronization)
	}

	if span.Start != secondOpen || span.End != secondOpen+2 {
		t.Fatalf("expected second empty group span [%d:%d], got [%d:%d]", secondOpen, secondOpen+2, span.Start, span.End)
	}
}

func waitForTestTokenAt(t *testing.T, query string, offset int) antlr.Token {
	t.Helper()
	lexer := fql.NewFqlLexer(antlr.NewInputStream(asciiUpper(query)))
	for {
		token := lexer.NextToken()
		if token == nil || token.GetTokenType() == antlr.TokenEOF {
			t.Fatalf("no token found at offset %d", offset)
		}

		if token.GetStart() <= offset && token.GetStop() >= offset {
			return token
		}
	}
}
