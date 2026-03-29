package fs_test

import (
	"context"
	"os"

	ferretfs "github.com/MontFerret/ferret/v2/pkg/fs"

	. "github.com/smartystreets/goconvey/convey"
)

type closeable interface {
	Close() error
}

func tempFileSystemContext() (context.Context, string, string, func()) {
	root, err := os.MkdirTemp("", "fstest")
	So(err, ShouldBeNil)

	filesystem, err := ferretfs.New(ferretfs.WithRoot(root))
	So(err, ShouldBeNil)

	ctx := ferretfs.WithFileSystem(context.Background(), filesystem)
	path := "test.txt"

	cleanup := func() {
		if c, ok := filesystem.(closeable); ok {
			_ = c.Close()
		}

		_ = os.RemoveAll(root)
	}

	return ctx, root, path, cleanup
}
