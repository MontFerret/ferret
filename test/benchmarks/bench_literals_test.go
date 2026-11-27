package benchmarks_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
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

func BenchmarkObjectDictionaryLiterals_None(b *testing.B) {
	// The 65th distinct field crosses the default VM dictionary threshold.
	fields := make([]string, 65)
	for i := range fields {
		fields[i] = "k" + strconv.Itoa(i) + ": n"
	}

	RunBenchmarkNone(b, "FOR n IN 1..16 LET obj = {"+strings.Join(fields, ",")+"} RETURN obj.k64")
}

func BenchmarkObjectDictionaryReads_None(b *testing.B) {
	const size = 128
	fields := make([]string, size)
	hits := make([]runtime.Value, size)
	misses := make([]runtime.Value, size)
	for i := range size {
		name := fmt.Sprintf("property_%03d_common_suffix", i)
		fields[i] = name + ":" + strconv.Itoa(i)
		hits[i] = runtime.String(strings.Clone(name))
		misses[i] = runtime.String(fmt.Sprintf("property_%03d_common_suffix", i+size))
	}
	// The runtime parameter keeps the literal on the VM construction path.
	fields[0] = hits[0].String() + ":@zero"
	query := "LET obj = {" + strings.Join(fields, ",") + "} RETURN FOR n IN 1..8 FOR key IN @keys RETURN obj?.[key]"
	program, err := spec.Compile(query, compiler.None)
	if err != nil {
		b.Fatal(err)
	}

	for _, tc := range []struct {
		name string
		keys []runtime.Value
	}{{"hit", hits}, {"miss", misses}} {
		b.Run(tc.name, func(b *testing.B) {
			opts := []vm.EnvironmentOption{WithParam("keys", runtime.NewArrayOf(tc.keys)), WithParam("zero", 0)}
			result, err := spec.Exec(program, false, opts...)
			if err != nil {
				b.Fatal(err)
			}
			values, ok := result.([]any)
			if !ok || len(values) != 8*size {
				b.Fatalf("expected %d dictionary reads, got %T with %d values", 8*size, result, len(values))
			}
			for i, value := range values {
				var want any
				if tc.name == "hit" {
					want = float64(i % size)
				}
				if value != want {
					b.Fatalf("read %d: got %v, want %v", i, value, want)
				}
			}

			RunBenchmarkNone(b, query, opts...)
		})
	}
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

func BenchmarkObjectShapeReads_None(b *testing.B) {
	const prefix = "LET obj = {a: @value} RETURN FOR i IN 1..1024 RETURN "
	for _, tc := range []struct{ name, read string }{{"constant", "obj?.a"}, {"dynamic", "obj?.[@key]"}} {
		b.Run(tc.name, func(b *testing.B) {
			RunBenchmarkNone(b, prefix+tc.read, WithParam("value", 17), WithParam("key", "a"))
		})
	}
}
