package core_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
)

func TestJoinValueTypes(t *testing.T) {
	tests := []struct {
		name  string
		left  core.ValueType
		right core.ValueType
		want  core.ValueType
	}{
		{name: "same stays exact", left: core.TypeArray, right: core.TypeArray, want: core.TypeArray},
		{name: "unknown adopts known left", left: core.TypeUnknown, right: core.TypeObject, want: core.TypeObject},
		{name: "unknown adopts known right", left: core.TypeList, right: core.TypeUnknown, want: core.TypeList},
		{name: "mixed known widens to any", left: core.TypeArray, right: core.TypeObject, want: core.TypeAny},
		{name: "any stays any", left: core.TypeAny, right: core.TypeArray, want: core.TypeAny},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := core.JoinValueTypes(tt.left, tt.right); got != tt.want {
				t.Fatalf("JoinValueTypes(%v, %v) = %v, want %v", tt.left, tt.right, got, tt.want)
			}
		})
	}
}
