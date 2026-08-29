package bytecode

import "testing"

func TestFunctionIDValid(t *testing.T) {
	tests := []struct {
		name  string
		id    FunctionID
		count int
		want  bool
	}{
		{name: "first", id: 0, count: 1, want: true},
		{name: "last", id: 2, count: 3, want: true},
		{name: "upper bound", id: 3, count: 3},
		{name: "top-level sentinel", id: NoFunction, count: 3},
		{name: "other negative", id: -2, count: 3},
		{name: "empty table", id: 0, count: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.id.Valid(tc.count); got != tc.want {
				t.Fatalf("FunctionID(%d).Valid(%d) = %t, want %t", tc.id, tc.count, got, tc.want)
			}
		})
	}
}
