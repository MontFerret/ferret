package compiler_test

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	pkgdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestSyntaxErrorsWaitfor(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(`RETURN WAITFOR ANY {}`, E{
			Kind:    parserd.SyntaxError,
			Message: "WAITFOR ANY group requires at least one arm",
			Hint:    "Add at least one expression or event subscription between the braces.",
		}, "WAITFOR ANY rejects an empty group"),
		Failure(`RETURN WAITFOR VALUE ALL {}`, E{
			Kind:    parserd.SyntaxError,
			Message: "WAITFOR VALUE ALL group requires at least one arm",
			Hint:    "Add at least one expression or event subscription between the braces.",
		}, "WAITFOR VALUE ALL rejects an empty group"),
		Failure(`RETURN WAITFOR EVENT ANY {}`, E{
			Kind:    parserd.SyntaxError,
			Message: "WAITFOR EVENT ANY group requires at least one arm",
			Hint:    "Add at least one expression or event subscription between the braces.",
		}, "WAITFOR EVENT ANY rejects an empty group"),
		Failure(`
			LET ok = WAITFOR EXISTS
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected expression after 'EXISTS' in WAITFOR predicate",
			Hint:    "Provide an expression after EXISTS, e.g. WAITFOR EXISTS x.",
		}, "Missing WAITFOR EXISTS expression"),
		Failure(`
			LET ok = WAITFOR NOT EXISTS
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected expression after 'NOT EXISTS' in WAITFOR predicate",
			Hint:    "Provide an expression after NOT EXISTS, e.g. WAITFOR NOT EXISTS x.",
		}, "Missing WAITFOR NOT EXISTS expression"),
		Failure(`
			LET ok = WAITFOR VALUE
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected expression after 'VALUE' in WAITFOR predicate",
			Hint:    "Provide an expression after VALUE, e.g. WAITFOR VALUE x.",
		}, "Missing WAITFOR VALUE expression"),
		Failure(`
			LET ok = WAITFOR TRUE TIMEOUT
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected value after 'TIMEOUT' in WAITFOR clause",
			Hint:    "Provide a duration, e.g. TIMEOUT 100ms.",
		}, "Missing WAITFOR TIMEOUT value"),
		Failure(`
			LET ok = WAITFOR TRUE EVERY
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected value after 'EVERY' in WAITFOR clause",
			Hint:    "Provide a duration, e.g. EVERY 100ms.",
		}, "Missing WAITFOR EVERY value"),
		Failure(`
			LET ok = WAITFOR TRUE BACKOFF
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected value after 'BACKOFF' in WAITFOR clause",
			Hint:    "Provide a backoff strategy, e.g. BACKOFF LINEAR.",
		}, "Missing WAITFOR BACKOFF strategy"),
		Failure(`
			LET ok = WAITFOR TRUE JITTER
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected value after 'JITTER' in WAITFOR clause",
			Hint:    "Provide a jitter value between 0 and 1, e.g. JITTER 0.2.",
		}, "Missing WAITFOR JITTER value"),
		Failure(`
			LET ok = WAITFOR TRUE EVERY 50ms,
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing WAITFOR EVERY cap value"),
		Failure(`
			LET ok = WAITFOR TRUE EVERY 50ms 2s
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing comma in WAITFOR EVERY cap clause"),
		Failure(`
			LET obs = {}
			LET ok = WAITFOR EVENT "test" IN obs FILTER .type == "match"
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Legacy FILTER keyword is invalid in WAITFOR EVENT"),
		Failure(`
			LET obs = {}
			LET ok = WAITFOR EVENT "test" IN obs TIMEOUT 5s EVERY 250ms
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "EVERY is not valid for WAITFOR EVENT",
			Hint:    "Remove EVERY; event waits subscribe to the event stream and use TIMEOUT as the wait deadline. Use WAITFOR VALUE ... EVERY ... for polling expressions.",
		}, "WAITFOR EVENT rejects EVERY after TIMEOUT"),
		Failure(`
			LET obs = {}
			LET btn = {}
			LET ok = WAITFOR EVENT "test" IN obs
				TRIGGER (
					btn <- "click"
				)
				TIMEOUT 5s
				EVERY 250ms
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "EVERY is not valid for WAITFOR EVENT",
			Hint:    "Remove EVERY; event waits subscribe to the event stream and use TIMEOUT as the wait deadline. Use WAITFOR VALUE ... EVERY ... for polling expressions.",
		}, "WAITFOR EVENT trigger block rejects EVERY after TIMEOUT"),
		Failure(`
			LET obs = {}
			LET ok = WAITFOR EVENT "test" IN obs
				TRIGGER
				TIMEOUT 5s
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected trigger statement after 'TRIGGER' in WAITFOR EVENT",
			Hint:    "Use a side-effect statement or TRIGGER (...), e.g. TRIGGER target <- \"click\".",
		}, "WAITFOR EVENT TRIGGER requires a body"),
		Failure(`
			LET obs = {}
			LET ok = WAITFOR EVENT "test" IN obs
				TRIGGER WAITFOR EVENT "inner" IN obs
				TIMEOUT 5s
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Nested WAITFOR in TRIGGER shorthand must use a parenthesized block",
			Hint:    "Use TRIGGER (...), e.g. TRIGGER (WAITFOR EVENT \"ready\" IN target).",
		}, "WAITFOR EVENT TRIGGER shorthand rejects nested WAITFOR"),
		Failure(`
			LET obs = {}
			LET ok = WAITFOR EVENT "test" IN obs
				TRIGGER 1 + 2
				TIMEOUT 5s
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected trigger statement after 'TRIGGER' in WAITFOR EVENT",
			Hint:    "Use a side-effect statement or TRIGGER (...), e.g. TRIGGER target <- \"click\".",
		}, "WAITFOR EVENT TRIGGER shorthand rejects arbitrary expressions"),
	})
}

func TestSyntaxErrorsWaitforEmptyGroupDoesNotHijackEarlierError(t *testing.T) {
	query := "RETURN [1,,2]\nRETURN WAITFOR ANY {}"
	_, err := compiler.New().Compile(source.NewAnonymous(query))
	if err == nil {
		t.Fatal("expected compilation error")
	}

	diag := firstCompilationError(err)
	if diag == nil {
		t.Fatalf("expected diagnostic, got %v", err)
	}

	if strings.Contains(diag.Message, "WAITFOR ANY group") {
		t.Fatalf("earlier syntax error was replaced by later empty-group diagnostic: %s", diag.Message)
	}
}

func TestSyntaxErrorsWaitforEmptyGroupsUseTheirOwnSpans(t *testing.T) {
	query := "LET first = WAITFOR ANY {}\nRETURN WAITFOR VALUE ALL {}"
	_, err := compiler.New().Compile(source.NewAnonymous(query))
	if err == nil {
		t.Fatal("expected compilation error")
	}

	wantLines := map[string]int{
		"WAITFOR ANY group requires at least one arm":       1,
		"WAITFOR VALUE ALL group requires at least one arm": 2,
	}
	found := make(map[string]bool, len(wantLines))
	set, ok := err.(*pkgdiagnostics.DiagnosticSet)
	if !ok {
		t.Fatalf("expected diagnostic set, got %T", err)
	}

	for _, diag := range set.Errors() {
		wantLine, relevant := wantLines[diag.Message]
		if !relevant {
			continue
		}

		if len(diag.Spans) == 0 {
			t.Fatalf("diagnostic %q has no span", diag.Message)
		}

		span := diag.Spans[0].Span
		if got := query[span.Start:span.End]; got != "{}" {
			t.Fatalf("diagnostic %q points at %q", diag.Message, got)
		}

		line, _ := diag.Source.LocationAt(span)
		if line != wantLine {
			t.Fatalf("diagnostic %q points at line %d, want %d", diag.Message, line, wantLine)
		}

		found[diag.Message] = true
	}
	for message := range wantLines {
		if !found[message] {
			t.Fatalf("missing diagnostic %q in %v", message, err)
		}
	}
}
