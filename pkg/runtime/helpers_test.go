package runtime_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"unsafe"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

type CustomValue struct {
	properties map[runtime.Value]runtime.Value
}

type DummyStruct struct{}

func TestIsNil(t *testing.T) {
	Convey("Should match", t, func() {
		// nil == invalid
		t := runtime.IsNil(nil)

		So(t, ShouldBeTrue)

		a := []string{}
		t = runtime.IsNil(a)

		So(t, ShouldBeFalse)

		b := make([]string, 1)
		t = runtime.IsNil(b)

		So(t, ShouldBeFalse)

		c := make(map[string]string)
		t = runtime.IsNil(c)

		So(t, ShouldBeFalse)

		var s struct {
			Test string
		}
		t = runtime.IsNil(s)

		So(t, ShouldBeFalse)

		f := func() {}
		t = runtime.IsNil(f)

		So(t, ShouldBeFalse)

		i := DummyStruct{}
		t = runtime.IsNil(i)

		So(t, ShouldBeFalse)

		ch := make(chan string)
		t = runtime.IsNil(ch)

		So(t, ShouldBeFalse)

		var y unsafe.Pointer
		var vy int
		y = unsafe.Pointer(&vy)
		t = runtime.IsNil(y)

		So(t, ShouldBeFalse)
	})
}

func (t *CustomValue) MarshalJSON() ([]byte, error) {
	return nil, runtime.ErrNotImplemented
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

func (t *CustomValue) Copy() runtime.Value {
	return runtime.None
}

func TestHelpers(t *testing.T) {
	Convey("Helpers", t, func() {
		Convey("Parse", func() {
			Convey("It should parse values", func() {
				inputs := []struct {
					Parsed runtime.Value
					Raw    interface{}
				}{
					{Parsed: runtime.NewInt(1), Raw: int(1)},
					{Parsed: runtime.NewInt(1), Raw: int8(1)},
					{Parsed: runtime.NewInt(1), Raw: int16(1)},
					{Parsed: runtime.NewInt(1), Raw: int32(1)},
					{Parsed: runtime.NewInt(1), Raw: int64(1)},
				}

				for _, input := range inputs {
					out := runtime.Parse(input.Raw)

					expectedType := reflect.TypeOf(input.Parsed)
					actualType := reflect.TypeOf(out)

					So(actualType, ShouldEqual, expectedType)
					So(out.Unwrap(), ShouldEqual, input.Parsed.Unwrap())
				}
			})
		})

		Convey("ToBoolean", func() {
			Convey("Should convert values", func() {
				inputs := [][]runtime.Value{
					{
						runtime.None,
						runtime.False,
					},
					{
						runtime.True,
						runtime.True,
					},
					{
						runtime.False,
						runtime.False,
					},
					{
						runtime.NewInt(1),
						runtime.True,
					},
					{
						runtime.NewInt(0),
						runtime.False,
					},
					{
						runtime.NewFloat(1),
						runtime.True,
					},
					{
						runtime.NewFloat(0),
						runtime.False,
					},
					{
						runtime.NewString("Foo"),
						runtime.True,
					},
					{
						runtime.EmptyString,
						runtime.False,
					},
					{
						runtime.NewCurrentDateTime(),
						runtime.True,
					},
					{
						runtime.NewArray(1),
						runtime.True,
					},
					{
						runtime.NewObject(),
						runtime.True,
					},
					{
						runtime.NewBinary([]byte("")),
						runtime.True,
					},
				}

				for _, pair := range inputs {
					actual := runtime.ToBoolean(pair[0])
					expected := pair[1]

					So(actual, ShouldEqual, expected)
				}
			})
		})

		ctx := context.Background()

		Convey("ToFloat", func() {
			Convey("Should convert Int", func() {
				input := runtime.NewInt(100)
				output, err := runtime.ToFloat(ctx, input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, runtime.NewFloat(100))
			})

			Convey("Should convert Float", func() {
				input := runtime.NewFloat(100)
				output, err := runtime.ToFloat(ctx, input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, runtime.NewFloat(100))
			})

			Convey("Should convert String", func() {
				input := runtime.NewString("100.1")
				output, err := runtime.ToFloat(ctx, input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, runtime.NewFloat(100.1))

				output2, err := runtime.ToFloat(ctx, runtime.NewString("foobar"))
				So(output2, ShouldEqual, runtime.ZeroFloat)
			})

			Convey("Should convert Boolean", func() {
				out, err := runtime.ToFloat(ctx, runtime.True)
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewFloat(1))

				out, err = runtime.ToFloat(ctx, runtime.False)
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewFloat(0))
			})

			Convey("Should convert Array with single item", func() {
				out, err := runtime.ToFloat(ctx, runtime.NewArrayWith(runtime.NewFloat(1)))
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewFloat(1))
			})

			Convey("Should convert Array with multiple items", func() {
				arg := runtime.NewArrayWith(runtime.NewFloat(1), runtime.NewFloat(1))

				out, err := runtime.ToFloat(ctx, arg)
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewFloat(2))
			})

			Convey("Should convert DateTime", func() {
				dt := runtime.NewCurrentDateTime()
				ts := dt.Time.Unix()

				out, err := runtime.ToFloat(ctx, dt)
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewFloat(float64(ts)))
			})

			Convey("Should NOT convert other types", func() {
				inputs := []runtime.Value{
					runtime.NewObject(),
					runtime.NewBinary([]byte("")),
				}

				for _, input := range inputs {
					out, err := runtime.ToFloat(ctx, input)
					So(err, ShouldNotBeNil)
					So(out, ShouldEqual, runtime.ZeroFloat)
				}
			})
		})

		Convey("ToInt", func() {
			Convey("Should convert Int", func() {
				input := runtime.NewInt(100)
				output, err := runtime.ToInt(ctx, input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, runtime.NewInt(100))
			})

			Convey("Should convert Float", func() {
				input := runtime.NewFloat(100.1)
				output, err := runtime.ToInt(ctx, input)
				So(err, ShouldBeNil)
				So(output, ShouldEqual, runtime.NewInt(100))
			})

			Convey("Should convert String", func() {
				input := runtime.NewString("100")
				output, err := runtime.ToInt(ctx, input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, runtime.NewInt(100))

				output2, err := runtime.ToInt(ctx, runtime.NewString("foobar"))
				So(err, ShouldNotBeNil)
				So(output2, ShouldEqual, runtime.ZeroInt)
			})

			Convey("Should convert Boolean", func() {
				out, err := runtime.ToInt(ctx, runtime.True)
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewInt(1))

				out, err = runtime.ToInt(ctx, runtime.False)
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewInt(0))
			})

			Convey("Should convert Array with single item", func() {
				out, err := runtime.ToInt(ctx, runtime.NewArrayWith(runtime.NewInt(1)))
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewInt(1))
			})

			Convey("Should convert Array with multiple items", func() {
				arg := runtime.NewArrayWith(runtime.NewFloat(1), runtime.NewFloat(1))
				out, err := runtime.ToInt(ctx, arg)

				So(err, ShouldBeNil)

				So(out, ShouldEqual, runtime.NewFloat(2))
			})

			Convey("Should convert DateTime", func() {
				dt := runtime.NewCurrentDateTime()
				ts := dt.Time.Unix()
				out, err := runtime.ToInt(ctx, dt)

				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewInt(int(ts)))
			})

			Convey("Should NOT convert other types", func() {
				inputs := []runtime.Value{
					runtime.NewObject(),
					runtime.NewBinary([]byte("")),
				}

				for _, input := range inputs {
					out, err := runtime.ToInt(ctx, input)
					So(err, ShouldNotBeNil)
					So(out, ShouldEqual, runtime.ZeroInt)
				}
			})
		})

		Convey("ToList", func() {
			Convey("Should convert primitives", func() {
				dt := runtime.NewCurrentDateTime()

				inputs := [][]runtime.Value{
					{
						runtime.None,
						runtime.NewArray(0),
					},
					{
						runtime.True,
						runtime.NewArrayWith(runtime.True),
					},
					{
						runtime.NewInt(1),
						runtime.NewArrayWith(runtime.NewInt(1)),
					},
					{
						runtime.NewFloat(1),
						runtime.NewArrayWith(runtime.NewFloat(1)),
					},
					{
						runtime.NewString("foo"),
						runtime.NewArrayWith(runtime.NewString("foo")),
					},
					{
						dt,
						runtime.NewArrayWith(dt),
					},
				}

				for _, pairs := range inputs {
					actual := runtime.ToList(context.Background(), pairs[0])
					expected := pairs[1]

					So(actual.Compare(expected), ShouldEqual, 0)
				}
			})

			//Convey("Should create a copy of a given array", func() {
			//	vals := []core.runtime.Value{
			//		values.runtime.NewInt(1),
			//		values.runtime.NewInt(2),
			//		values.runtime.NewInt(3),
			//		values.runtime.NewInt(4),
			//		values.runtime.NewArray(10),
			//		values.runtime.NewObject(),
			//	}
			//
			//	input := values.runtime.NewArrayWith(vals...)
			//	arr := values.ToList(context.Background(), input)
			//
			//	So(input == arr, ShouldBeFalse)
			//	So(arr.Length() == input.Length(), ShouldBeTrue)
			//
			//	for idx := range vals {
			//		expected := input.Get(values.runtime.NewInt(idx))
			//		actual := arr.Get(values.runtime.NewInt(idx))
			//
			//		// same ref
			//		So(actual == expected, ShouldBeTrue)
			//		So(actual.CompareValues(expected), ShouldEqual, 0)
			//	}
			//})

			//Convey("Should convert object to an array", func() {
			//	input := values.NewObjectWith(
			//		values.NewObjectProperty("foo", values.runtime.NewString("bar")),
			//		values.NewObjectProperty("baz", values.runtime.NewInt(1)),
			//		values.NewObjectProperty("qaz", values.runtime.NewObject()),
			//	)
			//
			//	arr := values.ToList(context.Background(), input).Sort()
			//
			//	So(arr.String(), ShouldEqual, "[1,\"bar\",{}]")
			//	So(arr.Get(values.runtime.NewInt(2)) == input.MustGet("qaz"), ShouldBeTrue)
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

				val, err := runtime.Unmarshal(json1)

				So(err, ShouldBeNil)

				json2, err := val.MarshalJSON()

				So(err, ShouldBeNil)
				So(json2, ShouldResemble, json1)
			})
		})
	})
}
