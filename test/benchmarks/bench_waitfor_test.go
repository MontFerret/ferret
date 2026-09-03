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

func BenchmarkWaitForValuePresent_None(b *testing.B) {
	RunBenchmarkNone(b, waitForValuePresentQuery, WithParam("candidate", []any{1}))
}

func BenchmarkWaitForValuePresent_Basic(b *testing.B) {
	RunBenchmarkBasic(b, waitForValuePresentQuery, WithParam("candidate", []any{1}))
}

func BenchmarkWaitForValuePresent_Full(b *testing.B) {
	RunBenchmarkFull(b, waitForValuePresentQuery, WithParam("candidate", []any{1}))
}

func BenchmarkWaitForEventPresent_None(b *testing.B) {
	RunBenchmarkNone(b, waitForEventPresentQuery, vm.WithParam("source", newBenchmarkWaitForObservable()))
}

func BenchmarkWaitForEventPresent_Basic(b *testing.B) {
	RunBenchmarkBasic(b, waitForEventPresentQuery, vm.WithParam("source", newBenchmarkWaitForObservable()))
}

func BenchmarkWaitForEventPresent_Full(b *testing.B) {
	RunBenchmarkFull(b, waitForEventPresentQuery, vm.WithParam("source", newBenchmarkWaitForObservable()))
}

func BenchmarkCompilerCompileWaitForEvent_None(b *testing.B) {
	benchmarkCompileQuery(b, waitForEventCompilerQuery, compiler.None)
}

func BenchmarkCompilerCompileWaitForEvent_Basic(b *testing.B) {
	benchmarkCompileQuery(b, waitForEventCompilerQuery, compiler.Basic)
}

func BenchmarkCompilerCompileWaitForEvent_Full(b *testing.B) {
	benchmarkCompileQuery(b, waitForEventCompilerQuery, compiler.Full)
}

func BenchmarkWaitForValueAnyPresent_None(b *testing.B) {
	RunBenchmarkNone(b, waitForValueAnyPresentQuery, WithParam("first", []any{1}), WithParam("second", []any{2}))
}

func BenchmarkWaitForValueAnyPresent_Basic(b *testing.B) {
	RunBenchmarkBasic(b, waitForValueAnyPresentQuery, WithParam("first", []any{1}), WithParam("second", []any{2}))
}

func BenchmarkWaitForValueAnyPresent_Full(b *testing.B) {
	RunBenchmarkFull(b, waitForValueAnyPresentQuery, WithParam("first", []any{1}), WithParam("second", []any{2}))
}

func BenchmarkWaitForValueAllPresent_None(b *testing.B) {
	RunBenchmarkNone(b, waitForValueAllPresentQuery, WithParam("first", []any{1}), WithParam("second", []any{2}))
}

func BenchmarkWaitForValueAllPresent_Basic(b *testing.B) {
	RunBenchmarkBasic(b, waitForValueAllPresentQuery, WithParam("first", []any{1}), WithParam("second", []any{2}))
}

func BenchmarkWaitForValueAllPresent_Full(b *testing.B) {
	RunBenchmarkFull(b, waitForValueAllPresentQuery, WithParam("first", []any{1}), WithParam("second", []any{2}))
}

func BenchmarkWaitForEventAnyPresent_None(b *testing.B) {
	RunBenchmarkNone(b, waitForEventAnyPresentQuery,
		vm.WithParam("first", newBenchmarkWaitForObservable()),
		vm.WithParam("second", newBenchmarkWaitForObservable()),
	)
}

func BenchmarkWaitForEventAnyPresent_Basic(b *testing.B) {
	RunBenchmarkBasic(b, waitForEventAnyPresentQuery,
		vm.WithParam("first", newBenchmarkWaitForObservable()),
		vm.WithParam("second", newBenchmarkWaitForObservable()),
	)
}

func BenchmarkWaitForEventAnyPresent_Full(b *testing.B) {
	RunBenchmarkFull(b, waitForEventAnyPresentQuery,
		vm.WithParam("first", newBenchmarkWaitForObservable()),
		vm.WithParam("second", newBenchmarkWaitForObservable()),
	)
}

func BenchmarkWaitForEventAllPresent_None(b *testing.B) {
	RunBenchmarkNone(b, waitForEventAllPresentQuery,
		vm.WithParam("first", newBenchmarkWaitForObservable()),
		vm.WithParam("second", newBenchmarkWaitForObservable()),
	)
}

func BenchmarkWaitForEventAllPresent_Basic(b *testing.B) {
	RunBenchmarkBasic(b, waitForEventAllPresentQuery,
		vm.WithParam("first", newBenchmarkWaitForObservable()),
		vm.WithParam("second", newBenchmarkWaitForObservable()),
	)
}

func BenchmarkWaitForEventAllPresent_Full(b *testing.B) {
	RunBenchmarkFull(b, waitForEventAllPresentQuery,
		vm.WithParam("first", newBenchmarkWaitForObservable()),
		vm.WithParam("second", newBenchmarkWaitForObservable()),
	)
}

func newBenchmarkWaitForObservable() *mock.Observable {
	return mock.NewObservable([]runtime.Value{mock.NewTestEventType("test")})
}
