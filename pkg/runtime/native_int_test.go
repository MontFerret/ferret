package runtime_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestToNativeInt(t *testing.T) {
	for _, value := range []runtime.Int{
		math.MinInt64, math.MinInt32 - 1, math.MinInt32, math.MinInt,
		-1, 0, 1, math.MaxInt32, math.MaxInt32 + 1, 1<<53 + 1,
		math.MaxInt, math.MaxInt64,
	} {
		t.Run(strconv.FormatInt(int64(value), 10), func(t *testing.T) {
			got, ok := runtime.ToNativeInt(value)
			if value < math.MinInt || value > math.MaxInt {
				if got != 0 || ok {
					t.Fatalf("ToNativeInt(%d) = %d, %v; want 0, false", value, got, ok)
				}

				return
			}

			if !ok || runtime.Int(got) != value {
				t.Fatalf("ToNativeInt(%d) = %d, %v; want exact native value", value, got, ok)
			}
		})
	}
}

func TestCapacityHint(t *testing.T) {
	for _, tc := range []struct {
		name                    string
		length                  runtime.Int
		multiplier, extra, want int
	}{
		{"scaled and extended", 4, 2, 1, 9},
		{"zero", 0, 2, 0, 0},
		{"zero with extra", 0, math.MaxInt, math.MaxInt, math.MaxInt},
		{"negative length", -1, 1, 0, 0},
		{"minimum length", math.MinInt64, 1, 0, 0},
		{"zero multiplier", 4, 0, 1, 0},
		{"negative multiplier", 4, -1, 0, 0},
		{"minimum multiplier", 4, math.MinInt, 0, 0},
		{"negative extra", 4, 1, -1, 0},
		{"minimum extra", 4, 1, math.MinInt, 0},
		{"native maximum", math.MaxInt, 1, 0, math.MaxInt},
		{"maximum multiplier", 1, math.MaxInt, 0, math.MaxInt},
		{"exact product and sum", math.MaxInt / 2, 2, 1, math.MaxInt},
		{"product overflow", runtime.Int(math.MaxInt/2) + 1, 2, 0, 0},
		{"sum overflow", math.MaxInt, 1, 1, 0},
		{"product fits but sum overflows", math.MaxInt / 2, 2, 2, 0},
		{"int64 product overflow", math.MaxInt64, 2, 0, 0},
		{"int64 sum overflow", math.MaxInt64, 1, 1, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := runtime.CapacityHint(tc.length, tc.multiplier, tc.extra); got != tc.want {
				t.Fatalf("CapacityHint(%d, %d, %d) = %d, want %d", tc.length, tc.multiplier, tc.extra, got, tc.want)
			}
		})
	}

	if strconv.IntSize == 32 {
		if got := runtime.CapacityHint(math.MaxInt32+1, 1, 0); got != 0 {
			t.Fatalf("unrepresentable native length = %d, want omitted hint", got)
		}
	}
}
