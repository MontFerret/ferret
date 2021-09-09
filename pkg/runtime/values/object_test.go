package values_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	. "github.com/smartystreets/goconvey/convey"
)

func TestObject(t *testing.T) {
	Convey("#constructor", t, func() {
		Convey("Should create an empty object", func() {
			obj := values.NewObject()

			So(obj.Length(), ShouldEqual, 0)
		})

		Convey("Should create an object, from passed values", func() {
			obj := values.NewObjectWith(
				values.NewObjectProperty("none", values.None),
				values.NewObjectProperty("boolean", values.False),
				values.NewObjectProperty("int", values.NewInt(1)),
				values.NewObjectProperty("float", values.Float(1)),
				values.NewObjectProperty("string", values.NewString("1")),
				values.NewObjectProperty("array", values.NewArray(10)),
				values.NewObjectProperty("object", values.NewObject()),
			)

			So(obj.Length(), ShouldEqual, 7)
		})
	})

	Convey(".MarshalJSON", t, func() {
		Convey("Should serialize an empty object", func() {
			obj := values.NewObject()
			marshaled, err := obj.MarshalJSON()

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "{}")
		})

		Convey("Should serialize full object", func() {
			obj := values.NewObjectWith(
				values.NewObjectProperty("none", values.None),
				values.NewObjectProperty("boolean", values.False),
				values.NewObjectProperty("int", values.NewInt(1)),
				values.NewObjectProperty("float", values.Float(1)),
				values.NewObjectProperty("string", values.NewString("1")),
				values.NewObjectProperty("array", values.NewArray(10)),
				values.NewObjectProperty("object", values.NewObject()),
			)
			marshaled, err := obj.MarshalJSON()

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "{\"array\":[],\"boolean\":false,\"float\":1,\"int\":1,\"none\":null,\"object\":{},\"string\":\"1\"}")
		})
	})

	Convey(".Type", t, func() {
		Convey("Should return type", func() {
			obj := values.NewObject()

			So(obj.Type().Equals(types.Object), ShouldBeTrue)
		})
	})

	Convey(".Unwrap", t, func() {
		Convey("Should return an unwrapped items", func() {
			obj := values.NewObjectWith(
				values.NewObjectProperty("foo", values.NewString("foo")),
				values.NewObjectProperty("bar", values.NewString("bar")),
			)

			for _, val := range obj.Unwrap().(map[string]interface{}) {
				So(val, ShouldHaveSameTypeAs, "")
			}
		})
	})

	Convey(".String", t, func() {
		Convey("Should return a string representation ", func() {
			obj := values.NewObjectWith(
				values.NewObjectProperty("foo", values.NewString("bar")),
			)

			So(obj.String(), ShouldEqual, "{\"foo\":\"bar\"}")
		})
	})

	Convey(".Compare", t, func() {
		Convey("It should return 1 for all non-object values", func() {
			arr := []core.Value{
				values.None,
				values.False,
				values.NewInt(1),
				values.Float(1),
				values.NewString("1"),
				values.NewArray(10),
			}
			obj := values.NewObject()

			for _, val := range arr {
				So(obj.Compare(val), ShouldEqual, 1)
			}
		})

		Convey("It should return -1 for all object values", func() {
			arr := values.NewArrayWith(values.ZeroInt, values.ZeroInt)
			obj := values.NewObject()

			So(arr.Compare(obj), ShouldEqual, -1)
		})

		Convey("It should return 0 when both objects are empty", func() {
			obj1 := values.NewObject()
			obj2 := values.NewObject()

			So(obj1.Compare(obj2), ShouldEqual, 0)
		})

		Convey("It should return 0 when both objects are equal (independent of key order)", func() {
			obj1 := values.NewObjectWith(
				values.NewObjectProperty("foo", values.NewString("foo")),
				values.NewObjectProperty("bar", values.NewString("bar")),
			)
			obj2 := values.NewObjectWith(
				values.NewObjectProperty("foo", values.NewString("foo")),
				values.NewObjectProperty("bar", values.NewString("bar")),
			)

			So(obj1.Compare(obj1), ShouldEqual, 0)
			So(obj2.Compare(obj2), ShouldEqual, 0)
			So(obj1.Compare(obj2), ShouldEqual, 0)
			So(obj2.Compare(obj1), ShouldEqual, 0)
		})

		Convey("It should return 1 when other array is empty", func() {
			obj1 := values.NewObjectWith(values.NewObjectProperty("foo", values.NewString("bar")))
			obj2 := values.NewObject()

			So(obj1.Compare(obj2), ShouldEqual, 1)
		})

		Convey("It should return 1 when values are bigger", func() {
			obj1 := values.NewObjectWith(values.NewObjectProperty("foo", values.NewFloat(3)))
			obj2 := values.NewObjectWith(values.NewObjectProperty("foo", values.NewFloat(2)))

			So(obj1.Compare(obj2), ShouldEqual, 1)
		})

		Convey("It should return 1 when values are less", func() {
			obj1 := values.NewObjectWith(values.NewObjectProperty("foo", values.NewFloat(1)))
			obj2 := values.NewObjectWith(values.NewObjectProperty("foo", values.NewFloat(2)))

			So(obj1.Compare(obj2), ShouldEqual, -1)
		})

		Convey("ArangoDB compatibility", func() {
			Convey("It should return 1 when {a:1} and {b:2}", func() {
				obj1 := values.NewObjectWith(values.NewObjectProperty("a", values.NewInt(1)))
				obj2 := values.NewObjectWith(values.NewObjectProperty("b", values.NewInt(2)))

				So(obj1.Compare(obj2), ShouldEqual, 1)
			})

			Convey("It should return 0 when {a:1} and {a:1}", func() {
				obj1 := values.NewObjectWith(values.NewObjectProperty("a", values.NewInt(1)))
				obj2 := values.NewObjectWith(values.NewObjectProperty("a", values.NewInt(1)))

				So(obj1.Compare(obj2), ShouldEqual, 0)
			})

			Convey("It should return 0 {a:1, c:2} and {c:2, a:1}", func() {
				obj1 := values.NewObjectWith(
					values.NewObjectProperty("a", values.NewInt(1)),
					values.NewObjectProperty("c", values.NewInt(2)),
				)
				obj2 := values.NewObjectWith(
					values.NewObjectProperty("c", values.NewInt(2)),
					values.NewObjectProperty("a", values.NewInt(1)),
				)

				So(obj1.Compare(obj2), ShouldEqual, 0)
			})

			Convey("It should return -1 when {a:1} and {a:2}", func() {
				obj1 := values.NewObjectWith(values.NewObjectProperty("a", values.NewInt(1)))
				obj2 := values.NewObjectWith(values.NewObjectProperty("a", values.NewInt(2)))

				So(obj1.Compare(obj2), ShouldEqual, -1)
			})

			Convey("It should return 1 when {a:1, c:2} and {c:2, b:2}", func() {
				obj1 := values.NewObjectWith(
					values.NewObjectProperty("a", values.NewInt(1)),
					values.NewObjectProperty("c", values.NewInt(2)),
				)
				obj2 := values.NewObjectWith(
					values.NewObjectProperty("c", values.NewInt(2)),
					values.NewObjectProperty("b", values.NewInt(2)),
				)

				So(obj1.Compare(obj2), ShouldEqual, 1)
			})

			Convey("It should return 1 {a:1, c:3} and {c:2, a:1}", func() {
				obj1 := values.NewObjectWith(
					values.NewObjectProperty("a", values.NewInt(1)),
					values.NewObjectProperty("c", values.NewInt(3)),
				)
				obj2 := values.NewObjectWith(
					values.NewObjectProperty("c", values.NewInt(2)),
					values.NewObjectProperty("a", values.NewInt(1)),
				)

				So(obj1.Compare(obj2), ShouldEqual, 1)
			})
		})
	})

	Convey(".Hash", t, func() {
		Convey("It should calculate hash of non-empty object", func() {
			v := values.NewObjectWith(
				values.NewObjectProperty("foo", values.NewString("bar")),
				values.NewObjectProperty("faz", values.NewInt(1)),
				values.NewObjectProperty("qaz", values.True),
			)

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)
		})

		Convey("It should calculate hash of empty object", func() {
			v := values.NewObject()

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)
		})

		Convey("Hash sum should be consistent", func() {
			v := values.NewObjectWith(
				values.NewObjectProperty("boolean", values.True),
				values.NewObjectProperty("int", values.NewInt(1)),
				values.NewObjectProperty("float", values.NewFloat(1.1)),
				values.NewObjectProperty("string", values.NewString("foobar")),
				values.NewObjectProperty("datetime", values.NewCurrentDateTime()),
				values.NewObjectProperty("array", values.NewArrayWith(values.NewInt(1), values.True)),
				values.NewObjectProperty("object", values.NewObjectWith(values.NewObjectProperty("foo", values.NewString("bar")))),
			)

			h1 := v.Hash()
			h2 := v.Hash()

			So(h1, ShouldEqual, h2)
		})
	})

	Convey(".Length", t, func() {
		Convey("Should return 0 when empty", func() {
			obj := values.NewObject()

			So(obj.Length(), ShouldEqual, 0)
		})

		Convey("Should return greater than 0 when not empty", func() {
			obj := values.NewObjectWith(
				values.NewObjectProperty("foo", values.ZeroInt),
				values.NewObjectProperty("bar", values.ZeroInt),
			)

			So(obj.Length(), ShouldEqual, 2)
		})
	})

	Convey(".ForEach", t, func() {
		Convey("Should iterate over elements", func() {
			obj := values.NewObjectWith(
				values.NewObjectProperty("foo", values.ZeroInt),
				values.NewObjectProperty("bar", values.ZeroInt),
			)
			counter := 0

			obj.ForEach(func(value core.Value, key string) bool {
				counter++

				return true
			})

			So(counter, ShouldEqual, obj.Length())
		})

		Convey("Should not iterate when empty", func() {
			obj := values.NewObject()
			counter := 0

			obj.ForEach(func(value core.Value, key string) bool {
				counter++

				return true
			})

			So(counter, ShouldEqual, obj.Length())
		})

		Convey("Should break iteration when false returned", func() {
			obj := values.NewObjectWith(
				values.NewObjectProperty("1", values.NewInt(1)),
				values.NewObjectProperty("2", values.NewInt(2)),
				values.NewObjectProperty("3", values.NewInt(3)),
				values.NewObjectProperty("4", values.NewInt(4)),
				values.NewObjectProperty("5", values.NewInt(5)),
			)
			threshold := 3
			counter := 0

			obj.ForEach(func(value core.Value, key string) bool {
				counter++

				return counter < threshold
			})

			So(counter, ShouldEqual, threshold)
		})
	})

	Convey(".Get", t, func() {
		Convey("Should return item by key", func() {
			obj := values.NewObjectWith(
				values.NewObjectProperty("foo", values.NewInt(1)),
				values.NewObjectProperty("bar", values.NewInt(2)),
				values.NewObjectProperty("qaz", values.NewInt(3)),
			)

			el, _ := obj.Get("foo")

			So(el.Compare(values.NewInt(1)), ShouldEqual, 0)
		})

		Convey("Should return None when no items", func() {
			obj := values.NewObject()

			el, _ := obj.Get("foo")

			So(el.Compare(values.None), ShouldEqual, 0)
		})
	})

	Convey(".Set", t, func() {
		Convey("Should set item by index", func() {
			obj := values.NewObject()

			obj.Set("foo", values.NewInt(1))

			So(obj.Length(), ShouldEqual, 1)

			v, _ := obj.Get("foo")
			So(v.Compare(values.NewInt(1)), ShouldEqual, 0)
		})
	})

	Convey(".Clone", t, func() {
		Convey("Cloned object should be equal to source object", func() {
			obj := values.NewObjectWith(
				values.NewObjectProperty("one", values.NewInt(1)),
				values.NewObjectProperty("two", values.NewInt(2)),
			)

			clone := obj.Clone().(*values.Object)

			So(obj.Compare(clone), ShouldEqual, 0)
		})

		Convey("Cloned object should be independent of the source object", func() {
			obj := values.NewObjectWith(
				values.NewObjectProperty("one", values.NewInt(1)),
				values.NewObjectProperty("two", values.NewInt(2)),
			)

			clone := obj.Clone().(*values.Object)

			obj.Remove(values.NewString("one"))

			So(obj.Compare(clone), ShouldNotEqual, 0)
		})

		Convey("Cloned object must contain copies of the nested objects", func() {
			obj := values.NewObjectWith(
				values.NewObjectProperty(
					"arr", values.NewArrayWith(values.NewInt(1)),
				),
			)

			clone := obj.Clone().(*values.Object)

			nestedInObj, _ := obj.Get(values.NewString("arr"))
			nestedInObjArr := nestedInObj.(*values.Array)
			nestedInObjArr.Push(values.NewInt(2))

			nestedInClone, _ := clone.Get(values.NewString("arr"))
			nestedInCloneArr := nestedInClone.(*values.Array)

			So(nestedInObjArr.Compare(nestedInCloneArr), ShouldNotEqual, 0)
		})
	})

	Convey(".GetIn", t, func() {

		ctx := context.Background()

		Convey("Should return the same as .Get when input is correct", func() {

			Convey("Should return item by key", func() {
				key := values.NewString("foo")
				obj := values.NewObjectWith(
					values.NewObjectProperty(key.String(), values.NewInt(1)),
				)

				el, err := obj.GetIn(ctx, []core.Value{key})
				elGet, _ := obj.Get(key)

				So(err, ShouldBeNil)
				So(el.Compare(elGet), ShouldEqual, 0)
			})

			Convey("Should return None when no items", func() {
				key := values.NewString("foo")
				obj := values.NewObject()

				el, err := obj.GetIn(ctx, []core.Value{key})
				elGet, _ := obj.Get(key)

				So(err, ShouldBeNil)
				So(el.Compare(elGet), ShouldEqual, 0)
			})
		})

		Convey("Should error when input is not correct", func() {

			Convey("Should return None when path[0] is not a string", func() {
				obj := values.NewObject()
				path := []core.Value{values.NewInt(0)}

				el, err := obj.GetIn(ctx, path)

				So(err, ShouldBeNil)
				So(el, ShouldNotBeNil)
				So(el.Type().String(), ShouldEqual, types.None.String())
			})

			Convey("Should error when first received item is not a Getter and len(path) > 1", func() {
				key := values.NewString("foo")
				obj := values.NewObjectWith(
					values.NewObjectProperty(key.String(), values.NewInt(1)),
				)
				path := []core.Value{key, key}

				el, err := obj.GetIn(ctx, path)

				So(err, ShouldBeError)
				So(el.Compare(values.None), ShouldEqual, 0)
			})
		})

		Convey("Should return None when len(path) == 0", func() {
			obj := values.NewObject()

			el, err := obj.GetIn(ctx, nil)

			So(err, ShouldBeNil)
			So(el.Compare(values.None), ShouldEqual, 0)
		})

		Convey("Should call the nested Getter", func() {
			key := values.NewString("foo")
			obj := values.NewObjectWith(
				values.NewObjectProperty(key.String(), values.NewArrayWith(key)),
			)

			el, err := obj.GetIn(ctx, []core.Value{
				key,              // obj.foo
				values.NewInt(0), // obj.foo[0]
			})

			So(err, ShouldBeNil)
			So(el.Compare(key), ShouldEqual, 0)
		})
	})
}
