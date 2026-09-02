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

func BenchmarkArrayLiterals_None(b *testing.B) {
	RunBenchmarkNone(b, arrayLiteralsQuery)
}

func BenchmarkArrayLiterals_Basic(b *testing.B) {
	RunBenchmarkBasic(b, arrayLiteralsQuery)
}

func BenchmarkArrayLiterals_Full(b *testing.B) {
	RunBenchmarkFull(b, arrayLiteralsQuery)
}

func BenchmarkArraySpread_None(b *testing.B) {
	RunBenchmarkNone(b, arraySpreadQuery)
}

func BenchmarkArraySpread_Basic(b *testing.B) {
	RunBenchmarkBasic(b, arraySpreadQuery)
}

func BenchmarkArraySpread_Full(b *testing.B) {
	RunBenchmarkFull(b, arraySpreadQuery)
}

func BenchmarkArraySpreadHostList_None(b *testing.B) {
	RunBenchmarkNone(b, arrayHostSpreadQuery, benchmarkSpreadSource(struct{ runtime.List }{
		List: runtime.NewArrayWith(
			runtime.NewString("foo"),
			runtime.NewString("bar"),
			runtime.NewString("qaz"),
			runtime.NewString("abc"),
		),
	}))
}

func BenchmarkArraySpreadHostList_Basic(b *testing.B) {
	RunBenchmarkBasic(b, arrayHostSpreadQuery, benchmarkSpreadSource(struct{ runtime.List }{
		List: runtime.NewArrayWith(
			runtime.NewString("foo"),
			runtime.NewString("bar"),
			runtime.NewString("qaz"),
			runtime.NewString("abc"),
		),
	}))
}

func BenchmarkArraySpreadHostList_Full(b *testing.B) {
	RunBenchmarkFull(b, arrayHostSpreadQuery, benchmarkSpreadSource(struct{ runtime.List }{
		List: runtime.NewArrayWith(
			runtime.NewString("foo"),
			runtime.NewString("bar"),
			runtime.NewString("qaz"),
			runtime.NewString("abc"),
		),
	}))
}

func BenchmarkArraySpreadHostSnapshot_None(b *testing.B) {
	values := runtime.NewArrayWith(
		runtime.NewString("foo"),
		runtime.NewString("bar"),
		runtime.NewString("qaz"),
		runtime.NewString("abc"),
	)

	RunBenchmarkNone(b, arrayHostSpreadQuery, benchmarkSpreadSource(&benchmarkSnapshotList{
		List:     values,
		snapshot: values,
	}))
}

func BenchmarkArraySpreadHostSnapshot_Basic(b *testing.B) {
	values := runtime.NewArrayWith(
		runtime.NewString("foo"),
		runtime.NewString("bar"),
		runtime.NewString("qaz"),
		runtime.NewString("abc"),
	)

	RunBenchmarkBasic(b, arrayHostSpreadQuery, benchmarkSpreadSource(&benchmarkSnapshotList{
		List:     values,
		snapshot: values,
	}))
}

func BenchmarkArraySpreadHostSnapshot_Full(b *testing.B) {
	values := runtime.NewArrayWith(
		runtime.NewString("foo"),
		runtime.NewString("bar"),
		runtime.NewString("qaz"),
		runtime.NewString("abc"),
	)

	RunBenchmarkFull(b, arrayHostSpreadQuery, benchmarkSpreadSource(&benchmarkSnapshotList{
		List:     values,
		snapshot: values,
	}))
}

func BenchmarkObjectLiterals_None(b *testing.B) {
	RunBenchmarkNone(b, objectLiteralsQuery)
}

func BenchmarkObjectLiterals_Basic(b *testing.B) {
	RunBenchmarkBasic(b, objectLiteralsQuery)
}

func BenchmarkObjectLiterals_Full(b *testing.B) {
	RunBenchmarkFull(b, objectLiteralsQuery)
}

func BenchmarkObjectSpread_None(b *testing.B) {
	RunBenchmarkNone(b, objectSpreadQuery)
}

func BenchmarkObjectSpread_Basic(b *testing.B) {
	RunBenchmarkBasic(b, objectSpreadQuery)
}

func BenchmarkObjectSpread_Full(b *testing.B) {
	RunBenchmarkFull(b, objectSpreadQuery)
}

func BenchmarkObjectSpreadHostMap_None(b *testing.B) {
	RunBenchmarkNone(b, objectHostSpreadQuery, benchmarkSpreadSource(&benchmarkMap{
		Map: benchmarkSpreadObject(),
	}))
}

func BenchmarkObjectSpreadHostMap_Basic(b *testing.B) {
	RunBenchmarkBasic(b, objectHostSpreadQuery, benchmarkSpreadSource(&benchmarkMap{
		Map: benchmarkSpreadObject(),
	}))
}

func BenchmarkObjectSpreadHostMap_Full(b *testing.B) {
	RunBenchmarkFull(b, objectHostSpreadQuery, benchmarkSpreadSource(&benchmarkMap{
		Map: benchmarkSpreadObject(),
	}))
}

func BenchmarkObjectSpreadHostSnapshot_None(b *testing.B) {
	object := benchmarkSpreadObject()

	RunBenchmarkNone(b, objectHostSpreadQuery, benchmarkSpreadSource(&benchmarkSnapshotObject{
		Map:      object,
		snapshot: object,
	}))
}

func BenchmarkObjectSpreadHostSnapshot_Basic(b *testing.B) {
	object := benchmarkSpreadObject()

	RunBenchmarkBasic(b, objectHostSpreadQuery, benchmarkSpreadSource(&benchmarkSnapshotObject{
		Map:      object,
		snapshot: object,
	}))
}

func BenchmarkObjectSpreadHostSnapshot_Full(b *testing.B) {
	object := benchmarkSpreadObject()

	RunBenchmarkFull(b, objectHostSpreadQuery, benchmarkSpreadSource(&benchmarkSnapshotObject{
		Map:      object,
		snapshot: object,
	}))
}

func BenchmarkObjectComputedLiterals_None(b *testing.B) {
	RunBenchmarkNone(b, objectComputedLiteralsQuery)
}

func BenchmarkObjectComputedLiterals_Basic(b *testing.B) {
	RunBenchmarkBasic(b, objectComputedLiteralsQuery)
}

func BenchmarkObjectComputedLiterals_Full(b *testing.B) {
	RunBenchmarkFull(b, objectComputedLiteralsQuery)
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
