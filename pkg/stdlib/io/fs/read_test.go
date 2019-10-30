package fs_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/MontFerret/ferret/pkg/stdlib/io/fs"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRead(t *testing.T) {

	Convey("Arguments passed", t, func() {

		Convey("No arguments passed", func() {
			out, err := fs.Read(context.Background())

			So(out, ShouldEqual, values.None)
			So(err, ShouldBeError)
		})

		Convey("Passed not a string", func() {
			args := []core.Value{values.NewInt(0)}
			out, err := fs.Read(context.Background(), args...)

			So(out, ShouldEqual, values.None)
			So(err, ShouldBeError)
		})

		Convey("Passed more that one argument", func() {
			args := []core.Value{
				values.NewString("filepath"),
				values.NewInt(0),
			}
			out, err := fs.Read(context.Background(), args...)

			So(out, ShouldEqual, values.None)
			So(err, ShouldBeError)
		})
	})

	Convey("Read from file", t, func() {

		Convey("File exists", func() {
			file, delFile := tempFile()
			defer delFile()

			text := "s string"
			file.WriteString(text)

			fname := values.NewString(file.Name())

			out, err := fs.Read(context.Background(), fname)
			So(err, ShouldBeNil)

			So(out.Type().ID(), ShouldEqual, types.Binary.ID())
			So(out.String(), ShouldEqual, text)
		})

		Convey("File does not exist", func() {
			fname := values.NewString("not_exist.file")

			out, err := fs.Read(context.Background(), fname)
			So(out, ShouldEqual, values.None)
			So(err, ShouldBeError)
		})
	})
}
