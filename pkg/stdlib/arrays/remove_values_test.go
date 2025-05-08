package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestRemoveValues(t *testing.T) {
	Convey("Should return a copy of an array without given elements", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		out, err := arrays.RemoveValues(
			context.Background(),
			arr,
			runtime.NewArrayWith(
				runtime.NewInt(3),
				runtime.NewInt(5),
				runtime.NewInt(6),
			),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,4]")
	})
}
