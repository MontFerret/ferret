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

const queryableConfiguredQuery = `
RETURN QUERY ".items" IN @doc USING css
	WITH { value: 1 }
	OPTIONS { timeout: 5000 }`

func BenchmarkQueryableModifiers_None(b *testing.B) {
	RunBenchmarkNone(b, queryableModifiersQuery, vm.WithParam("doc", newBenchmarkQueryable()))
}

func BenchmarkQueryableModifiers_Full(b *testing.B) {
	RunBenchmarkFull(b, queryableModifiersQuery, vm.WithParam("doc", newBenchmarkQueryable()))
}

func BenchmarkQueryableConfigured_None(b *testing.B) {
	RunBenchmarkNone(b, queryableConfiguredQuery, vm.WithParam("doc", newBenchmarkQueryable()))
}

func BenchmarkQueryableConfigured_Full(b *testing.B) {
	RunBenchmarkFull(b, queryableConfiguredQuery, vm.WithParam("doc", newBenchmarkQueryable()))
}
