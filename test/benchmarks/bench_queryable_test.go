package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/vm"
)

const queryableModifiersQuery = `
RETURN [
	QUERY EXISTS ".items" IN @doc USING css,
	QUERY COUNT ".items" IN @doc USING css,
	QUERY ONE ".items" IN @doc USING css,
]`

func BenchmarkQueryableModifiers_O0(b *testing.B) {
	RunBenchmarkO0(b, queryableModifiersQuery, vm.WithParam("doc", newBenchmarkQueryable()))
}

func BenchmarkQueryableModifiers_O1(b *testing.B) {
	RunBenchmarkO1(b, queryableModifiersQuery, vm.WithParam("doc", newBenchmarkQueryable()))
}
