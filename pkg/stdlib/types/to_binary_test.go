package types

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
)

func TestToBinary(t *testing.T) {
	Convey("TestToBinary", t, func() {
		expected := "some data"

		result, err := ToBinary(context.Background(), values.NewString("some data"))

		So(err, ShouldBeNil)
		So(result.String(), ShouldEqual, expected)
	})
}
