package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParam(t *testing.T) {
	Convey("Should be possible to use as a return value", t, func() {
		out := compiler.New().
			MustCompile(`
			RETURN @param
		`).
			MustRun(context.Background(), runtime.WithParam("param", "foobar"))

		So(string(out), ShouldEqual, `"foobar"`)
	})

	Convey("Should be possible to use as a FOR source", t, func() {
		out := compiler.New().
			MustCompile(`
			FOR i IN @values
			SORT i
			RETURN i
		`).
			MustRun(context.Background(), runtime.WithParam("values", []int{1, 2, 3, 4}))

		So(string(out), ShouldEqual, `[1,2,3,4]`)

		out2 := compiler.New().
			MustCompile(`
			FOR i IN @values
			SORT i
			RETURN i
		`).
			MustRun(context.Background(), runtime.WithParam("values", map[string]int{
				"foo": 1,
				"bar": 2,
				"faz": 3,
				"qaz": 4,
			}))

		So(string(out2), ShouldEqual, `[1,2,3,4]`)
	})

	Convey("Should be possible to use in range", t, func() {
		prog := compiler.New().
			MustCompile(`
			FOR i IN @start..@end
			SORT i
			RETURN i
		`)

		out := prog.MustRun(
			context.Background(),
			runtime.WithParam("start", 1),
			runtime.WithParam("end", 4),
		)

		So(string(out), ShouldEqual, `[1,2,3,4]`)

	})

	Convey("Should be possible to use in member expression", t, func() {
		prog := compiler.New().
			MustCompile(`
			RETURN @param.value
		`)

		out := prog.MustRun(
			context.Background(),
			runtime.WithParam("param", map[string]interface{}{
				"value": "foobar",
			}),
		)

		So(string(out), ShouldEqual, `"foobar"`)

	})

	Convey("Should be possible to use in member expression as a computed property", t, func() {
		prog := compiler.New().
			MustCompile(`
			LET obj = { foo: "bar" }
			RETURN obj[@param]
		`)

		out := prog.MustRun(
			context.Background(),
			runtime.WithParam("param", "foo"),
		)

		So(string(out), ShouldEqual, `"bar"`)
	})

	Convey("Should be possible to use in member expression as segments", t, func() {
		prog := compiler.New().
			MustCompile(`
			LET doc = { foo: { bar: "baz" } }

			RETURN doc.@attr.@subattr
		`)

		out := prog.MustRun(
			context.Background(),
			runtime.WithParam("attr", "foo"),
			runtime.WithParam("subattr", "bar"),
		)

		So(string(out), ShouldEqual, `"baz"`)
	})

	Convey("Should return an error if param values are not passed", t, func() {
		prog := compiler.New().
			MustCompile(`
			LET doc = { foo: { bar: "baz" } }

			RETURN doc.@attr.@subattr
		`)

		_, err := prog.Run(
			context.Background(),
		)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, runtime.ErrMissedParam.Error())
	})

	Convey("Should be possible to use in member expression as segments", t, func() {
		prog := compiler.New().
			MustCompile(`
			LET doc = { foo: { bar: "baz" } }

			RETURN doc.@attr.@subattr
		`)

		_, err := prog.Run(
			context.Background(),
			runtime.WithParam("attr", "foo"),
		)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "subattr")
	})

	Convey("Should be possible to use in struct with nested struct which nil", t, func() {
		type Some2 struct {
		}
		type Some struct {
			Some2 *Some2
		}

		someObj := &Some{}
		prog := compiler.New().
			MustCompile(`

			RETURN null
		`)

		panics := func() {
			_, _ = prog.Run(
				context.Background(),
				runtime.WithParam("struct", someObj),
			)
		}

		So(panics, ShouldNotPanic)
	})
}
