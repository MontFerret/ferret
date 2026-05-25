package internal

import (
	"bytes"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestClauseFormatter_TimeoutValueFormatsParam(t *testing.T) {
	input := "WAITFOR VALUE x TIMEOUT @t"
	program := parseProgram(t, input)
	timeout := mustFirst[*fql.TimeoutClauseContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, DefaultOptions())

	e.clause.formatTimeoutClause(timeout)
	if got := buf.String(); got != "TIMEOUT @t" {
		t.Fatalf("unexpected timeout formatting: %q", got)
	}
}

func TestClauseFormatter_EventFilterClauseUsesWhen(t *testing.T) {
	input := "WAITFOR EVENT \"test\" IN obs WHEN .type == \"match\" WHEN .visible"
	program := parseProgram(t, input)
	waitExpr := mustFirst[*fql.WaitForEventExpressionContext](t, program)
	filters := waitExpr.AllEventFilterClause()
	if len(filters) != 2 {
		t.Fatalf("expected two event filters, got %d", len(filters))
	}

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, DefaultOptions())

	e.clause.formatEventFilterClause(filters[1].(*fql.EventFilterClauseContext))
	if got := buf.String(); got != "WHEN .visible" {
		t.Fatalf("unexpected event filter formatting: %q", got)
	}
}

func TestClauseFormatter_WaitForPredicateWhenClause(t *testing.T) {
	input := "WAITFOR VALUE ready WHEN .state == \"ready\" WHEN .visible"
	program := parseProgram(t, input)
	waitExpr := mustFirst[*fql.WaitForPredicateExpressionContext](t, program)
	clauses := waitExpr.AllWaitForPredicateWhenClause()
	if len(clauses) != 2 {
		t.Fatalf("expected two predicate WHEN clauses, got %d", len(clauses))
	}

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, DefaultOptions())

	e.clause.formatWaitForPredicateWhenClause(clauses[1].(*fql.WaitForPredicateWhenClauseContext))
	if got := buf.String(); got != "WHEN .visible" {
		t.Fatalf("unexpected waitfor predicate WHEN formatting: %q", got)
	}
}
