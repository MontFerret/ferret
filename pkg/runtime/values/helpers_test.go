package values_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var CustomType = core.NewType("custom")

type CustomValue struct {
	properties map[core.Value]core.Value
}

func (t *CustomValue) MarshalJSON() ([]byte, error) {
	return nil, core.ErrNotImplemented
}

func (t *CustomValue) Type() core.Type {
	return CustomType
}

func (t *CustomValue) String() string {
	return ""
}

func (t *CustomValue) Compare(other core.Value) int64 {
	return other.Compare(t) * -1
}

func (t *CustomValue) Unwrap() interface{} {
	return t
}

func (t *CustomValue) Hash() uint64 {
	return 0
}

func (t *CustomValue) Copy() core.Value {
	return values.None
}

func (t *CustomValue) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	if path == nil || len(path) == 0 {
		return values.None, nil
	}

	propKey := path[0]
	propValue, ok := t.properties[propKey]

	if !ok {
		return values.None, nil
	}

	if len(path) == 1 {
		return propValue, nil
	}

	return values.GetIn(ctx, propValue, path[1:])
}

func (t *CustomValue) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	if path == nil || len(path) == 0 {
		return nil
	}

	propKey := path[0]
	propValue, ok := t.properties[propKey]

	if !ok {
		return nil
	}

	if len(path) == 1 {
		t.properties[propKey] = value

		return nil
	}

	return values.SetIn(ctx, propValue, path[1:], value)
}

func TestHelpers(t *testing.T) {
	Convey("Helpers", t, func() {
		Convey("Getter", func() {
			Convey("It should get a value by a given path", func() {
				ct := &CustomValue{
					properties: map[core.Value]core.Value{
						values.NewString("foo"): values.NewInt(1),
						values.NewString("bar"): &CustomValue{
							properties: map[core.Value]core.Value{
								values.NewString("qaz"): values.NewInt(2),
							},
						},
					},
				}

				ctx := context.Background()

				foo, err := values.GetIn(ctx, ct, []core.Value{
					values.NewString("foo"),
				})

				So(err, ShouldBeNil)
				So(foo, ShouldEqual, values.NewInt(1))

				qaz, err := values.GetIn(ctx, ct, []core.Value{
					values.NewString("bar"),
					values.NewString("qaz"),
				})

				So(err, ShouldBeNil)
				So(qaz, ShouldEqual, values.NewInt(2))
			})
		})

		Convey("Setter", func() {
			Convey("It should get a value by a given path", func() {
				ct := &CustomValue{
					properties: map[core.Value]core.Value{
						values.NewString("foo"): values.NewInt(1),
						values.NewString("bar"): &CustomValue{
							properties: map[core.Value]core.Value{
								values.NewString("qaz"): values.NewInt(2),
							},
						},
					},
				}

				ctx := context.Background()

				err := values.SetIn(ctx, ct, []core.Value{
					values.NewString("foo"),
				}, values.NewInt(2))

				So(err, ShouldBeNil)
				So(ct.properties[values.NewString("foo")], ShouldEqual, values.NewInt(2))

				err = values.SetIn(ctx, ct, []core.Value{
					values.NewString("bar"),
					values.NewString("qaz"),
				}, values.NewString("foobar"))

				So(err, ShouldBeNil)

				qaz, err := values.GetIn(ctx, ct, []core.Value{
					values.NewString("bar"),
					values.NewString("qaz"),
				})

				So(err, ShouldBeNil)
				So(qaz, ShouldEqual, values.NewString("foobar"))
			})
		})

		Convey("ToBoolean", func() {
			Convey("Should convert values", func() {
				inputs := [][]core.Value{
					{
						values.None,
						values.False,
					},
					{
						values.True,
						values.True,
					},
					{
						values.False,
						values.False,
					},
					{
						values.NewInt(1),
						values.True,
					},
					{
						values.NewInt(0),
						values.False,
					},
					{
						values.NewFloat(1),
						values.True,
					},
					{
						values.NewFloat(0),
						values.False,
					},
					{
						values.NewString("Foo"),
						values.True,
					},
					{
						values.EmptyString,
						values.False,
					},
					{
						values.NewCurrentDateTime(),
						values.True,
					},
					{
						values.NewArray(1),
						values.True,
					},
					{
						values.NewObject(),
						values.True,
					},
					{
						values.NewBinary([]byte("")),
						values.True,
					},
				}

				for _, pair := range inputs {
					actual := values.ToBoolean(pair[0])
					expected := pair[1]

					So(actual, ShouldEqual, expected)
				}
			})
		})

		Convey("ToFloat", func() {
			Convey("Should convert Int", func() {
				input := values.NewInt(100)
				output, err := values.ToFloat(input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, values.NewFloat(100))
			})

			Convey("Should convert Float", func() {
				input := values.NewFloat(100)
				output, err := values.ToFloat(input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, values.NewFloat(100))
			})

			Convey("Should convert String", func() {
				input := values.NewString("100.1")
				output, err := values.ToFloat(input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, values.NewFloat(100.1))
			})

			Convey("Should NOT convert other types", func() {
				inputs := []core.Value{
					values.NewBoolean(true),
					values.NewCurrentDateTime(),
					values.NewArray(1),
					values.NewObject(),
					values.NewBinary([]byte("")),
				}

				for _, input := range inputs {
					_, err := values.ToFloat(input)

					So(err, ShouldNotBeNil)
				}
			})
		})

		Convey("ToInt", func() {
			Convey("Should convert Int", func() {
				input := values.NewInt(100)
				output, err := values.ToInt(input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, values.NewInt(100))
			})

			Convey("Should convert Float", func() {
				input := values.NewFloat(100.1)
				output, err := values.ToInt(input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, values.NewInt(100))
			})

			Convey("Should convert String", func() {
				input := values.NewString("100")
				output, err := values.ToInt(input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, values.NewInt(100))
			})

			Convey("Should NOT convert other types", func() {
				inputs := []core.Value{
					values.NewBoolean(true),
					values.NewCurrentDateTime(),
					values.NewArray(1),
					values.NewObject(),
					values.NewBinary([]byte("")),
				}

				for _, input := range inputs {
					_, err := values.ToInt(input)

					So(err, ShouldNotBeNil)
				}
			})
		})

		Convey("ToArray", func() {
			Convey("Should convert primitives", func() {
				dt := values.NewCurrentDateTime()

				inputs := [][]core.Value{
					{
						values.None,
						values.NewArray(0),
					},
					{
						values.True,
						values.NewArrayWith(values.True),
					},
					{
						values.NewInt(1),
						values.NewArrayWith(values.NewInt(1)),
					},
					{
						values.NewFloat(1),
						values.NewArrayWith(values.NewFloat(1)),
					},
					{
						values.NewString("foo"),
						values.NewArrayWith(values.NewString("foo")),
					},
					{
						dt,
						values.NewArrayWith(dt),
					},
				}

				for _, pairs := range inputs {
					actual, err := values.ToArray(context.Background(), pairs[0])
					expected := pairs[1]

					So(err, ShouldBeNil)
					So(actual.Compare(expected), ShouldEqual, 0)
				}
			})

			Convey("Should create a copy of a given array", func() {
				vals := []core.Value{
					values.NewInt(1),
					values.NewInt(2),
					values.NewInt(3),
					values.NewInt(4),
					values.NewArray(10),
					values.NewObject(),
				}

				input := values.NewArrayWith(vals...)
				output, err := values.ToArray(context.Background(), input)

				So(err, ShouldBeNil)

				arr := output.(*values.Array)

				So(input == arr, ShouldBeFalse)
				So(arr.Length() == input.Length(), ShouldBeTrue)

				for idx := range vals {
					expected := input.Get(values.NewInt(idx))
					actual := arr.Get(values.NewInt(idx))

					// same ref
					So(actual == expected, ShouldBeTrue)
					So(actual.Compare(expected), ShouldEqual, 0)
				}
			})

			Convey("Should convert object to an array", func() {
				input := values.NewObjectWith(
					values.NewObjectProperty("foo", values.NewString("bar")),
					values.NewObjectProperty("baz", values.NewInt(1)),
					values.NewObjectProperty("qaz", values.NewObject()),
				)

				output, err := values.ToArray(context.Background(), input)

				So(err, ShouldBeNil)

				arr := output.(*values.Array).Sort()

				So(arr.String(), ShouldEqual, "[1,\"bar\",{}]")
				So(arr.Get(values.NewInt(2)) == input.MustGet("qaz"), ShouldBeTrue)
			})
		})
	})
}
