package parser_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func TestWaitForEventSourceAcceptsGeneralExpressions(t *testing.T) {
	tests := []string{
		"page",
		"@source",
		"session.page",
		"pages[@index]",
		"GET_PAGE(@id)",
		`QUERY ONE ".page" IN registry`,
		"(@page)",
		"@page ?? DEFAULT_PAGE()",
		"OPTIONS",
	}

	for _, expression := range tests {
		t.Run(expression, func(t *testing.T) {
			event := parseFirstWaitForEvent(t, `RETURN WAITFOR EVENT "message" IN `+expression+` TIMEOUT 1MS`)
			source := event.WaitForEventSource().(*fql.WaitForEventSourceContext).Expression()
			if source == nil {
				t.Fatal("expected source expression")
			}

			if got, want := source.GetText(), stripSpaces(expression); got != want {
				t.Fatalf("source expression = %q, want %q", got, want)
			}
		})
	}
}

func TestWaitForEventNameAcceptsGeneralExpressions(t *testing.T) {
	tests := []string{
		`"message"`,
		"@event",
		"session.event",
		"names[@index]",
		"GET_EVENT(@id)",
		`(@event ?? "message")`,
		`"network." + @phase`,
		"TIMEOUT",
	}

	for _, expression := range tests {
		t.Run(expression, func(t *testing.T) {
			event := parseFirstWaitForEvent(t, `RETURN WAITFOR EVENT `+expression+` IN @source`)
			name := event.WaitForEventName().(*fql.WaitForEventNameContext).Expression()
			if name == nil {
				t.Fatal("expected event-name expression")
			}

			if got, want := name.GetText(), stripSpaces(expression); got != want {
				t.Fatalf("event-name expression = %q, want %q", got, want)
			}
		})
	}
}

func TestWaitForEventInDelimiterUsesNaturalExpressionBoundary(t *testing.T) {
	event := parseFirstWaitForEvent(t, `RETURN WAITFOR EVENT @left IN @right IN @source`)
	name := event.WaitForEventName().(*fql.WaitForEventNameContext).Expression()
	source := event.WaitForEventSource().(*fql.WaitForEventSourceContext).Expression()

	if got, want := name.GetText(), "@leftIN@right"; got != want {
		t.Fatalf("event-name expression = %q, want %q", got, want)
	}

	if got, want := source.GetText(), "@source"; got != want {
		t.Fatalf("source expression = %q, want %q", got, want)
	}

	groupedSource := parseFirstWaitForEvent(t, `RETURN WAITFOR EVENT @event IN (@candidate IN @sources)`)
	grouped := groupedSource.WaitForEventSource().(*fql.WaitForEventSourceContext).Expression()
	if got, want := grouped.GetText(), "(@candidateIN@sources)"; got != want {
		t.Fatalf("grouped source expression = %q, want %q", got, want)
	}
}

func TestWaitForEventExpressionDoesNotConsumeClauses(t *testing.T) {
	const input = `RETURN WAITFOR EVENT (@kind ?? "message") IN (@page ?? fallback)
OPTIONS { capture: true }
WHEN .ok
TRIGGER ()
TIMEOUT 1MS
ON TIMEOUT RETURN NONE`

	program, errors := parseQueryPayloadProgram(input)
	if errors.HasErrors() {
		t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
	}

	wait := mustFindFirst[*fql.WaitForExpressionContext](t, program)
	event := wait.WaitForEventExpression().(*fql.WaitForEventExpressionContext)
	name := event.WaitForEventName().(*fql.WaitForEventNameContext).Expression()
	source := event.WaitForEventSource().(*fql.WaitForEventSourceContext).Expression()

	if got, want := name.GetText(), `(@kind??"message")`; got != want {
		t.Fatalf("event-name expression = %q, want %q", got, want)
	}

	if got, want := source.GetText(), "(@page??fallback)"; got != want {
		t.Fatalf("source expression = %q, want %q", got, want)
	}

	if event.OptionsClause() == nil || len(event.AllEventFilterClause()) != 1 {
		t.Fatal("expected OPTIONS and WHEN clauses")
	}

	tail := event.WaitForEventTail()
	if tail == nil || tail.WaitForTriggerClause() == nil || tail.TimeoutClause() == nil {
		t.Fatal("expected TRIGGER and TIMEOUT clauses")
	}

	if wait.RecoveryTails() == nil {
		t.Fatal("expected ON TIMEOUT recovery tail")
	}
}

func TestWaitForEventGroupEntriesAcceptGeneralExpressions(t *testing.T) {
	const input = `RETURN WAITFOR EVENT ANY {
    (@firstName ?? "first") IN (@first ?? fallback)
    "second" IN GET_PAGE(@id)
} TIMEOUT 1MS`

	program, errors := parseQueryPayloadProgram(input)
	if errors.HasErrors() {
		t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
	}

	group := mustFindFirst[*fql.WaitForEventGroupExpressionContext](t, program)
	entries := group.AllWaitForEventGroupEntry()
	if got, want := len(entries), 2; got != want {
		t.Fatalf("entry count = %d, want %d", got, want)
	}

	first := entries[0].(*fql.WaitForEventGroupEntryContext)
	firstName := first.WaitForEventName().(*fql.WaitForEventNameContext).Expression()
	firstSource := first.WaitForEventSource().(*fql.WaitForEventSourceContext).Expression()
	if got, want := firstName.GetText(), `(@firstName??"first")`; got != want {
		t.Fatalf("first event-name expression = %q, want %q", got, want)
	}

	if got, want := firstSource.GetText(), "(@first??fallback)"; got != want {
		t.Fatalf("first source expression = %q, want %q", got, want)
	}

	second := entries[1].(*fql.WaitForEventGroupEntryContext)
	secondSource := second.WaitForEventSource().(*fql.WaitForEventSourceContext).Expression()
	if got, want := secondSource.GetText(), "GET_PAGE(@id)"; got != want {
		t.Fatalf("second source expression = %q, want %q", got, want)
	}
}

func parseFirstWaitForEvent(t *testing.T, input string) *fql.WaitForEventExpressionContext {
	t.Helper()

	program, errors := parseQueryPayloadProgram(input)
	if errors.HasErrors() {
		t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
	}

	return mustFindFirst[*fql.WaitForEventExpressionContext](t, program)
}
