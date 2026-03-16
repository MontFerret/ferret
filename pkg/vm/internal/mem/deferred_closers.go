package mem

import (
	"errors"
	"io"
	"reflect"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DeferredClosers tracks discarded direct closers until a result or failed run
// drains them.
type DeferredClosers struct {
	seen    map[io.Closer]struct{}
	closers []io.Closer
}

func (d *DeferredClosers) Add(val runtime.Value) {
	closer, ok := trackedCloserOf(val)
	if !ok {
		return
	}

	d.AddCloser(closer)
}

func (d *DeferredClosers) AddCloser(closer io.Closer) {
	if closer == nil {
		return
	}

	typ := reflect.TypeOf(closer)
	if typ == nil || !typ.Comparable() {
		return
	}

	if d.seen == nil {
		d.seen = make(map[io.Closer]struct{})
	}

	if _, exists := d.seen[closer]; exists {
		return
	}

	d.seen[closer] = struct{}{}
	d.closers = append(d.closers, closer)
}

func (d *DeferredClosers) Merge(other *DeferredClosers) {
	if other == nil {
		return
	}

	for _, closer := range other.closers {
		d.AddCloser(closer)
	}

	other.Reset()
}

func (d *DeferredClosers) CloseAll() error {
	var err error

	for _, closer := range d.closers {
		if closeErr := closer.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}

	d.Reset()

	return err
}

func (d *DeferredClosers) ForEach(fn func(io.Closer)) {
	if fn == nil {
		return
	}

	for _, closer := range d.closers {
		fn(closer)
	}
}

func (d *DeferredClosers) Reset() {
	d.closers = nil
	d.seen = nil
}
