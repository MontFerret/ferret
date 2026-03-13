package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// WindowPool manages reusable register windows indexed by size.
type WindowPool struct {
	buckets [][][]runtime.Value
}

func NewWindowPool(maxSize int) WindowPool {
	if maxSize < 0 {
		maxSize = 0
	}

	return WindowPool{
		buckets: make([][][]runtime.Value, maxSize+1),
	}
}

func (p *WindowPool) Acquire(size int) []runtime.Value {
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

// Release clears and stores a register window for reuse.
func (p *WindowPool) Release(reg []runtime.Value) {
	if len(reg) == 0 {
		return
	}

	fillWithNone(reg)

	size := len(reg)
	if size < len(p.buckets) {
		p.buckets[size] = append(p.buckets[size], reg)
	}
}
