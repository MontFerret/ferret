package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestUnique(t *testing.T) {
	Convey("Should return only unique items", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(3),
			core.NewInt(5),
			core.NewInt(6),
			core.NewInt(5),
			core.NewInt(6),
		)

		res, err := arrays.Unique(
			context.Background(),
			arr,
		)

		So(err, ShouldBeNil)

		So(res.String(), ShouldEqual, `[1,2,3,4,5,6]`)
	})
}
