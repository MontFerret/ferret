package testing_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/testing"
)

func TestFail(t *t.T) {
	False := testing.NewPositive(testing.Fail)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := False(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("It should return an error", t, func() {
		_, err := False(context.Background(), values.False)

		So(err, ShouldBeError)
	})
}
