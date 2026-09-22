package json_test

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestJSONIntegerWidth(t *testing.T) {
	for _, value := range []int64{math.MinInt64, math.MinInt32 - 1, math.MinInt32, math.MaxInt32, math.MaxInt32 + 1, 1<<53 + 1, math.MaxInt64} {
		text := strconv.FormatInt(value, 10)
		t.Run(text, func(t *testing.T) {
			for _, document := range []string{text, fmt.Sprintf("{\"value\":[%s]}", text)} {
				decoded, err := json.Default.Decode([]byte(document))
				if err != nil {
					t.Fatal(err)
				}

				encoded, err := json.Default.Encode(decoded)
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
		decoded, err := json.Default.Decode([]byte(document))
		if _, ok := decoded.(runtime.Float); !ok || err != nil {
			t.Errorf("%s decoded as %T(%v), %v; want Float", document, decoded, decoded, err)
		}
	}
}
