package testing_test

import (
	"context"
	t "testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestFail(t *t.T) {
	Fail := base.NewPositiveAssertion(testing.Fail)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Fail(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("It should return an error", t, func() {
		_, err := Fail(context.Background(), runtime.False)

		So(err, ShouldBeError)
	})
}
