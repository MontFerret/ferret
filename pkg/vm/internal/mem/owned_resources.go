package mem

import (
	"io"
	"reflect"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// OwnedResources tracks only direct register-held closers for a single frame.
// Counts represent the number of live register aliases in the active frame for
// each owned closer. It ignores values that cannot be tracked safely and is not
// a general runtime lifetime manager.
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

// Extract removes all ownership tracking for val from the current frame without
// scheduling the closer for deferred cleanup. Use it when a value survives
// frame teardown, for example as a return value or final result root.
func (o *OwnedResources) Extract(val runtime.Value) bool {
	closer, ok := trackedCloserOf(val)
	if !ok || o.closers == nil {
		return false
	}

	_, exists := o.closers[closer]
	if !exists {
		return false
	}

	delete(o.closers, closer)
	if len(o.closers) == 0 {
		o.closers = nil
	}

	return true
}

// ExtractMany removes ownership tracking for the provided surviving values from
// the current frame and reassigns the matching alias counts to dst. Duplicate
// values retain their multiplicity in the destination frame, while dead aliases
// in the retiring frame are dropped instead of being deferred.
func (o *OwnedResources) ExtractMany(values []runtime.Value, dst *OwnedResources) {
	if o.closers == nil || len(values) == 0 {
		return
	}

	survivors := make(map[io.Closer]int)

	for _, val := range values {
		closer, ok := trackedCloserOf(val)
		if !ok {
			continue
		}

		if _, exists := o.closers[closer]; !exists {
			continue
		}

		survivors[closer]++
	}

	if len(survivors) == 0 {
		return
	}

	if dst != nil && dst.closers == nil {
		dst.closers = make(map[io.Closer]int)
	}

	for closer, count := range survivors {
		if ownedCount, exists := o.closers[closer]; exists && count > ownedCount {
			count = ownedCount
		}

		delete(o.closers, closer)

		if dst != nil {
			dst.closers[closer] += count
		}
	}

	if len(o.closers) == 0 {
		o.closers = nil
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

	// Release is terminal for explicit close, not a generic ref-count decrement:
	// once one alias is closed, the closer is retired from ownership tracking so
	// stale aliases cannot close it again during later cleanup.
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
