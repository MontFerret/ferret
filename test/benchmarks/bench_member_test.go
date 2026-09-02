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

func BenchmarkMemberAccess_Short_None(b *testing.B) {
	RunBenchmarkNone(b, memberAccessShort)
}

func BenchmarkMemberAccess_Short_Full(b *testing.B) {
	RunBenchmarkFull(b, memberAccessShort)
}

func BenchmarkMemberAccess_Long_None(b *testing.B) {
	RunBenchmarkNone(b, memberAccessLong)
}

func BenchmarkMemberAccess_Long_Full(b *testing.B) {
	RunBenchmarkFull(b, memberAccessLong)
}

func BenchmarkUnknownMemberAccess_Short_None(b *testing.B) {
	RunBenchmarkNone(b, unknownMemberAccessShort, WithParam("obj", map[string]any{"foo": "bar"}))
}

func BenchmarkUnknownMemberAccess_Short_Full(b *testing.B) {
	RunBenchmarkFull(b, unknownMemberAccessShort, WithParam("obj", map[string]any{"foo": "bar"}))
}

func BenchmarkUnknownMemberAccess_Long_None(b *testing.B) {
	RunBenchmarkNone(b, unknownMemberAccessLong, WithParam("obj", map[string]any{"foo": "bar", "bar": map[string]any{"qaz": map[string]any{"abc": 42}}}))
}

func BenchmarkUnknownMemberAccess_Long_Full(b *testing.B) {
	RunBenchmarkFull(b, unknownMemberAccessLong, WithParam("obj", map[string]any{"foo": "bar", "bar": map[string]any{"qaz": map[string]any{"abc": 42}}}))
}
