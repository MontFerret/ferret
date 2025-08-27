package fs_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/io/fs"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWrite(t *testing.T) {

	Convey("Invalid arguments", t, func() {

		path := runtime.NewString("path.txt")
		data := runtime.NewBinary([]byte("3timeslazy"))
		params := runtime.NewObjectWith(
			runtime.NewObjectProperty("mode", runtime.NewString("w")),
		)
		someInt := runtime.NewInt(0)

		Convey("All arguments", func() {

			testCases := []struct {
				Name string
				Args []runtime.Value
			}{
				{
					Name: "Arguments Number: No arguments passed",
					Args: []runtime.Value{},
				},
				{
					Name: "Arguments Number: Only `path` passed",
					Args: []runtime.Value{path},
				},
				{
					Name: "Arguments Number: more than 3 arguments passed",
					Args: []runtime.Value{path, data, params, someInt},
				},
				{
					Name: "Arguments Type: `path` not a string",
					Args: []runtime.Value{someInt},
				},
				{
					Name: "Arguments Type: `params` not an object",
					Args: []runtime.Value{path, data, someInt},
				},
			}

			for _, tC := range testCases {
				tC := tC

				Convey(tC.Name, func() {
					none, err := fs.Write(context.Background(), tC.Args...)
					So(err, ShouldBeError)
					So(none, ShouldResemble, runtime.None)
				})
			}
		})

		Convey("Argument `params`", func() {

			Convey("First `mode`", func() {

				testCases := []runtime.Value{
					// empty mode string
					runtime.NewString(""),

					// `a` and `w` cannot be used together
					runtime.NewString("aw"),

					// two equal letters
					runtime.NewString("ww"),

					// mode string is too long
					runtime.NewString("awx"),

					// the `x` mode only
					runtime.NewString("x"),

					// mode is not a string
					runtime.NewInt(1),
				}

				for _, mode := range testCases {
					mode := mode
					name := fmt.Sprintf("mode `%s`", mode)

					Convey(name, func() {
						params := runtime.NewObjectWith(
							runtime.NewObjectProperty("mode", mode),
						)

						none, err := fs.Write(context.Background(), path, data, params)
						So(err, ShouldBeError)
						So(none, ShouldResemble, runtime.None)
					})
				}
			})
		})
	})

	Convey("Error cases", t, func() {

		Convey("Write into existing file with `x` mode", func() {
			file, delFile := tempFile()
			defer delFile()

			none, err := fs.Write(
				context.Background(),
				runtime.NewString(file.Name()),
				runtime.NewBinary([]byte("3timeslazy")),
				runtime.NewObjectWith(
					runtime.NewObjectProperty("mode", runtime.NewString("wx")),
				),
			)
			So(none, ShouldResemble, runtime.None)
			So(err, ShouldBeError)
		})

		Convey("Filepath is empty", func() {
			none, err := fs.Write(
				context.Background(),
				runtime.NewString(""),
				runtime.NewBinary([]byte("3timeslazy")),
			)
			So(none, ShouldResemble, runtime.None)
			So(err, ShouldBeError)
		})
	})

	Convey("Success cases", t, func() {

		Convey("Mode `w` should truncate file", func() {
			file, delFile := tempFile()
			defer delFile()

			data := runtime.NewBinary([]byte("3timeslazy"))
			fpath := runtime.NewString(file.Name())
			params := runtime.NewObjectWith(
				runtime.NewObjectProperty("mode", runtime.NewString("w")),
			)

			for _ = range [2]struct{}{} {
				// at first iteration check that `Write` creates file and writes `data`.
				// At second iteration check that `Write` truncates the file and
				// writes `data` again

				_, err := fs.Write(context.Background(), fpath, data, params)
				So(err, ShouldBeNil)

				read, err := fs.Read(context.Background(), fpath)
				So(err, ShouldBeNil)
				So(read, ShouldResemble, data)
			}
		})

		Convey("Mode `a` should append into file", func() {
			file, delFile := tempFile()
			defer delFile()

			data := runtime.NewBinary([]byte("3timeslazy"))
			fpath := runtime.NewString(file.Name())
			params := runtime.NewObjectWith(
				runtime.NewObjectProperty("mode", runtime.NewString("a")),
			)

			for i := range [2]struct{}{} {
				// at first iteration check that `Write` creates file and appends `data`
				// into it.
				// At second iteration check that `Write` appends `data` into file
				// one more time using bytes.Repeat

				_, err := fs.Write(context.Background(), fpath, data, params)
				So(err, ShouldBeNil)

				read, err := fs.Read(context.Background(), fpath)
				So(err, ShouldBeNil)

				readBytes := read.Unwrap().([]byte)
				So(readBytes, ShouldResemble, bytes.Repeat(data, i+1))
			}
		})
	})
}
