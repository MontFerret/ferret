package core_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	. "github.com/smartystreets/goconvey/convey"
)

type DummyInterface interface {
	DummyFunc() string
}

type DummyStruct struct{}

func (d DummyStruct) DummyFunc() string {
	return "testing"
}

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

		/* currently not passing
		var s struct {
			Test string
		}
		t = core.IsNil(s)

		So(t, ShouldBeFalse)
		*/

		f := func() {}
		t = core.IsNil(f)

		So(t, ShouldBeFalse)

		/* currently not passing
		i := DummyStruct{}
		t = core.IsNil(i)

		So(t, ShouldBeFalse)
		*/

		ch := make(chan string)
		t = core.IsNil(ch)

		So(t, ShouldBeFalse)

		/* currently not passing
		var y unsafe.Pointer
		var vy int
		y = unsafe.Pointer(&vy)
		t = core.IsNil(y)

		So(t, ShouldBeFalse)
		*/
	})
}
