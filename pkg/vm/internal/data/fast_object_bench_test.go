package data

import (
	"strconv"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var benchmarkFastObjectResult runtime.Value

func BenchmarkFastObject(b *testing.B) {
	const threshold = 64
	keys := make([]string, threshold+1)
	for i := range keys {
		keys[i] = "field" + strconv.Itoa(i)
	}

	cache := NewShapeCache(threshold)
	build := func(count int) *FastObject {
		object := NewFastObject(cache, threshold)
		for _, key := range keys[:count] {
			object.SetStringCached(key, runtime.True)
		}
		return object
	}

	shape := build(threshold)
	dictionary := build(threshold + 1)
	if shape.Shape() == nil || dictionary.Shape() != nil {
		b.Fatal("fixtures must straddle the dictionary transition")
	}

	b.Run("shape_build_4", func(b *testing.B) {
		for b.Loop() {
			benchmarkFastObjectResult = build(4)
		}
	})
	b.Run("shape_overwrite", func(b *testing.B) {
		for b.Loop() {
			shape.SetStringCached(keys[0], runtime.False)
		}
	})
	b.Run("transition", func(b *testing.B) {
		for b.Loop() {
			// The threshold-crossing write reads the old slots and detaches.
			object := *shape
			object.SetStringCached(keys[threshold], runtime.True)
		}
	})
	b.Run("dictionary_build_65", func(b *testing.B) {
		for b.Loop() {
			benchmarkFastObjectResult = build(threshold + 1)
		}
	})
	b.Run("dictionary_overwrite", func(b *testing.B) {
		for b.Loop() {
			dictionary.SetStringCached(keys[0], runtime.False)
		}
	})
	b.Run("dictionary_insert_remove", func(b *testing.B) {
		object := build(threshold + 1)
		for b.Loop() {
			object.SetStringCached("extra", runtime.True)
			object.removeString("extra")
		}
	})
	b.Run("dictionary_copy", func(b *testing.B) {
		for b.Loop() {
			benchmarkFastObjectResult = dictionary.Copy()
		}
	})
	b.Run("dictionary_clone", func(b *testing.B) {
		ctx := b.Context()
		for b.Loop() {
			cloned, err := dictionary.Clone(ctx)
			if err != nil {
				b.Fatal(err)
			}
			benchmarkFastObjectResult = cloned
		}
	})
	b.Run("dictionary_clear_reuse", func(b *testing.B) {
		ctx := b.Context()
		object := build(threshold + 1)
		for b.Loop() {
			if err := object.Clear(ctx); err != nil {
				b.Fatal(err)
			}
			object.SetStringCached(keys[0], runtime.True)
		}
	})
}
