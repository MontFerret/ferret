package fs_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/io/fs"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRead(t *testing.T) {
	Convey("Arguments passed", t, func() {
		Convey("Passed not a string", func() {
			out, err := fs.Read(context.Background(), runtime.NewInt(0))

			So(out, ShouldEqual, runtime.None)
			So(err, ShouldBeError)
		})
	})

	Convey("Read from file", t, func() {
		Convey("File exists", func() {
			ctx, root, path, cleanup := tempFileSystemContext()
			defer cleanup()

			text := "s string"
			err := os.WriteFile(filepath.Join(root, path), []byte(text), 0o666)
			So(err, ShouldBeNil)

			fname := runtime.NewString(path)

			out, err := fs.Read(ctx, fname)
			So(err, ShouldBeNil)

			SoMsg("Output should be binary", runtime.AssertBinary(out), ShouldBeNil)

			So(out.String(), ShouldEqual, text)
		})

		Convey("File does not exist", func() {
			ctx, _, _, cleanup := tempFileSystemContext()
			defer cleanup()

			fname := runtime.NewString("not_exist.file")

			out, err := fs.Read(ctx, fname)
			So(out, ShouldEqual, runtime.None)
			So(err, ShouldBeError)
		})
	})
}
