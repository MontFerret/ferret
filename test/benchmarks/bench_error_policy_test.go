package benchmarks_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

const suppressedHostCallQuery = `
RETURN FAIL() ON ERROR RETURN NONE`

const retriedHostCallQuery = `
RETURN STEP() ON ERROR RETRY 2 DELAY 0 OR RETURN NONE`

const retriedHostCallFallbackQuery = `
RETURN STEP() ON ERROR RETRY 2 DELAY 0 OR RETURN NONE`

const waitForTimeoutReturnNoneQuery = `
RETURN WAITFOR VALUE NONE TIMEOUT 1 EVERY 0 ON TIMEOUT RETURN NONE`

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

func BenchmarkRetriedHostCall_O0(b *testing.B) {
	callCount := 0

	RunBenchmarkO0(b, retriedHostCallQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		callCount++
		if callCount%3 != 0 {
			return runtime.None, errors.New("boom")
		}

		return runtime.NewInt(1), nil
	}))
}

func BenchmarkRetriedHostCall_O1(b *testing.B) {
	callCount := 0

	RunBenchmarkO1(b, retriedHostCallQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		callCount++
		if callCount%3 != 0 {
			return runtime.None, errors.New("boom")
		}

		return runtime.NewInt(1), nil
	}))
}

func BenchmarkRetriedHostCallFallback_O0(b *testing.B) {
	RunBenchmarkO0(b, retriedHostCallFallbackQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, errors.New("boom")
	}))
}

func BenchmarkRetriedHostCallFallback_O1(b *testing.B) {
	RunBenchmarkO1(b, retriedHostCallFallbackQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, errors.New("boom")
	}))
}

func BenchmarkWaitForTimeoutReturnNone_O0(b *testing.B) {
	RunBenchmarkO0(b, waitForTimeoutReturnNoneQuery)
}

func BenchmarkWaitForTimeoutReturnNone_O1(b *testing.B) {
	RunBenchmarkO1(b, waitForTimeoutReturnNoneQuery)
}
