package benchmarks_test

import (
	"testing"
)

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

	unknownMemberAccessShort = `
LET obj = @obj

RETURN obj.foo
	`

	unknownMemberAccessLong = `
LET obj = @obj

RETURN obj.bar.qaz.abc
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

func BenchmarkUnknownMemberAccess_Short_O0(b *testing.B) {
	RunBenchmarkO0(b, unknownMemberAccessShort, WithParam("obj", map[string]any{"foo": "bar"}))
}

func BenchmarkUnknownMemberAccess_Short_O1(b *testing.B) {
	RunBenchmarkO1(b, unknownMemberAccessShort, WithParam("obj", map[string]any{"foo": "bar"}))
}

func BenchmarkUnknownMemberAccess_Long_O0(b *testing.B) {
	RunBenchmarkO0(b, unknownMemberAccessLong, WithParam("obj", map[string]any{"foo": "bar", "bar": map[string]any{"qaz": map[string]any{"abc": 42}}}))
}

func BenchmarkUnknownMemberAccess_Long_O1(b *testing.B) {
	RunBenchmarkO1(b, unknownMemberAccessLong, WithParam("obj", map[string]any{"foo": "bar", "bar": map[string]any{"qaz": map[string]any{"abc": 42}}}))
}
