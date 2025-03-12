package arrays_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestRemoveValues(t *testing.T) {
	Convey("Should return a copy of an array without given elements", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
			core.NewInt(6),
		)

		out, err := arrays.RemoveValues(
			context.Background(),
			arr,
			internal.NewArrayWith(
				core.NewInt(3),
				core.NewInt(5),
				core.NewInt(6),
			),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,4]")
	})
}
