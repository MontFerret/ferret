package fs_test

import (
	"context"
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
			file, delFile := tempFile()
			defer delFile()

			text := "s string"
			file.WriteString(text)

			fname := runtime.NewString(file.Name())

			out, err := fs.Read(context.Background(), fname)
			So(err, ShouldBeNil)

			SoMsg("Output should be binary", runtime.AssertBinary(out), ShouldBeNil)

			So(out.String(), ShouldEqual, text)
		})

		Convey("File does not exist", func() {
			fname := runtime.NewString("not_exist.file")

			out, err := fs.Read(context.Background(), fname)
			So(out, ShouldEqual, runtime.None)
			So(err, ShouldBeError)
		})
	})
}
