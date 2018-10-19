package collections_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestResultSet(t *testing.T) {
	Convey(".Hash", t, func() {
		Convey("Should return same result for same value", func() {
			set := collections.ResultSet{values.NewInt(1), values.NewInt(0)}
			expected := set.Hash()

			for i := 0; i < 200; i++ {
				actual := set.Hash()

				So(actual, ShouldEqual, expected)
			}
		})

		Convey("Should return different result for different values", func() {
			set1 := collections.ResultSet{values.NewInt(1), values.NewInt(0)}
			set2 := collections.ResultSet{values.NewInt(0), values.NewInt(1)}

			So(set1.Hash(), ShouldNotEqual, set2.Hash())
		})
	})
}
