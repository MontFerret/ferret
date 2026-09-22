package sdk_test

import (
	"context"
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

type oversizedSliceSource struct {
	runtime.List
	failure error
	length  runtime.Int
}

func (v oversizedSliceSource) Length(context.Context) (runtime.Int, error) {
	return v.length, v.failure
}

func TestToSliceOversizedHint(t *testing.T) {
	if math.MaxInt != math.MaxInt32 {
		t.Skip("an int64 hint exceeds native int only on 32-bit hosts")
	}
	input := oversizedSliceSource{List: runtime.NewArrayWith(runtime.Int(7)), length: math.MaxInt32 + 1}
	mapper := func(_ context.Context, value, _ runtime.Value) (int64, error) { return int64(value.(runtime.Int)), nil }
	got, err := sdk.ToSlice(t.Context(), input, mapper)
	if err != nil || !reflect.DeepEqual(got, []int64{7}) {
		t.Fatalf("ToSlice = %v, %v", got, err)
	}

	failure := errors.New("length failed")
	input.failure = failure
	if _, err := sdk.ToSlice(t.Context(), input, mapper); !errors.Is(err, failure) {
		t.Fatalf("lost length failure: %v", err)
	}
	input.failure = nil
	if _, err := sdk.ToSlice(t.Context(), input, func(context.Context, runtime.Value, runtime.Value) (int64, error) { return 0, failure }); !errors.Is(err, failure) {
		t.Fatalf("lost mapper failure: %v", err)
	}
}
