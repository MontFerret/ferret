package internal

import (
	"bytes"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func TestStatementFormatter_DispatchEventNameString(t *testing.T) {
	input := "DISPATCH \"evt\" IN target"
	program := parseProgram(t, input)
	eventName := mustFirst[*fql.DispatchEventNameContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymousSource(input), &buf, DefaultOptions())

	e.statement.formatDispatchEventName(eventName)
	if got := buf.String(); got != "\"evt\"" {
		t.Fatalf("unexpected dispatch event name formatting: %q", got)
	}
}
