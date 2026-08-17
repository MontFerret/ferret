package fs

import (
	"os"
	"sync"
	"testing"
)

func TestRootFSCloseIsConcurrentAndIdempotent(t *testing.T) {
	root, err := os.OpenRoot(t.TempDir())
	if err != nil {
		t.Fatalf("open root: %v", err)
	}

	filesystem := &rootFS{root: root}

	const callers = 32
	start := make(chan struct{})
	errs := make(chan error, callers)
	var callersDone sync.WaitGroup
	callersDone.Add(callers)

	for range callers {
		go func() {
			defer callersDone.Done()
			<-start
			errs <- filesystem.Close()
		}()
	}

	close(start)
	callersDone.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("close filesystem: %v", err)
		}
	}

	if _, err := filesystem.Stat("."); err == nil {
		t.Fatal("expected root to be closed")
	}
}
