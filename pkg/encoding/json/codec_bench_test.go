package json_test

import (
	"strconv"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var benchmarkJSONEncoded []byte

func benchmarkFlatArray(size int) runtime.Value {
	values := make([]runtime.Value, size)

	for i := range size {
		values[i] = runtime.NewInt(i)
	}

	return runtime.NewArrayOf(values)
}

func benchmarkFlatObject(size int) runtime.Value {
	props := make(map[string]runtime.Value, size)

	for i := range size {
		props[strconv.Itoa(i)] = runtime.NewInt(i)
	}

	return runtime.NewObjectWith(props)
}

func BenchmarkJSONCodecEncode(b *testing.B) {
	codec := json.Default

	cases := []struct {
		name  string
		value runtime.Value
	}{
		{name: "flat_array_1024", value: benchmarkFlatArray(1024)},
		{name: "flat_object_256", value: benchmarkFlatObject(256)},
		{name: "nested_array_10000", value: nestedArray(10_000)},
		{name: "nested_object_5000", value: nestedObject(5_000)},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			out, err := codec.Encode(tc.value)
			if err != nil {
				b.Fatalf("setup encode failed: %v", err)
			}

			b.ReportAllocs()
			b.SetBytes(int64(len(out)))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				out, err := codec.Encode(tc.value)
				if err != nil {
					b.Fatalf("encode failed: %v", err)
				}

				benchmarkJSONEncoded = out
			}
		})
	}
}
