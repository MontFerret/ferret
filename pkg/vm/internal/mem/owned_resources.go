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
	closers map[io.Closer]struct{}
}

func (o *OwnedResources) Track(val runtime.Value) {
	closer, ok := trackedCloserOf(val)
	if !ok {
		return
	}

	if o.closers == nil {
		o.closers = make(map[io.Closer]struct{})
	}

	o.closers[closer] = struct{}{}
}

func (o *OwnedResources) Forget(val runtime.Value) {
	closer, ok := trackedCloserOf(val)
	if !ok || o.closers == nil {
		return
	}

	delete(o.closers, closer)
	if len(o.closers) == 0 {
		o.closers = nil
	}
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
			dst.closers = make(map[io.Closer]struct{})
		}

		dst.closers[closer] = struct{}{}
	}
}

func (o *OwnedResources) TransferMany(values []runtime.Value, dst *OwnedResources) {
	for _, val := range values {
		o.Transfer(val, dst)
	}
}

func (o *OwnedResources) CloseAll() {
	for closer := range o.closers {
		_ = closer.Close()
	}

	o.closers = nil
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
