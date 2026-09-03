package benchmarks_test

import "testing"

// Preserve the legacy query bytes without mixed space/tab source indentation.
const regexpLoopQuery = `
LET users = [
  {
` + "  \t" + `name: "Alice",
  },
{
	name: "Bob",
},
{
	name: "Charlie",
},
{
	name: "Dave",
},
{
	name: "Eve",
}
]
RETURN FOR i IN users
	FILTER i.name =~ "^[A-D].*"
` + "  \t" + `RETURN i.name
`

func BenchmarkRegexp_Loop_None(b *testing.B) {
	RunBenchmarkNone(b, regexpLoopQuery)
}

func BenchmarkRegexp_Loop_Basic(b *testing.B) {
	RunBenchmarkBasic(b, regexpLoopQuery)
}

func BenchmarkRegexp_Loop_Full(b *testing.B) {
	RunBenchmarkFull(b, regexpLoopQuery)
}
