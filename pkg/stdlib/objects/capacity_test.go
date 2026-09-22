package objects

import (
	"context"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type oversizedCopyList struct{ runtime.List }

func (v oversizedCopyList) Length(context.Context) (runtime.Int, error) {
	return math.MaxInt32 + 1, nil
}

func TestCopyListOversizedHint(t *testing.T) {
	if math.MaxInt != math.MaxInt32 {
		t.Skip("an int64 hint exceeds native int only on 32-bit hosts")
	}
	input := oversizedCopyList{runtime.NewArrayWith(runtime.Int(7))}
	got, err := copyListValues(t.Context(), input)
	if err != nil || got.String() != "[7]" {
		t.Fatalf("copy = %v, %v", got, err)
	}
	if got == input.List {
		t.Fatal("copy reused the source")
	}
}
