package engine

import (
	"errors"
	"sync/atomic"
	"testing"

	ferretfs "github.com/MontFerret/ferret/v2/pkg/fs"
)

type countingCloseFileSystem struct {
	ferretfs.FileSystem
	closeErr   error
	closeCalls atomic.Int32
}

func (f *countingCloseFileSystem) Close() error {
	f.closeCalls.Add(1)

	return errors.Join(f.closeErr, f.FileSystem.Close())
}

func newCountingCloseFileSystem(t testing.TB, root string, closeErr error) *countingCloseFileSystem {
	t.Helper()

	filesystem, err := ferretfs.New(ferretfs.WithRoot(root))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = filesystem.Close() })

	return &countingCloseFileSystem{FileSystem: filesystem, closeErr: closeErr}
}
