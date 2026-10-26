package json_test

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"
	"math"
	"strconv"
	"testing"

	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestJSONNumericBufferCapacity(t *testing.T) {
	const prefix = "prefix:"

	items := make([]int, 4096)
	for i := range items {
		items[i] = i
	}

	rangeJSON, err := stdjson.Marshal(items)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		value runtime.Value
		name  string
		want  string
	}{
		{
			name:  "mixed",
			value: runtime.NewArrayWith(runtime.Int(math.MinInt64), runtime.Int(math.MaxInt64), runtime.Float(1.25), runtime.Float(1e100), runtime.True, runtime.None),
			want:  "[-9223372036854775808,9223372036854775807,1.25,1e+100,true,null]",
		},
		{name: "range_ascending", value: runtime.NewRange(math.MinInt64, math.MinInt64+1), want: "[-9223372036854775808,-9223372036854775807]"},
		{name: "range_descending", value: runtime.NewRange(math.MaxInt64, math.MaxInt64-1), want: "[9223372036854775807,9223372036854775806]"},
		{name: "range_chunks", value: runtime.NewRange(0, 4095), want: string(rangeJSON)},
	} {
		for _, spare := range []int{0, 1, 16, 31, 32, 33, 64, 4096} {
			t.Run(tc.name+"/"+strconv.Itoa(spare), func(t *testing.T) {
				data := make([]byte, len(prefix), len(prefix)+spare)
				copy(data, prefix)
				dst := bytes.NewBuffer(data)
				want := prefix + tc.want
				if err := json.Default.Encode(t.Context(), dst, tc.value); err != nil || dst.String() != want {
					t.Fatalf("encode = %s, %v; want %s", dst.Bytes(), err, want)
				}
			})
		}
	}
}

func TestJSONIntegerWidth(t *testing.T) {
	for _, value := range []int64{math.MinInt64, math.MinInt32 - 1, math.MinInt32, math.MaxInt32, math.MaxInt32 + 1, 1<<53 + 1, math.MaxInt64} {
		text := strconv.FormatInt(value, 10)
		t.Run(text, func(t *testing.T) {
			for _, document := range []string{text, fmt.Sprintf("{\"value\":[%s]}", text)} {
				decoded, err := ferretencoding.DecodeBytes(t.Context(), json.Default, []byte(document))
				if err != nil {
					t.Fatal(err)
				}

				encoded, err := ferretencoding.EncodeBytes(t.Context(), json.Default, decoded)
				if err != nil || string(encoded) != document {
					t.Fatalf("round trip = %s, %v; want %s", encoded, err, document)
				}

				scalar := decoded
				if object, ok := decoded.(*runtime.Object); ok {
					nested, err := object.Get(t.Context(), runtime.String("value"))
					if err != nil {
						t.Fatal(err)
					}

					scalar, err = nested.(*runtime.Array).At(t.Context(), 0)
					if err != nil {
						t.Fatal(err)
					}
				}

				if scalar != runtime.NewInt64(value) {
					t.Fatalf("decoded = %T(%v), want Int(%d)", scalar, scalar, value)
				}
			}
		})
	}

	for _, document := range []string{"1.5", "1e3", "9223372036854775808", "-9223372036854775809"} {
		decoded, err := ferretencoding.DecodeBytes(t.Context(), json.Default, []byte(document))
		if _, ok := decoded.(runtime.Float); !ok || err != nil {
			t.Errorf("%s decoded as %T(%v), %v; want Float", document, decoded, decoded, err)
		}
	}
}
