package collections_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMapIterator(t *testing.T) {
	Convey("Should iterate over a map", t, func() {
		m := map[string]core.Value{
			"one":   values.NewInt(1),
			"two":   values.NewInt(2),
			"three": values.NewInt(3),
			"four":  values.NewInt(4),
			"five":  values.NewInt(5),
		}

		iter := collections.NewMapIterator(m)

		res := make([]core.Value, 0, len(m))

		for iter.HasNext() {
			set, err := iter.Next()

			So(err, ShouldBeNil)
			So(set, ShouldHaveLength, 2)

			expected, exists := m[set[1].String()]

			So(exists, ShouldBeTrue)
			So(expected, ShouldEqual, set[0])

			res = append(res, set[0])
		}

		So(res, ShouldHaveLength, len(m))
	})

	Convey("Should return an error when exhausted", t, func() {
		m := map[string]core.Value{
			"one":   values.NewInt(1),
			"two":   values.NewInt(2),
			"three": values.NewInt(3),
			"four":  values.NewInt(4),
			"five":  values.NewInt(5),
		}

		iter := collections.NewMapIterator(m)

		res := make([]core.Value, 0, len(m))

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
		m := make(map[string]core.Value)

		iter := collections.NewMapIterator(m)

		var iterated bool

		for iter.HasNext() {
			iterated = true
		}

		So(iterated, ShouldBeFalse)
	})
}
