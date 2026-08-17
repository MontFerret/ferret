package ferret

import "testing"

func BenchmarkEngineFSRootLifecycle(b *testing.B) {
	root := b.TempDir()

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		engine, err := New(WithFSRoot(root))
		if err != nil {
			b.Fatalf("new engine: %v", err)
		}

		if err := engine.Close(); err != nil {
			b.Fatalf("close engine: %v", err)
		}
	}
}
