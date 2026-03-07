package runtime_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

func capturePanic(fn func()) (recovered any) {
	defer func() {
		recovered = recover()
	}()

	fn()

	return nil
}

func TestParams(t *testing.T) {
	Convey("Params", t, func() {
		Convey("NewParams and read methods", func() {
			params := runtime.NewParams()

			So(params, ShouldNotBeNil)
			So(len(params), ShouldEqual, 0)
			So(params.Has("missing"), ShouldBeFalse)

			value, exists := params.Get("missing")
			So(exists, ShouldBeFalse)
			So(value, ShouldEqual, runtime.None)

			So(params.GetOr("missing", runtime.NewString("fallback")), ShouldEqual, runtime.NewString("fallback"))
			So(params.GetOr("missing", nil), ShouldEqual, runtime.None)

			params.SetValue("nil-value", nil)

			value, exists = params.Get("nil-value")
			So(exists, ShouldBeTrue)
			So(value, ShouldEqual, runtime.None)
			So(params.Has("nil-value"), ShouldBeTrue)
			So(params.GetOr("nil-value", runtime.NewInt(99)), ShouldEqual, runtime.None)

			params.SetValue("present", runtime.NewString("ok"))
			So(params.MustGet("present"), ShouldEqual, runtime.NewString("ok"))

			recovered := capturePanic(func() {
				_ = params.MustGet("unknown")
			})

			So(recovered, ShouldNotBeNil)

			err, ok := recovered.(error)
			So(ok, ShouldBeTrue)
			So(errors.Is(err, runtime.ErrNotFound), ShouldBeTrue)
			So(strings.Contains(err.Error(), `param "unknown"`), ShouldBeTrue)
		})

		Convey("SetValue, SetAllValues and Delete", func() {
			params := runtime.NewParams()
			answer := runtime.NewInt(42)

			out := params.SetValue("answer", answer)
			So(out, ShouldResemble, params)
			So(params.Has("answer"), ShouldBeTrue)
			So(params.MustGet("answer"), ShouldEqual, answer)

			out = params.SetAllValues(map[string]runtime.Value{
				"name": runtime.NewString("ferret"),
				"none": nil,
			})

			So(out, ShouldResemble, params)
			So(params.MustGet("name"), ShouldEqual, runtime.NewString("ferret"))

			value, exists := params.Get("none")
			So(exists, ShouldBeTrue)
			So(value, ShouldEqual, runtime.None)

			params.Delete("name")
			So(params.Has("name"), ShouldBeFalse)
			So(params.Has("answer"), ShouldBeTrue)
		})

		Convey("Set and SetAll", func() {
			params := runtime.NewParams()

			err := params.Set("count", 7)
			So(err, ShouldBeNil)
			So(params.MustGet("count"), ShouldEqual, runtime.NewInt(7))

			err = params.Set("none", nil)
			So(err, ShouldBeNil)
			So(params.MustGet("none"), ShouldEqual, runtime.None)

			err = params.Set("invalid", make(chan int))
			So(err, ShouldNotBeNil)
			So(strings.Contains(err.Error(), `param "invalid"`), ShouldBeTrue)
			So(errors.Is(err, runtime.ErrInvalidType), ShouldBeTrue)
			So(params.Has("invalid"), ShouldBeFalse)

			err = params.SetAll(map[string]any{
				"title": "test",
				"size":  2,
			})
			So(err, ShouldBeNil)
			So(params.MustGet("title"), ShouldEqual, runtime.NewString("test"))
			So(params.MustGet("size"), ShouldEqual, runtime.NewInt(2))

			err = params.SetAll(map[string]any{
				"broken": make(chan int),
			})
			So(err, ShouldNotBeNil)
			So(strings.Contains(err.Error(), `param "broken"`), ShouldBeTrue)
			So(errors.Is(err, runtime.ErrInvalidType), ShouldBeTrue)
			So(params.Has("broken"), ShouldBeFalse)
		})

		Convey("MustSet", func() {
			params := runtime.NewParams()

			out := params.MustSet("ok", 1)
			So(out.Has("ok"), ShouldBeTrue)
			So(params.MustGet("ok"), ShouldEqual, runtime.NewInt(1))

			out.MustSet("chained", "value")
			So(params.MustGet("chained"), ShouldEqual, runtime.NewString("value"))

			recovered := capturePanic(func() {
				params.MustSet("bad", make(chan int))
			})

			So(recovered, ShouldNotBeNil)

			err, ok := recovered.(error)
			So(ok, ShouldBeTrue)
			So(errors.Is(err, runtime.ErrInvalidType), ShouldBeTrue)
			So(strings.Contains(err.Error(), `param "bad"`), ShouldBeTrue)
		})

		Convey("Clone", func() {
			var nilParams runtime.Params
			So(nilParams.Clone(), ShouldBeNil)

			ctx := context.Background()
			shared := runtime.NewArrayWith(runtime.NewInt(1))

			original := runtime.NewParams().
				SetValue("arr", shared).
				SetValue("name", runtime.NewString("ferret"))

			cloned := original.Clone()
			So(cloned, ShouldNotBeNil)
			So(cloned, ShouldResemble, original)

			cloned.SetValue("new", runtime.True)
			So(cloned.Has("new"), ShouldBeTrue)
			So(original.Has("new"), ShouldBeFalse)

			cloned.Delete("name")
			So(cloned.Has("name"), ShouldBeFalse)
			So(original.Has("name"), ShouldBeTrue)

			originalValue, exists := original.Get("arr")
			So(exists, ShouldBeTrue)

			clonedValue, exists := cloned.Get("arr")
			So(exists, ShouldBeTrue)
			So(clonedValue, ShouldEqual, originalValue)

			err := shared.Append(ctx, runtime.NewInt(2))
			So(err, ShouldBeNil)

			clonedArray, ok := clonedValue.(*runtime.Array)
			So(ok, ShouldBeTrue)

			size, err := clonedArray.Length(ctx)
			So(err, ShouldBeNil)
			So(size, ShouldEqual, runtime.NewInt(2))
		})
	})
}
