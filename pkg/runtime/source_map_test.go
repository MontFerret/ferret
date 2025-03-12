package runtime_test

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSourceMap(t *testing.T) {
	Convey("Should match", t, func() {
		s := runtime.NewSourceMap("test", 1, 100)
		sFmt := fmt.Sprintf("%s at %d:%d", "test", 1, 100)

		So(s, ShouldNotBeNil)

		So(s.Line(), ShouldEqual, 1)

		So(s.Column(), ShouldEqual, 100)

		So(s.String(), ShouldEqual, sFmt)
	})
}
