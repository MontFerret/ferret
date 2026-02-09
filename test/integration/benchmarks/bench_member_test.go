package benchmarks_test

import "testing"

const (
	memberAccessShort = `
LET obj = {
	"abc": 42
}

RETURN obj.abc
	`

	memberAccessLong = `
LET obj = {
	"foo": { "bar": { "qaz": { "abc": 42 } } }
}

RETURN obj.foo.bar.qaz.abc
	`
)

func BenchmarkMemberAccess_Short_O0(b *testing.B) {
	RunBenchmarkO0(b, memberAccessShort)
}

func BenchmarkMemberAccess_Short_O1(b *testing.B) {
	RunBenchmarkO1(b, memberAccessShort)
}

func BenchmarkMemberAccess_Long_O0(b *testing.B) {
	RunBenchmarkO0(b, memberAccessLong)
}

func BenchmarkMemberAccess_Long_O1(b *testing.B) {
	RunBenchmarkO1(b, memberAccessLong)
}
