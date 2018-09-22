package strings_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestReverse(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Reverse(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("Should reverse a text with right encoding", t, func() {
		out, _ := strings.Reverse(
			context.Background(),
			values.NewString("The quick brown 狐 jumped over the lazy 犬"),
		)

		So(out, ShouldEqual, "犬 yzal eht revo depmuj 狐 nworb kciuq ehT")
	})
}
