package formatter_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterWaitForGroups(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`return waitfor value all { foo when .ready when .valid bar } timeout 10s every 100ms backoff LINEAR jitter 0.2 on error return none on timeout return []`, `return waitfor value all {
    foo
        when .ready
        when .valid
    bar
}
timeout 10s
every 100ms
backoff LINEAR
jitter 0.2
on timeout return []
on error return none`),
		S(`return waitfor event any { "navigation" in page "dialog" in page when .type == "confirm" } trigger (page <- "click") timeout 10s`, `return waitfor event any {
    "navigation" in page
    "dialog" in page
        when .type == "confirm"
}
trigger (
    page <- "click"
)
timeout 10s`),
		S(`return waitfor event all { // subscriptions
"response" in page // server response
when .status == 200
// browser download
"download" in browser
}`, `return waitfor event all { // subscriptions
    "response" in page // server response
        when .status == 200
    // browser download
    "download" in browser
}`),
		S(`return waitfor any { ready }`, `return waitfor any {
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
			input: `return waitfor event any {
"ready" in source
} // after group
trigger source <- "go" // after trigger
timeout 1s`,
			want: `return waitfor event any {
    "ready" in source
} // after group
trigger source <- "go" // after trigger
timeout 1s`,
		},
		"predicate schedule comments": {
			input: `return waitfor value all {
ready
} // before timeout
timeout 1s // before every
every 10ms // before backoff
backoff LINEAR // before jitter
jitter 0.1`,
			want: `return waitfor value all {
    ready
} // before timeout
timeout 1s // before every
every 10ms // before backoff
backoff LINEAR // before jitter
jitter 0.1`,
		},
		"blank line after inline comment": {
			input: `return waitfor any {
ready
} // before timeout

on timeout return "timeout"`,
			want: `return waitfor any {
    ready
} // before timeout

on timeout return "timeout"`,
		},
		"recovery comments retain source order": {
			input: `return waitfor any {
ready
} // before error
on error return "error" // before timeout recovery
on timeout return "timeout"`,
			want: `return waitfor any {
    ready
} // before error
on error return "error" // before timeout recovery
on timeout return "timeout"`,
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
