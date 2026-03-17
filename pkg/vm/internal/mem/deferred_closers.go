package mem

import (
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DeferredClosers tracks discarded direct closers until a result or failed run
// drains them.
type DeferredClosers struct {
	set CloserSet
}

func (d *DeferredClosers) Add(val runtime.Value) {
	closer, ok := val.(io.Closer)
	if !ok {
		return
	}

	d.AddCloser(closer)
}

func (d *DeferredClosers) AddCloser(closer io.Closer) {
	d.set.Add(closer)
}

func (d *DeferredClosers) Merge(other *DeferredClosers) {
	if other == nil {
		return
	}

	d.set.Merge(&other.set)
}

func (d *DeferredClosers) CloseAll() error {
	return d.set.CloseAll()
}

func (d *DeferredClosers) ForEach(fn func(io.Closer)) {
	d.set.ForEach(fn)
}

func (d *DeferredClosers) Reset() {
	d.set.Reset()
}
