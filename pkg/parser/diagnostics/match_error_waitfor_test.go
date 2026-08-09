package diagnostics

import (
	"testing"

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
			mode, synchronization, span, ok := waitForEmptyGroupInSource(source.NewAnonymous(test.query))
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
