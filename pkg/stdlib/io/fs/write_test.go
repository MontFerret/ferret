package fs_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/io/fs"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWrite(t *testing.T) {

	Convey("Invalid arguments", t, func() {

		path := values.NewString("path.txt")
		data := values.NewBinary([]byte("3timeslazy"))
		params := values.NewObjectWith(
			values.NewObjectProperty("mode", values.NewString("w")),
		)
		someInt := values.NewInt(0)

		Convey("All arguments", func() {

			testCases := []struct {
				Name string
				Args []core.Value
			}{
				{
					Name: "Arguments Number: No arguments passed",
					Args: []core.Value{},
				},
				{
					Name: "Arguments Number: Only `path` passed",
					Args: []core.Value{path},
				},
				{
					Name: "Arguments Number: more than 3 arguments passed",
					Args: []core.Value{path, data, params, someInt},
				},
				{
					Name: "Arguments Type: `path` not a string",
					Args: []core.Value{someInt},
				},
				{
					Name: "Arguments Type: `params` not an object",
					Args: []core.Value{path, data, someInt},
				},
			}

			for _, tC := range testCases {
				tC := tC

				Convey(tC.Name, func() {
					none, err := fs.Write(context.Background(), tC.Args...)
					So(err, ShouldBeError)
					So(none, ShouldResemble, values.None)
				})
			}
		})

		Convey("Argument `params`", func() {

			Convey("Key `mode`", func() {

				testCases := []core.Value{
					// empty mode string
					values.NewString(""),

					// `a` and `w` cannot be used together
					values.NewString("aw"),

					// two equal letters
					values.NewString("ww"),

					// mode string is too long
					values.NewString("awx"),

					// the `x` mode only
					values.NewString("x"),

					// mode is not a string
					values.NewInt(1),
				}

				for _, mode := range testCases {
					mode := mode
					name := fmt.Sprintf("mode `%s`", mode)

					Convey(name, func() {
						params := values.NewObjectWith(
							values.NewObjectProperty("mode", mode),
						)

						none, err := fs.Write(context.Background(), path, data, params)
						So(err, ShouldBeError)
						So(none, ShouldResemble, values.None)
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
				values.NewString(file.Name()),
				values.NewBinary([]byte("3timeslazy")),
				values.NewObjectWith(
					values.NewObjectProperty("mode", values.NewString("wx")),
				),
			)
			So(none, ShouldResemble, values.None)
			So(err, ShouldBeError)
		})

		Convey("Filepath is empty", func() {
			none, err := fs.Write(
				context.Background(),
				values.NewString(""),
				values.NewBinary([]byte("3timeslazy")),
			)
			So(none, ShouldResemble, values.None)
			So(err, ShouldBeError)
		})
	})

	Convey("Success cases", t, func() {

		Convey("Mode `w` should truncate file", func() {
			file, delFile := tempFile()
			defer delFile()

			data := values.NewBinary([]byte("3timeslazy"))
			fpath := values.NewString(file.Name())
			params := values.NewObjectWith(
				values.NewObjectProperty("mode", values.NewString("w")),
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

			data := values.NewBinary([]byte("3timeslazy"))
			fpath := values.NewString(file.Name())
			params := values.NewObjectWith(
				values.NewObjectProperty("mode", values.NewString("a")),
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
