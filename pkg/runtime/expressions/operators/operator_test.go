package operators_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/operators"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAdd(t *testing.T) {
	Convey("Add", t, func() {
		Convey("Integers", func() {
			arg1 := values.NewInt(1)
			arg2 := values.NewInt(2)

			So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(3))
		})

		Convey("Floats", func() {
			arg1 := values.NewInt(1)
			arg2 := values.NewFloat(2.1)

			So(operators.Add(arg1, arg2), ShouldEqual, values.NewFloat(3.1))
		})

		Convey("Number and String", func() {
			arg1 := values.NewInt(1)
			arg2 := values.NewString("2")

			So(operators.Add(arg1, arg2), ShouldEqual, values.NewFloat(3))

			arg3 := values.NewString("abc")

			So(operators.Add(arg1, arg3), ShouldEqual, values.NewInt(1))
		})

		Convey("Strings", func() {
			arg1 := values.NewString("1")
			arg2 := values.NewString("2")

			So(operators.Add(arg1, arg2), ShouldEqual, values.NewString("12"))

			arg3 := values.NewArrayWith(values.NewInt(1))

			So(operators.Add(arg1, arg3), ShouldEqual, values.NewString("1[1]"))
		})

		Convey("Boolean", func() {
			arg1 := values.NewInt(1)
			arg2 := values.True

			So(operators.Add(arg1, arg2), ShouldEqual, values.NewFloat(2))

			arg3 := values.False

			So(operators.Add(arg1, arg3), ShouldEqual, values.NewInt(1))
		})

		Convey("Array(0)", func() {
			arg1 := values.NewInt(1)
			arg2 := values.NewArray(0)

			So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(1))
		})

		Convey("Array(1)", func() {
			arg1 := values.NewInt(1)
			arg2 := values.NewArrayWith(values.NewInt(2))

			So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(3))
		})

		Convey("Array(2)", func() {
			arg1 := values.NewInt(1)
			arg2 := values.NewArrayWith(values.NewInt(2), values.NewFloat(2))

			So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(1))
		})

		Convey("Any", func() {
			arg1 := values.NewInt(1)
			args := []core.Value{
				values.NewArray(10),
				values.NewObject(),
				values.NewBinary([]byte("1")),
				values.NewCurrentDateTime(),
			}

			for _, argN := range args {
				So(operators.Add(arg1, argN), ShouldEqual, values.NewInt(1))
			}
		})
	})
}

func TestSubtract(t *testing.T) {
	Convey("Add", t, func() {
		Convey("Integers", func() {
			arg1 := values.NewInt(3)
			arg2 := values.NewInt(2)

			So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewInt(1))
		})

		Convey("Floats", func() {
			arg1 := values.NewInt(3)
			arg2 := values.NewFloat(2)

			So(operators.Subtract(arg1, arg2), ShouldEqual, 1)
		})

		Convey("Strings", func() {
			arg1 := values.NewInt(3)
			arg2 := values.NewString("2")

			So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewFloat(1))

			arg3 := values.NewString("abc")

			So(operators.Subtract(arg1, arg3), ShouldEqual, values.NewInt(3))
		})

		Convey("Boolean", func() {
			arg1 := values.NewInt(3)
			arg2 := values.True

			So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewFloat(2))

			arg3 := values.False

			So(operators.Subtract(arg1, arg3), ShouldEqual, values.NewInt(3))
		})

		Convey("Array(0)", func() {
			arg1 := values.NewInt(3)
			arg2 := values.NewArray(0)

			So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewInt(3))
		})

		Convey("Array(1)", func() {
			arg1 := values.NewInt(3)
			arg2 := values.NewArrayWith(values.NewInt(2))

			So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewInt(1))
		})

		Convey("Array(2)", func() {
			arg1 := values.NewInt(3)
			arg2 := values.NewArrayWith(values.NewInt(2), values.NewFloat(2))

			So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewInt(3))
		})

		Convey("Any", func() {
			arg1 := values.NewInt(3)
			args := []core.Value{
				values.NewArray(10),
				values.NewObject(),
				values.NewBinary([]byte("1")),
				values.NewCurrentDateTime(),
			}

			for _, argN := range args {
				So(operators.Subtract(arg1, argN), ShouldEqual, values.NewInt(3))
			}
		})
	})
}