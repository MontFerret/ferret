package data

import (
	"context"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type oversizedObject struct{ runtime.Map }

func (oversizedObject) ObjectLike() {}

func (v oversizedObject) Length(context.Context) (runtime.Int, error) { return math.MaxInt32 + 1, nil }

func TestObjectSnapshotOversizedHint(t *testing.T) {
	if math.MaxInt != math.MaxInt32 {
		t.Skip("an int64 hint exceeds native int only on 32-bit hosts")
	}
	left := oversizedObject{runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.Int(7)})}
	right := oversizedObject{runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.Int(7)})}
	equal, err := equalObjectLike(t.Context(), left, right)
	if err != nil || !equal {
		t.Fatalf("equality = %v, %v", equal, err)
	}
	ordering, err := compareObjectLike(t.Context(), left, right)
	if err != nil || ordering != runtime.Equal {
		t.Fatalf("ordering = %v, %v", ordering, err)
	}
}
