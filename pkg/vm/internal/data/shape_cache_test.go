package data

import (
	"reflect"
	"testing"
)

var benchmarkShapeCacheResult *Shape

func TestShapeCacheTransition(t *testing.T) {
	t.Parallel()

	cache := NewShapeCache(8)
	root := cache.Root()
	first := cache.Transition(root, "first")
	second := cache.Transition(first, "second")

	if actual, expected := first.names, []string{"first"}; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("first names = %v, want %v", actual, expected)
	}

	if actual, expected := second.names, []string{"first", "second"}; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("second names = %v, want %v", actual, expected)
	}

	if slot := second.fields["first"]; slot != 0 {
		t.Fatalf("first slot = %d, want 0", slot)
	}

	if slot := second.fields["second"]; slot != 1 {
		t.Fatalf("second slot = %d, want 1", slot)
	}

	if _, exists := first.fields["second"]; exists {
		t.Fatal("transition mutated the previous shape")
	}

	if cached := cache.Transition(root, "first"); cached != first {
		t.Fatal("repeated transition did not reuse the cached shape")
	}
}

func TestShapeCacheTransitionAfterLimitIsNotCached(t *testing.T) {
	t.Parallel()

	cache := NewShapeCache(1)
	root := cache.Root()
	cache.Transition(root, "cached")

	first := cache.Transition(root, "uncached")
	second := cache.Transition(root, "uncached")

	if first == second {
		t.Fatal("transition after the cache limit was unexpectedly reused")
	}

	if actual, expected := first.names, []string{"uncached"}; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("uncached names = %v, want %v", actual, expected)
	}
}

func BenchmarkShapeCacheTransitionChain(b *testing.B) {
	keys := []string{"first", "second", "third", "fourth"}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cache := NewShapeCache(len(keys))
		shape := cache.Root()

		for _, key := range keys {
			shape = cache.Transition(shape, key)
		}

		benchmarkShapeCacheResult = shape
	}
}
