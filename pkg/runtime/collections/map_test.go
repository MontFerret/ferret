package collections_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func mapIterator(m map[string]core.Value) collections.Iterator {
	return collections.NewMapIterator(valVar, keyVar, m)
}

func TestMapIterator(t *testing.T) {
	Convey("Should iterate over a map", t, func() {
		m := map[string]core.Value{
			"one":   values.NewInt(1),
			"two":   values.NewInt(2),
			"three": values.NewInt(3),
			"four":  values.NewInt(4),
			"five":  values.NewInt(5),
		}

		iter := mapIterator(m)

		res := make([]core.Value, 0, len(m))

		for iter.HasNext() {
			item, key, err := next(iter)

			So(err, ShouldBeNil)

			expected, exists := m[key.String()]

			So(exists, ShouldBeTrue)
			So(expected, ShouldEqual, item)

			res = append(res, item)
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

		iter := mapIterator(m)

		res := make([]core.Value, 0, len(m))

		for iter.HasNext() {
			item, _, err := next(iter)

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		item, _, err := next(iter)

		So(item, ShouldBeNil)
		So(err, ShouldBeError)
	})

	Convey("Should NOT iterate over a empty map", t, func() {
		m := make(map[string]core.Value)

		iter := mapIterator(m)

		var iterated bool

		for iter.HasNext() {
			iterated = true
		}

		So(iterated, ShouldBeFalse)
	})
}
