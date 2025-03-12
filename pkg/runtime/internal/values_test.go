package internal_test

import (
	"context"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type CustomValue struct {
	properties map[Value]Value
}

func (t *CustomValue) MarshalJSON() ([]byte, error) {
	return nil, ErrNotImplemented
}

func (t *CustomValue) String() string {
	return ""
}

func (t *CustomValue) Unwrap() interface{} {
	return t
}

func (t *CustomValue) Hash() uint64 {
	return 0
}

func (t *CustomValue) Copy() Value {
	return None
}

func TestHelpers(t *testing.T) {
	Convey("Helpers", t, func() {
		Convey("Parse", func() {
			Convey("It should parse values", func() {
				inputs := []struct {
					Parsed Value
					Raw    interface{}
				}{
					{Parsed: NewInt(1), Raw: int(1)},
					{Parsed: NewInt(1), Raw: int8(1)},
					{Parsed: NewInt(1), Raw: int16(1)},
					{Parsed: NewInt(1), Raw: int32(1)},
					{Parsed: NewInt(1), Raw: int64(1)},
				}

				for _, input := range inputs {
					out := internal.Parse(input.Raw)

					expectedType := reflect.TypeOf(input.Parsed)
					actualType := reflect.TypeOf(out)

					So(actualType, ShouldEqual, expectedType)
					So(out.Unwrap(), ShouldEqual, input.Parsed.Unwrap())
				}
			})
		})

		Convey("ToBoolean", func() {
			Convey("Should convert values", func() {
				inputs := [][]Value{
					{
						None,
						False,
					},
					{
						True,
						True,
					},
					{
						False,
						False,
					},
					{
						NewInt(1),
						True,
					},
					{
						NewInt(0),
						False,
					},
					{
						NewFloat(1),
						True,
					},
					{
						NewFloat(0),
						False,
					},
					{
						NewString("Foo"),
						True,
					},
					{
						EmptyString,
						False,
					},
					{
						NewCurrentDateTime(),
						True,
					},
					{
						internal.NewArray(1),
						True,
					},
					{
						internal.NewObject(),
						True,
					},
					{
						NewBinary([]byte("")),
						True,
					},
				}

				for _, pair := range inputs {
					actual := internal.ToBoolean(pair[0])
					expected := pair[1]

					So(actual, ShouldEqual, expected)
				}
			})
		})

		Convey("ToFloat", func() {
			Convey("Should convert Int", func() {
				input := NewInt(100)
				output := internal.ToFloat(input)

				So(output, ShouldEqual, NewFloat(100))
			})

			Convey("Should convert Float", func() {
				input := NewFloat(100)
				output := internal.ToFloat(input)

				So(output, ShouldEqual, NewFloat(100))
			})

			Convey("Should convert String", func() {
				input := NewString("100.1")
				output := internal.ToFloat(input)

				So(output, ShouldEqual, NewFloat(100.1))

				output2 := internal.ToFloat(NewString("foobar"))
				So(output2, ShouldEqual, ZeroFloat)
			})

			Convey("Should convert Boolean", func() {
				So(internal.ToFloat(True), ShouldEqual, NewFloat(1))
				So(internal.ToFloat(False), ShouldEqual, NewFloat(0))
			})

			Convey("Should convert Array with single item", func() {
				So(internal.ToFloat(internal.NewArrayWith(NewFloat(1))), ShouldEqual, NewFloat(1))
			})

			Convey("Should convert Array with multiple items", func() {
				arg := internal.NewArrayWith(NewFloat(1), NewFloat(1))

				So(internal.ToFloat(arg), ShouldEqual, NewFloat(2))
			})

			Convey("Should convert DateTime", func() {
				dt := NewCurrentDateTime()
				ts := dt.Time.Unix()

				So(internal.ToFloat(dt), ShouldEqual, NewFloat(float64(ts)))
			})

			Convey("Should NOT convert other types", func() {
				inputs := []Value{
					internal.NewObject(),
					NewBinary([]byte("")),
				}

				for _, input := range inputs {
					So(internal.ToFloat(input), ShouldEqual, ZeroFloat)
				}
			})
		})

		Convey("ToInt", func() {
			Convey("Should convert Int", func() {
				input := NewInt(100)
				output := internal.ToInt(input)

				So(output, ShouldEqual, NewInt(100))
			})

			Convey("Should convert Float", func() {
				input := NewFloat(100.1)
				output := internal.ToInt(input)

				So(output, ShouldEqual, NewInt(100))
			})

			Convey("Should convert String", func() {
				input := NewString("100")
				output := internal.ToInt(input)

				So(output, ShouldEqual, NewInt(100))

				output2 := internal.ToInt(NewString("foobar"))
				So(output2, ShouldEqual, ZeroInt)
			})

			Convey("Should convert Boolean", func() {
				So(internal.ToInt(True), ShouldEqual, NewInt(1))
				So(internal.ToInt(False), ShouldEqual, NewInt(0))
			})

			Convey("Should convert Array with single item", func() {
				So(internal.ToInt(internal.NewArrayWith(NewFloat(1))), ShouldEqual, NewInt(1))
			})

			Convey("Should convert Array with multiple items", func() {
				arg := internal.NewArrayWith(NewFloat(1), NewFloat(1))

				So(internal.ToInt(arg), ShouldEqual, NewFloat(2))
			})

			Convey("Should convert DateTime", func() {
				dt := NewCurrentDateTime()
				ts := dt.Time.Unix()

				So(internal.ToInt(dt), ShouldEqual, NewInt(int(ts)))
			})

			Convey("Should NOT convert other types", func() {
				inputs := []Value{
					internal.NewObject(),
					NewBinary([]byte("")),
				}

				for _, input := range inputs {
					So(internal.ToInt(input), ShouldEqual, ZeroInt)
				}
			})
		})

		Convey("ToArray", func() {
			Convey("Should convert primitives", func() {
				dt := NewCurrentDateTime()

				inputs := [][]Value{
					{
						None,
						internal.NewArray(0),
					},
					{
						True,
						internal.NewArrayWith(True),
					},
					{
						NewInt(1),
						internal.NewArrayWith(NewInt(1)),
					},
					{
						NewFloat(1),
						internal.NewArrayWith(NewFloat(1)),
					},
					{
						NewString("foo"),
						internal.NewArrayWith(NewString("foo")),
					},
					{
						dt,
						internal.NewArrayWith(dt),
					},
				}

				for _, pairs := range inputs {
					actual := internal.ToArray(context.Background(), pairs[0])
					expected := pairs[1]

					So(actual.Compare(expected), ShouldEqual, 0)
				}
			})

			//Convey("Should create a copy of a given array", func() {
			//	vals := []core.Value{
			//		values.NewInt(1),
			//		values.NewInt(2),
			//		values.NewInt(3),
			//		values.NewInt(4),
			//		values.NewArray(10),
			//		values.NewObject(),
			//	}
			//
			//	input := values.NewArrayWith(vals...)
			//	arr := values.ToArray(context.Background(), input)
			//
			//	So(input == arr, ShouldBeFalse)
			//	So(arr.Length() == input.Length(), ShouldBeTrue)
			//
			//	for idx := range vals {
			//		expected := input.Get(values.NewInt(idx))
			//		actual := arr.Get(values.NewInt(idx))
			//
			//		// same ref
			//		So(actual == expected, ShouldBeTrue)
			//		So(actual.CompareValues(expected), ShouldEqual, 0)
			//	}
			//})

			//Convey("Should convert object to an array", func() {
			//	input := values.NewObjectWith(
			//		values.NewObjectProperty("foo", values.NewString("bar")),
			//		values.NewObjectProperty("baz", values.NewInt(1)),
			//		values.NewObjectProperty("qaz", values.NewObject()),
			//	)
			//
			//	arr := values.ToArray(context.Background(), input).Sort()
			//
			//	So(arr.String(), ShouldEqual, "[1,\"bar\",{}]")
			//	So(arr.Get(values.NewInt(2)) == input.MustGet("qaz"), ShouldBeTrue)
			//})
		})

		Convey("Unmarshal", func() {
			Convey("Should deserialize object", func() {
				input := map[string]interface{}{
					"foo": []string{
						"bar",
						"qaz",
					},
				}
				json1, err := json.Marshal(input)

				So(err, ShouldBeNil)

				val, err := internal.Unmarshal(json1)

				So(err, ShouldBeNil)

				json2, err := val.MarshalJSON()

				So(err, ShouldBeNil)
				So(json2, ShouldResemble, json1)
			})
		})
	})
}
