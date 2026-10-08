package codecutil

import (
	"context"
	"io"
	"reflect"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Iterator owns only a distinct iterator acquired for encoding. An iterator
// sharing the borrowed source's identity remains with the source's owner.
type Iterator struct {
	runtime.Iterator
	closer io.Closer
}

// NewIterator acquires an iterator without taking ownership of the source.
func NewIterator(ctx context.Context, source runtime.Iterable) (Iterator, error) {
	iter, err := source.Iterate(ctx)
	if err != nil {
		return Iterator{}, err
	}

	handle := Iterator{Iterator: iter}
	if closer, ok := iter.(io.Closer); ok && !handle.borrows(source) {
		handle.closer = closer
	}

	return handle, nil
}

// Close releases an owned iterator once. Codecs call it after child encoding so
// iterator cleanup failures cannot skip those children's ownership hooks.
func (iter *Iterator) Close() error {
	closer := iter.closer
	iter.closer = nil
	if closer == nil {
		return nil
	}

	return closer.Close()
}

func (iter *Iterator) borrows(source runtime.Iterable) bool {
	if resource, ok := iter.Iterator.(runtime.Resource); ok {
		if owner, ok := source.(runtime.Resource); ok {
			return resource.ResourceID() == owner.ResourceID()
		}
	}

	if _, ok := source.(io.Closer); !ok {
		return false
	}

	// Pointers compare safely against any dynamic type. Checking their kind
	// avoids the general value-comparability check on this common path.
	if reflect.TypeOf(source).Kind() == reflect.Pointer {
		return any(source) == any(iter.Iterator)
	}

	// Value.Comparable also checks interface fields; Type.Comparable alone can
	// still allow a comparison that panics on a dynamically non-comparable field.
	return reflect.ValueOf(source).Comparable() && reflect.ValueOf(iter.Iterator).Comparable() && any(source) == any(iter.Iterator)
}
