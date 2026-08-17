package ferret

import (
	"io"
	"testing"
)

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

		// Keep a direct close outside the measured interval so a lifecycle
		// regression cannot make the benchmark accumulate descriptors.
		b.StopTimer()
		if closer, ok := engine.host.fs.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				b.Fatalf("close filesystem: %v", err)
			}
		}
		b.StartTimer()
	}
}
