package runtime

import (
	"strconv"
	"testing"
)

var benchmarkObjectValue Value

func BenchmarkObjectOperations(b *testing.B) {
	ctx := b.Context()
	keys := make([]String, 256)
	for i := range keys {
		keys[i] = String("field" + strconv.Itoa(i))
	}

	build := func(hint, size int) *Object {
		obj := NewObjectOf(hint)
		for _, key := range keys[:size] {
			obj.SetString(string(key), True)
		}

		return obj
	}

	for _, size := range []int{1, 2, 8, 16, 64, 128, 256} {
		name := strconv.Itoa(size)
		b.Run("grow_"+name, func(b *testing.B) {
			for b.Loop() {
				benchmarkObjectValue = build(0, size)
			}
		})
		b.Run("reserved_"+name, func(b *testing.B) {
			for b.Loop() {
				benchmarkObjectValue = build(size, size)
			}
		})
	}

	for _, size := range []int{1, 16, 256} {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			obj := build(size, size)
			key := Value(keys[size-1])
			name := string(keys[size-1])
			b.Run("overwrite", func(b *testing.B) {
				for b.Loop() {
					obj.SetString(name, False)
				}
			})
			b.Run("set", func(b *testing.B) {
				for b.Loop() {
					_ = obj.Set(ctx, key, True)
				}
			})
			b.Run("lookup", func(b *testing.B) {
				for b.Loop() {
					benchmarkObjectValue, _, _ = obj.Lookup(ctx, key)
				}
			})
			b.Run("insert_remove", func(b *testing.B) {
				obj := build(size+1, size)
				for b.Loop() {
					obj.SetString("extra", True)
					_ = obj.RemoveKey(ctx, String("extra"))
				}
			})
			b.Run("iterate", func(b *testing.B) {
				for b.Loop() {
					iter, _ := obj.Iterate(ctx)
					for {
						v, _, err := iter.Next(ctx)
						if err != nil {
							break
						}
						benchmarkObjectValue = v
					}
				}
			})
			b.Run("copy", func(b *testing.B) {
				for b.Loop() {
					benchmarkObjectValue = obj.Copy()
				}
			})
			b.Run("clone", func(b *testing.B) {
				for b.Loop() {
					cloned, err := obj.Clone(ctx)
					if err != nil {
						b.Fatal(err)
					}
					benchmarkObjectValue = cloned.(Value)
				}
			})
		})
	}

	b.Run("clear_reuse", func(b *testing.B) {
		obj := NewObject()
		for b.Loop() {
			_ = obj.Clear(ctx)
			obj.SetString("x", True)
		}
	})
}
