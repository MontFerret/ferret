package core_test

import (
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func TestSourceError(t *testing.T) {
	Convey("Should match", t, func() {
		sm := core.NewSourceMap("test", 1, 1)

		msg := "test at 1:1"
		cause := errors.New("cause")
		e := errors.Errorf("%s: %s", cause.Error(), msg)

		cse := core.SourceError(sm, cause)
		So(cse, ShouldNotBeNil)
		So(cse.Error(), ShouldEqual, e.Error())
	})
}

func TestTypeError(t *testing.T) {
	Convey("Should match", t, func() {
		e := core.TypeError(TypeA{})
		So(e, ShouldNotBeNil)

		e = core.TypeError(TypeA{}, TypeB{})
		So(e, ShouldNotBeNil)

		cause := errors.New("invalid type: expected type_b or type_c, but got type_a")
		e = core.TypeError(TypeA{}, TypeB{}, TypeC{})
		So(e.Error(), ShouldEqual, cause.Error())
	})
}

func TestError(t *testing.T) {
	Convey("Should match", t, func() {
		msg := "test message"
		cause := errors.New("cause")
		e := errors.Errorf("%s: %s", cause.Error(), msg)

		ce := core.Error(cause, msg)
		So(ce, ShouldNotBeNil)
		So(ce.Error(), ShouldEqual, e.Error())
	})
}
