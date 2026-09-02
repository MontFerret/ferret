package benchmarks_test

import "testing"

func BenchmarkRegexp_Loop_None(b *testing.B) {
	RunBenchmarkNone(b, `
LET users = [
  {
  	name: "Alice",
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
  	RETURN i.name
`)
}

func BenchmarkRegexp_Loop_Full(b *testing.B) {
	RunBenchmarkFull(b, `
LET users = [
  {
  	name: "Alice",
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
  	RETURN i.name
`)
}
