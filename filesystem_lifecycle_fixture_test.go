package ferret

import (
	"errors"
	"sync/atomic"

	ferretfs "github.com/MontFerret/ferret/v2/pkg/fs"
)

type failingCloseFileSystem struct {
	ferretfs.FileSystem
	closeErr error
	closes   atomic.Int32
}

func (f *failingCloseFileSystem) Close() error {
	f.closes.Add(1)

	return errors.Join(f.closeErr, f.FileSystem.Close())
}
