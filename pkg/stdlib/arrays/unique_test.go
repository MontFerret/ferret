package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestUnique(t *testing.T) {
	Convey("Should return only unique items", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(3),
			values.NewInt(5),
			values.NewInt(6),
			values.NewInt(5),
			values.NewInt(6),
		)

		res, err := arrays.Unique(
			context.Background(),
			arr,
		)

		So(err, ShouldBeNil)

		So(res.String(), ShouldEqual, `[1,2,3,4,5,6]`)
	})
}
