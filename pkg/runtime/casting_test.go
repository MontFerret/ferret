package runtime_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestCasting(t *testing.T) {
	Convey("Casting Builder", t, func() {
		Convey("Cast", func() {
			Convey("Should cast concrete and interface types", func() {
				val, err := runtime.Cast[runtime.Int](runtime.NewInt(42))
				So(err, ShouldBeNil)
				So(val, ShouldEqual, runtime.NewInt(42))

				arr := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
				list, err := runtime.Cast[runtime.List](arr)
				So(err, ShouldBeNil)
				So(list, ShouldEqual, arr)
			})

			Convey("Should return error with expected interface type", func() {
				_, err := runtime.Cast[runtime.List](runtime.True)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "expected List")
				So(err.Error(), ShouldContainSubstring, "got Boolean")
			})
		})

		Convey("CastArg", func() {
			Convey("Should cast argument to interface type", func() {
				arr := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
				list, err := runtime.CastArg[runtime.List](arr, 0)
				So(err, ShouldBeNil)
				So(list, ShouldEqual, arr)
			})

			Convey("Should return error with expected type and position", func() {
				_, err := runtime.CastArg[runtime.List](runtime.True, 0)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "expected List")
				So(err.Error(), ShouldContainSubstring, "position 1")
			})
		})

		Convey("CastOr", func() {
			Convey("Should cast T to T", func() {
				input := runtime.NewInt(1)
				fallback := runtime.NewInt(0)
				actual := runtime.CastOr(input, fallback)
				So(actual, ShouldEqual, input)
			})

			Convey("Should return fallback for non-T types", func() {
				input := runtime.NewString("1")
				fallback := runtime.NewInt(0)
				actual := runtime.CastOr(input, fallback)
				So(actual, ShouldEqual, fallback)
			})
		})

		Convey("CastBoolean", func() {
			Convey("Should cast Boolean to Boolean", func() {
				result, err := runtime.CastBoolean(runtime.True)
				So(err, ShouldBeNil)
				So(result, ShouldEqual, runtime.True)

				result, err = runtime.CastBoolean(runtime.False)
				So(err, ShouldBeNil)
				So(result, ShouldEqual, runtime.False)
			})

			Convey("Should return error for non-Boolean types", func() {
				_, err := runtime.CastBoolean(runtime.NewInt(1))
				So(err, ShouldNotBeNil)

				_, err = runtime.CastBoolean(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				_, err = runtime.CastBoolean(runtime.NewFloat(3.14))
				So(err, ShouldNotBeNil)
			})
		})

		Convey("CastInt", func() {
			Convey("Should cast Int to Int", func() {
				result, err := runtime.CastInt(runtime.NewInt(42))
				So(err, ShouldBeNil)
				So(result, ShouldEqual, runtime.NewInt(42))
			})

			Convey("Should return error for non-Int types", func() {
				_, err := runtime.CastInt(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				_, err = runtime.CastInt(runtime.NewFloat(3.14))
				So(err, ShouldNotBeNil)

				_, err = runtime.CastInt(runtime.True)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("CastFloat", func() {
			Convey("Should cast Float to Float", func() {
				result, err := runtime.CastFloat(runtime.NewFloat(3.14))
				So(err, ShouldBeNil)
				So(result, ShouldEqual, runtime.NewFloat(3.14))
			})

			Convey("Should return error for non-Float types", func() {
				_, err := runtime.CastFloat(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				_, err = runtime.CastFloat(runtime.NewInt(42))
				So(err, ShouldNotBeNil)

				_, err = runtime.CastFloat(runtime.True)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("CastString", func() {
			Convey("Should cast String to String", func() {
				result, err := runtime.CastString(runtime.NewString("hello"))
				So(err, ShouldBeNil)
				So(result, ShouldEqual, runtime.NewString("hello"))
			})

			Convey("Should return error for non-String types", func() {
				_, err := runtime.CastString(runtime.NewInt(42))
				So(err, ShouldNotBeNil)

				_, err = runtime.CastString(runtime.NewFloat(3.14))
				So(err, ShouldNotBeNil)

				_, err = runtime.CastString(runtime.True)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("CastBinary", func() {
			Convey("Should cast Binary to Binary", func() {
				binary := runtime.NewBinary([]byte("test"))
				result, err := runtime.CastBinary(binary)
				So(err, ShouldBeNil)
				So(result, ShouldEqual, binary)
			})

			Convey("Should return error for non-Binary types", func() {
				_, err := runtime.CastBinary(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				_, err = runtime.CastBinary(runtime.NewInt(42))
				So(err, ShouldNotBeNil)
			})
		})

		Convey("CastList", func() {
			Convey("Should cast Array to List", func() {
				arr := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
				result, err := runtime.CastList(arr)
				So(err, ShouldBeNil)
				So(result, ShouldEqual, arr)
			})

			Convey("Should return error for non-List types", func() {
				_, err := runtime.CastList(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				_, err = runtime.CastList(runtime.NewInt(42))
				So(err, ShouldNotBeNil)
			})
		})

		Convey("CastMap", func() {
			Convey("Should cast Object to Map", func() {
				obj := runtime.NewObject()
				result, err := runtime.CastMap(obj)
				So(err, ShouldBeNil)
				So(result, ShouldEqual, obj)
			})

			Convey("Should return error for non-Map types", func() {
				_, err := runtime.CastMap(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				_, err = runtime.CastMap(runtime.NewInt(42))
				So(err, ShouldNotBeNil)
			})
		})
	})
}
