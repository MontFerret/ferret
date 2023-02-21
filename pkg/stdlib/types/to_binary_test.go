package types

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
)

func TestToBinary(t *testing.T) {
	Convey("TestToBinary", t, func() {
		value := "abc"

		result, err := ToBinary(context.Background(), values.NewString(value))
		So(err, ShouldBeNil)

		wasBinary, err := IsBinary(context.Background(), result)
		So(err, ShouldBeNil)
		So(wasBinary, ShouldEqual, values.True)

		So(result.String(), ShouldEqual, value)
	})
}
