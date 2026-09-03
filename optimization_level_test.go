package ferret

import "testing"

func TestOptimizationLevelString(t *testing.T) {
	tests := []struct {
		want  string
		level OptimizationLevel
	}{
		{level: OptimizationNone, want: "none"},
		{level: OptimizationBasic, want: "basic"},
		{level: OptimizationFull, want: "full"},
		{level: OptimizationLevel(-1), want: "unknown"},
	}

	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("OptimizationLevel(%d).String() = %q, want %q", tt.level, got, tt.want)
		}
	}
}
