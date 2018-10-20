package collections_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestObjectIterator(t *testing.T) {
	Convey("Should iterate over a map", t, func() {
		m := values.NewObjectWith(
			values.NewObjectProperty("one", values.NewInt(1)),
			values.NewObjectProperty("two", values.NewInt(2)),
			values.NewObjectProperty("three", values.NewInt(3)),
			values.NewObjectProperty("four", values.NewInt(4)),
			values.NewObjectProperty("five", values.NewInt(5)),
		)

		iter := collections.NewObjectIterator(m)

		res := make([]core.Value, 0, m.Length())

		for iter.HasNext() {
			item, key, err := iter.Next()

			So(err, ShouldBeNil)

			expected, exists := m.Get(values.NewString(key.String()))

			So(bool(exists), ShouldBeTrue)
			So(expected, ShouldEqual, item)

			res = append(res, item)
		}

		So(res, ShouldHaveLength, m.Length())
	})

	Convey("Should return an error when exhausted", t, func() {
		m := values.NewObjectWith(
			values.NewObjectProperty("one", values.NewInt(1)),
			values.NewObjectProperty("two", values.NewInt(2)),
			values.NewObjectProperty("three", values.NewInt(3)),
			values.NewObjectProperty("four", values.NewInt(4)),
			values.NewObjectProperty("five", values.NewInt(5)),
		)

		iter := collections.NewObjectIterator(m)

		res := make([]core.Value, 0, m.Length())

		for iter.HasNext() {
			item, _, err := iter.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		item, _, err := iter.Next()

		So(item, ShouldEqual, values.None)
		So(err, ShouldBeError)
	})

	Convey("Should NOT iterate over a empty map", t, func() {
		m := values.NewObject()

		iter := collections.NewObjectIterator(m)

		var iterated bool

		for iter.HasNext() {
			iterated = true
		}

		So(iterated, ShouldBeFalse)
	})
}
