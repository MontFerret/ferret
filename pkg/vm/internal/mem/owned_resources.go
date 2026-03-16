package mem

import (
	"io"
	"reflect"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// OwnedResources tracks only direct register-held closers for a single frame.
// It ignores values that cannot be tracked safely and is not a general runtime
// lifetime manager.
type OwnedResources struct {
	closers map[io.Closer]int
}

func (o *OwnedResources) Track(val runtime.Value) {
	closer, ok := trackedCloserOf(val)
	if !ok {
		return
	}

	if o.closers == nil {
		o.closers = make(map[io.Closer]int)
	}

	o.closers[closer]++
}

func (o *OwnedResources) Forget(val runtime.Value) {
	closer, ok := trackedCloserOf(val)
	if !ok || o.closers == nil {
		return
	}

	count, exists := o.closers[closer]
	if !exists {
		return
	}

	if count <= 1 {
		delete(o.closers, closer)
		if len(o.closers) == 0 {
			o.closers = nil
		}

		return
	}

	o.closers[closer] = count - 1
}

func (o *OwnedResources) Owns(val runtime.Value) bool {
	closer, ok := trackedCloserOf(val)
	if !ok || o.closers == nil {
		return false
	}

	_, ok = o.closers[closer]
	return ok
}

func (o *OwnedResources) Transfer(val runtime.Value, dst *OwnedResources) {
	closer, ok := trackedCloserOf(val)
	if !ok || o.closers == nil {
		return
	}

	if _, exists := o.closers[closer]; !exists {
		return
	}

	delete(o.closers, closer)
	if len(o.closers) == 0 {
		o.closers = nil
	}

	if dst != nil {
		if dst.closers == nil {
			dst.closers = make(map[io.Closer]int)
		}

		dst.closers[closer]++
	}
}

func (o *OwnedResources) TransferMany(values []runtime.Value, dst *OwnedResources) {
	for _, val := range values {
		o.Transfer(val, dst)
	}
}

func (o *OwnedResources) Discard(val runtime.Value, deferred *DeferredClosers) {
	closer, ok := trackedCloserOf(val)
	if !ok || o.closers == nil {
		return
	}

	count, exists := o.closers[closer]
	if !exists {
		return
	}

	if count > 1 {
		o.closers[closer] = count - 1
		return
	}

	delete(o.closers, closer)
	if len(o.closers) == 0 {
		o.closers = nil
	}

	if deferred != nil {
		deferred.AddCloser(closer)
	}
}

func (o *OwnedResources) DrainTo(deferred *DeferredClosers) {
	if o.closers == nil {
		return
	}

	if deferred != nil {
		for closer := range o.closers {
			deferred.AddCloser(closer)
		}
	}

	o.closers = nil
}

func (o *OwnedResources) Release(val runtime.Value) (io.Closer, bool) {
	closer, ok := trackedCloserOf(val)
	if !ok || o.closers == nil {
		return nil, false
	}

	if _, exists := o.closers[closer]; !exists {
		return nil, false
	}

	delete(o.closers, closer)
	if len(o.closers) == 0 {
		o.closers = nil
	}

	return closer, true
}

func (o *OwnedResources) CloseAll() {
	deferred := DeferredClosers{}
	o.DrainTo(&deferred)
	_ = deferred.CloseAll()
}

func (o *OwnedResources) ForEach(fn func(io.Closer)) {
	if fn == nil {
		return
	}

	for closer := range o.closers {
		fn(closer)
	}
}

func trackedCloserOf(val runtime.Value) (io.Closer, bool) {
	closer, ok := val.(io.Closer)
	if !ok || closer == nil {
		return nil, false
	}

	typ := reflect.TypeOf(closer)
	if typ == nil || !typ.Comparable() {
		return nil, false
	}

	return closer, true
}

func SameTrackedCloser(left, right runtime.Value) bool {
	leftCloser, leftOK := trackedCloserOf(left)
	if !leftOK {
		return false
	}

	rightCloser, rightOK := trackedCloserOf(right)
	if !rightOK {
		return false
	}

	return leftCloser == rightCloser
}
