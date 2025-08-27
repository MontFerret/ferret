package fs_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/io/fs"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRead(t *testing.T) {

	Convey("Arguments passed", t, func() {

		Convey("No arguments passed", func() {
			out, err := fs.Read(context.Background())

			So(out, ShouldEqual, runtime.None)
			So(err, ShouldBeError)
		})

		Convey("Passed not a string", func() {
			args := []runtime.Value{runtime.NewInt(0)}
			out, err := fs.Read(context.Background(), args...)

			So(out, ShouldEqual, runtime.None)
			So(err, ShouldBeError)
		})

		Convey("Passed more that one argument", func() {
			args := []runtime.Value{
				runtime.NewString("filepath"),
				runtime.NewInt(0),
			}
			out, err := fs.Read(context.Background(), args...)

			So(out, ShouldEqual, runtime.None)
			So(err, ShouldBeError)
		})
	})

	Convey("Read from file", t, func() {

		//Convey("File exists", func() {
		//	file, delFile := tempFile()
		//	defer delFile()
		//
		//	text := "s string"
		//	file.WriteString(text)
		//
		//	fname := runtime.NewString(file.Name())
		//
		//	out, err := fs.Read(context.Background(), fname)
		//	So(err, ShouldBeNil)
		//
		//	So(out.Type().ID(), ShouldEqual, types.Binary.ID())
		//	So(out.String(), ShouldEqual, text)
		//})

		Convey("File does not exist", func() {
			fname := runtime.NewString("not_exist.file")

			out, err := fs.Read(context.Background(), fname)
			So(out, ShouldEqual, runtime.None)
			So(err, ShouldBeError)
		})
	})
}
