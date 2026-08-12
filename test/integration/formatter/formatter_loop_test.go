package formatter_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterLoopBindings(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
for _ while i < 2
return i
`, `for while i < 2
    return i`),
		S(`
for _ do while false
return 1
`, `for do while false
    return 1`),
		S(`
for n while i < 2
return n
`, "for n while i < 2\n    return n"),
		S(`
for while i < 1
	dispatch "click" in @d
	return i
	`, "for while i < 1\n    dispatch \"click\" in @d\n    return i"),
		S(`
for while ready
waitfor event "navigation" in doc when .type == "match" timeout 10s
return ready
`, "for while ready\n    waitfor event \"navigation\" in doc when .type == \"match\" timeout 10s\n    return ready"),
		S(`
for order in query "/orders" in api with { query: { status: "open" } }
return order
`, "for order in query \"/orders\" in api with { query: { status: \"open\" } }\n    return order"),
		S(`
for order in (query "/orders" in api with { query: { status: "open" } })
return order
`, "for order in (query \"/orders\" in api with { query: { status: \"open\" } })\n    return order"),
	})
}
