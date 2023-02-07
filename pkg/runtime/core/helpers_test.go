package core_test

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type DummyStruct struct{}

func TestIsNil(t *testing.T) {
	Convey("Should match", t, func() {
		// nil == invalid
		t := core.IsNil(nil)

		So(t, ShouldBeTrue)

		a := []string{}
		t = core.IsNil(a)

		So(t, ShouldBeFalse)

		b := make([]string, 1)
		t = core.IsNil(b)

		So(t, ShouldBeFalse)

		c := make(map[string]string)
		t = core.IsNil(c)

		So(t, ShouldBeFalse)

		var s struct {
			Test string
		}
		t = core.IsNil(s)

		So(t, ShouldBeFalse)

		f := func() {}
		t = core.IsNil(f)

		So(t, ShouldBeFalse)

		i := DummyStruct{}
		t = core.IsNil(i)

		So(t, ShouldBeFalse)

		ch := make(chan string)
		t = core.IsNil(ch)

		So(t, ShouldBeFalse)

		var y unsafe.Pointer
		var vy int
		y = unsafe.Pointer(&vy)
		t = core.IsNil(y)

		So(t, ShouldBeFalse)
	})
}
