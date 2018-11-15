package compiler_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	. "github.com/smartystreets/goconvey/convey"
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
		`).MustRun(context.Background(), runtime.WithParam("values", []int64{1, 2, 3, 4}))

		So(string(out), ShouldEqual, `[1,2,3,4]`)

		out2 := compiler.New().
			MustCompile(`
			FOR i IN @values
			SORT i
			RETURN i
		`).MustRun(context.Background(), runtime.WithParam("values", map[string]int64{
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
			runtime.WithParam("start", int64(1)),
			runtime.WithParam("end", int64(4)),
		)

		So(string(out), ShouldEqual, `[1,2,3,4]`)

	})
}
