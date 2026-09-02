package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

const (
	optionalMemberAccessShort = `
LET obj = NONE

RETURN obj?.abc
	`

	optionalMemberAccessShort2 = `
LET obj = {}

RETURN obj.abc
	`

	optionalMemberAccessLong = `
LET obj = NONE

RETURN obj?.foo?.["bar"]?.qaz.abc
	`

	optionalUnknownMemberAccessShort = `
LET obj = @obj

RETURN obj?.foo
	`

	optionalUnknownMemberAccessLong = `
LET obj = @obj

RETURN obj?.bar?.qaz?.abc
	`
)

func BenchmarkOptionalMemberAccess_Short_None(b *testing.B) {
	RunBenchmarkNone(b, optionalMemberAccessShort)
}

func BenchmarkOptionalMemberAccess_Short_Full(b *testing.B) {
	RunBenchmarkFull(b, optionalMemberAccessShort)
}

func BenchmarkOptionalMemberAccess_Short2_None(b *testing.B) {
	RunBenchmarkNone(b, optionalMemberAccessShort2)
}

func BenchmarkOptionalMemberAccess_Short2_Full(b *testing.B) {
	RunBenchmarkFull(b, optionalMemberAccessShort2)
}

func BenchmarkOptionalMemberAccess_Long_None(b *testing.B) {
	RunBenchmarkNone(b, optionalMemberAccessLong)
}

func BenchmarkOptionalMemberAccess_Long_Full(b *testing.B) {
	RunBenchmarkFull(b, optionalMemberAccessLong)
}

func BenchmarkOptionalUnknownMemberAccess_Short_None(b *testing.B) {
	RunBenchmarkNone(b, optionalUnknownMemberAccessShort, vm.WithParam("obj", runtime.None))
}

func BenchmarkOptionalUnknownMemberAccess_Short_Full(b *testing.B) {
	RunBenchmarkFull(b, optionalUnknownMemberAccessShort, vm.WithParam("obj", runtime.None))
}

func BenchmarkOptionalUnknownMemberAccess_Long_None(b *testing.B) {
	RunBenchmarkNone(b, optionalUnknownMemberAccessLong, vm.WithParam("obj", runtime.None))
}

func BenchmarkOptionalUnknownMemberAccess_Long_Full(b *testing.B) {
	RunBenchmarkFull(b, optionalUnknownMemberAccessLong, vm.WithParam("obj", runtime.None))
}
