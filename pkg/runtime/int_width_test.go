package runtime_test

import (
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestIntWidthConversions(t *testing.T) {
	for _, value := range []int64{math.MinInt64, math.MinInt32 - 1, math.MinInt32, 0, 42, math.MaxInt32, math.MaxInt32 + 1, 1<<53 + 1, math.MaxInt64} {
		t.Run(strconv.FormatInt(value, 10), func(t *testing.T) {
			integer := runtime.NewInt64(value)
			unwrapped, ok := integer.Unwrap().(int64)
			if !ok || unwrapped != value {
				t.Fatalf("Unwrap = %T(%v), want int64(%d)", integer.Unwrap(), integer.Unwrap(), value)
			}

			if got, ok := runtime.UnwrapAs[int64](integer); !ok || got != value {
				t.Fatalf("UnwrapAs[int64] = %d, %v", got, ok)
			}

			if _, ok := runtime.UnwrapAs[int](integer); ok {
				t.Fatal("UnwrapAs[int] must not narrow an Int")
			}

			parsed, err := runtime.ParseInt(strconv.FormatInt(value, 10))
			if err != nil || parsed != integer {
				t.Fatalf("ParseInt = %v, %v, want %d", parsed, err, value)
			}

			host, err := runtime.ValueOf(value)
			if err != nil || host != integer {
				t.Fatalf("ValueOf = %T(%v), %v", host, host, err)
			}

			floating, err := runtime.ParseFloat(strconv.FormatInt(value, 10))
			if err != nil || floating != runtime.Float(value) {
				t.Fatalf("ParseFloat integer string = %v, %v", floating, err)
			}
		})
	}

	for _, text := range []string{"9223372036854775808", "-9223372036854775809", "1.5", "1e3", " 42"} {
		if _, err := runtime.ParseInt(text); err == nil {
			t.Errorf("ParseInt(%q) succeeded", text)
		}

		if _, err := runtime.ParseFloat(text); err == nil {
			t.Errorf("ParseFloat(%q) expanded its integer-string grammar", text)
		}
	}
}

func TestToIntDateTimeWidth(t *testing.T) {
	for _, seconds := range []int64{math.MinInt32 - 1, math.MinInt32, math.MaxInt32, math.MaxInt32 + 1, 1<<53 + 1} {
		value := runtime.NewDateTime(time.Unix(seconds, 0))
		got, err := runtime.ToInt(t.Context(), value)
		if err != nil || got != runtime.Int(seconds) {
			t.Errorf("ToInt(%d) = %d, %v", seconds, got, err)
		}
	}

	if got, err := runtime.ToInt(t.Context(), runtime.ZeroDateTime); err != nil || got != runtime.ZeroInt {
		t.Fatalf("zero date = %d, %v", got, err)
	}
}

func TestIsInfSignWidth(t *testing.T) {
	for _, sign := range []runtime.Int{math.MinInt64, math.MinInt32 - 1, -1, 0, 1, math.MaxInt32 + 1, 1 << 32, math.MaxInt64} {
		for _, polarity := range []int{-1, 1} {
			want := sign == 0 || (sign < 0 && polarity < 0) || (sign > 0 && polarity > 0)
			if got := runtime.IsInf(runtime.Float(math.Inf(polarity)), sign); bool(got) != want {
				t.Errorf("IsInf(%d, %d) = %v, want %v", polarity, sign, got, want)
			}
		}

		if runtime.IsInf(runtime.Float(1), sign) {
			t.Errorf("finite value identified as infinity for sign %d", sign)
		}
	}
}
