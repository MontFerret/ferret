package internal_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func TestObject(t *testing.T) {
	Convey("#constructor", t, func() {
		Convey("Should create an empty object", func() {
			obj := internal.NewObject()

			So(obj.Length(), ShouldEqual, 0)
		})

		Convey("Should create an object, from passed values", func() {
			obj := internal.NewObjectWith(
				internal.NewObjectProperty("none", core.None),
				internal.NewObjectProperty("boolean", core.False),
				internal.NewObjectProperty("int", core.NewInt(1)),
				internal.NewObjectProperty("float", core.Float(1)),
				internal.NewObjectProperty("string", core.NewString("1")),
				internal.NewObjectProperty("array", internal.NewArray(10)),
				internal.NewObjectProperty("object", internal.NewObject()),
			)

			So(obj.Length(), ShouldEqual, 7)
		})
	})

	Convey(".MarshalJSON", t, func() {
		Convey("Should serialize an empty object", func() {
			obj := internal.NewObject()
			marshaled, err := obj.MarshalJSON()

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "{}")
		})

		Convey("Should serialize full object", func() {
			obj := internal.NewObjectWith(
				internal.NewObjectProperty("none", core.None),
				internal.NewObjectProperty("boolean", core.False),
				internal.NewObjectProperty("int", core.NewInt(1)),
				internal.NewObjectProperty("float", core.Float(1)),
				internal.NewObjectProperty("string", core.NewString("1")),
				internal.NewObjectProperty("array", internal.NewArray(10)),
				internal.NewObjectProperty("object", internal.NewObject()),
			)
			marshaled, err := obj.MarshalJSON()

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "{\"array\":[],\"boolean\":false,\"float\":1,\"int\":1,\"none\":null,\"object\":{},\"string\":\"1\"}")
		})
	})

	Convey(".Unwrap", t, func() {
		Convey("Should return an unwrapped items", func() {
			obj := internal.NewObjectWith(
				internal.NewObjectProperty("foo", core.NewString("foo")),
				internal.NewObjectProperty("bar", core.NewString("bar")),
			)

			for _, val := range obj.Unwrap().(map[string]interface{}) {
				So(val, ShouldHaveSameTypeAs, "")
			}
		})
	})

	Convey(".String", t, func() {
		Convey("Should return a string representation ", func() {
			obj := internal.NewObjectWith(
				internal.NewObjectProperty("foo", core.NewString("bar")),
			)

			So(obj.String(), ShouldEqual, "{\"foo\":\"bar\"}")
		})
	})

	Convey(".CompareValues", t, func() {
		Convey("It should return 1 for all non-object values", func() {
			arr := []core.Value{
				core.None,
				core.False,
				core.NewInt(1),
				core.Float(1),
				core.NewString("1"),
				internal.NewArray(10),
			}
			obj := internal.NewObject()

			for _, val := range arr {
				So(obj.Compare(val), ShouldEqual, 1)
			}
		})

		Convey("It should return -1 for all object values", func() {
			arr := internal.NewArrayWith(core.ZeroInt, core.ZeroInt)
			obj := internal.NewObject()

			So(arr.Compare(obj), ShouldEqual, -1)
		})

		Convey("It should return 0 when both objects are empty", func() {
			obj1 := internal.NewObject()
			obj2 := internal.NewObject()

			So(obj1.Compare(obj2), ShouldEqual, 0)
		})

		Convey("It should return 0 when both objects are equal (independent of key order)", func() {
			obj1 := internal.NewObjectWith(
				internal.NewObjectProperty("foo", core.NewString("foo")),
				internal.NewObjectProperty("bar", core.NewString("bar")),
			)
			obj2 := internal.NewObjectWith(
				internal.NewObjectProperty("foo", core.NewString("foo")),
				internal.NewObjectProperty("bar", core.NewString("bar")),
			)

			So(obj1.Compare(obj1), ShouldEqual, 0)
			So(obj2.Compare(obj2), ShouldEqual, 0)
			So(obj1.Compare(obj2), ShouldEqual, 0)
			So(obj2.Compare(obj1), ShouldEqual, 0)
		})

		Convey("It should return 1 when other array is empty", func() {
			obj1 := internal.NewObjectWith(internal.NewObjectProperty("foo", core.NewString("bar")))
			obj2 := internal.NewObject()

			So(obj1.Compare(obj2), ShouldEqual, 1)
		})

		Convey("It should return 1 when values are bigger", func() {
			obj1 := internal.NewObjectWith(internal.NewObjectProperty("foo", core.NewFloat(3)))
			obj2 := internal.NewObjectWith(internal.NewObjectProperty("foo", core.NewFloat(2)))

			So(obj1.Compare(obj2), ShouldEqual, 1)
		})

		Convey("It should return 1 when values are less", func() {
			obj1 := internal.NewObjectWith(internal.NewObjectProperty("foo", core.NewFloat(1)))
			obj2 := internal.NewObjectWith(internal.NewObjectProperty("foo", core.NewFloat(2)))

			So(obj1.Compare(obj2), ShouldEqual, -1)
		})

		Convey("ArangoDB compatibility", func() {
			Convey("It should return 1 when {a:1} and {b:2}", func() {
				obj1 := internal.NewObjectWith(internal.NewObjectProperty("a", core.NewInt(1)))
				obj2 := internal.NewObjectWith(internal.NewObjectProperty("b", core.NewInt(2)))

				So(obj1.Compare(obj2), ShouldEqual, 1)
			})

			Convey("It should return 0 when {a:1} and {a:1}", func() {
				obj1 := internal.NewObjectWith(internal.NewObjectProperty("a", core.NewInt(1)))
				obj2 := internal.NewObjectWith(internal.NewObjectProperty("a", core.NewInt(1)))

				So(obj1.Compare(obj2), ShouldEqual, 0)
			})

			Convey("It should return 0 {a:1, c:2} and {c:2, a:1}", func() {
				obj1 := internal.NewObjectWith(
					internal.NewObjectProperty("a", core.NewInt(1)),
					internal.NewObjectProperty("c", core.NewInt(2)),
				)
				obj2 := internal.NewObjectWith(
					internal.NewObjectProperty("c", core.NewInt(2)),
					internal.NewObjectProperty("a", core.NewInt(1)),
				)

				So(obj1.Compare(obj2), ShouldEqual, 0)
			})

			Convey("It should return -1 when {a:1} and {a:2}", func() {
				obj1 := internal.NewObjectWith(internal.NewObjectProperty("a", core.NewInt(1)))
				obj2 := internal.NewObjectWith(internal.NewObjectProperty("a", core.NewInt(2)))

				So(obj1.Compare(obj2), ShouldEqual, -1)
			})

			Convey("It should return 1 when {a:1, c:2} and {c:2, b:2}", func() {
				obj1 := internal.NewObjectWith(
					internal.NewObjectProperty("a", core.NewInt(1)),
					internal.NewObjectProperty("c", core.NewInt(2)),
				)
				obj2 := internal.NewObjectWith(
					internal.NewObjectProperty("c", core.NewInt(2)),
					internal.NewObjectProperty("b", core.NewInt(2)),
				)

				So(obj1.Compare(obj2), ShouldEqual, 1)
			})

			Convey("It should return 1 {a:1, c:3} and {c:2, a:1}", func() {
				obj1 := internal.NewObjectWith(
					internal.NewObjectProperty("a", core.NewInt(1)),
					internal.NewObjectProperty("c", core.NewInt(3)),
				)
				obj2 := internal.NewObjectWith(
					internal.NewObjectProperty("c", core.NewInt(2)),
					internal.NewObjectProperty("a", core.NewInt(1)),
				)

				So(obj1.Compare(obj2), ShouldEqual, 1)
			})
		})
	})

	Convey(".Hash", t, func() {
		Convey("It should calculate hash of non-empty object", func() {
			v := internal.NewObjectWith(
				internal.NewObjectProperty("foo", core.NewString("bar")),
				internal.NewObjectProperty("faz", core.NewInt(1)),
				internal.NewObjectProperty("qaz", core.True),
			)

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)
		})

		Convey("It should calculate hash of empty object", func() {
			v := internal.NewObject()

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)
		})

		Convey("Hash sum should be consistent", func() {
			v := internal.NewObjectWith(
				internal.NewObjectProperty("boolean", core.True),
				internal.NewObjectProperty("int", core.NewInt(1)),
				internal.NewObjectProperty("float", core.NewFloat(1.1)),
				internal.NewObjectProperty("string", core.NewString("foobar")),
				internal.NewObjectProperty("datetime", core.NewCurrentDateTime()),
				internal.NewObjectProperty("array", internal.NewArrayWith(core.NewInt(1), core.True)),
				internal.NewObjectProperty("object", internal.NewObjectWith(internal.NewObjectProperty("foo", core.NewString("bar")))),
			)

			h1 := v.Hash()
			h2 := v.Hash()

			So(h1, ShouldEqual, h2)
		})
	})

	Convey(".Length", t, func() {
		Convey("Should return 0 when empty", func() {
			obj := internal.NewObject()

			So(obj.Length(), ShouldEqual, 0)
		})

		Convey("Should return greater than 0 when not empty", func() {
			obj := internal.NewObjectWith(
				internal.NewObjectProperty("foo", core.ZeroInt),
				internal.NewObjectProperty("bar", core.ZeroInt),
			)

			So(obj.Length(), ShouldEqual, 2)
		})
	})

	Convey(".ForEach", t, func() {
		Convey("Should iterate over elements", func() {
			obj := internal.NewObjectWith(
				internal.NewObjectProperty("foo", core.ZeroInt),
				internal.NewObjectProperty("bar", core.ZeroInt),
			)
			counter := 0

			obj.ForEach(func(value core.Value, key string) bool {
				counter++

				return true
			})

			So(counter, ShouldEqual, obj.Length())
		})

		Convey("Should not iterate when empty", func() {
			obj := internal.NewObject()
			counter := 0

			obj.ForEach(func(value core.Value, key string) bool {
				counter++

				return true
			})

			So(counter, ShouldEqual, obj.Length())
		})

		Convey("Should break iteration when false returned", func() {
			obj := internal.NewObjectWith(
				internal.NewObjectProperty("1", core.NewInt(1)),
				internal.NewObjectProperty("2", core.NewInt(2)),
				internal.NewObjectProperty("3", core.NewInt(3)),
				internal.NewObjectProperty("4", core.NewInt(4)),
				internal.NewObjectProperty("5", core.NewInt(5)),
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
		//Convey("Should return item by key", func() {
		//	obj := values.NewObjectWith(
		//		values.NewObjectProperty("foo", values.NewInt(1)),
		//		values.NewObjectProperty("bar", values.NewInt(2)),
		//		values.NewObjectProperty("qaz", values.NewInt(3)),
		//	)
		//
		//	el, _ := obj.Get("foo")
		//
		//	So(el.CompareValues(values.NewInt(1)), ShouldEqual, 0)
		//})

		//Convey("Should return None when no items", func() {
		//	obj := values.NewObject()
		//
		//	el, _ := obj.Get("foo")
		//
		//	So(el.CompareValues(values.None), ShouldEqual, 0)
		//})
	})

	Convey(".Set", t, func() {
		//Convey("Should set item by index", func() {
		//	obj := values.NewObject()
		//
		//	obj.Set("foo", values.NewInt(1))
		//
		//	So(obj.Length(), ShouldEqual, 1)
		//
		//	v, _ := obj.Get("foo")
		//	So(v.CompareValues(values.NewInt(1)), ShouldEqual, 0)
		//})
	})

	Convey(".Clone", t, func() {
		Convey("Cloned object should be equal to source object", func() {
			obj := internal.NewObjectWith(
				internal.NewObjectProperty("one", core.NewInt(1)),
				internal.NewObjectProperty("two", core.NewInt(2)),
			)

			clone := obj.Clone().(*internal.Object)

			So(obj.Compare(clone), ShouldEqual, 0)
		})

		Convey("Cloned object should be independent of the source object", func() {
			obj := internal.NewObjectWith(
				internal.NewObjectProperty("one", core.NewInt(1)),
				internal.NewObjectProperty("two", core.NewInt(2)),
			)

			clone := obj.Clone().(*internal.Object)

			obj.Remove(core.NewString("one"))

			So(obj.Compare(clone), ShouldNotEqual, 0)
		})

		Convey("Cloned object must contain copies of the nested objects", func() {
			obj := internal.NewObjectWith(
				internal.NewObjectProperty(
					"arr", internal.NewArrayWith(core.NewInt(1)),
				),
			)

			clone := obj.Clone().(*internal.Object)

			nestedInObj, _ := obj.Get(core.NewString("arr"))
			nestedInObjArr := nestedInObj.(*internal.Array)
			nestedInObjArr.Push(core.NewInt(2))

			nestedInClone, _ := clone.Get(core.NewString("arr"))
			nestedInCloneArr := nestedInClone.(*internal.Array)

			So(nestedInObjArr.Compare(nestedInCloneArr), ShouldNotEqual, 0)
		})
	})
}
