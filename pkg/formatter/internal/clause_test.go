package internal

import (
	"bytes"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func TestClauseFormatter_TimeoutValueFormatsParam(t *testing.T) {
	input := "WAITFOR VALUE x TIMEOUT @t"
	program := parseProgram(t, input)
	timeout := mustFirst[*fql.TimeoutClauseContext](t, program)

	var buf bytes.Buffer
	e := newEngine(file.NewAnonymousSource(input), &buf, DefaultOptions())

	e.clause.formatTimeoutClause(timeout)
	if got := buf.String(); got != "TIMEOUT @t" {
		t.Fatalf("unexpected timeout formatting: %q", got)
	}
}
