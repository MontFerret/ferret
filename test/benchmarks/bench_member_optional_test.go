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

func BenchmarkOptionalMemberAccess_Short_O0(b *testing.B) {
	RunBenchmarkO0(b, optionalMemberAccessShort)
}

func BenchmarkOptionalMemberAccess_Short_O1(b *testing.B) {
	RunBenchmarkO1(b, optionalMemberAccessShort)
}

func BenchmarkOptionalMemberAccess_Short2_O0(b *testing.B) {
	RunBenchmarkO0(b, optionalMemberAccessShort2)
}

func BenchmarkOptionalMemberAccess_Short2_O1(b *testing.B) {
	RunBenchmarkO1(b, optionalMemberAccessShort2)
}

func BenchmarkOptionalMemberAccess_Long_O0(b *testing.B) {
	RunBenchmarkO0(b, optionalMemberAccessLong)
}

func BenchmarkOptionalMemberAccess_Long_O1(b *testing.B) {
	RunBenchmarkO1(b, optionalMemberAccessLong)
}

func BenchmarkOptionalUnknownMemberAccess_Short_O0(b *testing.B) {
	RunBenchmarkO0(b, optionalUnknownMemberAccessShort, vm.WithParam("obj", runtime.None))
}

func BenchmarkOptionalUnknownMemberAccess_Short_O1(b *testing.B) {
	RunBenchmarkO1(b, optionalUnknownMemberAccessShort, vm.WithParam("obj", runtime.None))
}

func BenchmarkOptionalUnknownMemberAccess_Long_O0(b *testing.B) {
	RunBenchmarkO0(b, optionalUnknownMemberAccessLong, vm.WithParam("obj", runtime.None))
}

func BenchmarkOptionalUnknownMemberAccess_Long_O1(b *testing.B) {
	RunBenchmarkO1(b, optionalUnknownMemberAccessLong, vm.WithParam("obj", runtime.None))
}
