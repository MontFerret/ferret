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

func TestDefaultEntropyDeferred(t *testing.T) {
	calls := uint64(0)
	entropy := func() [16]byte {
		calls++
		var seed [16]byte
		binary.LittleEndian.PutUint64(seed[:8], calls)
		binary.LittleEndian.PutUint64(seed[8:], calls+10)

		return seed
	}
	first, second := newSource(entropy), newSource(entropy)
	if calls != 0 || first == second {
		t.Fatal("construction must allocate independent sources without reading entropy")
	}

	ctx := WithContext(t.Context(), first)
	for range 10 {
		if got, ok := FromContext(ctx); !ok || got != first {
			t.Fatal("context did not retain the pending source")
		}

		if got := first.Int64(5, 5); got != 5 {
			t.Fatalf("equal bounds returned %d", got)
		}
	}

	if calls != 0 {
		t.Fatal("transport, lookup, or equal bounds read entropy")
	}

	sources := []*Source{first, second}
	controls := []*rand.Rand{rand.New(rand.NewPCG(1, 11)), rand.New(rand.NewPCG(2, 12))}
	for range 64 {
		for index, src := range sources {
			control := controls[index]
			if src.Float64() != control.Float64() {
				t.Fatal("initialization lost entropy, reseeded, or shared mutable state")
			}

			if src.entropy != nil {
				t.Fatal("initialized source retained its entropy callback")
			}
		}
	}

	if calls != 2 {
		t.Fatal("each source must obtain entropy exactly once")
	}

	if first, second := New(), New(); first == nil || second == nil || first == second {
		t.Fatal("default constructors did not allocate independent sources")
	}
}

func TestDefaultSourceFirstDraw(t *testing.T) {
	for _, test := range []struct {
		draw func(*Source) any
		want func(*rand.Rand) any
		name string
	}{
		{name: "float", draw: func(src *Source) any { return src.Float64() }, want: func(control *rand.Rand) any { return control.Float64() }},
		{name: "int", draw: func(src *Source) any { return src.Int64(-3, 7) }, want: func(control *rand.Rand) any { return int64(control.Uint64N(11)) - 3 }},
		{name: "int full", draw: func(src *Source) any { return src.Int64(math.MinInt64, math.MaxInt64) }, want: func(control *rand.Rand) any { return int64(control.Uint64() + uint64(1)<<63) }},
		{name: "bool", draw: func(src *Source) any { return src.Bool() }, want: func(control *rand.Rand) any { return control.Uint64()&1 != 0 }},
		{name: "legacy", draw: func(src *Source) any { return src.LegacyFloat64(8.5, -2.5) }, want: func(control *rand.Rand) any { return math.Floor(control.Float64()*12) - 2.5 }},
		{name: "legacy equal", draw: func(src *Source) any { return src.LegacyFloat64(5, 5) }, want: func(control *rand.Rand) any { return math.Floor(control.Float64()) + 5 }},
	} {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			src := newSource(func() [16]byte {
				calls++

				return [16]byte{42}
			})
			control := rand.New(rand.NewPCG(42, 0))
			for range 32 {
				if got, want := test.draw(src), test.want(control); got != want {
					t.Fatalf("draw = %v, want %v", got, want)
				}

				if calls != 1 {
					t.Fatalf("entropy reads = %d, want one", calls)
				}
			}

			if src.Float64() != control.Float64() {
				t.Fatal("initialization changed sequence consumption")
			}
		})
	}
}

func TestSourceFirstDrawAllocations(t *testing.T) {
	for name, draw := range map[string]func(*Source){
		"float":    func(src *Source) { src.Float64() },
		"int":      func(src *Source) { src.Int64(-3, 7) },
		"int full": func(src *Source) { src.Int64(math.MinInt64, math.MaxInt64) },
		"bool":     func(src *Source) { src.Bool() },
		"legacy":   func(src *Source) { src.LegacyFloat64(8.5, -2.5) },
	} {
		t.Run(name, func(t *testing.T) {
			// AllocsPerRun performs one warmup call. Supply a fresh, already-owned
			// source for every call so the measured draws all initialize a source.
			// Fixed entropy isolates our allocations from crypto/rand, which can
			// move its buffer to the heap under race instrumentation on Linux.
			const runs = 100
			sources := make([]*Source, runs+1)
			for index := range sources {
				sources[index] = newSource(func() [16]byte { return [16]byte{42} })
			}

			next := 0
			allocations := testing.AllocsPerRun(runs, func() {
				draw(sources[next])
				next++
			})
			if allocations != 0 {
				t.Fatalf("first draw allocated %v times", allocations)
			}
		})
	}
}

func TestSourceSteadyStateDrawAllocations(t *testing.T) {
	for name, draw := range map[string]func(*Source){
		"float":    func(src *Source) { src.Float64() },
		"int":      func(src *Source) { src.Int64(-3, 7) },
		"int full": func(src *Source) { src.Int64(math.MinInt64, math.MaxInt64) },
		"bool":     func(src *Source) { src.Bool() },
		"legacy":   func(src *Source) { src.LegacyFloat64(8.5, -2.5) },
	} {
		t.Run(name, func(t *testing.T) {
			// Exercise the real entropy reader before measuring initialized draws.
			src := New()
			draw(src)
			if src.entropy != nil {
				t.Fatal("first draw did not finish initialization")
			}

			allocations := testing.AllocsPerRun(100, func() {
				draw(src)
			})
			if allocations != 0 {
				t.Fatalf("steady-state draw allocated %v times", allocations)
			}
		})
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
