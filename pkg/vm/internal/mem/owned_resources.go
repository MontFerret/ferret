package mem

import (
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// OwnedResources tracks only direct register-held closers for a single frame.
// Ownership is per resource, not per register slot: duplicate aliases in the
// same frame still map to a single owned closer. The VM is responsible for
// deciding when the last live register alias disappears before calling Discard.
// OwnedResources ignores values that cannot be tracked safely and is not a
// general runtime lifetime manager.
type OwnedResources struct {
	closers map[ResourceKey]io.Closer
}

func (o *OwnedResources) Track(val runtime.Value) {
	key, closer, ok := ResourceKeyOf(val)
	if !ok {
		return
	}

	o.TrackResolved(key, closer)
}

func (o *OwnedResources) TrackResolved(key ResourceKey, closer io.Closer) {
	if closer == nil {
		return
	}

	if o.closers == nil {
		o.closers = make(map[ResourceKey]io.Closer)
	}

	o.closers[key] = closer
}

func (o *OwnedResources) Owns(val runtime.Value) bool {
	key, _, ok := ResourceKeyOf(val)
	if !ok {
		return false
	}

	return o.OwnsKey(key)
}

func (o *OwnedResources) OwnsKey(key ResourceKey) bool {
	if o.closers == nil {
		return false
	}

	_, exists := o.closers[key]
	return exists
}

// Extract removes all ownership tracking for val from the current frame without
// scheduling the closer for deferred cleanup. Use it when a value survives
// frame teardown, for example as a return value or final result root.
func (o *OwnedResources) Extract(val runtime.Value) bool {
	key, _, ok := ResourceKeyOf(val)
	if !ok {
		return false
	}

	return o.ExtractByKey(key)
}

func (o *OwnedResources) ExtractByKey(key ResourceKey) bool {
	if o.closers == nil {
		return false
	}

	if _, exists := o.closers[key]; !exists {
		return false
	}

	delete(o.closers, key)
	if len(o.closers) == 0 {
		o.closers = nil
	}

	return true
}

// ExtractMany removes ownership tracking for the provided surviving values from
// the current frame and reassigns the matching owned closers to dst. Duplicate
// values transfer ownership only once per resource, while dead aliases in the
// retiring frame are dropped instead of being deferred.
func (o *OwnedResources) ExtractMany(values []runtime.Value, dst *OwnedResources) {
	if o.closers == nil || len(values) == 0 {
		return
	}

	survivors := make(map[ResourceKey]io.Closer)

	for _, val := range values {
		key, _, ok := ResourceKeyOf(val)
		if !ok {
			continue
		}

		closer, exists := o.closers[key]
		if !exists {
			continue
		}

		survivors[key] = closer
	}

	if len(survivors) == 0 {
		return
	}

	if dst != nil && dst.closers == nil {
		dst.closers = make(map[ResourceKey]io.Closer)
	}

	for key, closer := range survivors {
		delete(o.closers, key)

		if dst != nil {
			dst.closers[key] = closer
		}
	}

	if len(o.closers) == 0 {
		o.closers = nil
	}
}

func (o *OwnedResources) Discard(val runtime.Value, deferred *DeferredClosers) {
	key, _, ok := ResourceKeyOf(val)
	if !ok || o.closers == nil {
		return
	}

	closer, exists := o.closers[key]
	if !exists {
		return
	}

	delete(o.closers, key)
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
		for _, closer := range o.closers {
			deferred.AddCloser(closer)
		}
	}

	o.closers = nil
}

func (o *OwnedResources) Release(val runtime.Value) (io.Closer, bool) {
	key, _, ok := ResourceKeyOf(val)
	if !ok || o.closers == nil {
		return nil, false
	}

	closer, exists := o.closers[key]
	if !exists {
		return nil, false
	}

	// Release is terminal for explicit close, not a generic ref-count decrement:
	// once one alias is closed, the closer is retired from ownership tracking so
	// stale aliases cannot close it again during later cleanup.
	delete(o.closers, key)
	if len(o.closers) == 0 {
		o.closers = nil
	}

	return closer, true
}

// DiscardByKey removes ownership of a resource by its pre-resolved key and
// schedules the closer for deferred cleanup. Unlike Discard, this avoids a
// redundant ResourceKeyOf call when the caller already holds the key.
func (o *OwnedResources) DiscardByKey(key ResourceKey, closer io.Closer, deferred *DeferredClosers) {
	if o.closers == nil {
		return
	}

	if _, exists := o.closers[key]; !exists {
		return
	}

	delete(o.closers, key)
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

func (o *OwnedResources) Empty() bool {
	return len(o.closers) == 0
}

func (o *OwnedResources) ForEach(fn func(io.Closer)) {
	if fn == nil {
		return
	}

	for _, closer := range o.closers {
		fn(closer)
	}
}

// TrackedCloserOf extracts an io.Closer from val if it implements the
// interface. Tracked closers are pointer-comparable by construction
// (pointer-receiver io.Closer or runtime.Resource implementations).
//
// Deprecated: prefer ResourceKeyOf when the ResourceKey is also needed.
func TrackedCloserOf(val runtime.Value) (io.Closer, bool) {
	closer, ok := val.(io.Closer)
	if !ok || closer == nil {
		return nil, false
	}

	return closer, true
}
