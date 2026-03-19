package mem

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var (
	benchCellHandle CellHandle
	benchCellValue  runtime.Value
)

func BenchmarkCellStore_New(b *testing.B) {
	store := newBenchmarkCellStore(b.N)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchCellHandle = store.New(runtime.True)
	}
}

func BenchmarkCellStore_Get(b *testing.B) {
	store := newBenchmarkCellStore(1)
	handle := store.New(runtime.True)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		val, ok := store.Get(handle)
		if !ok {
			b.Fatal("expected benchmark handle to stay valid")
		}

		benchCellValue = val
	}
}

func BenchmarkCellStore_Set(b *testing.B) {
	store := newBenchmarkCellStore(1)
	handle := store.New(runtime.True)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		val := runtime.Value(runtime.True)
		if i&1 == 1 {
			val = runtime.False
		}

		if !store.Set(handle, val) {
			b.Fatal("expected benchmark handle to stay valid")
		}
	}
}

func BenchmarkCellStore_DeleteThenNew(b *testing.B) {
	store := newBenchmarkCellStore(1)
	handle := store.New(runtime.True)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		val, ok := store.Delete(handle)
		if !ok {
			b.Fatal("expected benchmark handle to stay valid")
		}

		benchCellValue = val
		handle = store.New(runtime.True)
	}

	benchCellHandle = handle
}

func BenchmarkCellStore_ResetThenNew(b *testing.B) {
	store := newBenchmarkCellStore(1)
	handle := store.New(runtime.True)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Reset()
		handle = store.New(runtime.True)
	}

	benchCellHandle = handle
}

func newBenchmarkCellStore(capacity int) CellStore {
	if capacity < 1 {
		capacity = 1
	}

	return CellStore{
		values:     make(map[CellHandle]runtime.Value, capacity),
		generation: 1,
		nextSlot:   1,
		nextToken:  1,
	}
}
