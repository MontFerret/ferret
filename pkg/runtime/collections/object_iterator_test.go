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
			set, err := iter.Next()

			So(err, ShouldBeNil)
			So(set, ShouldHaveLength, 2)

			expected, exists := m.Get(values.NewString(set[1].String()))

			So(bool(exists), ShouldBeTrue)
			So(expected, ShouldEqual, set[0])

			res = append(res, set[0])
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
			set, err := iter.Next()

			So(err, ShouldBeNil)
			So(set, ShouldHaveLength, 2)

			res = append(res, set[0])
		}

		set, err := iter.Next()

		So(set, ShouldBeNil)
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
