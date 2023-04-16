package core_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func TestNewSourceMap(t *testing.T) {
	Convey("Should match", t, func() {
		s := core.NewSourceMap("test", 1, 100)
		sFmt := fmt.Sprintf("%s at %d:%d", "test", 1, 100)

		So(s, ShouldNotBeNil)

		So(s.Line(), ShouldEqual, 1)

		So(s.Column(), ShouldEqual, 100)

		So(s.String(), ShouldEqual, sFmt)
	})
}
