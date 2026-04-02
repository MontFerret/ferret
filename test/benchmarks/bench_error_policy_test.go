package benchmarks_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

const suppressedHostCallQuery = `
RETURN FAIL() ON ERROR SUPPRESS`

func BenchmarkSuppressedHostCall_O0(b *testing.B) {
	boom := errors.New("boom")

	RunBenchmarkO0(b, suppressedHostCallQuery, vm.WithFunction("FAIL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, boom
	}))
}

func BenchmarkSuppressedHostCall_O1(b *testing.B) {
	boom := errors.New("boom")

	RunBenchmarkO1(b, suppressedHostCallQuery, vm.WithFunction("FAIL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, boom
	}))
}
