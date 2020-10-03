package encoder_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/values/encoder"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestJSON(t *testing.T) {
	Convey(".EncodeJSON", t, func() {
		Convey("Should not HTML escape", func() {
			expectedValue := []byte(`"name=Jane&age=38"`)

			marshalled, err := encoder.EncodeJSON("name=Jane&age=38")
			So(err, ShouldBeNil)
			So(string(marshalled), ShouldEqual, string(expectedValue))
		})
	})
}
