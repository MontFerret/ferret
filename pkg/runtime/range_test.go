package runtime_test

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var _ runtime.Measurable = runtime.NewRange(0, 0)

func TestRangeTypeRemainsIterable(t *testing.T) {
	if got := runtime.TypeOf(runtime.NewRange(1, 3)); got != runtime.TypeIterable {
		t.Fatalf("TypeOf(Range) = %s, want %s", got, runtime.TypeIterable)
	}
}

func TestRangeLength(t *testing.T) {
	tests := []struct {
		name       string
		start, end runtime.Int
		want       runtime.Int
	}{
		{name: "ascending", start: 1, end: 3, want: 3},
		{name: "descending", start: 3, end: 1, want: 3},
		{name: "singleton", start: 7, end: 7, want: 1},
		{name: "negative ascending", start: -3, end: -1, want: 3},
		{name: "negative descending", start: -1, end: -3, want: 3},
		{name: "cross zero", start: -2, end: 2, want: 5},
		{name: "largest representable", start: math.MinInt64, end: -2, want: math.MaxInt64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runtime.NewRange(tt.start, tt.end).Length(context.Background())
			if err != nil {
				t.Fatalf("Length() error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("Length() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestRangeLengthRejectsOverflow(t *testing.T) {
	tests := []struct {
		name       string
		start, end runtime.Int
	}{
		{name: "ascending", start: math.MinInt64, end: -1},
		{name: "descending", start: math.MaxInt64, end: 0},
		{name: "full domain", start: math.MinInt64, end: math.MaxInt64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runtime.NewRange(tt.start, tt.end).Length(context.Background())
			if !errors.Is(err, runtime.ErrRange) {
				t.Fatalf("Length() error = %v, want ErrRange", err)
			}
		})
	}
}

func TestRangeMarshalJSONNegativeRanges(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		start runtime.Int
		end   runtime.Int
	}{
		{name: "ascending", start: -3, end: -1, want: "[-3,-2,-1]"},
		{name: "descending", start: -1, end: -3, want: "[-1,-2,-3]"},
		{name: "cross zero", start: -1, end: 1, want: "[-1,0,1]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runtime.NewRange(tt.start, tt.end).MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() error: %v", err)
			}

			if string(got) != tt.want {
				t.Fatalf("MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}

var benchmarkRangeLength runtime.Int

func BenchmarkRangeLength(b *testing.B) {
	ctx := context.Background()
	r := runtime.NewRange(math.MinInt64, -2)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		length, err := r.Length(ctx)
		if err != nil {
			b.Fatal(err)
		}

		benchmarkRangeLength = length
	}
}
