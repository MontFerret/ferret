package bytecode

import "testing"

func TestFunctionIDValid(t *testing.T) {
	tests := []struct {
		name string
		id   FunctionID
		want bool
	}{
		{name: "top-level sentinel", id: NoFunction},
		{name: "other negative", id: -2},
		{name: "zero", id: 0, want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.id.Valid(); got != tc.want {
				t.Fatalf("FunctionID(%d).Valid() = %t, want %t", tc.id, got, tc.want)
			}
		})
	}
}

func TestFunctionIDInRange(t *testing.T) {
	tests := []struct {
		name  string
		id    FunctionID
		count int
		want  bool
	}{
		{name: "first", id: 0, count: 3, want: true},
		{name: "last", id: 2, count: 3, want: true},
		{name: "upper bound", id: 3, count: 3},
		{name: "top-level sentinel", id: NoFunction, count: 3},
		{name: "other negative", id: -2, count: 3},
		{name: "empty table", id: 0, count: 0},
		{name: "negative count", id: 0, count: -1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.id.InRange(tc.count); got != tc.want {
				t.Fatalf("FunctionID(%d).InRange(%d) = %t, want %t", tc.id, tc.count, got, tc.want)
			}
		})
	}
}
