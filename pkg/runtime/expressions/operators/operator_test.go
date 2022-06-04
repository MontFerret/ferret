package operators_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/operators"
	"github.com/MontFerret/ferret/pkg/runtime/values"

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

		Convey("String", func() {
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
				argN := argN
				Convey(argN.Type().String(), func() {
					So(operators.Add(arg1, argN), ShouldEqual, values.NewInt(1))
				})
			}
		})
	})
}

func TestSubtract(t *testing.T) {
	Convey("Subtract", t, func() {
		Convey("Integer", func() {
			Convey("3 - 2 = 1", func() {
				arg1 := values.NewInt(3)
				arg2 := values.NewInt(2)

				So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewInt(1))
			})
		})

		Convey("Float", func() {
			Convey("3.4 - 2.2 = 1.2", func() {
				arg1 := values.NewFloat(3.40)
				arg2 := values.NewFloat(2.20)

				So(operators.Subtract(arg1, arg2), ShouldEqual, arg1-arg2)
			})
		})

		Convey("Integer & Float", func() {
			Convey("3.1 - 2 = 1.1", func() {
				arg1 := values.NewFloat(3.1)
				arg2 := values.NewInt(2)

				So(operators.Subtract(arg1, arg2), ShouldEqual, 1.1)
			})

			Convey("3 - 2.1 = 0.9", func() {
				arg1 := values.NewInt(3)
				arg2 := values.NewFloat(2.1)

				So(operators.Subtract(arg1, arg2), ShouldEqual, values.Float(arg1)-arg2)
			})
		})

		Convey("String", func() {
			Convey("'a' - 'b' = 0", func() {
				arg1 := values.NewString("a")
				arg2 := values.NewString("b")

				out := operators.Subtract(arg1, arg2)

				So(out, ShouldBeZeroValue)
			})

			Convey("'1' - 'b' = 1", func() {
				arg1 := values.NewString("1")
				arg2 := values.NewString("b")

				out := operators.Subtract(arg1, arg2)

				So(out, ShouldEqual, values.NewInt(1))
			})

			Convey("'a' - '1' = NaN", func() {
				arg1 := values.NewString("a")
				arg2 := values.NewString("1")

				out := operators.Subtract(arg1, arg2)

				So(out, ShouldEqual, values.NewInt(-1))
			})

			Convey("'2' - '1' = 1", func() {
				arg1 := values.NewString("2")
				arg2 := values.NewString("1")

				out := operators.Subtract(arg1, arg2)

				So(out, ShouldEqual, values.NewInt(1))
			})
		})

		Convey("String & Int", func() {
			Convey("1 - 'b' = 1", func() {
				arg1 := values.NewInt(1)
				arg2 := values.NewString("b")

				out := operators.Subtract(arg1, arg2)

				So(out, ShouldEqual, values.NewInt(1))
			})

			Convey("'2' - 1 = 1", func() {
				arg1 := values.NewString("2")
				arg2 := values.NewInt(1)

				out := operators.Subtract(arg1, arg2)

				So(out, ShouldEqual, values.NewInt(1))
			})

			Convey("'a' - '1' = -1", func() {
				arg1 := values.NewString("a")
				arg2 := values.NewString("1")

				out := operators.Subtract(arg1, arg2)

				So(out, ShouldEqual, values.NewInt(-1))
			})

			Convey("'2' - '1' = 1", func() {
				arg1 := values.NewString("2")
				arg2 := values.NewString("1")

				out := operators.Subtract(arg1, arg2)

				So(out, ShouldEqual, values.NewInt(1))
			})
		})

		Convey("Boolean", func() {
			Convey("TRUE - TRUE = 0", func() {
				arg1 := values.True
				arg2 := values.True

				So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewFloat(0))
			})

			Convey("TRUE - FALSE = 1", func() {
				arg1 := values.True
				arg2 := values.False

				So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewFloat(1))
			})
		})

		Convey("Boolean & Int", func() {
			Convey("2 - TRUE = 1", func() {
				arg1 := values.NewInt(2)
				arg2 := values.True

				So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewFloat(1))
			})

			Convey("1 - FALSE = 1", func() {
				arg1 := values.NewInt(1)
				arg2 := values.False

				So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewFloat(1))
			})
		})

		Convey("Boolean & String", func() {
			Convey("'2' - TRUE = 1", func() {
				arg1 := values.NewString("2")
				arg2 := values.True

				So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewFloat(1))
			})

			Convey("a - FALSE = 1", func() {
				arg1 := values.NewString("a")
				arg2 := values.False

				So(operators.Subtract(arg1, arg2), ShouldBeZeroValue)
			})
		})

		Convey("Array", func() {
			Convey("[2] - [1]", func() {
				arg1 := values.NewArrayWith(values.NewInt(2))
				arg2 := values.NewArrayWith(values.NewInt(1))

				So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewInt(1))
			})

			Convey("[2] - [1, 1]", func() {
				arg1 := values.NewArrayWith(values.NewInt(2))
				arg2 := values.NewArrayWith(values.NewInt(1), values.NewInt(1))

				So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("Array & Int", func() {
			Convey("2 - [1]", func() {
				arg1 := values.NewInt(2)
				arg2 := values.NewArrayWith(values.NewInt(1))

				So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewInt(1))
			})

			Convey("[1, '1'] - 2", func() {
				arg1 := values.NewArrayWith(values.NewInt(1), values.NewString("1"))
				arg2 := values.NewInt(2)

				So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("Array & Float", func() {
			Convey("2 - [1, 1.1]", func() {
				arg1 := values.NewInt(2)
				arg2 := values.NewArrayWith(values.NewInt(1), values.NewFloat(1.1))

				actual := operators.Subtract(arg1, arg2)
				expected := values.Float(arg1) - (operators.ToNumberOnly(arg2).(values.Float))

				So(actual, ShouldEqual, expected)
			})

			Convey("[1.1, '1.1'] - 2", func() {
				arg1 := values.NewArrayWith(values.NewFloat(1.1), values.NewString("1.1"))
				arg2 := values.NewInt(2)

				actual := operators.Subtract(arg1, arg2)
				expected := (operators.ToNumberOnly(arg1).(values.Float)) - values.Float(arg2)

				So(actual, ShouldEqual, expected)
			})
		})

		Convey("Array & String", func() {
			Convey("'2' - ['1']", func() {
				arg1 := values.NewInt(2)
				arg2 := values.NewArrayWith(values.NewInt(1))

				So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewInt(1))
			})

			Convey("[1, '1'] - 2", func() {
				arg1 := values.NewArrayWith(values.NewInt(1), values.NewString("1"))
				arg2 := values.NewInt(2)

				So(operators.Subtract(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("Datetime", func() {
			Convey("NOW() - NOW() = 0", func() {
				arg1 := values.NewCurrentDateTime()

				So(operators.Subtract(arg1, arg1), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("Any", func() {
			arg1 := values.NewInt(3)
			args := []core.Value{
				values.NewObject(),
				values.NewBinary([]byte("1")),
			}

			for _, argN := range args {
				So(operators.Subtract(arg1, argN), ShouldEqual, values.NewInt(3))
			}
		})
	})
}

func TestMultiply(t *testing.T) {
	Convey("Multiply", t, func() {
		Convey("Int", func() {
			Convey("2 * 2 = 4", func() {
				arg1 := values.NewInt(2)
				arg2 := values.NewInt(2)

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(4))
				So(operators.Multiply(arg2, arg1), ShouldEqual, values.NewInt(4))
			})
		})

		Convey("Float", func() {
			Convey("2 * 2.1 = 4.2", func() {
				arg1 := values.NewFloat(2.0)
				arg2 := values.NewFloat(2.1)

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewFloat(4.2))
				So(operators.Multiply(arg2, arg1), ShouldEqual, values.NewFloat(4.2))
			})
		})

		Convey("Float & Int", func() {
			Convey("2 * 2.2 = 4.4", func() {
				arg1 := values.NewInt(2)
				arg2 := values.NewFloat(2.2)

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewFloat(4.4))
			})

			Convey("2.2 * 2 = 4.4", func() {
				arg1 := values.NewFloat(2.2)
				arg2 := values.NewInt(2)

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewFloat(4.4))
			})
		})

		Convey("Int & String", func() {
			Convey("1 * 'a' = 0", func() {
				arg1 := values.NewInt(1)
				arg2 := values.NewString("a")

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(0))
			})

			Convey("'a' * 1 = 0", func() {
				arg1 := values.NewString("a")
				arg2 := values.NewInt(1)

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("Float & String", func() {
			Convey("1.1 * 'a' = 0", func() {
				arg1 := values.NewFloat(1.1)
				arg2 := values.NewString("a")

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(0))
			})

			Convey("'a' * 1.1 = 0", func() {
				arg1 := values.NewString("a")
				arg2 := values.NewFloat(1.1)

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("String", func() {
			Convey("'2' * '2' = 4", func() {
				arg1 := values.NewString("2")
				arg2 := values.NewString("2")

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(4))
			})

			Convey("'a' * 'b' = 0", func() {
				arg1 := values.NewString("a")
				arg2 := values.NewString("b")

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("Boolean", func() {
			Convey("TRUE * TRUE = 1", func() {
				arg1 := values.True
				arg2 := values.True

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(1))
			})

			Convey("TRUE * FALSE = 0", func() {
				arg1 := values.True
				arg2 := values.False

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("Boolean & Int", func() {
			Convey("TRUE * 2 = 2", func() {
				arg1 := values.True
				arg2 := values.NewInt(2)

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(2))
			})

			Convey("1 * FALSE = 0", func() {
				arg1 := values.NewInt(1)
				arg2 := values.False

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("Boolean & Float", func() {
			Convey("TRUE * 1.2 = 1.2", func() {
				arg1 := values.True
				arg2 := values.NewFloat(1.2)

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewFloat(1.2))
			})

			Convey("1.2 + FALSE = 0", func() {
				arg1 := values.NewFloat(1.2)
				arg2 := values.False

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("None", func() {
			Convey("NONE * NONE = 0", func() {
				arg1 := values.None
				arg2 := values.None

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("None & Int", func() {
			Convey("NONE * 1 = 0", func() {
				arg1 := values.None
				arg2 := values.NewInt(1)

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(0))
			})

			Convey("1 * NONE = 0", func() {
				arg1 := values.NewInt(1)
				arg2 := values.None

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("None & Float", func() {
			Convey("NONE * 1.2 = 0", func() {
				arg1 := values.None
				arg2 := values.NewFloat(1.2)

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewFloat(0))
			})

			Convey("1.2 + NONE = 0", func() {
				arg1 := values.NewFloat(1.2)
				arg2 := values.None

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewFloat(0))
			})
		})

		Convey("Array", func() {
			Convey("[2] * [2] = 4", func() {
				arg1 := values.NewArrayWith(values.NewInt(2))
				arg2 := values.NewArrayWith(values.NewInt(2))

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(4))
			})

			Convey("[2] * [2, 3] = 10", func() {
				arg1 := values.NewArrayWith(values.NewInt(2))
				arg2 := values.NewArrayWith(values.NewInt(2), values.NewInt(3))

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(10))
			})

			Convey("[1, 2] * [1, 1] = 6", func() {
				arg1 := values.NewArrayWith(values.NewInt(1), values.NewInt(2))
				arg2 := values.NewArrayWith(values.NewInt(1), values.NewInt(1))

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(6))
			})
		})

		Convey("Datetime", func() {
			Convey("NOW() * NOW() = now*now", func() {
				arg1 := values.NewCurrentDateTime()
				arg2 := values.NewCurrentDateTime()

				expected := arg1.Time.Unix() * arg2.Time.Unix()

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(int(expected)))
			})
		})

		Convey("Datetime & Int", func() {
			Convey("NOW() * 1 = unix", func() {
				arg1 := values.NewCurrentDateTime()
				arg2 := values.NewArrayWith(values.NewInt(1))

				expected := arg1.Time.Unix() * 1

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(int(expected)))
			})

			Convey("1 * NOW() = unix", func() {
				arg1 := values.NewArrayWith(values.NewInt(1))
				arg2 := values.NewCurrentDateTime()

				expected := arg2.Time.Unix() * 1

				So(operators.Multiply(arg1, arg2), ShouldEqual, values.NewInt(int(expected)))
			})
		})

		Convey("Other types", func() {
			arg1 := values.NewInt(1)
			args := []core.Value{
				values.NewObject(),
				values.NewBinary([]byte("1")),
			}

			for _, argN := range args {
				argN := argN
				Convey(argN.Type().String(), func() {
					So(operators.Multiply(arg1, argN), ShouldEqual, values.NewInt(0))
				})
			}
		})
	})
}

func TestDivide(t *testing.T) {
	Convey("Divide", t, func() {
		Convey("Int", func() {
			Convey("4 / 2 = 2", func() {
				arg1 := values.NewInt(4)
				arg2 := values.NewInt(2)

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewInt(2))
			})

			Convey("5 / 2 = 2.5", func() {
				arg1 := values.NewInt(5)
				arg2 := values.NewInt(2)

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewFloat(2.5))
			})
		})

		Convey("Float", func() {
			Convey("4.4 / 2.2 = 2", func() {
				arg1 := values.NewFloat(4.4)
				arg2 := values.NewFloat(2.2)

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewFloat(2))
			})
		})

		Convey("Float & Int", func() {
			Convey("5 / 2.5 = 2", func() {
				arg1 := values.NewInt(5)
				arg2 := values.NewFloat(2.5)

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewFloat(2))
			})

			Convey("4.4 / 2 = 2.2", func() {
				arg1 := values.NewFloat(4.4)
				arg2 := values.NewInt(2)

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewFloat(2.2))
			})
		})

		Convey("Int & String", func() {
			Convey("1 / 'a' = PANIC", func() {
				arg1 := values.NewInt(1)
				arg2 := values.NewString("a")

				So(func() {
					operators.Divide(arg1, arg2)
				}, ShouldPanic)
			})

			Convey("'a' / 1 = 0", func() {
				arg1 := values.NewString("a")
				arg2 := values.NewInt(1)

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("Float & String", func() {
			Convey("1.1 / 'a' = PANIC", func() {
				arg1 := values.NewFloat(1.1)
				arg2 := values.NewString("a")

				So(func() {
					operators.Divide(arg1, arg2)
				}, ShouldPanic)
			})

			Convey("'a' / 1.1 = 0", func() {
				arg1 := values.NewString("a")
				arg2 := values.NewFloat(1.1)

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewInt(0))
			})
		})

		Convey("String", func() {
			Convey("'2' / '2' = 1", func() {
				arg1 := values.NewString("2")
				arg2 := values.NewString("2")

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewInt(1))
			})

			Convey("'a' / 'b' = PANIC", func() {
				arg1 := values.NewString("a")
				arg2 := values.NewString("b")

				So(func() {
					operators.Divide(arg1, arg2)
				}, ShouldPanic)
			})
		})

		Convey("Boolean", func() {
			Convey("TRUE / TRUE = 1", func() {
				arg1 := values.True
				arg2 := values.True

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewInt(1))
			})

			Convey("TRUE / FALSE = PANIC", func() {
				arg1 := values.True
				arg2 := values.False

				So(func() {
					operators.Divide(arg1, arg2)
				}, ShouldPanic)
			})
		})

		Convey("Boolean & Int", func() {
			Convey("TRUE / 2 = 0.5", func() {
				arg1 := values.True
				arg2 := values.NewInt(2)

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewFloat(0.5))
			})

			Convey("1 / FALSE = PANIC", func() {
				arg1 := values.NewInt(1)
				arg2 := values.False

				So(func() {
					operators.Divide(arg1, arg2)
				}, ShouldPanic)
			})
		})

		Convey("Boolean & Float", func() {
			Convey("TRUE / 0.5 = 2", func() {
				arg1 := values.True
				arg2 := values.NewFloat(0.5)

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewFloat(2))
			})

			Convey("1.2 / FALSE = PANIC", func() {
				arg1 := values.NewFloat(1.2)
				arg2 := values.False

				So(func() {
					operators.Divide(arg1, arg2)
				}, ShouldPanic)
			})
		})

		Convey("None", func() {
			Convey("NONE / NONE = PANIC", func() {
				arg1 := values.None
				arg2 := values.None

				So(func() {
					operators.Divide(arg1, arg2)
				}, ShouldPanic)
			})
		})

		Convey("None & Int", func() {
			Convey("NONE / 1 = 0", func() {
				arg1 := values.None
				arg2 := values.NewInt(1)

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewInt(0))
			})

			Convey("1 / NONE = PANIC", func() {
				arg1 := values.NewInt(1)
				arg2 := values.None

				So(func() {
					operators.Divide(arg1, arg2)
				}, ShouldPanic)
			})
		})

		Convey("None & Float", func() {
			Convey("NONE / 1.2 = 0", func() {
				arg1 := values.None
				arg2 := values.NewFloat(1.2)

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewFloat(0))
			})

			Convey("1.2 / NONE = PANIC", func() {
				arg1 := values.NewFloat(1.2)
				arg2 := values.None

				So(func() {
					operators.Divide(arg1, arg2)
				}, ShouldPanic)
			})
		})

		Convey("Array", func() {
			Convey("[2] / [2] = 1", func() {
				arg1 := values.NewArrayWith(values.NewInt(2))
				arg2 := values.NewArrayWith(values.NewInt(2))

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewInt(1))
			})

			Convey("[4, 2] / [2] = 3", func() {
				arg1 := values.NewArrayWith(values.NewInt(4), values.NewInt(2))
				arg2 := values.NewArrayWith(values.NewInt(2))

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewInt(3))
			})

			Convey("[4, 2] / [1, 1] = 3", func() {
				arg1 := values.NewArrayWith(values.NewInt(4), values.NewInt(2))
				arg2 := values.NewArrayWith(values.NewInt(1), values.NewInt(1))

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewInt(3))
			})
		})

		Convey("Datetime", func() {
			Convey("NOW() / NOW() = 1", func() {
				arg1 := values.NewCurrentDateTime()
				arg2 := values.NewCurrentDateTime()

				expected := arg1.Time.Unix() / arg2.Time.Unix()

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewInt(int(expected)))
			})
		})

		Convey("Datetime & Int", func() {
			Convey("NOW() / 1 = unix", func() {
				arg1 := values.NewCurrentDateTime()
				arg2 := values.NewArrayWith(values.NewInt(1))

				expected := arg1.Time.Unix() / 1

				So(operators.Divide(arg1, arg2), ShouldEqual, values.NewInt(int(expected)))
			})

			SkipConvey("1 / NOW() = unix", func() {
				arg1 := values.NewArrayWith(values.NewInt(1))
				arg2 := values.NewCurrentDateTime()

				expected := arg2.Time.Unix() / 1

				So(operators.Divide(arg1, arg2), ShouldEqual, expected)
			})
		})

		Convey("Other types", func() {
			arg1 := values.NewInt(1)
			args := []core.Value{
				values.NewObject(),
				values.NewBinary([]byte("1")),
			}

			for _, argN := range args {
				argN := argN
				Convey(argN.Type().String(), func() {
					So(func() {
						operators.Divide(arg1, argN)
					}, ShouldPanic)
				})
			}
		})
	})
}

func TestPositive(t *testing.T) {
	Convey("converts value to positive", t, func() {
		Convey("when int", func() {
			Convey("postive", func() {
				So(operators.Positive(values.NewInt(1), nil), ShouldEqual, 1)
			})
			Convey("negative", func() {
				So(operators.Positive(values.NewInt(-1), nil), ShouldEqual, -1)
			})
		})

		Convey("when float", func() {
			Convey("postive", func() {
				So(operators.Positive(values.NewFloat(1.0), nil), ShouldEqual, 1)
			})
			Convey("negative", func() {
				So(operators.Positive(values.NewFloat(-1.0), nil), ShouldEqual, -1.0)
			})
		})

		Convey("non numeric", func() {
			Convey("string", func() {
				So(operators.Positive(values.String("1"), nil), ShouldEqual, "1")
			})
		})
	})
}

func TestNegative(t *testing.T) {
	Convey("converts value to negative", t, func() {
		Convey("when int", func() {
			Convey("postive", func() {
				So(operators.Negative(values.NewInt(1), nil), ShouldEqual, -1)
			})
			Convey("negative", func() {
				So(operators.Negative(values.NewInt(-1), nil), ShouldEqual, 1)
			})
		})

		Convey("when float", func() {
			Convey("postive", func() {
				So(operators.Negative(values.NewFloat(1.0), nil), ShouldEqual, -1.0)
			})
			Convey("negative", func() {
				So(operators.Negative(values.NewFloat(-1.0), nil), ShouldEqual, 1.0)
			})
		})

		Convey("non numeric", func() {
			Convey("string", func() {
				So(operators.Negative(values.String("1"), nil), ShouldEqual, "1")
			})
		})
	})
}

func TestDecrement(t *testing.T) {
	Convey("decrement values", t, func() {
		Convey("decrement int", func() {
			Convey("positives", func() {
				So(operators.Decrement(values.NewInt(1), nil), ShouldEqual, 0)
			})
			Convey("negatives", func() {
				So(operators.Decrement(values.NewInt(-1), nil), ShouldEqual, -2)
			})
		})

		Convey("decrement float", func() {
			Convey("positives", func() {
				So(operators.Decrement(values.NewFloat(1.1), nil), ShouldEqual, 0.10000000000000009)
			})
			Convey("negatives", func() {
				So(operators.Decrement(values.NewFloat(-1.1), nil), ShouldEqual, -2.1)
			})
		})

		Convey("other values", func() {
			So(operators.Decrement(values.None, nil), ShouldEqual, -1)
		})
	})
}

func TestIncrement(t *testing.T) {
	Convey("increment values", t, func() {
		Convey("increment int", func() {
			Convey("positives", func() {
				So(operators.Increment(values.NewInt(1), nil), ShouldEqual, 2)
			})
			Convey("negatives", func() {
				So(operators.Increment(values.NewInt(-1), nil), ShouldEqual, 0)
			})
		})

		Convey("Increment float", func() {
			Convey("positives", func() {
				So(operators.Increment(values.NewFloat(1.1), nil), ShouldEqual, 2.1)
			})
			Convey("negatives", func() {
				So(operators.Increment(values.NewFloat(-1.1), nil), ShouldEqual, -0.10000000000000009)
			})
		})

		Convey("other values", func() {
			So(operators.Increment(values.None, nil), ShouldEqual, 1)
		})
	})
}

func TestEqual(t *testing.T) {
	Convey("Equality in values", t, func() {
		Convey("1 == 1", func() {
			So(operators.Equal(values.NewInt(1), values.NewInt(1)), ShouldEqual, values.True)
		})

		Convey("1 != 2", func() {
			So(operators.Equal(values.NewInt(1), values.NewInt(2)), ShouldEqual, values.False)
		})

		Convey("'hello' == 'hello", func() {
			So(operators.Equal(values.String("hello"), values.String("hello")), ShouldEqual, values.True)
		})

		Convey("'foo' != 'bar", func() {
			So(operators.Equal(values.String("foo"), values.String("bar")), ShouldEqual, values.False)
		})
	})
}

func TestNotEqual(t *testing.T) {
	Convey("Inequality in values", t, func() {
		Convey("1 == 1", func() {
			So(operators.NotEqual(values.NewInt(1), values.NewInt(1)), ShouldEqual, values.False)
		})

		Convey("1 != 2", func() {
			So(operators.NotEqual(values.NewInt(1), values.NewInt(2)), ShouldEqual, values.True)
		})

		Convey("'hello' == 'hello", func() {
			So(operators.NotEqual(values.String("hello"), values.String("hello")), ShouldEqual, values.False)
		})

		Convey("'foo' != 'bar", func() {
			So(operators.NotEqual(values.String("foo"), values.String("bar")), ShouldEqual, values.True)
		})
	})
}

func TestModulus(t *testing.T) {
	Convey("Modulus", t, func() {
		Convey("5 % 2", func() {
			So(operators.Modulus(values.NewInt(5), values.NewInt(2)), ShouldEqual, 1)
		})

		Convey("5.0 % 2", func() {
			So(operators.Modulus(values.NewFloat(5.0), values.NewInt(2)), ShouldEqual, 1)
		})

		Convey("5 % 2.0", func() {
			So(operators.Modulus(values.NewInt(5), values.NewFloat(2.0)), ShouldEqual, 1)
		})

		Convey("5.1 % 3.2", func() {
			So(operators.Modulus(values.NewFloat(5.1), values.NewFloat(3.2)), ShouldEqual, 2)
		})

		Convey("Non int or float", func() {
			So(operators.Modulus(values.String("foo"), values.NewFloat(3.2)), ShouldEqual, 0)
		})
	})
}

func TestGreater(t *testing.T) {
	Convey("Greater than value", t, func() {
		Convey("5 > 2", func() {
			So(operators.Greater(values.NewInt(5), values.NewInt(2)), ShouldEqual, values.True)
		})

		Convey("2 > 5", func() {
			So(operators.Greater(values.NewInt(2), values.NewInt(5)), ShouldEqual, values.False)
		})
	})
}

func TestGreaterOrEqual(t *testing.T) {
	Convey("Greater or equal than value", t, func() {
		Convey("5 >= 5", func() {
			So(operators.GreaterOrEqual(values.NewInt(5), values.NewInt(5)), ShouldEqual, values.True)
		})

		Convey("2 >= 5", func() {
			So(operators.GreaterOrEqual(values.NewInt(2), values.NewInt(5)), ShouldEqual, values.False)
		})
	})
}

func TestLess(t *testing.T) {
	Convey("Less than value", t, func() {
		Convey("1 < 5", func() {
			So(operators.Less(values.NewInt(1), values.NewInt(5)), ShouldEqual, values.True)
		})
		Convey("5 < 2", func() {
			So(operators.Less(values.NewInt(5), values.NewInt(2)), ShouldEqual, values.False)
		})
	})
}

func TestLessOrEqual(t *testing.T) {
	Convey("Less than value", t, func() {
		Convey("1 < 5", func() {
			So(operators.LessOrEqual(values.NewInt(1), values.NewInt(5)), ShouldEqual, values.True)
		})
		Convey("5 < 2", func() {
			So(operators.LessOrEqual(values.NewInt(5), values.NewInt(2)), ShouldEqual, values.False)
		})
		Convey("5 <= 5", func() {
			So(operators.LessOrEqual(values.NewInt(5), values.NewInt(5)), ShouldEqual, values.True)
		})
	})
}

func TestNot(t *testing.T) {
	Convey("Invert truthiness", t, func() {
		Convey("true turns false", func() {
			So(operators.Not(values.NewBoolean(true), nil), ShouldEqual, values.False)
		})
		Convey("false turns true", func() {
			So(operators.Not(values.NewBoolean(false), nil), ShouldEqual, values.True)
		})
		Convey("'' turns true", func() {
			So(operators.Not(values.NewString(""), nil), ShouldEqual, values.True)
		})
		Convey("'foo' turns false", func() {
			So(operators.Not(values.NewString("foo"), nil), ShouldEqual, values.False)
		})
		Convey("1 turns false", func() {
			So(operators.Not(values.NewInt(1), nil), ShouldEqual, values.False)
		})
		Convey("0 turns true", func() {
			So(operators.Not(values.NewInt(0), nil), ShouldEqual, values.True)
		})
		Convey("1.0 turns false", func() {
			So(operators.Not(values.NewFloat(1), nil), ShouldEqual, values.False)
		})
		Convey("0.0 turns true", func() {
			So(operators.Not(values.NewFloat(0.0), nil), ShouldEqual, values.True)
		})
		Convey("current turns false", func() {
			So(operators.Not(values.NewDateTime(values.NewCurrentDateTime().Time), nil), ShouldEqual, values.False)
		})
		Convey("zerotime turns true", func() {
			So(operators.Not(values.NewDateTime(values.ZeroDateTime.Time), nil), ShouldEqual, values.True)
		})
	})
}
