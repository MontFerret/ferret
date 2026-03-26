package benchmarks_test

import "testing"

func BenchmarkRegexp_Loop_O0(b *testing.B) {
	RunBenchmarkO0(b, `
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
FOR i IN users
	FILTER i.name =~ "^[A-D].*"
  	RETURN i.name
`)
}

func BenchmarkRegexp_Loop_O1(b *testing.B) {
	RunBenchmarkO1(b, `
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
FOR i IN users
	FILTER i.name =~ "^[A-D].*"
  	RETURN i.name
`)
}
