package runtime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTypeError(t *testing.T) {
	Convey("Should match", t, func() {
		e := runtime.TypeError(runtime.True, runtime.TypeList)
		So(e, ShouldNotBeNil)

		e = runtime.TypeError(runtime.True, runtime.TypeList, runtime.TypeString)
		So(e, ShouldNotBeNil)

		cause := errors.New("invalid type: expected foo or bar, but got boolean")
		e = runtime.TypeError(runtime.True, "foo", "bar")
		So(e.Error(), ShouldEqual, cause.Error())
	})
}

func TestError(t *testing.T) {
	Convey("Should match", t, func() {
		msg := "test message"
		cause := errors.New("cause")
		e := errors.Errorf("%s: %s", cause.Error(), msg)

		ce := runtime.Error(cause, msg)
		So(ce, ShouldNotBeNil)
		So(ce.Error(), ShouldEqual, e.Error())
	})
}
