package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestUnique(t *testing.T) {
	Convey("Should return only unique items", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(3),
			runtime.NewInt(5),
			runtime.NewInt(6),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		res, err := arrays.Unique(
			context.Background(),
			arr,
		)

		So(err, ShouldBeNil)

		So(res.String(), ShouldEqual, `[1,2,3,4,5,6]`)
	})
}
