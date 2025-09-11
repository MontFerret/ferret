package runtime_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"errors"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTypeError(t *testing.T) {
	Convey("Should match", t, func() {
		e := runtime.TypeErrorOf(runtime.True, runtime.TypeList)
		So(e, ShouldNotBeNil)

		e = runtime.TypeErrorOf(runtime.True, runtime.TypeList, runtime.TypeString)
		So(e, ShouldNotBeNil)

		cause := errors.New("invalid type: expected foo or bar, but got boolean")
		e = runtime.TypeErrorOf(runtime.True, "foo", "bar")
		So(e.Error(), ShouldEqual, cause.Error())
	})
}

func TestError(t *testing.T) {
	Convey("Should match", t, func() {
		msg := "test message"
		cause := errors.New("cause")
		e := fmt.Errorf("%w: %s", cause, msg)

		ce := runtime.Error(cause, msg)
		So(ce, ShouldNotBeNil)
		So(ce.Error(), ShouldEqual, e.Error())
	})
}
