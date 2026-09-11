package collections

import (
	"errors"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestIncrementCountBoundary(t *testing.T) {
	got, err := incrementCount(math.MaxInt64 - 1)
	if err != nil || got != math.MaxInt64 {
		t.Fatalf("last valid count: %v, %v", got, err)
	}

	got, err = incrementCount(got)
	if got != runtime.ZeroInt || !errors.Is(err, runtime.ErrRange) {
		t.Fatalf("overflow: %v, %v", got, err)
	}
}
