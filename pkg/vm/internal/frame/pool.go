package frame

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// Pool manages reusable register windows indexed by size.
type Pool struct {
	buckets [][][]runtime.Value
}

// Init prepares the pool with buckets up to maxSize.
func (p *Pool) Init(maxSize int) {
	if maxSize < 0 {
		maxSize = 0
	}

	p.buckets = make([][][]runtime.Value, maxSize+1)
}

func (p *Pool) Get(size int) []runtime.Value {
	if size <= 0 {
		return nil
	}

	if size < len(p.buckets) {
		bucket := p.buckets[size]
		n := len(bucket)
		if n > 0 {
			reg := bucket[n-1]
			p.buckets[size] = bucket[:n-1]
			return reg
		}
	}

	reg := make([]runtime.Value, size)
	fillWithNone(reg)
	return reg
}

// Put clears and stores a register window for reuse.
func (p *Pool) Put(reg []runtime.Value) {
	if len(reg) == 0 {
		return
	}

	// Normalize to runtime.None and avoid retaining references across calls.
	fillWithNone(reg)

	size := len(reg)
	if size < len(p.buckets) {
		p.buckets[size] = append(p.buckets[size], reg)
	}
}

func fillWithNone(reg []runtime.Value) {
	for i := range reg {
		reg[i] = runtime.None
	}
}
