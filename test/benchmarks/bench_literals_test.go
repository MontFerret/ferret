package benchmarks_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

const (
	arrayLiteralsQuery          = `RETURN ["foo", "bar", "qaz", "abc"]`
	arraySpreadQuery            = `LET left = ["foo", "bar"] LET right = ["qaz", "abc"] RETURN [...left, ...right]`
	arrayHostSpreadQuery        = `RETURN [...SOURCE()]`
	objectLiteralsQuery         = `RETURN { "foo": 1, "bar": 2, "qaz": 3, "abc": 4 }`
	objectSpreadQuery           = `LET left = { "foo": 1, "bar": 2 } LET right = { "qaz": 3, "abc": 4 } RETURN { ...left, ...right }`
	objectHostSpreadQuery       = `RETURN { ...SOURCE() }`
	objectComputedLiteralsQuery = `RETURN { ["foo"]: 1, ["bar"]: 2, [1]: 3 }`
)

func BenchmarkArrayLiterals_O0(b *testing.B) {
	RunBenchmarkO0(b, arrayLiteralsQuery)
}

func BenchmarkArrayLiterals_O1(b *testing.B) {
	RunBenchmarkO1(b, arrayLiteralsQuery)
}

func BenchmarkArraySpread_O0(b *testing.B) {
	RunBenchmarkO0(b, arraySpreadQuery)
}

func BenchmarkArraySpread_O1(b *testing.B) {
	RunBenchmarkO1(b, arraySpreadQuery)
}

func BenchmarkArraySpreadHostList_O0(b *testing.B) {
	RunBenchmarkO0(b, arrayHostSpreadQuery, benchmarkSpreadSource(struct{ runtime.List }{
		List: runtime.NewArrayWith(
			runtime.NewString("foo"),
			runtime.NewString("bar"),
			runtime.NewString("qaz"),
			runtime.NewString("abc"),
		),
	}))
}

func BenchmarkArraySpreadHostList_O1(b *testing.B) {
	RunBenchmarkO1(b, arrayHostSpreadQuery, benchmarkSpreadSource(struct{ runtime.List }{
		List: runtime.NewArrayWith(
			runtime.NewString("foo"),
			runtime.NewString("bar"),
			runtime.NewString("qaz"),
			runtime.NewString("abc"),
		),
	}))
}

func BenchmarkArraySpreadHostSnapshot_O0(b *testing.B) {
	values := runtime.NewArrayWith(
		runtime.NewString("foo"),
		runtime.NewString("bar"),
		runtime.NewString("qaz"),
		runtime.NewString("abc"),
	)

	RunBenchmarkO0(b, arrayHostSpreadQuery, benchmarkSpreadSource(&benchmarkSnapshotList{
		List:     values,
		snapshot: values,
	}))
}

func BenchmarkArraySpreadHostSnapshot_O1(b *testing.B) {
	values := runtime.NewArrayWith(
		runtime.NewString("foo"),
		runtime.NewString("bar"),
		runtime.NewString("qaz"),
		runtime.NewString("abc"),
	)

	RunBenchmarkO1(b, arrayHostSpreadQuery, benchmarkSpreadSource(&benchmarkSnapshotList{
		List:     values,
		snapshot: values,
	}))
}

func BenchmarkObjectLiterals_O0(b *testing.B) {
	RunBenchmarkO0(b, objectLiteralsQuery)
}

func BenchmarkObjectLiterals_O1(b *testing.B) {
	RunBenchmarkO1(b, objectLiteralsQuery)
}

func BenchmarkObjectSpread_O0(b *testing.B) {
	RunBenchmarkO0(b, objectSpreadQuery)
}

func BenchmarkObjectSpread_O1(b *testing.B) {
	RunBenchmarkO1(b, objectSpreadQuery)
}

func BenchmarkObjectSpreadHostObjectLike_O0(b *testing.B) {
	RunBenchmarkO0(b, objectHostSpreadQuery, benchmarkSpreadSource(&benchmarkObjectLike{
		Map: benchmarkSpreadObject(),
	}))
}

func BenchmarkObjectSpreadHostObjectLike_O1(b *testing.B) {
	RunBenchmarkO1(b, objectHostSpreadQuery, benchmarkSpreadSource(&benchmarkObjectLike{
		Map: benchmarkSpreadObject(),
	}))
}

func BenchmarkObjectSpreadHostSnapshot_O0(b *testing.B) {
	object := benchmarkSpreadObject()

	RunBenchmarkO0(b, objectHostSpreadQuery, benchmarkSpreadSource(&benchmarkSnapshotObject{
		Map:      object,
		snapshot: object,
	}))
}

func BenchmarkObjectSpreadHostSnapshot_O1(b *testing.B) {
	object := benchmarkSpreadObject()

	RunBenchmarkO1(b, objectHostSpreadQuery, benchmarkSpreadSource(&benchmarkSnapshotObject{
		Map:      object,
		snapshot: object,
	}))
}

func BenchmarkObjectComputedLiterals_O0(b *testing.B) {
	RunBenchmarkO0(b, objectComputedLiteralsQuery)
}

func BenchmarkObjectComputedLiterals_O1(b *testing.B) {
	RunBenchmarkO1(b, objectComputedLiteralsQuery)
}

func benchmarkSpreadSource(value runtime.Value) vm.EnvironmentOption {
	return vm.WithFunction("SOURCE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return value, nil
	})
}

func benchmarkSpreadObject() *runtime.Object {
	return runtime.NewObjectWith(map[string]runtime.Value{
		"foo": runtime.NewInt(1),
		"bar": runtime.NewInt(2),
		"qaz": runtime.NewInt(3),
		"abc": runtime.NewInt(4),
	})
}
