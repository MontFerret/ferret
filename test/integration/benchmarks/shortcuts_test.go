package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/integration/base"
)

// Benchmark gate commands:
//
//	go test ./test/integration/benchmarks -run ^$ -bench UdfCalls -benchmem
//	go test ./test/integration/benchmarks -run ^$ -bench 'FunctionCall[0-4]?(Fallback)?_O[01]$' -benchmem
//	go test ./test/integration/benchmarks -run ^$ -bench 'MemberAccess|UnknownMemberAccess|OptionalMemberAccess' -benchmem
//	go test ./test/integration/benchmarks -run ^$ -bench 'AddConstString|TemplateLiteral|AddNumeric|AddConstNumeric' -benchmem
func RunBenchmarkO0(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	base.RunBenchmarkWithOptimization(b, expression, compiler.O0, opts...)
}

func RunBenchmarkO1(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	base.RunBenchmarkWithOptimization(b, expression, compiler.O1, opts...)
}
