package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestRemoveValues(t *testing.T) {
	Convey("Should return a copy of an array without given elements", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
			values.NewInt(6),
		)

		out, err := arrays.RemoveValues(
			context.Background(),
			arr,
			values.NewArrayWith(
				values.NewInt(3),
				values.NewInt(5),
				values.NewInt(6),
			),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,4]")
	})
}
