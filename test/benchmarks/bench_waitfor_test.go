package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec/mock"
)

const waitForValuePresentQuery = `
RETURN WAITFOR VALUE @candidate`

const waitForEventPresentQuery = `
LET source = @source
RETURN WAITFOR EVENT "test" IN source`

const waitForEventCompilerQuery = `
LET eventName = @eventName
LET source = @source
RETURN WAITFOR EVENT eventName IN source
	WHEN .type == eventName
	TIMEOUT 10ms`

const waitForValueAnyPresentQuery = `
RETURN WAITFOR VALUE ANY {
	@first
	@second
}`

const waitForValueAllPresentQuery = `
RETURN WAITFOR VALUE ALL {
	@first
	@second
}`

const waitForEventAnyPresentQuery = `
LET first = @first
LET second = @second
RETURN WAITFOR EVENT ANY {
	"first" IN first
	"second" IN second
}`

const waitForEventAllPresentQuery = `
LET first = @first
LET second = @second
RETURN WAITFOR EVENT ALL {
	"first" IN first
	"second" IN second
}`

func BenchmarkWaitForValuePresent_O0(b *testing.B) {
	RunBenchmarkO0(b, waitForValuePresentQuery, WithParam("candidate", []any{1}))
}

func BenchmarkWaitForValuePresent_O1(b *testing.B) {
	RunBenchmarkO1(b, waitForValuePresentQuery, WithParam("candidate", []any{1}))
}

func BenchmarkWaitForEventPresent_O0(b *testing.B) {
	RunBenchmarkO0(b, waitForEventPresentQuery, vm.WithParam("source", newBenchmarkWaitForObservable()))
}

func BenchmarkWaitForEventPresent_O1(b *testing.B) {
	RunBenchmarkO1(b, waitForEventPresentQuery, vm.WithParam("source", newBenchmarkWaitForObservable()))
}

func BenchmarkCompilerCompileWaitForEvent_O0(b *testing.B) {
	benchmarkCompileQuery(b, waitForEventCompilerQuery, compiler.OptimizationNone)
}

func BenchmarkCompilerCompileWaitForEvent_O1(b *testing.B) {
	benchmarkCompileQuery(b, waitForEventCompilerQuery, compiler.OptimizationFull)
}

func BenchmarkWaitForValueAnyPresent_O0(b *testing.B) {
	RunBenchmarkO0(b, waitForValueAnyPresentQuery, WithParam("first", []any{1}), WithParam("second", []any{2}))
}

func BenchmarkWaitForValueAnyPresent_O1(b *testing.B) {
	RunBenchmarkO1(b, waitForValueAnyPresentQuery, WithParam("first", []any{1}), WithParam("second", []any{2}))
}

func BenchmarkWaitForValueAllPresent_O0(b *testing.B) {
	RunBenchmarkO0(b, waitForValueAllPresentQuery, WithParam("first", []any{1}), WithParam("second", []any{2}))
}

func BenchmarkWaitForValueAllPresent_O1(b *testing.B) {
	RunBenchmarkO1(b, waitForValueAllPresentQuery, WithParam("first", []any{1}), WithParam("second", []any{2}))
}

func BenchmarkWaitForEventAnyPresent_O0(b *testing.B) {
	RunBenchmarkO0(b, waitForEventAnyPresentQuery,
		vm.WithParam("first", newBenchmarkWaitForObservable()),
		vm.WithParam("second", newBenchmarkWaitForObservable()),
	)
}

func BenchmarkWaitForEventAnyPresent_O1(b *testing.B) {
	RunBenchmarkO1(b, waitForEventAnyPresentQuery,
		vm.WithParam("first", newBenchmarkWaitForObservable()),
		vm.WithParam("second", newBenchmarkWaitForObservable()),
	)
}

func BenchmarkWaitForEventAllPresent_O0(b *testing.B) {
	RunBenchmarkO0(b, waitForEventAllPresentQuery,
		vm.WithParam("first", newBenchmarkWaitForObservable()),
		vm.WithParam("second", newBenchmarkWaitForObservable()),
	)
}

func BenchmarkWaitForEventAllPresent_O1(b *testing.B) {
	RunBenchmarkO1(b, waitForEventAllPresentQuery,
		vm.WithParam("first", newBenchmarkWaitForObservable()),
		vm.WithParam("second", newBenchmarkWaitForObservable()),
	)
}

func newBenchmarkWaitForObservable() *mock.Observable {
	return mock.NewObservable([]runtime.Value{mock.NewTestEventType("test")})
}
