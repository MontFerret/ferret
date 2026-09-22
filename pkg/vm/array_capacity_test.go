package vm

import (
	"context"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type oversizedArrayList struct {
	runtime.List
	length runtime.Int
}

func (v oversizedArrayList) Length(context.Context) (runtime.Int, error) { return v.length, nil }

func TestVMArrayCapacityHints(t *testing.T) {
	input := oversizedArrayList{List: runtime.NewArrayWith(runtime.Int(7), runtime.Int(7)), length: math.MaxInt64}
	flattened, err := arrayFlatten(t.Context(), input, 1)
	if err != nil || flattened.String() != "[7,7]" {
		t.Fatalf("flatten = %v, %v", flattened, err)
	}

	if math.MaxInt == math.MaxInt32 {
		distinct, err := arrayDistinct(t.Context(), input)
		if err != nil || distinct.String() != "[7]" {
			t.Fatalf("distinct = %v, %v", distinct, err)
		}
	}
}
