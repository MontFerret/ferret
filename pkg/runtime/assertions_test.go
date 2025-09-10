package runtime_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func TestAssertions(t *testing.T) {
	Convey("Assertion Builder", t, func() {

		Convey("AssertString", func() {
			Convey("Should pass for String type", func() {
				str := runtime.NewString("test")
				err := runtime.AssertString(str)
				So(err, ShouldBeNil)
			})

			Convey("Should fail for non-String types", func() {
				err := runtime.AssertString(runtime.NewInt(42))
				So(err, ShouldNotBeNil)

				err = runtime.AssertString(runtime.NewFloat(3.14))
				So(err, ShouldNotBeNil)

				err = runtime.AssertString(runtime.True)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("AssertInt", func() {
			Convey("Should pass for Int type", func() {
				val := runtime.NewInt(42)
				err := runtime.AssertInt(val)
				So(err, ShouldBeNil)
			})

			Convey("Should fail for non-Int types", func() {
				err := runtime.AssertInt(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				err = runtime.AssertInt(runtime.NewFloat(3.14))
				So(err, ShouldNotBeNil)

				err = runtime.AssertInt(runtime.True)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("AssertFloat", func() {
			Convey("Should pass for Float type", func() {
				val := runtime.NewFloat(3.14)
				err := runtime.AssertFloat(val)
				So(err, ShouldBeNil)
			})

			Convey("Should fail for non-Float types", func() {
				err := runtime.AssertFloat(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				err = runtime.AssertFloat(runtime.NewInt(42))
				So(err, ShouldNotBeNil)

				err = runtime.AssertFloat(runtime.True)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("AssertNumber", func() {
			Convey("Should pass for numeric types", func() {
				err := runtime.AssertNumber(runtime.NewInt(42))
				So(err, ShouldBeNil)

				err = runtime.AssertNumber(runtime.NewFloat(3.14))
				So(err, ShouldBeNil)
			})

			Convey("Should fail for non-numeric types", func() {
				err := runtime.AssertNumber(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				err = runtime.AssertNumber(runtime.True)
				So(err, ShouldNotBeNil)

				err = runtime.AssertNumber(runtime.NewObject())
				So(err, ShouldNotBeNil)
			})
		})

		Convey("AssertBoolean", func() {
			Convey("Should pass for Boolean type", func() {
				err := runtime.AssertBoolean(runtime.True)
				So(err, ShouldBeNil)

				err = runtime.AssertBoolean(runtime.False)
				So(err, ShouldBeNil)
			})

			Convey("Should fail for non-Boolean types", func() {
				err := runtime.AssertBoolean(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				err = runtime.AssertBoolean(runtime.NewInt(42))
				So(err, ShouldNotBeNil)

				err = runtime.AssertBoolean(runtime.NewFloat(3.14))
				So(err, ShouldNotBeNil)
			})
		})

		Convey("AssertCollection", func() {
			Convey("Should pass for Collection types", func() {
				arr := runtime.NewArray(5)
				err := runtime.AssertCollection(arr)
				So(err, ShouldBeNil)

				obj := runtime.NewObject()
				err = runtime.AssertCollection(obj)
				So(err, ShouldBeNil)
			})

			Convey("Should fail for non-Collection types", func() {
				err := runtime.AssertCollection(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				err = runtime.AssertCollection(runtime.NewInt(42))
				So(err, ShouldNotBeNil)

				err = runtime.AssertCollection(runtime.True)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("AssertArray", func() {
			Convey("Should pass for Array type", func() {
				arr := runtime.NewArray(5)
				err := runtime.AssertArray(arr)
				So(err, ShouldBeNil)
			})

			Convey("Should fail for non-Array types", func() {
				err := runtime.AssertArray(runtime.NewObject())
				So(err, ShouldNotBeNil)

				err = runtime.AssertArray(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				err = runtime.AssertArray(runtime.NewInt(42))
				So(err, ShouldNotBeNil)
			})
		})

		Convey("AssertList", func() {
			Convey("Should pass for List type", func() {
				arr := runtime.NewArray(5)
				err := runtime.AssertList(arr)
				So(err, ShouldBeNil)
			})

			Convey("Should fail for non-List types", func() {
				err := runtime.AssertList(runtime.NewObject())
				So(err, ShouldNotBeNil)

				err = runtime.AssertList(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				err = runtime.AssertList(runtime.NewInt(42))
				So(err, ShouldNotBeNil)
			})
		})

		Convey("AssertObject", func() {
			Convey("Should pass for Object type", func() {
				obj := runtime.NewObject()
				err := runtime.AssertObject(obj)
				So(err, ShouldBeNil)
			})

			Convey("Should fail for non-Object types", func() {
				err := runtime.AssertObject(runtime.NewArray(5))
				So(err, ShouldNotBeNil)

				err = runtime.AssertObject(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				err = runtime.AssertObject(runtime.NewInt(42))
				So(err, ShouldNotBeNil)
			})
		})

		Convey("AssertMap", func() {
			Convey("Should pass for Map type", func() {
				obj := runtime.NewObject()
				err := runtime.AssertMap(obj)
				So(err, ShouldBeNil)
			})

			Convey("Should fail for non-Map types", func() {
				err := runtime.AssertMap(runtime.NewArray(5))
				So(err, ShouldNotBeNil)

				err = runtime.AssertMap(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				err = runtime.AssertMap(runtime.NewInt(42))
				So(err, ShouldNotBeNil)
			})
		})

		Convey("AssertBinary", func() {
			Convey("Should pass for Binary type", func() {
				bin := runtime.NewBinary([]byte("test"))
				err := runtime.AssertBinary(bin)
				So(err, ShouldBeNil)
			})

			Convey("Should fail for non-Binary types", func() {
				err := runtime.AssertBinary(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				err = runtime.AssertBinary(runtime.NewInt(42))
				So(err, ShouldNotBeNil)

				err = runtime.AssertBinary(runtime.NewObject())
				So(err, ShouldNotBeNil)
			})
		})

		Convey("AssertDateTime", func() {
			Convey("Should pass for DateTime type", func() {
				dt := runtime.NewCurrentDateTime()
				err := runtime.AssertDateTime(dt)
				So(err, ShouldBeNil)
			})

			Convey("Should fail for non-DateTime types", func() {
				err := runtime.AssertDateTime(runtime.NewString("test"))
				So(err, ShouldNotBeNil)

				err = runtime.AssertDateTime(runtime.NewInt(42))
				So(err, ShouldNotBeNil)

				err = runtime.AssertDateTime(runtime.NewObject())
				So(err, ShouldNotBeNil)
			})
		})

		Convey("AssertMeasurable", func() {
			Convey("Should pass for Measurable types", func() {
				// Array is Measurable
				arr := runtime.NewArray(5)
				err := runtime.AssertMeasurable(arr)
				So(err, ShouldBeNil)

				// String is Measurable
				str := runtime.NewString("test")
				err = runtime.AssertMeasurable(str)
				So(err, ShouldBeNil)

				// Binary is Measurable
				bin := runtime.NewBinary([]byte("test"))
				err = runtime.AssertMeasurable(bin)
				So(err, ShouldBeNil)
			})

			Convey("Should fail for non-Measurable types", func() {
				err := runtime.AssertMeasurable(runtime.NewInt(42))
				So(err, ShouldNotBeNil)

				err = runtime.AssertMeasurable(runtime.NewFloat(3.14))
				So(err, ShouldNotBeNil)

				err = runtime.AssertMeasurable(runtime.True)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
