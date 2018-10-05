package core_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidateArgs(t *testing.T) {
	Convey("Should match", t, func() {
		a := []core.Value{values.NewInt(1), values.NewInt(2)}

		e := core.ValidateArgs(a, 1, 2)
		So(e, ShouldBeNil)

		e = core.ValidateArgs(a, 3, 4)
		So(e, ShouldNotBeNil)
	})
}
