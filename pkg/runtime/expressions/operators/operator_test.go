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
		Convey("Int", func() {
			Convey("1 + 2 = 3", func() {
				arg1 := values.NewInt(1)
				arg2 := values.NewInt(2)

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(3))
				So(operators.Add(arg2, arg1), ShouldEqual, values.NewInt(3))
			})
		})

		Convey("Float", func() {
			Convey("1.1 + 2.1 = 3.2", func() {
				arg1 := values.NewFloat(1.1)
				arg2 := values.NewFloat(2.1)

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewFloat(3.2))
				So(operators.Add(arg2, arg1), ShouldEqual, values.NewFloat(3.2))
			})
		})

		Convey("Float & Int", func() {
			Convey("1 + 2.1 = 3.1", func() {
				arg1 := values.NewInt(1)
				arg2 := values.NewFloat(2.1)

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewFloat(3.1))
			})

			Convey("1.1 + 2 = 3.1", func() {
				arg1 := values.NewFloat(1.1)
				arg2 := values.NewInt(2)

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewFloat(3.1))
			})
		})

		Convey("Int & String", func() {
			Convey("1 + 'a' = '1a'", func() {
				arg1 := values.NewInt(1)
				arg2 := values.NewString("a")

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewString("1a"))
			})

			Convey("'a' + 1 = 'a1'", func() {
				arg1 := values.NewString("a")
				arg2 := values.NewInt(1)

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewString("a1"))
			})
		})

		Convey("Float & String", func() {
			Convey("1.1 + 'a' = '1.1a'", func() {
				arg1 := values.NewFloat(1.1)
				arg2 := values.NewString("a")

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewString("1.1a"))
			})

			Convey("'a' + 1.1 = 'a1.1'", func() {
				arg1 := values.NewString("a")
				arg2 := values.NewFloat(1.1)

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewString("a1.1"))
			})
		})

		Convey("Strings", func() {
			Convey("'1' + '2' = '12'", func() {
				arg1 := values.NewString("1")
				arg2 := values.NewString("2")

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewString("12"))
			})

			Convey("'a' + 'b' = 'ab'", func() {
				arg1 := values.NewString("a")
				arg2 := values.NewString("b")

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewString("ab"))
			})
		})

		Convey("Boolean", func() {
			Convey("TRUE + TRUE = 2", func() {
				arg1 := values.True
				arg2 := values.True

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(2))
			})

			Convey("TRUE + FALSE = 1", func() {
				arg1 := values.True
				arg2 := values.False

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(1))
			})
		})

		Convey("Boolean & Int", func() {
			Convey("TRUE + 1 = 2", func() {
				arg1 := values.True
				arg2 := values.NewInt(1)

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(2))
			})

			Convey("1 + FALSE = 1", func() {
				arg1 := values.NewInt(1)
				arg2 := values.False

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(1))
			})
		})

		Convey("Boolean & Float", func() {
			Convey("TRUE + 1.2 = 2.2", func() {
				arg1 := values.True
				arg2 := values.NewFloat(1.2)

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewFloat(2.2))
			})

			Convey("1.2 + FALSE = 1.2", func() {
				arg1 := values.NewFloat(1.2)
				arg2 := values.False

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewFloat(1.2))
			})
		})

		Convey("None", func() {
			Convey("NONE + NONE = 0", func() {
				arg1 := values.None
				arg2 := values.None

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("None & Int", func() {
			Convey("NONE + 1 = 1", func() {
				arg1 := values.None
				arg2 := values.NewInt(1)

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(1))
			})

			Convey("1 + NONE = 1", func() {
				arg1 := values.NewInt(1)
				arg2 := values.None

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(1))
			})
		})

		Convey("None & Float", func() {
			Convey("NONE + 1.2 = 1.2", func() {
				arg1 := values.None
				arg2 := values.NewFloat(1.2)

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewFloat(1.2))
			})

			Convey("1.2 + NONE = 1.2", func() {
				arg1 := values.NewFloat(1.2)
				arg2 := values.None

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewFloat(1.2))
			})
		})

		Convey("Array", func() {
			Convey("[1] + [2] = 3", func() {
				arg1 := values.NewArrayWith(values.NewInt(1))
				arg2 := values.NewArrayWith(values.NewInt(2))

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(3))
			})

			Convey("[1] + [1, 1] = 3", func() {
				arg1 := values.NewArrayWith(values.NewInt(1))
				arg2 := values.NewArrayWith(values.NewInt(1), values.NewInt(1))

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(3))
			})

			Convey("[1, 2] + [1, 1] = 5", func() {
				arg1 := values.NewArrayWith(values.NewInt(1), values.NewInt(2))
				arg2 := values.NewArrayWith(values.NewInt(1), values.NewInt(1))

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(5))
			})
		})

		Convey("Datetime", func() {
			Convey("NOW() + NOW() = now+now", func() {
				arg1 := values.NewCurrentDateTime()
				arg2 := values.NewCurrentDateTime()

				expected := arg1.Time.Unix() + arg2.Time.Unix()

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(int(expected)))
			})
		})

		Convey("Datetime & Int", func() {
			Convey("NOW() + 1 = unix+1", func() {
				arg1 := values.NewCurrentDateTime()
				arg2 := values.NewArrayWith(values.NewInt(1))

				expected := arg1.Time.Unix() + 1

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(int(expected)))
			})

			Convey("1 + NOW() = 1+unix", func() {
				arg1 := values.NewArrayWith(values.NewInt(1))
				arg2 := values.NewCurrentDateTime()

				expected := arg2.Time.Unix() + 1

				So(operators.Add(arg1, arg2), ShouldEqual, values.NewInt(int(expected)))
			})
		})

		Convey("Other types", func() {
			arg1 := values.NewInt(1)
			args := []core.Value{
				values.NewObject(),
				values.NewBinary([]byte("1")),
			}

			for _, argN := range args {
				Convey(argN.Type().String(), func() {
					So(operators.Add(arg1, argN), ShouldEqual, values.NewInt(1))
				})
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