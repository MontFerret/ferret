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
RETURN STEP() ON ERROR RETRY 2 DELAY 0s OR RETURN NONE`

const waitForTimeoutReturnNoneQuery = `
RETURN WAITFOR VALUE NONE TIMEOUT 1ms EVERY 0ms ON TIMEOUT RETURN NONE`

const groupedForRetryQuery = `
LET xs = (FOR i IN [1, 2] LET y = STEP() RETURN y + i) ON ERROR RETRY 1 OR RETURN []
RETURN xs`

const waitForEventRetryQuery = `
RETURN WAITFOR EVENT "test" IN SOURCE() TIMEOUT 20ms ON TIMEOUT RETURN "timeout" ON ERROR RETRY 2 DELAY 0s OR RETURN "error"`

func BenchmarkSuppressedHostCall_None(b *testing.B) {
	boom := errors.New("boom")

	RunBenchmarkNone(b, suppressedHostCallQuery, vm.WithFunction("FAIL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, boom
	}))
}

func BenchmarkSuppressedHostCall_Full(b *testing.B) {
	boom := errors.New("boom")

	RunBenchmarkFull(b, suppressedHostCallQuery, vm.WithFunction("FAIL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, boom
	}))
}

func BenchmarkRetriedHostCall_None(b *testing.B) {
	callCount := 0

	RunBenchmarkNone(b, retriedHostCallQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		callCount++
		if callCount%3 != 0 {
			return runtime.None, errors.New("boom")
		}

		return runtime.NewInt(1), nil
	}))
}

func BenchmarkRetriedHostCall_Full(b *testing.B) {
	callCount := 0

	RunBenchmarkFull(b, retriedHostCallQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		callCount++
		if callCount%3 != 0 {
			return runtime.None, errors.New("boom")
		}

		return runtime.NewInt(1), nil
	}))
}

func BenchmarkRetriedHostCallFallback_None(b *testing.B) {
	RunBenchmarkNone(b, retriedHostCallQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, errors.New("boom")
	}))
}

func BenchmarkRetriedHostCallFallback_Full(b *testing.B) {
	RunBenchmarkFull(b, retriedHostCallQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, errors.New("boom")
	}))
}

func BenchmarkWaitForTimeoutReturnNone_None(b *testing.B) {
	RunBenchmarkNone(b, waitForTimeoutReturnNoneQuery)
}

func BenchmarkWaitForTimeoutReturnNone_Full(b *testing.B) {
	RunBenchmarkFull(b, waitForTimeoutReturnNoneQuery)
}

func BenchmarkGroupedForRetry_None(b *testing.B) {
	callCount := 0

	RunBenchmarkNone(b, groupedForRetryQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		callCount++
		if callCount%3 == 1 {
			return runtime.None, errors.New("boom")
		}

		return runtime.NewInt(10), nil
	}))
}

func BenchmarkGroupedForRetry_Full(b *testing.B) {
	callCount := 0

	RunBenchmarkFull(b, groupedForRetryQuery, vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		callCount++
		if callCount%3 == 1 {
			return runtime.None, errors.New("boom")
		}

		return runtime.NewInt(10), nil
	}))
}

func BenchmarkWaitForEventRetry_None(b *testing.B) {
	sourceCalls := 0

	RunBenchmarkNone(b, waitForEventRetryQuery, vm.WithFunction("SOURCE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		sourceCalls++
		if sourceCalls%3 != 0 {
			return runtime.None, errors.New("boom")
		}

		return mock.NewObservable([]runtime.Value{mock.NewTestEventType("match")}), nil
	}))
}

func BenchmarkWaitForEventRetry_Full(b *testing.B) {
	sourceCalls := 0

	RunBenchmarkFull(b, waitForEventRetryQuery, vm.WithFunction("SOURCE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		sourceCalls++
		if sourceCalls%3 != 0 {
			return runtime.None, errors.New("boom")
		}

		return mock.NewObservable([]runtime.Value{mock.NewTestEventType("match")}), nil
	}))
}
