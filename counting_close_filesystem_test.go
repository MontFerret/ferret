package ferret

import (
	"errors"
	"sync/atomic"

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
