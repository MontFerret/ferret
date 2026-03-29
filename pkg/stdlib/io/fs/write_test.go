package fs_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/io/fs"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWrite(t *testing.T) {

	Convey("Invalid arguments", t, func() {

		path := runtime.NewString("path.txt")
		data := runtime.NewBinary([]byte("3timeslazy"))
		params := runtime.NewObjectWith(
			map[string]runtime.Value{
				"mode": runtime.NewString("foo"),
			},
		)
		someInt := runtime.NewInt(0)
		someStr := runtime.NewString("test")

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
					Name: "Arguments Type: `path` not a string",
					Args: []runtime.Value{someInt, data},
				},
				{
					Name: "Arguments Type: `data` not a binary",
					Args: []runtime.Value{path, someStr},
				},
				{
					Name: "Arguments Type: `params` not an object",
					Args: []runtime.Value{path, data, someInt},
				},
				{
					Name: "Arguments Type: `params` contains invalid `mode`",
					Args: []runtime.Value{path, data, params},
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
							map[string]runtime.Value{
								"mode": mode,
							},
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
			ctx, root, path, cleanup := tempFileSystemContext()
			defer cleanup()

			err := os.WriteFile(filepath.Join(root, path), []byte("existing"), 0o666)
			So(err, ShouldBeNil)

			none, err := fs.Write(
				ctx,
				runtime.NewString(path),
				runtime.NewBinary([]byte("3timeslazy")),
				runtime.NewObjectWith(
					map[string]runtime.Value{
						"mode": runtime.NewString("wx"),
					},
				),
			)
			So(none, ShouldResemble, runtime.None)
			So(err, ShouldBeError)
		})

		Convey("Filepath is empty", func() {
			ctx, _, _, cleanup := tempFileSystemContext()
			defer cleanup()

			none, err := fs.Write(
				ctx,
				runtime.NewString(""),
				runtime.NewBinary([]byte("3timeslazy")),
			)
			So(none, ShouldResemble, runtime.None)
			So(err, ShouldBeError)
		})
	})

	Convey("Success cases", t, func() {
		Convey("Mode `w` should truncate file", func() {
			ctx, _, path, cleanup := tempFileSystemContext()
			defer cleanup()

			data := runtime.NewBinary([]byte("3timeslazy"))
			fpath := runtime.NewString(path)
			params := runtime.NewObjectWith(
				map[string]runtime.Value{
					"mode": runtime.NewString("w"),
				},
			)

			for range [2]struct{}{} {
				// at first iteration check that `Write` creates file and writes `data`.
				// At second iteration check that `Write` truncates the file and
				// writes `data` again

				_, err := fs.Write(ctx, fpath, data, params)
				So(err, ShouldBeNil)

				read, err := fs.Read(ctx, fpath)
				So(err, ShouldBeNil)
				So(read, ShouldResemble, data)
			}
		})

		Convey("Mode `a` should append into file", func() {
			ctx, _, path, cleanup := tempFileSystemContext()
			defer cleanup()

			data := runtime.NewBinary([]byte("3timeslazy"))
			fpath := runtime.NewString(path)
			params := runtime.NewObjectWith(
				map[string]runtime.Value{
					"mode": runtime.NewString("a"),
				},
			)

			for i := range [2]struct{}{} {
				// at first iteration check that `Write` creates file and appends `data`
				// into it.
				// At second iteration check that `Write` appends `data` into file
				// one more time using bytes.Repeat

				_, err := fs.Write(ctx, fpath, data, params)
				So(err, ShouldBeNil)

				read, err := fs.Read(ctx, fpath)
				So(err, ShouldBeNil)

				readBytes := runtime.Unwrap(read)
				So(readBytes, ShouldResemble, bytes.Repeat(data, i+1))
			}
		})

		Convey("Write string data should error", func() {
			ctx, _, path, cleanup := tempFileSystemContext()
			defer cleanup()

			text := "test string data"
			fpath := runtime.NewString(path)

			none, err := fs.Write(ctx, fpath, runtime.NewString(text))
			So(err, ShouldBeError)
			So(none, ShouldResemble, runtime.None)
		})
	})
}
