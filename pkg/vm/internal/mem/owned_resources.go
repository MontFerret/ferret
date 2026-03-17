package mem

import (
	"io"
	"reflect"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// OwnedResources tracks only direct register-held closers for a single frame.
// Ownership is per resource, not per register slot: duplicate aliases in the
// same frame still map to a single owned closer. The VM is responsible for
// deciding when the last live register alias disappears before calling Discard.
// OwnedResources ignores values that cannot be tracked safely and is not a
// general runtime lifetime manager.
type OwnedResources struct {
	closers map[io.Closer]struct{}
}

func (o *OwnedResources) Track(val runtime.Value) {
	closer, ok := TrackedCloserOf(val)
	if !ok {
		return
	}

	if o.closers == nil {
		o.closers = make(map[io.Closer]struct{})
	}

	o.closers[closer] = struct{}{}
}

func (o *OwnedResources) Owns(val runtime.Value) bool {
	closer, ok := TrackedCloserOf(val)
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
	closer, ok := TrackedCloserOf(val)
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
// the current frame and reassigns the matching owned closers to dst. Duplicate
// values transfer ownership only once per closer, while dead aliases in the
// retiring frame are dropped instead of being deferred.
func (o *OwnedResources) ExtractMany(values []runtime.Value, dst *OwnedResources) {
	if o.closers == nil || len(values) == 0 {
		return
	}

	survivors := make(map[io.Closer]struct{})

	for _, val := range values {
		closer, ok := TrackedCloserOf(val)
		if !ok {
			continue
		}

		if _, exists := o.closers[closer]; !exists {
			continue
		}

		survivors[closer] = struct{}{}
	}

	if len(survivors) == 0 {
		return
	}

	if dst != nil && dst.closers == nil {
		dst.closers = make(map[io.Closer]struct{})
	}

	for closer := range survivors {
		delete(o.closers, closer)

		if dst != nil {
			dst.closers[closer] = struct{}{}
		}
	}

	if len(o.closers) == 0 {
		o.closers = nil
	}
}

func (o *OwnedResources) Discard(val runtime.Value, deferred *DeferredClosers) {
	closer, ok := TrackedCloserOf(val)
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
	closer, ok := TrackedCloserOf(val)
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

// DiscardCloser removes ownership of an already-extracted io.Closer and
// schedules it for deferred cleanup. Unlike Discard, this avoids a
// redundant TrackedCloserOf call when the caller already holds the closer.
func (o *OwnedResources) DiscardCloser(closer io.Closer, deferred *DeferredClosers) {
	if o.closers == nil {
		return
	}

	if _, exists := o.closers[closer]; !exists {
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

func TrackedCloserOf(val runtime.Value) (io.Closer, bool) {
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

