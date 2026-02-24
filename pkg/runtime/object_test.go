package runtime_test

import (
	c "context"
	"testing"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"

	. "github.com/MontFerret/ferret/v2/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

func TestObject(t *testing.T) {
	Convey("#constructor", t, func() {
		Convey("Should create an empty object", func() {
			obj := NewObject()
			size, err := obj.Length(c.Background())

			So(err, ShouldBeNil)
			So(size, ShouldEqual, 0)
		})

		Convey("Should create an object, from passed values", func() {
			values := map[string]Value{
				"none":    None,
				"boolean": False,
				"int":     NewInt(1),
				"float":   Float(1),
				"string":  NewString("1"),
				"array":   NewArray(10),
				"object":  NewObject(),
			}
			obj := NewObjectWith(values)

			size, err := obj.Length(c.Background())

			So(err, ShouldBeNil)
			So(size, ShouldEqual, 7)

			values["int"] = ZeroInt

			current, _ := obj.Get(c.Background(), NewString("int"))

			So(current.String(), ShouldNotEqual, values["int"].String())
		})

		Convey("Should implement runtime.Map interface", func() {
			var obj Value
			obj = NewObject()

			_, ok := obj.(Map)

			So(ok, ShouldBeTrue)
		})
	})

	Convey(".EncodeJSON", t, func() {
		Convey("Should serialize an empty object", func() {
			obj := NewObject()
			marshaled, err := encodingjson.Default.Encode(obj)

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "{}")
		})

		Convey("Should serialize full object", func() {
			obj := NewObjectWith(
				map[string]Value{
					"none":    None,
					"boolean": False,
					"int":     NewInt(1),
					"float":   Float(1),
					"string":  NewString("1"),
					"array":   NewArray(10),
					"object":  NewObject(),
				},
			)
			marshaled, err := encodingjson.Default.Encode(obj)

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "{\"array\":[],\"boolean\":false,\"float\":1,\"int\":1,\"none\":null,\"object\":{},\"string\":\"1\"}")
		})
	})

	Convey(".String", t, func() {
		Convey("Should return a string representation ", func() {
			obj := NewObjectWith(
				map[string]Value{
					"foo": NewString("bar"),
				},
			)

			So(obj.String(), ShouldEqual, "{\"foo\":\"bar\"}")
		})
	})

	Convey(".CompareValues", t, func() {
		Convey("It should return 1 for all non-object values", func() {
			arr := []Value{
				None,
				False,
				NewInt(1),
				Float(1),
				NewString("1"),
				NewArray(10),
			}
			obj := NewObject()

			for _, val := range arr {
				So(obj.Compare(val), ShouldEqual, 1)
			}
		})

		Convey("It should return -1 for all object values", func() {
			arr := NewArrayWith(ZeroInt, ZeroInt)
			obj := NewObject()

			So(arr.Compare(obj), ShouldEqual, -1)
		})

		Convey("It should return 0 when both objects are empty", func() {
			obj1 := NewObject()
			obj2 := NewObject()

			So(obj1.Compare(obj2), ShouldEqual, 0)
		})

		Convey("It should return 0 when both objects are equal (independent of key order)", func() {
			obj1 := NewObjectWith(
				map[string]Value{
					"foo": NewString("foo"),
					"bar": NewString("bar"),
				},
			)
			obj2 := NewObjectWith(
				map[string]Value{
					"foo": NewString("foo"),
					"bar": NewString("bar"),
				},
			)

			So(obj1.Compare(obj1), ShouldEqual, 0)
			So(obj2.Compare(obj2), ShouldEqual, 0)
			So(obj1.Compare(obj2), ShouldEqual, 0)
			So(obj2.Compare(obj1), ShouldEqual, 0)
		})

		Convey("It should return 1 when other array is empty", func() {
			obj1 := NewObjectWith(
				map[string]Value{
					"foo": NewString("bar"),
				},
			)
			obj2 := NewObject()

			So(obj1.Compare(obj2), ShouldEqual, 1)
		})

		Convey("It should return 1 when values are bigger", func() {
			obj1 := NewObjectWith(
				map[string]Value{
					"foo": NewFloat(3),
				},
			)
			obj2 := NewObjectWith(
				map[string]Value{
					"foo": NewFloat(2),
				},
			)

			So(obj1.Compare(obj2), ShouldEqual, 1)
		})

		Convey("It should return 1 when values are less", func() {
			obj1 := NewObjectWith(
				map[string]Value{
					"foo": NewFloat(1),
				},
			)
			obj2 := NewObjectWith(
				map[string]Value{
					"foo": NewFloat(2),
				},
			)

			So(obj1.Compare(obj2), ShouldEqual, -1)
		})

		Convey("ArangoDB compatibility", func() {
			Convey("It should return 1 when {a:1} and {b:2}", func() {
				obj1 := NewObjectWith(
					map[string]Value{
						"a": NewInt(1),
					},
				)
				obj2 := NewObjectWith(
					map[string]Value{
						"b": NewInt(2),
					},
				)

				So(obj1.Compare(obj2), ShouldEqual, 1)
			})

			Convey("It should return 0 when {a:1} and {a:1}", func() {
				obj1 := NewObjectWith(
					map[string]Value{
						"a": NewInt(1),
					},
				)
				obj2 := NewObjectWith(
					map[string]Value{
						"a": NewInt(1),
					},
				)

				So(obj1.Compare(obj2), ShouldEqual, 0)
			})

			Convey("It should return 0 {a:1, c:2} and {c:2, a:1}", func() {
				obj1 := NewObjectWith(
					map[string]Value{
						"a": NewInt(1),
						"c": NewInt(2),
					},
				)
				obj2 := NewObjectWith(
					map[string]Value{
						"c": NewInt(2),
						"a": NewInt(1),
					},
				)

				So(obj1.Compare(obj2), ShouldEqual, 0)
			})

			Convey("It should return -1 when {a:1} and {a:2}", func() {
				obj1 := NewObjectWith(
					map[string]Value{
						"a": NewInt(1),
					},
				)
				obj2 := NewObjectWith(
					map[string]Value{
						"a": NewInt(2),
					},
				)

				So(obj1.Compare(obj2), ShouldEqual, -1)
			})

			Convey("It should return 1 when {a:1, c:2} and {c:2, b:2}", func() {
				obj1 := NewObjectWith(
					map[string]Value{
						"a": NewInt(1),
						"c": NewInt(2),
					},
				)
				obj2 := NewObjectWith(
					map[string]Value{
						"c": NewInt(2),
						"b": NewInt(2),
					},
				)

				So(obj1.Compare(obj2), ShouldEqual, 1)
			})

			Convey("It should return 1 {a:1, c:3} and {c:2, a:1}", func() {
				obj1 := NewObjectWith(
					map[string]Value{
						"a": NewInt(1),
						"c": NewInt(3),
					},
				)
				obj2 := NewObjectWith(
					map[string]Value{
						"c": NewInt(2),
						"a": NewInt(1),
					},
				)

				So(obj1.Compare(obj2), ShouldEqual, 1)
			})
		})
	})

	Convey(".Hash", t, func() {
		Convey("It should calculate hash of non-empty object", func() {
			v := NewObjectWith(
				map[string]Value{
					"foo": NewString("bar"),
					"faz": NewInt(1),
					"qaz": True,
				},
			)

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)
		})

		Convey("It should calculate hash of empty object", func() {
			v := NewObject()

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)
		})

		Convey("Hash sum should be consistent", func() {
			v := NewObjectWith(
				map[string]Value{
					"boolean":  True,
					"int":      NewInt(1),
					"float":    NewFloat(1.1),
					"string":   NewString("foobar"),
					"datetime": NewCurrentDateTime(),
					"array":    NewArrayWith(NewInt(1), True),
					"object":   NewObjectWith(map[string]Value{"foo": NewString("bar")}),
				},
			)

			h1 := v.Hash()
			h2 := v.Hash()

			So(h1, ShouldEqual, h2)
		})
	})

	Convey(".Length", t, func() {
		Convey("Should return 0 when empty", func() {
			obj := NewObject()
			size, err := obj.Length(c.Background())

			So(err, ShouldBeNil)
			So(size, ShouldEqual, 0)
		})

		Convey("Should return greater than 0 when not empty", func() {
			obj := NewObjectWith(
				map[string]Value{
					"foo": ZeroInt,
					"bar": ZeroInt,
				},
			)

			size, err := obj.Length(c.Background())
			So(err, ShouldBeNil)
			So(size, ShouldEqual, 2)
		})
	})

	Convey(".ForEachIter", t, func() {
		Convey("Should iterate over elements", func() {
			obj := NewObjectWith(
				map[string]Value{
					"bar": ZeroInt,
					"foo": ZeroInt,
				},
			)
			counter := 0
			ctx := c.Background()

			obj.ForEach(ctx, func(_ c.Context, value, key Value) (Boolean, error) {
				counter++

				return true, nil
			})

			size, err := obj.Length(ctx)

			So(err, ShouldBeNil)
			So(counter, ShouldEqual, size)
		})

		Convey("Should not iterate when empty", func() {
			obj := NewObject()
			counter := 0

			ctx := c.Background()

			obj.ForEach(ctx, func(_ c.Context, value, key Value) (Boolean, error) {
				counter++

				return true, nil
			})

			size, err := obj.Length(ctx)

			So(err, ShouldBeNil)
			So(counter, ShouldEqual, size)
		})

		Convey("Should break iteration when false returned", func() {
			obj := NewObjectWith(
				map[string]Value{
					"1": NewInt(1),
					"2": NewInt(2),
					"3": NewInt(3),
					"4": NewInt(4),
					"5": NewInt(5),
				},
			)
			threshold := 3
			counter := 0
			ctx := c.Background()

			obj.ForEach(ctx, func(_ c.Context, value, key Value) (Boolean, error) {
				counter++

				return counter < threshold, nil
			})

			So(counter, ShouldEqual, threshold)
		})
	})

	Convey(".At", t, func() {
		//Convey("Should return item by key", func() {
		//	obj := values.NewObjectWith(
		//		values.NewObjectProperty("foo", values.NewInt(1)),
		//		values.NewObjectProperty("bar", values.NewInt(2)),
		//		values.NewObjectProperty("qaz", values.NewInt(3)),
		//	)
		//
		//	el, _ := obj.At("foo")
		//
		//	So(el.CompareValues(values.NewInt(1)), ShouldEqual, 0)
		//})

		//Convey("Should return None when no items", func() {
		//	obj := values.NewObject()
		//
		//	el, _ := obj.At("foo")
		//
		//	So(el.CompareValues(values.None), ShouldEqual, 0)
		//})
	})

	Convey(".SetAt", t, func() {
		//Convey("Should set item by index", func() {
		//	obj := values.NewObject()
		//
		//	obj.SetAt("foo", values.NewInt(1))
		//
		//	So(obj.Length(), ShouldEqual, 1)
		//
		//	v, _ := obj.At("foo")
		//	So(v.CompareValues(values.NewInt(1)), ShouldEqual, 0)
		//})
	})

	Convey(".Clone", t, func() {
		Convey("Cloned object should be equal to source object", func() {
			obj := NewObjectWith(
				map[string]Value{
					"one": NewInt(1),
					"two": NewInt(2),
				},
			)

			ctx := c.Background()
			cloned, _ := obj.Clone(ctx)
			clone := cloned.(*Object)

			So(obj.Compare(clone), ShouldEqual, 0)
		})

		Convey("Cloned object should be independent of the source object", func() {
			obj := NewObjectWith(
				map[string]Value{
					"one": NewInt(1),
					"two": NewInt(2),
				},
			)

			ctx := c.Background()
			cloned, _ := obj.Clone(ctx)
			clone := cloned.(*Object)

			obj.RemoveKey(ctx, NewString("one"))

			So(obj.Compare(clone), ShouldNotEqual, 0)
		})

		Convey("Cloned object must contain copies of the nested objects", func() {
			obj := NewObjectWith(
				map[string]Value{
					"arr": NewArrayWith(NewInt(1)),
				},
			)

			ctx := c.Background()
			cloned, _ := obj.Clone(ctx)
			clone := cloned.(*Object)

			nestedInObj, _ := obj.Get(ctx, NewString("arr"))
			nestedInObjArr := nestedInObj.(*Array)
			nestedInObjArr.Append(ctx, NewInt(2))

			nestedInClone, _ := clone.Get(ctx, NewString("arr"))
			nestedInCloneArr := nestedInClone.(*Array)

			So(nestedInObjArr.Compare(nestedInCloneArr), ShouldNotEqual, 0)
		})
	})
}
