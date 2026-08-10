package formatter_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterWaitForGroups(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`RETURN WAITFOR VALUE ALL { foo WHEN .ready WHEN .valid bar } TIMEOUT 10s EVERY 100ms BACKOFF LINEAR JITTER 0.2 ON ERROR RETURN NONE ON TIMEOUT RETURN []`, `RETURN WAITFOR VALUE ALL {
    foo
        WHEN .ready
        WHEN .valid
    bar
}
TIMEOUT 10s
EVERY 100ms
BACKOFF LINEAR
JITTER 0.2
ON TIMEOUT RETURN []
ON ERROR RETURN NONE`),
		S(`RETURN WAITFOR EVENT ANY { "navigation" IN page "dialog" IN page WHEN .type == "confirm" } TRIGGER (page <- "click") TIMEOUT 10s`, `RETURN WAITFOR EVENT ANY {
    "navigation" IN page
    "dialog" IN page
        WHEN .type == "confirm"
}
TRIGGER (
    page <- "click"
)
TIMEOUT 10s`),
		S(`RETURN WAITFOR EVENT ALL { // subscriptions
"response" IN page // server response
WHEN .status == 200
// browser download
"download" IN browser
}`, `RETURN WAITFOR EVENT ALL { // subscriptions
    "response" IN page // server response
        WHEN .status == 200
    // browser download
    "download" IN browser
}`),
		S(`RETURN WAITFOR ANY { ready }`, `RETURN WAITFOR ANY {
    ready
}`),
	})
}

func TestFormatterWaitForGroupComments(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"event tail comments": {
			input: `RETURN WAITFOR EVENT ANY {
"ready" IN source
} // after group
TRIGGER source <- "go" // after trigger
TIMEOUT 1s`,
			want: `RETURN WAITFOR EVENT ANY {
    "ready" IN source
} // after group
TRIGGER source <- "go" // after trigger
TIMEOUT 1s`,
		},
		"predicate schedule comments": {
			input: `RETURN WAITFOR VALUE ALL {
ready
} // before timeout
TIMEOUT 1s // before every
EVERY 10ms // before backoff
BACKOFF LINEAR // before jitter
JITTER 0.1`,
			want: `RETURN WAITFOR VALUE ALL {
    ready
} // before timeout
TIMEOUT 1s // before every
EVERY 10ms // before backoff
BACKOFF LINEAR // before jitter
JITTER 0.1`,
		},
		"blank line after inline comment": {
			input: `RETURN WAITFOR ANY {
ready
} // before timeout

ON TIMEOUT RETURN "timeout"`,
			want: `RETURN WAITFOR ANY {
    ready
} // before timeout

ON TIMEOUT RETURN "timeout"`,
		},
		"recovery comments retain source order": {
			input: `RETURN WAITFOR ANY {
ready
} // before error
ON ERROR RETURN "error" // before timeout recovery
ON TIMEOUT RETURN "timeout"`,
			want: `RETURN WAITFOR ANY {
    ready
} // before error
ON ERROR RETURN "error" // before timeout recovery
ON TIMEOUT RETURN "timeout"`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			first := formatSource(t, test.input)
			if first != test.want {
				t.Fatalf("formatted output mismatch\nwant:\n%q\ngot:\n%q", test.want, first)
			}

			second := formatSource(t, first)
			if second != first {
				t.Fatalf("formatter is not idempotent\nfirst:\n%q\nsecond:\n%q", first, second)
			}
		})
	}
}
