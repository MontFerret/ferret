package runtime_test

import (
	"encoding/json"
	"testing"

	c "context"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArray(t *testing.T) {
	ctx := c.Background()
	Convey("#constructor", t, func() {
		Convey("Should create an empty array", func() {
			arr := runtime.NewArray(10)
			size, _ := arr.Length(ctx)

			So(size, ShouldEqual, 0)
		})

		Convey("Should create an array, from passed values", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			size, _ := arr.Length(ctx)

			So(size, ShouldEqual, 3)
		})
	})

	Convey(".MarshalJSON", t, func() {
		Convey("Should serialize empty array", func() {
			arr := runtime.NewArray(10)
			marshaled, err := arr.MarshalJSON()

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "[]")
		})

		Convey("Should serialize full array", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			marshaled, err := json.Marshal(arr)

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "[1,2,3]")
		})
	})

	Convey(".Unwrap", t, func() {
		Convey("Should return a an array of unwrapped values", func() {
			arr := runtime.NewArrayWith(
				runtime.ZeroInt,
				runtime.ZeroInt,
			)

			for _, val := range arr.Unwrap().([]interface{}) {
				So(val, ShouldHaveSameTypeAs, 0)
			}
		})
	})

	Convey(".String", t, func() {
		Convey("Should return a string representation ", func() {
			arr := runtime.NewArrayWith(runtime.ZeroInt, runtime.ZeroInt)

			So(arr.String(), ShouldEqual, "[0,0]")
		})
	})

	Convey(".Compareruntime.Values", t, func() {
		Convey("It should return 1 for all non-array and non-object values", func() {
			arr := runtime.NewArrayWith(runtime.ZeroInt, runtime.ZeroInt)

			So(arr.Compare(runtime.None), ShouldEqual, 1)
			So(arr.Compare(runtime.ZeroInt), ShouldEqual, 1)
			So(arr.Compare(runtime.ZeroFloat), ShouldEqual, 1)
			So(arr.Compare(runtime.EmptyString), ShouldEqual, 1)
		})

		Convey("It should return -1 for all object values", func() {
			arr := runtime.NewArrayWith(runtime.ZeroInt, runtime.ZeroInt)
			obj := runtime.NewObject()

			So(arr.Compare(obj), ShouldEqual, -1)
		})

		Convey("It should return 0 when both arrays are empty", func() {
			arr1 := runtime.NewArray(1)
			arr2 := runtime.NewArray(1)

			So(arr1.Compare(arr2), ShouldEqual, 0)
		})

		Convey("It should return 1 when other array is empty", func() {
			arr1 := runtime.NewArrayWith(runtime.ZeroFloat)
			arr2 := runtime.NewArray(1)

			So(arr1.Compare(arr2), ShouldEqual, 1)
		})

		Convey("It should return 1 when values are bigger", func() {
			arr1 := runtime.NewArrayWith(runtime.NewInt(1))
			arr2 := runtime.NewArrayWith(runtime.ZeroInt)

			So(arr1.Compare(arr2), ShouldEqual, 1)
		})

		Convey("It should return 0 when arrays are equal", func() {
			Convey("When only simple types are nested", func() {
				arr1 := runtime.NewArrayWith(
					runtime.NewInt(0), runtime.NewString("str"),
				)
				arr2 := runtime.NewArrayWith(
					runtime.NewInt(0), runtime.NewString("str"),
				)

				So(arr1.Compare(arr2), ShouldEqual, 0)
			})

			Convey("When object and array are nested at the same time", func() {
				arr1 := runtime.NewArrayWith(
					runtime.NewObjectWith(
						runtime.NewObjectProperty("one", runtime.NewInt(1)),
					),
					runtime.NewArrayWith(
						runtime.NewInt(2),
					),
				)
				arr2 := runtime.NewArrayWith(
					runtime.NewObjectWith(
						runtime.NewObjectProperty("one", runtime.NewInt(1)),
					),
					runtime.NewArrayWith(
						runtime.NewInt(2),
					),
				)

				So(arr1.Compare(arr2), ShouldEqual, 0)
			})

			Convey("When only objects are nested", func() {
				arr1 := runtime.NewArrayWith(
					runtime.NewObjectWith(
						runtime.NewObjectProperty("one", runtime.NewInt(1)),
					),
				)
				arr2 := runtime.NewArrayWith(
					runtime.NewObjectWith(
						runtime.NewObjectProperty("one", runtime.NewInt(1)),
					),
				)

				So(arr1.Compare(arr2), ShouldEqual, 0)
			})

			Convey("When only arrays are nested", func() {
				arr1 := runtime.NewArrayWith(
					runtime.NewArrayWith(
						runtime.NewInt(2),
					),
				)
				arr2 := runtime.NewArrayWith(
					runtime.NewArrayWith(
						runtime.NewInt(2),
					),
				)

				So(arr1.Compare(arr2), ShouldEqual, 0)
			})

			Convey("When simple and complex types at the same time", func() {
				arr1 := runtime.NewArrayWith(
					runtime.NewInt(0),
					runtime.NewObjectWith(
						runtime.NewObjectProperty("one", runtime.NewInt(1)),
					),
					runtime.NewArrayWith(
						runtime.NewInt(2),
					),
				)
				arr2 := runtime.NewArrayWith(
					runtime.NewInt(0),
					runtime.NewObjectWith(
						runtime.NewObjectProperty("one", runtime.NewInt(1)),
					),
					runtime.NewArrayWith(
						runtime.NewInt(2),
					),
				)

				So(arr1.Compare(arr2), ShouldEqual, 0)
			})

			Convey("When custom complex type", func() {
				arr1 := runtime.NewArrayWith(
					runtime.NewObjectWith(
						runtime.NewObjectProperty(
							"arr", runtime.NewArrayWith(runtime.NewObject()),
						),
					),
				)
				arr2 := runtime.NewArrayWith(
					runtime.NewObjectWith(
						runtime.NewObjectProperty(
							"arr", runtime.NewArrayWith(runtime.NewObject()),
						),
					),
				)

				So(arr1.Compare(arr2), ShouldEqual, 0)
			})
		})
	})

	Convey(".Hash", t, func() {
		Convey("It should calculate hash of non-empty array", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)

			h := arr.Hash()

			So(h, ShouldBeGreaterThan, 0)
		})

		Convey("It should calculate hash of empty array", func() {
			arr := runtime.NewArrayWith()

			h := arr.Hash()

			So(h, ShouldBeGreaterThan, 0)
		})

		Convey("Hash sum should be consistent", func() {
			arr := runtime.NewArrayWith(
				runtime.True,
				runtime.NewInt(1),
				runtime.NewFloat(1.1),
				runtime.NewString("foobar"),
				runtime.NewCurrentDateTime(),
				runtime.NewArrayWith(runtime.NewInt(1), runtime.True),
				runtime.NewObjectWith(runtime.NewObjectProperty("foo", runtime.NewString("bar"))),
			)

			h1 := arr.Hash()
			h2 := arr.Hash()

			So(h1, ShouldEqual, h2)
		})
	})

	Convey(".Length", t, func() {
		Convey("Should return 0 when empty", func() {
			arr := runtime.NewArray(1)

			size, _ := arr.Length(ctx)

			So(size, ShouldEqual, 0)
		})

		Convey("Should return greater than 0 when not empty", func() {
			arr := runtime.NewArrayWith(runtime.ZeroInt, runtime.ZeroInt)

			size, _ := arr.Length(ctx)

			So(size, ShouldEqual, 2)
		})
	})

	Convey(".ForEachIter", t, func() {
		Convey("Should iterate over elements", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			counter := 0

			_ = arr.ForEach(ctx, func(_ c.Context, val runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
				counter++

				return true, nil
			})

			size, _ := arr.Length(ctx)

			So(counter, ShouldEqual, size)
		})

		Convey("Should not iterate when empty", func() {
			arr := runtime.NewArrayWith()
			counter := 0

			_ = arr.ForEach(ctx, func(_ c.Context, val runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
				counter++

				return true, nil
			})

			size, _ := arr.Length(ctx)

			So(counter, ShouldEqual, size)
		})

		//Convey("Should break iteration when false returned", func() {
		//	arr := values.NewArrayWith(
		//		values.runtime.NewInt(1),
		//		values.runtime.NewInt(2),
		//		values.runtime.NewInt(3),
		//		values.runtime.NewInt(4),
		//		values.runtime.NewInt(5),
		//	)
		//	threshold := 3
		//	counter := 0
		//
		//	arr.ForEachIter(func(value core.runtime.Value, idx int) bool {
		//		counter++
		//
		//		return value.Compareruntime.Values(values.runtime.NewInt(threshold)) == -1
		//	})
		//
		//	So(counter, ShouldEqual, threshold)
		//})
	})

	//Convey(".Get", t, func() {
	//	Convey("Should return item by index", func() {
	//		arr := values.NewArrayWith(
	//			values.runtime.NewInt(1),
	//			values.runtime.NewInt(2),
	//			values.runtime.NewInt(3),
	//			values.runtime.NewInt(4),
	//			values.runtime.NewInt(5),
	//		)
	//
	//		el := arr.Get(1)
	//
	//		So(el.Compareruntime.Values(values.runtime.NewInt(2)), ShouldEqual, 0)
	//	})
	//
	//	Convey("Should return runtime.None when no items", func() {
	//		arr := values.NewArrayWith()
	//
	//		el := arr.Get(1)
	//
	//		So(el.Compareruntime.Values(values.runtime.None), ShouldEqual, 0)
	//	})
	//})

	Convey(".Set", t, func() {
		//Convey("Should set item by index", func() {
		//	arr := values.NewArrayWith(values.runtime.ZeroInt)
		//
		//	err := arr.Set(0, values.runtime.NewInt(1))
		//
		//	So(err, ShouldBeNil)
		//	So(arr.Length(), ShouldEqual, 1)
		//	So(arr.Get(0).CompareValues(values.NewInt(1)), ShouldEqual, 0)
		//})

		Convey("Should return an error when index is out of bounds", func() {
			arr := runtime.NewArray(10)

			err := arr.Set(ctx, 0, runtime.NewInt(1))

			So(err, ShouldNotBeNil)
			size, _ := arr.Length(ctx)

			So(size, ShouldEqual, 0)
		})
	})

	Convey(".Push", t, func() {
		Convey("Should add an item", func() {
			arr := runtime.NewArray(10)

			src := []runtime.Value{
				runtime.ZeroInt,
				runtime.ZeroInt,
				runtime.ZeroInt,
				runtime.ZeroInt,
				runtime.ZeroInt,
			}

			for _, val := range src {
				arr.Add(ctx, val)
			}

			size, _ := arr.Length(ctx)

			So(size, ShouldEqual, len(src))
		})
	})

	//Convey(".Slice", t, func() {
	//	Convey("Should return a slice", func() {
	//		arr := values.NewArrayWith(
	//			values.runtime.NewInt(0),
	//			values.runtime.NewInt(1),
	//			values.runtime.NewInt(2),
	//			values.runtime.NewInt(3),
	//			values.runtime.NewInt(4),
	//			values.runtime.NewInt(5),
	//		)
	//
	//		s := arr.Slice(0, 1)
	//
	//		So(s.Length(ctx), ShouldEqual, 1)
	//		So(s.Get(0).Compareruntime.Values(values.runtime.ZeroInt), ShouldEqual, 0)
	//
	//		s2 := arr.Slice(2, arr.Length(ctx))
	//
	//		So(s2.Length(ctx), ShouldEqual, arr.Length(ctx)-2)
	//	})
	//})

	Convey(".Insert", t, func() {
		Convey("Should insert an item in the middle of an array", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(0),
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
				runtime.NewInt(4),
				runtime.NewInt(5),
			)

			lenBefore, _ := arr.Length(ctx)

			arr.Insert(ctx, 3, runtime.NewInt(100))

			lenAfter, _ := arr.Length(ctx)

			act, _ := arr.Get(ctx, 3)

			So(lenAfter, ShouldBeGreaterThan, lenBefore)
			So(act, ShouldEqual, 100)
		})
	})

	Convey(".RemoveAt", t, func() {
		Convey("Should remove an item from the middle", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(0),
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
				runtime.NewInt(4),
				runtime.NewInt(5),
			)

			lenBefore, _ := arr.Length(ctx)

			arr.RemoveAt(ctx, 3)

			lenAfter, _ := arr.Length(ctx)

			val, _ := arr.Get(ctx, 3)

			So(lenAfter, ShouldBeLessThan, lenBefore)
			So(val, ShouldEqual, 4)
		})

		Convey("Should remove an item from the end", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(0),
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
				runtime.NewInt(4),
				runtime.NewInt(5),
			)

			lenBefore, _ := arr.Length(ctx)

			arr.RemoveAt(ctx, 5)

			lenAfter, _ := arr.Length(ctx)

			val, _ := arr.Get(ctx, 4)

			So(lenAfter, ShouldBeLessThan, lenBefore)
			So(lenAfter, ShouldEqual, 5)
			So(val, ShouldEqual, 4)
		})

		Convey("Should remove an item from the beginning", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(0),
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
				runtime.NewInt(4),
				runtime.NewInt(5),
			)

			lenBefore, _ := arr.Length(ctx)

			arr.RemoveAt(ctx, 0)

			lenAfter, _ := arr.Length(ctx)

			val, _ := arr.Get(ctx, 0)

			So(lenAfter, ShouldBeLessThan, lenBefore)
			So(val, ShouldEqual, 1)
		})
	})

	Convey(".Clone", t, func() {
		Convey("Cloned array should be equal to source array", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(0),
				runtime.NewObjectWith(
					runtime.NewObjectProperty("one", runtime.NewInt(1)),
				),
				runtime.NewArrayWith(
					runtime.NewInt(2),
				),
			)

			cloned, _ := arr.Clone(ctx)
			clone := cloned.(*runtime.Array)

			size, _ := arr.Length(ctx)
			cloneSize, _ := clone.Length(ctx)

			So(size, ShouldEqual, cloneSize)
			So(arr.Compare(clone), ShouldEqual, 0)
		})

		Convey("Cloned array should be independent of the source array", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(0),
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
				runtime.NewInt(4),
				runtime.NewInt(5),
			)

			cloned, _ := arr.Clone(ctx)
			clone := cloned.(*runtime.Array)

			arr.Add(ctx, runtime.NewInt(6))

			size, _ := arr.Length(ctx)
			cloneSize, _ := clone.Length(ctx)

			So(size, ShouldNotEqual, cloneSize)
			So(arr.Compare(clone), ShouldNotEqual, 0)
		})

		//Convey("Cloned array must contain copies of the nested objects", func() {
		//	arr := values.NewArrayWith(
		//		values.NewArrayWith(
		//			values.runtime.NewInt(0),
		//			values.runtime.NewInt(1),
		//			values.runtime.NewInt(2),
		//			values.runtime.NewInt(3),
		//			values.runtime.NewInt(4),
		//		),
		//	)
		//
		//	clone := arr.Clone().(*values.Array)
		//
		//	nestedInArr := arr.Get(values.runtime.NewInt(0)).(*values.Array)
		//	nestedInArr.Add(ctx, values.runtime.NewInt(5))
		//
		//	nestedInClone := clone.Get(values.runtime.NewInt(0)).(*values.Array)
		//
		//	So(nestedInArr.Compareruntime.Values(nestedInClone), ShouldNotEqual, 0)
		//})

		Convey(".Contains", func() {
			Convey("Should return true when element exists", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
				)

				contains, err := arr.Contains(ctx, runtime.NewInt(1))
				So(err, ShouldBeNil)
				So(contains, ShouldEqual, runtime.True)

				contains, err = arr.Contains(ctx, runtime.NewString("test"))
				So(err, ShouldBeNil)
				So(contains, ShouldEqual, runtime.True)

				contains, err = arr.Contains(ctx, runtime.NewFloat(3.14))
				So(err, ShouldBeNil)
				So(contains, ShouldEqual, runtime.True)
			})

			Convey("Should return false when element does not exist", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
				)

				contains, err := arr.Contains(ctx, runtime.NewInt(2))
				So(err, ShouldBeNil)
				So(contains, ShouldEqual, runtime.False)

				contains, err = arr.Contains(ctx, runtime.NewString("other"))
				So(err, ShouldBeNil)
				So(contains, ShouldEqual, runtime.False)
			})

			Convey("Should return false for empty array", func() {
				arr := runtime.NewArray(0)

				contains, err := arr.Contains(ctx, runtime.NewInt(1))
				So(err, ShouldBeNil)
				So(contains, ShouldEqual, runtime.False)
			})
		})

		Convey(".IndexOf", func() {
			Convey("Should return correct index when element exists", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
				)

				index, err := arr.IndexOf(ctx, runtime.NewInt(1))
				So(err, ShouldBeNil)
				So(index, ShouldEqual, runtime.NewInt(0))

				index, err = arr.IndexOf(ctx, runtime.NewString("test"))
				So(err, ShouldBeNil)
				So(index, ShouldEqual, runtime.NewInt(1))

				index, err = arr.IndexOf(ctx, runtime.NewFloat(3.14))
				So(err, ShouldBeNil)
				So(index, ShouldEqual, runtime.NewInt(2))
			})

			Convey("Should return -1 when element does not exist", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
				)

				index, err := arr.IndexOf(ctx, runtime.NewInt(2))
				So(err, ShouldBeNil)
				So(index, ShouldEqual, runtime.NewInt(-1))

				index, err = arr.IndexOf(ctx, runtime.NewString("other"))
				So(err, ShouldBeNil)
				So(index, ShouldEqual, runtime.NewInt(-1))
			})

			Convey("Should return -1 for empty array", func() {
				arr := runtime.NewArray(0)

				index, err := arr.IndexOf(ctx, runtime.NewInt(1))
				So(err, ShouldBeNil)
				So(index, ShouldEqual, runtime.NewInt(-1))
			})
		})

		Convey(".First", func() {
			Convey("Should return first element", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
				)

				first, err := arr.First(ctx)
				So(err, ShouldBeNil)
				So(first, ShouldEqual, runtime.NewInt(1))
			})

			Convey("Should return None for empty array", func() {
				arr := runtime.NewArray(0)

				first, err := arr.First(ctx)
				So(err, ShouldBeNil)
				So(first, ShouldEqual, runtime.None)
			})
		})

		Convey(".Last", func() {
			Convey("Should return last element", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
				)

				last, err := arr.Last(ctx)
				So(err, ShouldBeNil)
				So(last, ShouldEqual, runtime.NewFloat(3.14))
			})

			Convey("Should return None for empty array", func() {
				arr := runtime.NewArray(0)

				last, err := arr.Last(ctx)
				So(err, ShouldBeNil)
				So(last, ShouldEqual, runtime.None)
			})
		})

		Convey(".IsEmpty", func() {
			Convey("Should return true for empty array", func() {
				arr := runtime.NewArray(0)

				isEmpty, err := arr.IsEmpty(ctx)
				So(err, ShouldBeNil)
				So(isEmpty, ShouldEqual, runtime.True)
			})

			Convey("Should return false for non-empty array", func() {
				arr := runtime.NewArrayWith(runtime.NewInt(1))

				isEmpty, err := arr.IsEmpty(ctx)
				So(err, ShouldBeNil)
				So(isEmpty, ShouldEqual, runtime.False)
			})
		})

		Convey(".Clear", func() {
			Convey("Should clear all elements", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
				)

				length, _ := arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(3))

				err := arr.Clear(ctx)
				So(err, ShouldBeNil)

				length, _ = arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(0))

				isEmpty, _ := arr.IsEmpty(ctx)
				So(isEmpty, ShouldEqual, runtime.True)
			})

			Convey("Should work on empty array", func() {
				arr := runtime.NewArray(0)

				err := arr.Clear(ctx)
				So(err, ShouldBeNil)

				length, _ := arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(0))
			})
		})

		Convey(".Copy", func() {
			Convey("Should create a shallow copy", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
				)

				copied := arr.Copy()
				copyArr := copied.(*runtime.Array)

				// Should be equal but different instances
				So(arr.Compare(copyArr), ShouldEqual, 0)

				// Modifying copy should not affect original
				copyArr.Add(ctx, runtime.NewInt(42))

				arrLength, _ := arr.Length(ctx)
				copyLength, _ := copyArr.Length(ctx)

				So(arrLength, ShouldEqual, runtime.NewInt(3))
				So(copyLength, ShouldEqual, runtime.NewInt(4))
			})

			Convey("Should copy empty array", func() {
				arr := runtime.NewArray(0)

				copied := arr.Copy()
				copyArr := copied.(*runtime.Array)

				So(arr.Compare(copyArr), ShouldEqual, 0)
				
				isEmpty, _ := copyArr.IsEmpty(ctx)
				So(isEmpty, ShouldEqual, runtime.True)
			})
		})

		Convey(".Remove", func() {
			Convey("Should remove existing element", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
					runtime.NewString("test"),
				)

				err := arr.Remove(ctx, runtime.NewString("test"))
				So(err, ShouldBeNil)

				length, _ := arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(3))

				// Should remove first occurrence
				val0, _ := arr.Get(ctx, runtime.NewInt(0))
				So(val0, ShouldEqual, runtime.NewInt(1))

				val1, _ := arr.Get(ctx, runtime.NewInt(1))
				So(val1, ShouldEqual, runtime.NewFloat(3.14))

				val2, _ := arr.Get(ctx, runtime.NewInt(2))
				So(val2, ShouldEqual, runtime.NewString("test"))
			})

			Convey("Should return nil when element does not exist", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
				)

				err := arr.Remove(ctx, runtime.NewString("notfound"))
				So(err, ShouldBeNil)  // Returns nil even when not found

				length, _ := arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(2))
			})

			Convey("Should work with empty array", func() {
				arr := runtime.NewArray(0)

				err := arr.Remove(ctx, runtime.NewInt(1))
				So(err, ShouldBeNil)  // Returns nil even when not found

				length, _ := arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(0))
			})
		})

		Convey(".SortAsc", func() {
			Convey("Should sort integers in ascending order", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(3),
					runtime.NewInt(1),
					runtime.NewInt(4),
					runtime.NewInt(2),
				)

				err := arr.SortAsc(ctx)
				So(err, ShouldBeNil)

				val0, _ := arr.Get(ctx, runtime.NewInt(0))
				So(val0, ShouldEqual, runtime.NewInt(1))

				val1, _ := arr.Get(ctx, runtime.NewInt(1))
				So(val1, ShouldEqual, runtime.NewInt(2))

				val2, _ := arr.Get(ctx, runtime.NewInt(2))
				So(val2, ShouldEqual, runtime.NewInt(3))

				val3, _ := arr.Get(ctx, runtime.NewInt(3))
				So(val3, ShouldEqual, runtime.NewInt(4))
			})

			Convey("Should work with empty array", func() {
				arr := runtime.NewArray(0)

				err := arr.SortAsc(ctx)
				So(err, ShouldBeNil)

				length, _ := arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(0))
			})
		})
	})
}
