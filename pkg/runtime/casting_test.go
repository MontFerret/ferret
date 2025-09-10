package runtime_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func TestCasting(t *testing.T) {
	Convey("Casting Builder", t, func() {

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

		Convey("SafeCastBoolean", func() {
			Convey("Should cast Boolean to Boolean", func() {
				result := runtime.SafeCastBoolean(runtime.True, runtime.False)
				So(result, ShouldEqual, runtime.True)
			})

			Convey("Should return fallback for non-Boolean types", func() {
				result := runtime.SafeCastBoolean(runtime.NewInt(1), runtime.True)
				So(result, ShouldEqual, runtime.True)

				result = runtime.SafeCastBoolean(runtime.NewString("test"), runtime.False)
				So(result, ShouldEqual, runtime.False)
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

		Convey("SafeCastInt", func() {
			Convey("Should cast Int to Int", func() {
				result := runtime.SafeCastInt(runtime.NewInt(42), runtime.NewInt(0))
				So(result, ShouldEqual, runtime.NewInt(42))
			})

			Convey("Should return fallback for non-Int types", func() {
				fallback := runtime.NewInt(99)
				result := runtime.SafeCastInt(runtime.NewString("test"), fallback)
				So(result, ShouldEqual, fallback)
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

		Convey("SafeCastFloat", func() {
			Convey("Should cast Float to Float", func() {
				result := runtime.SafeCastFloat(runtime.NewFloat(3.14), runtime.NewFloat(0))
				So(result, ShouldEqual, runtime.NewFloat(3.14))
			})

			Convey("Should return fallback for non-Float types", func() {
				fallback := runtime.NewFloat(99.5)
				result := runtime.SafeCastFloat(runtime.NewString("test"), fallback)
				So(result, ShouldEqual, fallback)
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

		Convey("SafeCastString", func() {
			Convey("Should cast String to String", func() {
				result := runtime.SafeCastString(runtime.NewString("hello"), runtime.NewString("fallback"))
				So(result, ShouldEqual, runtime.NewString("hello"))
			})

			Convey("Should return fallback for non-String types", func() {
				fallback := runtime.NewString("fallback")
				result := runtime.SafeCastString(runtime.NewInt(42), fallback)
				So(result, ShouldEqual, fallback)
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

		Convey("SafeCastBinary", func() {
			Convey("Should cast Binary to Binary", func() {
				binary := runtime.NewBinary([]byte("test"))
				fallback := runtime.NewBinary([]byte("fallback"))
				result := runtime.SafeCastBinary(binary, fallback)
				So(result, ShouldEqual, binary)
			})

			Convey("Should return fallback for non-Binary types", func() {
				fallback := runtime.NewBinary([]byte("fallback"))
				result := runtime.SafeCastBinary(runtime.NewString("test"), fallback)
				So(result, ShouldEqual, fallback)
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

		Convey("SafeCastList", func() {
			Convey("Should cast Array to List", func() {
				arr := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
				fallback := runtime.NewArray(0)
				result := runtime.SafeCastList(arr, fallback)
				So(result, ShouldEqual, arr)
			})

			Convey("Should return fallback for non-List types", func() {
				fallback := runtime.NewArray(0)
				result := runtime.SafeCastList(runtime.NewString("test"), fallback)
				So(result, ShouldEqual, fallback)
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

		Convey("SafeCastMap", func() {
			Convey("Should cast Object to Map", func() {
				obj := runtime.NewObject()
				fallback := runtime.NewObject()
				result := runtime.SafeCastMap(obj, fallback)
				So(result, ShouldEqual, obj)
			})

			Convey("Should return fallback for non-Map types", func() {
				fallback := runtime.NewObject()
				result := runtime.SafeCastMap(runtime.NewString("test"), fallback)
				So(result, ShouldEqual, fallback)
			})
		})
	})
}
