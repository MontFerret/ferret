package sdk_test

import (
	"errors"
	"math"
	"strconv"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

func TestDecodeIntWidth(t *testing.T) {
	for _, value := range []int64{math.MinInt64, math.MinInt32 - 1, math.MinInt32, math.MaxInt32, math.MaxInt32 + 1, 1<<53 + 1, math.MaxInt64} {
		t.Run(strconv.FormatInt(value, 10), func(t *testing.T) {
			source := runtime.NewInt64(value)
			target64 := int64(17)
			if err := sdk.Decode(t.Context(), source, &target64); err != nil || target64 != value {
				t.Fatalf("int64 decode = %d, %v", target64, err)
			}

			target := 17
			err := sdk.Decode(t.Context(), source, &target)
			if value < math.MinInt || value > math.MaxInt {
				var decodeError *sdk.DecodeError
				if !errors.As(err, &decodeError) || decodeError.Kind() != sdk.DecodeErrorKindRange || decodeError.Path() != "$" || target != 17 {
					t.Fatalf("overflow decode = %d, %v; want unchanged destination and DecodeError", target, err)
				}
			} else if err != nil || int64(target) != value {
				t.Fatalf("int decode = %d, %v", target, err)
			}

			encoded, err := sdk.Encode(t.Context(), value)
			if err != nil || encoded != source {
				t.Fatalf("Encode = %T(%v), %v", encoded, encoded, err)
			}
		})
	}
}
