package core_test

import (
	"context"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"reflect"
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"
)

type CustomValue struct {
	properties map[core.Value]core.Value
}

type DummyStruct struct{}

func TestIsNil(t *testing.T) {
	Convey("Should match", t, func() {
		// nil == invalid
		t := core.IsNil(nil)

		So(t, ShouldBeTrue)

		a := []string{}
		t = core.IsNil(a)

		So(t, ShouldBeFalse)

		b := make([]string, 1)
		t = core.IsNil(b)

		So(t, ShouldBeFalse)

		c := make(map[string]string)
		t = core.IsNil(c)

		So(t, ShouldBeFalse)

		var s struct {
			Test string
		}
		t = core.IsNil(s)

		So(t, ShouldBeFalse)

		f := func() {}
		t = core.IsNil(f)

		So(t, ShouldBeFalse)

		i := DummyStruct{}
		t = core.IsNil(i)

		So(t, ShouldBeFalse)

		ch := make(chan string)
		t = core.IsNil(ch)

		So(t, ShouldBeFalse)

		var y unsafe.Pointer
		var vy int
		y = unsafe.Pointer(&vy)
		t = core.IsNil(y)

		So(t, ShouldBeFalse)
	})
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
					out := core.Parse(input.Raw)

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
						NewArray(1),
						True,
					},
					{
						NewObject(),
						True,
					},
					{
						NewBinary([]byte("")),
						True,
					},
				}

				for _, pair := range inputs {
					actual := core.ToBoolean(pair[0])
					expected := pair[1]

					So(actual, ShouldEqual, expected)
				}
			})
		})

		Convey("ToFloat", func() {
			Convey("Should convert Int", func() {
				input := NewInt(100)
				output := core.ToFloat(input)

				So(output, ShouldEqual, NewFloat(100))
			})

			Convey("Should convert Float", func() {
				input := NewFloat(100)
				output := core.ToFloat(input)

				So(output, ShouldEqual, NewFloat(100))
			})

			Convey("Should convert String", func() {
				input := NewString("100.1")
				output := core.ToFloat(input)

				So(output, ShouldEqual, NewFloat(100.1))

				output2 := core.ToFloat(NewString("foobar"))
				So(output2, ShouldEqual, ZeroFloat)
			})

			Convey("Should convert Boolean", func() {
				So(core.ToFloat(True), ShouldEqual, NewFloat(1))
				So(core.ToFloat(False), ShouldEqual, NewFloat(0))
			})

			Convey("Should convert Array with single item", func() {
				So(core.ToFloat(NewArrayWith(NewFloat(1))), ShouldEqual, NewFloat(1))
			})

			Convey("Should convert Array with multiple items", func() {
				arg := NewArrayWith(NewFloat(1), NewFloat(1))

				So(core.ToFloat(arg), ShouldEqual, NewFloat(2))
			})

			Convey("Should convert DateTime", func() {
				dt := NewCurrentDateTime()
				ts := dt.Time.Unix()

				So(core.ToFloat(dt), ShouldEqual, NewFloat(float64(ts)))
			})

			Convey("Should NOT convert other types", func() {
				inputs := []Value{
					NewObject(),
					NewBinary([]byte("")),
				}

				for _, input := range inputs {
					So(core.ToFloat(input), ShouldEqual, ZeroFloat)
				}
			})
		})

		Convey("ToInt", func() {
			Convey("Should convert Int", func() {
				input := NewInt(100)
				output := core.ToInt(input)

				So(output, ShouldEqual, NewInt(100))
			})

			Convey("Should convert Float", func() {
				input := NewFloat(100.1)
				output := core.ToInt(input)

				So(output, ShouldEqual, NewInt(100))
			})

			Convey("Should convert String", func() {
				input := NewString("100")
				output := core.ToInt(input)

				So(output, ShouldEqual, NewInt(100))

				output2 := core.ToInt(NewString("foobar"))
				So(output2, ShouldEqual, ZeroInt)
			})

			Convey("Should convert Boolean", func() {
				So(core.ToInt(True), ShouldEqual, NewInt(1))
				So(core.ToInt(False), ShouldEqual, NewInt(0))
			})

			Convey("Should convert Array with single item", func() {
				So(core.ToInt(NewArrayWith(NewFloat(1))), ShouldEqual, NewInt(1))
			})

			Convey("Should convert Array with multiple items", func() {
				arg := NewArrayWith(NewFloat(1), NewFloat(1))

				So(core.ToInt(arg), ShouldEqual, NewFloat(2))
			})

			Convey("Should convert DateTime", func() {
				dt := NewCurrentDateTime()
				ts := dt.Time.Unix()

				So(core.ToInt(dt), ShouldEqual, NewInt(int(ts)))
			})

			Convey("Should NOT convert other types", func() {
				inputs := []Value{
					NewObject(),
					NewBinary([]byte("")),
				}

				for _, input := range inputs {
					So(core.ToInt(input), ShouldEqual, ZeroInt)
				}
			})
		})

		Convey("ToList", func() {
			Convey("Should convert primitives", func() {
				dt := NewCurrentDateTime()

				inputs := [][]Value{
					{
						None,
						NewArray(0),
					},
					{
						True,
						NewArrayWith(True),
					},
					{
						NewInt(1),
						NewArrayWith(NewInt(1)),
					},
					{
						NewFloat(1),
						NewArrayWith(NewFloat(1)),
					},
					{
						NewString("foo"),
						NewArrayWith(NewString("foo")),
					},
					{
						dt,
						NewArrayWith(dt),
					},
				}

				for _, pairs := range inputs {
					actual := core.ToList(context.Background(), pairs[0])
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
			//	arr := values.ToList(context.Background(), input)
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
			//	arr := values.ToList(context.Background(), input).Sort()
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

				val, err := core.Unmarshal(json1)

				So(err, ShouldBeNil)

				json2, err := val.MarshalJSON()

				So(err, ShouldBeNil)
				So(json2, ShouldResemble, json1)
			})
		})
	})
}
