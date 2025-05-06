package runtime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

func TestObject(t *testing.T) {
	Convey("#constructor", t, func() {
		Convey("Should create an empty object", func() {
			obj := runtime.NewObject()

			So(obj.Length(), ShouldEqual, 0)
		})

		Convey("Should create an object, from passed values", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("none", None),
				runtime.NewObjectProperty("boolean", False),
				runtime.NewObjectProperty("int", NewInt(1)),
				runtime.NewObjectProperty("float", Float(1)),
				runtime.NewObjectProperty("string", NewString("1")),
				runtime.NewObjectProperty("array", NewArray(10)),
				runtime.NewObjectProperty("object", runtime.NewObject()),
			)

			So(obj.Length(), ShouldEqual, 7)
		})
	})

	Convey(".MarshalJSON", t, func() {
		Convey("Should serialize an empty object", func() {
			obj := runtime.NewObject()
			marshaled, err := obj.MarshalJSON()

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "{}")
		})

		Convey("Should serialize full object", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("none", None),
				runtime.NewObjectProperty("boolean", False),
				runtime.NewObjectProperty("int", NewInt(1)),
				runtime.NewObjectProperty("float", Float(1)),
				runtime.NewObjectProperty("string", NewString("1")),
				runtime.NewObjectProperty("array", NewArray(10)),
				runtime.NewObjectProperty("object", runtime.NewObject()),
			)
			marshaled, err := obj.MarshalJSON()

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "{\"array\":[],\"boolean\":false,\"float\":1,\"int\":1,\"none\":null,\"object\":{},\"string\":\"1\"}")
		})
	})

	Convey(".Unwrap", t, func() {
		Convey("Should return an unwrapped items", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("foo", NewString("foo")),
				runtime.NewObjectProperty("bar", NewString("bar")),
			)

			for _, val := range obj.Unwrap().(map[string]interface{}) {
				So(val, ShouldHaveSameTypeAs, "")
			}
		})
	})

	Convey(".String", t, func() {
		Convey("Should return a string representation ", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("foo", NewString("bar")),
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
			obj := runtime.NewObject()

			for _, val := range arr {
				So(obj.Compare(val), ShouldEqual, 1)
			}
		})

		Convey("It should return -1 for all object values", func() {
			arr := NewArrayWith(ZeroInt, ZeroInt)
			obj := runtime.NewObject()

			So(arr.Compare(obj), ShouldEqual, -1)
		})

		Convey("It should return 0 when both objects are empty", func() {
			obj1 := runtime.NewObject()
			obj2 := runtime.NewObject()

			So(obj1.Compare(obj2), ShouldEqual, 0)
		})

		Convey("It should return 0 when both objects are equal (independent of key order)", func() {
			obj1 := runtime.NewObjectWith(
				runtime.NewObjectProperty("foo", NewString("foo")),
				runtime.NewObjectProperty("bar", NewString("bar")),
			)
			obj2 := runtime.NewObjectWith(
				runtime.NewObjectProperty("foo", NewString("foo")),
				runtime.NewObjectProperty("bar", NewString("bar")),
			)

			So(obj1.Compare(obj1), ShouldEqual, 0)
			So(obj2.Compare(obj2), ShouldEqual, 0)
			So(obj1.Compare(obj2), ShouldEqual, 0)
			So(obj2.Compare(obj1), ShouldEqual, 0)
		})

		Convey("It should return 1 when other array is empty", func() {
			obj1 := runtime.NewObjectWith(runtime.NewObjectProperty("foo", NewString("bar")))
			obj2 := runtime.NewObject()

			So(obj1.Compare(obj2), ShouldEqual, 1)
		})

		Convey("It should return 1 when values are bigger", func() {
			obj1 := runtime.NewObjectWith(runtime.NewObjectProperty("foo", NewFloat(3)))
			obj2 := runtime.NewObjectWith(runtime.NewObjectProperty("foo", NewFloat(2)))

			So(obj1.Compare(obj2), ShouldEqual, 1)
		})

		Convey("It should return 1 when values are less", func() {
			obj1 := runtime.NewObjectWith(runtime.NewObjectProperty("foo", NewFloat(1)))
			obj2 := runtime.NewObjectWith(runtime.NewObjectProperty("foo", NewFloat(2)))

			So(obj1.Compare(obj2), ShouldEqual, -1)
		})

		Convey("ArangoDB compatibility", func() {
			Convey("It should return 1 when {a:1} and {b:2}", func() {
				obj1 := runtime.NewObjectWith(runtime.NewObjectProperty("a", NewInt(1)))
				obj2 := runtime.NewObjectWith(runtime.NewObjectProperty("b", NewInt(2)))

				So(obj1.Compare(obj2), ShouldEqual, 1)
			})

			Convey("It should return 0 when {a:1} and {a:1}", func() {
				obj1 := runtime.NewObjectWith(runtime.NewObjectProperty("a", NewInt(1)))
				obj2 := runtime.NewObjectWith(runtime.NewObjectProperty("a", NewInt(1)))

				So(obj1.Compare(obj2), ShouldEqual, 0)
			})

			Convey("It should return 0 {a:1, c:2} and {c:2, a:1}", func() {
				obj1 := runtime.NewObjectWith(
					runtime.NewObjectProperty("a", NewInt(1)),
					runtime.NewObjectProperty("c", NewInt(2)),
				)
				obj2 := runtime.NewObjectWith(
					runtime.NewObjectProperty("c", NewInt(2)),
					runtime.NewObjectProperty("a", NewInt(1)),
				)

				So(obj1.Compare(obj2), ShouldEqual, 0)
			})

			Convey("It should return -1 when {a:1} and {a:2}", func() {
				obj1 := runtime.NewObjectWith(runtime.NewObjectProperty("a", NewInt(1)))
				obj2 := runtime.NewObjectWith(runtime.NewObjectProperty("a", NewInt(2)))

				So(obj1.Compare(obj2), ShouldEqual, -1)
			})

			Convey("It should return 1 when {a:1, c:2} and {c:2, b:2}", func() {
				obj1 := runtime.NewObjectWith(
					runtime.NewObjectProperty("a", NewInt(1)),
					runtime.NewObjectProperty("c", NewInt(2)),
				)
				obj2 := runtime.NewObjectWith(
					runtime.NewObjectProperty("c", NewInt(2)),
					runtime.NewObjectProperty("b", NewInt(2)),
				)

				So(obj1.Compare(obj2), ShouldEqual, 1)
			})

			Convey("It should return 1 {a:1, c:3} and {c:2, a:1}", func() {
				obj1 := runtime.NewObjectWith(
					runtime.NewObjectProperty("a", NewInt(1)),
					runtime.NewObjectProperty("c", NewInt(3)),
				)
				obj2 := runtime.NewObjectWith(
					runtime.NewObjectProperty("c", NewInt(2)),
					runtime.NewObjectProperty("a", NewInt(1)),
				)

				So(obj1.Compare(obj2), ShouldEqual, 1)
			})
		})
	})

	Convey(".Hash", t, func() {
		Convey("It should calculate hash of non-empty object", func() {
			v := runtime.NewObjectWith(
				runtime.NewObjectProperty("foo", NewString("bar")),
				runtime.NewObjectProperty("faz", NewInt(1)),
				runtime.NewObjectProperty("qaz", True),
			)

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)
		})

		Convey("It should calculate hash of empty object", func() {
			v := runtime.NewObject()

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)
		})

		Convey("Hash sum should be consistent", func() {
			v := runtime.NewObjectWith(
				runtime.NewObjectProperty("boolean", True),
				runtime.NewObjectProperty("int", NewInt(1)),
				runtime.NewObjectProperty("float", NewFloat(1.1)),
				runtime.NewObjectProperty("string", NewString("foobar")),
				runtime.NewObjectProperty("datetime", NewCurrentDateTime()),
				runtime.NewObjectProperty("array", NewArrayWith(NewInt(1), True)),
				runtime.NewObjectProperty("object", runtime.NewObjectWith(runtime.NewObjectProperty("foo", NewString("bar")))),
			)

			h1 := v.Hash()
			h2 := v.Hash()

			So(h1, ShouldEqual, h2)
		})
	})

	Convey(".Length", t, func() {
		Convey("Should return 0 when empty", func() {
			obj := runtime.NewObject()

			So(obj.Length(), ShouldEqual, 0)
		})

		Convey("Should return greater than 0 when not empty", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("foo", ZeroInt),
				runtime.NewObjectProperty("bar", ZeroInt),
			)

			So(obj.Length(), ShouldEqual, 2)
		})
	})

	Convey(".ForEach", t, func() {
		Convey("Should iterate over elements", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("foo", ZeroInt),
				runtime.NewObjectProperty("bar", ZeroInt),
			)
			counter := 0

			obj.ForEach(func(value Value, key string) bool {
				counter++

				return true
			})

			So(counter, ShouldEqual, obj.Length())
		})

		Convey("Should not iterate when empty", func() {
			obj := runtime.NewObject()
			counter := 0

			obj.ForEach(func(value Value, key string) bool {
				counter++

				return true
			})

			So(counter, ShouldEqual, obj.Length())
		})

		Convey("Should break iteration when false returned", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("1", NewInt(1)),
				runtime.NewObjectProperty("2", NewInt(2)),
				runtime.NewObjectProperty("3", NewInt(3)),
				runtime.NewObjectProperty("4", NewInt(4)),
				runtime.NewObjectProperty("5", NewInt(5)),
			)
			threshold := 3
			counter := 0

			obj.ForEach(func(value Value, key string) bool {
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
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("one", NewInt(1)),
				runtime.NewObjectProperty("two", NewInt(2)),
			)

			clone := obj.Clone().(*runtime.Object)

			So(obj.Compare(clone), ShouldEqual, 0)
		})

		Convey("Cloned object should be independent of the source object", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("one", NewInt(1)),
				runtime.NewObjectProperty("two", NewInt(2)),
			)

			clone := obj.Clone().(*runtime.Object)

			obj.Remove(NewString("one"))

			So(obj.Compare(clone), ShouldNotEqual, 0)
		})

		Convey("Cloned object must contain copies of the nested objects", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty(
					"arr", NewArrayWith(NewInt(1)),
				),
			)

			clone := obj.Clone().(*runtime.Object)

			nestedInObj, _ := obj.Get(NewString("arr"))
			nestedInObjArr := nestedInObj.(*Array)
			nestedInObjArr.Push(NewInt(2))

			nestedInClone, _ := clone.Get(NewString("arr"))
			nestedInCloneArr := nestedInClone.(*Array)

			So(nestedInObjArr.Compare(nestedInCloneArr), ShouldNotEqual, 0)
		})
	})
}
