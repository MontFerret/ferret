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
