package types_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/types"
)

func TestToBinary(t *testing.T) {
	Convey("TestToBinary", t, func() {
		value := "abc"

		result, err := types.ToBinary(context.Background(), values.NewString(value))
		So(err, ShouldBeNil)

		wasBinary, err := types.IsBinary(context.Background(), result)
		So(err, ShouldBeNil)
		So(wasBinary, ShouldEqual, values.True)

		So(result.String(), ShouldEqual, value)
	})
}
