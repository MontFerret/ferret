package rnd

import (
	"context"
	"encoding/binary"
	"math"
	"math/rand/v2"
	"testing"
)

func TestSeededSharedSequence(t *testing.T) {
	for _, seed := range []int64{0, 42, -1, math.MinInt64, math.MaxInt64} {
		src := NewSeed(seed)
		control := rand.New(rand.NewPCG(uint64(seed), 0))
		for range 64 {
			if got, want := src.Float64(), control.Float64(); got != want {
				t.Fatalf("seed %d: float = %v, want %v", seed, got, want)
			}

			if got, want := src.Int64(-3, 7), int64(control.Uint64N(11))-3; got != want {
				t.Fatalf("seed %d: int = %v, want %v", seed, got, want)
			}

			if got, want := src.Bool(), control.Uint64()&1 != 0; got != want {
				t.Fatalf("seed %d: bool = %v, want %v", seed, got, want)
			}

			if got, want := src.LegacyFloat64(8.5, -2.5), math.Floor(control.Float64()*12)-2.5; got != want {
				t.Fatalf("seed %d: legacy = %v, want %v", seed, got, want)
			}
		}
	}
}

func TestInt64EntireDomain(t *testing.T) {
	for _, bounds := range [][2]int64{
		{1, 6}, {-10, -2}, {-10, 10},
		{math.MinInt64, math.MaxInt64}, {math.MinInt64, math.MaxInt64 - 1},
		{math.MinInt64 + 1, math.MaxInt64}, {math.MinInt64, 0}, {-1, math.MaxInt64},
		{math.MinInt64, math.MinInt64 + 1}, {math.MaxInt64 - 1, math.MaxInt64},
		{math.MinInt64, math.MinInt64}, {math.MaxInt64, math.MaxInt64},
	} {
		src, control := NewSeed(42), rand.New(rand.NewPCG(42, 0))
		min, max := bounds[0], bounds[1]
		for range 512 {
			got := src.Int64(min, max)
			if got < min || got > max {
				t.Fatalf("%v: out-of-range value %d", bounds, got)
			}

			want := min
			if min != max {
				width := uint64(max) - uint64(min) + 1
				if width == 0 {
					want = int64(uint64(min) + control.Uint64())
				} else {
					want = int64(uint64(min) + control.Uint64N(width))
				}
			}

			if got != want {
				t.Fatalf("%v: value %d, want unbiased draw %d", bounds, got, want)
			}
		}

		if src.Float64() != control.Float64() {
			t.Fatalf("%v: unexpected sequence consumption", bounds)
		}
	}
}

func TestSourcesDoNotShareState(t *testing.T) {
	first, second, control := NewSeed(42), NewSeed(42), NewSeed(42)
	if first == second {
		t.Fatal("constructors shared a source")
	}

	for range 30 {
		first.Float64()
	}

	for range 20 {
		if second.Float64() != control.Float64() {
			t.Fatal("advancing the first source changed the second")
		}
	}
}

func TestDefaultEntropyConstruction(t *testing.T) {
	calls := uint64(0)
	fill := func(seed []byte) {
		calls++
		if len(seed) != 16 {
			t.Fatalf("entropy request = %d bytes", len(seed))
		}

		binary.LittleEndian.PutUint64(seed[:8], calls)
		binary.LittleEndian.PutUint64(seed[8:], calls+10)
	}
	first, second := newSource(fill), newSource(fill)
	if calls != 2 || first == second {
		t.Fatal("construction must obtain independent entropy once per source")
	}

	for index, src := range []*Source{first, second} {
		control := rand.New(rand.NewPCG(uint64(index+1), uint64(index+11)))
		for range 64 {
			if src.Float64() != control.Float64() {
				t.Fatal("construction lost entropy or reseeded during draws")
			}
		}
	}

	if calls != 2 {
		t.Fatal("draws requested additional entropy")
	}

	if first, second := New(), New(); first == nil || second == nil || first == second {
		t.Fatal("default constructors did not allocate independent sources")
	}
}

func TestContextTransport(t *testing.T) {
	for _, ctx := range []context.Context{nil, context.Background(), WithContext(context.Background(), nil)} {
		if src, ok := FromContext(ctx); src != nil || ok {
			t.Fatal("missing source lookup constructed a fallback")
		}
	}

	src, control := NewSeed(42), NewSeed(42)
	parent, cancel := context.WithCancel(t.Context())
	ctx := WithContext(parent, src)
	for range 10 {
		if got, ok := FromContext(ctx); !ok || got != src {
			t.Fatal("context did not retain source identity")
		}
	}

	cancel()
	if ctx.Err() != context.Canceled || src.Float64() != control.Float64() {
		t.Fatal("transport lost cancellation or consumed randomness")
	}

	if _, ok := FromContext(parent); ok {
		t.Fatal("transport mutated parent context")
	}
}
