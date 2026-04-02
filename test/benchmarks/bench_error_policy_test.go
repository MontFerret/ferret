package benchmarks_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	mock "github.com/MontFerret/ferret/v2/test/spec/mock"
)

const suppressedHostCallQuery = `
RETURN FAIL() ON ERROR RETURN NONE`

const retriedHostCallQuery = `
RETURN STEP() ON ERROR RETRY 2 DELAY 0 OR RETURN NONE`

const waitForTimeoutReturnNoneQuery = `
RETURN WAITFOR VALUE NONE TIMEOUT 1 EVERY 0 ON TIMEOUT RETURN NONE`

const groupedForRetryQuery = `
LET xs = (FOR i IN [1, 2] LET y = STEP() RETURN y + i) ON ERROR RETRY 1 OR RETURN []
RETURN xs`

const waitForEventRetryQuery = `
RETURN WAITFOR EVENT "test" IN SOURCE() TIMEOUT 20ms ON TIMEOUT RETURN "timeout" ON ERROR RETRY 2 DELAY 0 OR RETURN "error"`

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
	RunBenchmarkO0(b, retriedHostCallQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, errors.New("boom")
	}))
}

func BenchmarkRetriedHostCallFallback_O1(b *testing.B) {
	RunBenchmarkO1(b, retriedHostCallQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, errors.New("boom")
	}))
}

func BenchmarkWaitForTimeoutReturnNone_O0(b *testing.B) {
	RunBenchmarkO0(b, waitForTimeoutReturnNoneQuery)
}

func BenchmarkWaitForTimeoutReturnNone_O1(b *testing.B) {
	RunBenchmarkO1(b, waitForTimeoutReturnNoneQuery)
}

func BenchmarkGroupedForRetry_O0(b *testing.B) {
	callCount := 0

	RunBenchmarkO0(b, groupedForRetryQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		callCount++
		if callCount%3 == 1 {
			return runtime.None, errors.New("boom")
		}

		return runtime.NewInt(10), nil
	}))
}

func BenchmarkGroupedForRetry_O1(b *testing.B) {
	callCount := 0

	RunBenchmarkO1(b, groupedForRetryQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		callCount++
		if callCount%3 == 1 {
			return runtime.None, errors.New("boom")
		}

		return runtime.NewInt(10), nil
	}))
}

func BenchmarkWaitForEventRetry_O0(b *testing.B) {
	sourceCalls := 0

	RunBenchmarkO0(b, waitForEventRetryQuery, vm.WithFunction("SOURCE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		sourceCalls++
		if sourceCalls%3 != 0 {
			return runtime.None, errors.New("boom")
		}

		return mock.NewObservable([]runtime.Value{mock.NewTestEventType("match")}), nil
	}))
}

func BenchmarkWaitForEventRetry_O1(b *testing.B) {
	sourceCalls := 0

	RunBenchmarkO1(b, waitForEventRetryQuery, vm.WithFunction("SOURCE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		sourceCalls++
		if sourceCalls%3 != 0 {
			return runtime.None, errors.New("boom")
		}

		return mock.NewObservable([]runtime.Value{mock.NewTestEventType("match")}), nil
	}))
}
